package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"crs-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IPWhitelistManager مدیریت لیست سفید آی‌پی‌ها و به‌روزرسانی خودکار آن‌ها از دیتابیس را بر عهده دارد.
type IPWhitelistManager struct {
	allowedIPs   map[string]bool
	whitelist    []*net.IPNet
	mu           sync.RWMutex
	lastUpdate   time.Time
	lastReload   time.Time
	reloadPeriod time.Duration
	logger       *logrus.Logger
	DB           *gorm.DB
}

// NewIPWhitelistManager ایجاد و مقداردهی اولیه IPWhitelistManager.
// در صورت نبود اتصال به DB، به‌روزرسانی خودکار فعال نخواهد شد.
func NewIPWhitelistManager(allowedIPs []string, logger *logrus.Logger, db *gorm.DB, reloadPeriod time.Duration) *IPWhitelistManager {
	iw := &IPWhitelistManager{
		allowedIPs:   make(map[string]bool),
		whitelist:    make([]*net.IPNet, 0),
		lastUpdate:   time.Now(),
		logger:       logger,
		DB:           db,
		reloadPeriod: reloadPeriod,
	}

	// مقداردهی اولیه آی‌پی‌های مجاز و تبدیل به CIDR
	for _, ip := range allowedIPs {
		iw.allowedIPs[ip] = true
		if !strings.Contains(ip, "/") {
			parsedIP := net.ParseIP(ip)
			if parsedIP != nil {
				if parsedIP.To4() != nil {
					ip = ip + "/32"
				} else {
					ip = ip + "/128"
				}
			}
		}
		_, ipnet, err := net.ParseCIDR(ip)
		if err != nil {
			if logger != nil {
				logger.Warnf("Invalid CIDR format for allowed IP: %s", ip)
			}
			continue
		}
		iw.whitelist = append(iw.whitelist, ipnet)
	}

	// شروع به‌روزرسانی خودکار در صورت وجود DB و reloadPeriod معتبر.
	if reloadPeriod > 0 && db != nil {
		go iw.autoReload()
	}

	return iw
}

// Middleware یک middleware Gin برای بررسی IP درخواست‌کننده بر اساس whitelist فراهم می‌کند.
func (w *IPWhitelistManager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := w.getRealIP(c)
		if !w.isAllowed(clientIP) {
			if w.logger != nil {
				w.logger.WithFields(logrus.Fields{
					"ip":     clientIP,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"client": c.Request.UserAgent(),
				}).Warn("Unauthorized access attempt from IP")
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			return
		}
		c.Next()
	}
}

// getRealIP استخراج آی‌پی واقعی کاربر از هدرهای مختلف است.
func (w *IPWhitelistManager) getRealIP(c *gin.Context) string {
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			parsedIP := net.ParseIP(ip)
			if parsedIP != nil && !parsedIP.IsPrivate() {
				return ip
			}
		}
	}
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		if net.ParseIP(realIP) != nil {
			return realIP
		}
	}
	return c.ClientIP()
}

// isAllowed بررسی می‌کند که آیا آی‌پی داده‌شده در whitelist مجاز است یا خیر.
func (w *IPWhitelistManager) isAllowed(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, network := range w.whitelist {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// reloadWhitelist به‌روزرسانی whitelist از دیتابیس است.
func (w *IPWhitelistManager) reloadWhitelist() error {
	var entries []models.IPWhitelist
	result := w.DB.Where("active = ?", true).Find(&entries)
	if result.Error != nil {
		return fmt.Errorf("error reloading whitelist from database: %v", result.Error)
	}

	newWhitelist := make([]*net.IPNet, 0)
	for _, entry := range entries {
		cidr := fmt.Sprintf("%s/%d", entry.IP, entry.CIDR)
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			if w.logger != nil {
				w.logger.Warnf("Invalid IP record: %s/%d - %v", entry.IP, entry.CIDR, err)
			}
			continue
		}
		newWhitelist = append(newWhitelist, ipnet)
	}

	w.mu.Lock()
	w.whitelist = newWhitelist
	w.lastReload = time.Now()
	w.mu.Unlock()

	if w.logger != nil {
		w.logger.Infof("IP whitelist updated with %d valid records", len(newWhitelist))
	}

	return nil
}

// autoReload به‌روزرسانی خودکار whitelist در فواصل زمانی مشخص است.
func (w *IPWhitelistManager) autoReload() {
	ticker := time.NewTicker(w.reloadPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.reloadWhitelist(); err != nil && w.logger != nil {
				w.logger.Errorf("Error auto-reloading IP whitelist: %v", err)
			}
		}
	}
}

// RouteSpecificMiddleware اعمال محدودیت‌های IP به‌صورت اختصاصی برای مسیرهای حساس است.
func (w *IPWhitelistManager) RouteSpecificMiddleware(allowedCIDRs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := w.getRealIP(c)
		valid := false
		for _, cidr := range allowedCIDRs {
			_, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if ipnet.Contains(net.ParseIP(clientIP)) {
				valid = true
				break
			}
		}

		if !valid {
			if w.logger != nil {
				w.logger.WithFields(logrus.Fields{
					"ip":      clientIP,
					"route":   c.FullPath(),
					"allowed": allowedCIDRs,
				}).Error("Access denied due to IP restriction")
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access restricted on this route",
				"code":  "ROUTE_IP_RESTRICTED",
			})
			return
		}
		c.Next()
	}
}

package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SecurityHeadersConfig تنظیمات مربوط به هدرهای امنیتی را نگه می‌دارد.
type SecurityHeadersConfig struct {
	CSPDirectives         map[string][]string
	FeaturePolicy         map[string][]string
	PermissionsPolicy     map[string][]string
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	HSTSPreload           bool
	FrameOptions          string
	ContentTypeOptions    string
	XSSProtection         string
	ReferrerPolicy        string
	ExpectCTMaxAge        int
	ExpectCTEnforce       bool
	ReportURI             string
	Logger                *logrus.Logger
}

// SecurityHeadersManager مدیریت هدرهای امنیتی و اعمال آن‌ها بر روی درخواست‌ها است.
type SecurityHeadersManager struct {
	config         *SecurityHeadersConfig
	cspNonce       string
	nonceGenerator func() string
}

// NewSecurityHeadersManager ایجاد و مقداردهی اولیه SecurityHeadersManager.
func NewSecurityHeadersManager(config *SecurityHeadersConfig) *SecurityHeadersManager {
	if config == nil {
		config = &SecurityHeadersConfig{}
	}
	sh := &SecurityHeadersManager{
		config:         config,
		nonceGenerator: generateNonce,
	}
	// تنظیم پیش‌فرض برای CSP در صورت نبود تنظیمات
	if len(config.CSPDirectives) == 0 {
		config.CSPDirectives = map[string][]string{
			"default-src": {"'self'"},
			"script-src":  {"'self'", "'unsafe-inline'", "https://cdn.example.com"},
			"style-src":   {"'self'", "'unsafe-inline'"},
			"img-src":     {"'self'", "data:", "https://*.example.com"},
			"connect-src": {"'self'", "https://api.example.com"},
			"font-src":    {"'self'", "https://fonts.gstatic.com"},
			"object-src":  {"'none'"},
			"frame-src":   {"'none'"},
			"report-uri":  {config.ReportURI},
		}
	}
	return sh
}

// Apply یک middleware Gin برای اعمال هدرهای امنیتی فراهم می‌کند.
func (sh *SecurityHeadersManager) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := sh.nonceGenerator()
		sh.cspNonce = nonce

		sh.applyHeaders(c, nonce)

		// جلوگیری از کش شدن پاسخ‌های حساس برای مسیرهای API
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Header("Cache-Control", "no-store, max-age=0")
		}

		c.Next()

		// پاکسازی nonce پس از پردازش درخواست
		sh.cspNonce = ""
	}
}

// GetNonce دریافت nonce برای استفاده در templateها.
func (sh *SecurityHeadersManager) GetNonce() string {
	return sh.cspNonce
}

// applyHeaders هدرها را بر اساس تنظیمات در پاسخ اعمال می‌کند.
func (sh *SecurityHeadersManager) applyHeaders(c *gin.Context, nonce string) {
	// Content Security Policy (CSP)
	cspHeader := sh.buildCSPHeader(nonce)
	c.Header("Content-Security-Policy", cspHeader)
	c.Header("Content-Security-Policy-Report-Only", "")

	// Strict-Transport-Security (HSTS)
	if sh.config.HSTSMaxAge > 0 {
		hstsValue := fmt.Sprintf("max-age=%d", sh.config.HSTSMaxAge)
		if sh.config.HSTSIncludeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if sh.config.HSTSPreload {
			hstsValue += "; preload"
		}
		c.Header("Strict-Transport-Security", hstsValue)
	}

	// سایر هدرهای امنیتی
	c.Header("X-Frame-Options", sh.config.FrameOptions)
	c.Header("X-Content-Type-Options", sh.config.ContentTypeOptions)
	c.Header("X-XSS-Protection", sh.config.XSSProtection)
	c.Header("Referrer-Policy", sh.config.ReferrerPolicy)
	c.Header("Permissions-Policy", sh.buildPermissionsPolicyHeader())
	c.Header("Expect-CT", sh.buildExpectCTHeader())

	// هدرهای مربوط به Cross-Origin
	c.Header("Cross-Origin-Embedder-Policy", "require-corp")
	c.Header("Cross-Origin-Opener-Policy", "same-origin")
	c.Header("Cross-Origin-Resource-Policy", "same-site")
}

// buildCSPHeader ساخت رشته CSP با اضافه کردن nonce به script-src.
func (sh *SecurityHeadersManager) buildCSPHeader(nonce string) string {
	directives := []string{}

	// اضافه کردن nonce به script-src در صورت وجود
	if sources, ok := sh.config.CSPDirectives["script-src"]; ok {
		sh.config.CSPDirectives["script-src"] = append(sources, fmt.Sprintf("'nonce-%s'", nonce))
	}

	for directive, sources := range sh.config.CSPDirectives {
		directives = append(directives,
			fmt.Sprintf("%s %s", directive, strings.Join(sources, " ")),
		)
	}

	return strings.Join(directives, "; ")
}

// buildPermissionsPolicyHeader ساخت رشته Permissions-Policy.
func (sh *SecurityHeadersManager) buildPermissionsPolicyHeader() string {
	policies := []string{}
	for feature, origins := range sh.config.PermissionsPolicy {
		policies = append(policies,
			fmt.Sprintf("%s=(%s)", feature, strings.Join(origins, " ")),
		)
	}
	return strings.Join(policies, ", ")
}

// buildExpectCTHeader ساخت رشته Expect-CT.
func (sh *SecurityHeadersManager) buildExpectCTHeader() string {
	if sh.config.ExpectCTMaxAge == 0 {
		return ""
	}
	value := fmt.Sprintf("max-age=%d", sh.config.ExpectCTMaxAge)
	if sh.config.ExpectCTEnforce {
		value += ", enforce"
	}
	if sh.config.ReportURI != "" {
		value += fmt.Sprintf(", report-uri=\"%s\"", sh.config.ReportURI)
	}
	return value
}

// generateNonce تولید یک رشته تصادفی (nonce) برای امنیت CSP است.
func generateNonce() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// ReportHandler هندلر دریافت گزارش تخلفات CSP است.
func (sh *SecurityHeadersManager) ReportHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var report struct {
			CSPReport struct {
				DocumentURI        string `json:"document-uri"`
				Referrer           string `json:"referrer"`
				ViolatedDirective  string `json:"violated-directive"`
				EffectiveDirective string `json:"effective-directive"`
				OriginalPolicy     string `json:"original-policy"`
				BlockedURI         string `json:"blocked-uri"`
				StatusCode         int    `json:"status-code"`
			} `json:"csp-report"`
		}

		if err := c.BindJSON(&report); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid report format"})
			return
		}

		if sh.config.Logger != nil {
			sh.config.Logger.WithFields(logrus.Fields{
				"document_uri":        report.CSPReport.DocumentURI,
				"violated_directive":  report.CSPReport.ViolatedDirective,
				"blocked_uri":         report.CSPReport.BlockedURI,
				"status_code":         report.CSPReport.StatusCode,
				"effective_directive": report.CSPReport.EffectiveDirective,
				"client_ip":           c.ClientIP(),
				"user_agent":          c.Request.UserAgent(),
			}).Warn("CSP Policy Violation")
		}

		c.Status(http.StatusNoContent)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
	"crs-backend/internal/models"
	"crs-backend/internal/utils"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
)

type TwoFactorHandler struct {
	DB         *gorm.DB
	JWTSecret  string
	Env        string
}

// ----------- توابع اصلی -----------

func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "دسترسی غیرمجاز"})
		return
	}

	currentUser := user.(models.User)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CRS Backend",
		AccountName: currentUser.Email,
		Period:      30,
		Digits:      6,
		//Algorithm:   totp.AlgorithmSHA1, // تغییر به SHA1 برای نسخه 1.4.0
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در تولید کلید امنیتی"})
		return
	}

	tx := h.DB.Model(&currentUser).Updates(map[string]interface{}{
		"two_factor_secret":   key.Secret(),
		"two_factor_enabled":  false,
	})

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در بروزرسانی اطلاعات کاربر"})
		return
	}

	recoveryCodes := generateRecoveryCodes(8)
	recoveryCodesJSON, _ := json.Marshal(recoveryCodes)

	h.DB.Model(&currentUser).Update("recovery_codes", string(recoveryCodesJSON))

	c.JSON(http.StatusOK, gin.H{
		"qr_code_url":    key.URL(),
		"manual_code":    key.Secret(),
		"recovery_codes": recoveryCodes,
		"message":        "اسکن QR Code یا وارد کردن دستی کلید در برنامه احراز هویت",
	})
}

func (h *TwoFactorHandler) Verify2FA(c *gin.Context) {
	var input struct {
		Code string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت داده نامعتبر"})
		return
	}

	user, _ := c.Get("current_user")
	currentUser := user.(models.User)

	if !totp.Validate(input.Code, currentUser.TwoFactorSecret) {
		go h.logFailed2FAAttempt(currentUser.ID, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "کد احراز نامعتبر"})
		return
	}

	h.DB.Model(&currentUser).Update("two_factor_enabled", true)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    currentUser.ID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"2fa_valid":  true,
		"session_id": generateSessionID(),
	})

	tokenString, _ := token.SignedString([]byte(h.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenString,
		"token_type":   "bearer",
		"expires_in":   86400,
		"message":      "احراز هویت دو مرحله‌ای با موفقیت فعال شد",
	})
}

func (h *TwoFactorHandler) LoginWith2FA(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Code     string `json:"code,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت داده نامعتبر"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "اطلاعات ورود نامعتبر"})
		return
	}

	if !checkPasswordHash(input.Password, user.PasswordHash) { // اصلاح تابع
		c.JSON(http.StatusUnauthorized, gin.H{"error": "اطلاعات ورود نامعتبر"})
		return
	}

	if user.TwoFactorEnabled {
		if input.Code == "" {
			tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": user.ID,
				"exp":     time.Now().Add(5 * time.Minute).Unix(),
				"2fa":     true,
			})

			tempTokenString, _ := tempToken.SignedString([]byte(h.JWTSecret))

			c.JSON(http.StatusOK, gin.H{
				"2fa_required":  true,
				"temp_token":    tempTokenString,
				"message":       "نیاز به کد احراز دو مرحله‌ای",
			})
			return
		}

		if !totp.Validate(input.Code, user.TwoFactorSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "کد احراز نامعتبر"})
			return
		}
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"2fa_valid":  user.TwoFactorEnabled,
		"session_id": generateSessionID(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(h.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   "bearer",
		"expires_in":   86400,
	})
}

// ----------- توابع کمکی -----------

func checkPasswordHash(password, hash string) bool { // اضافه شد
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateRecoveryCodes(count int) []string { // پیاده‌سازی واقعی
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		codes[i] = utils.GenerateRandomString(12) // نیاز به پیاده‌سازی در utils
	}
	return codes
}

func generateSessionID() string { // اضافه شد
	return utils.GenerateUUID() // نیاز به پیاده‌سازی در utils
}

func (h *TwoFactorHandler) logFailed2FAAttempt(userID uint, ip string) {
	h.DB.Create(&models.SecurityLog{
		UserID:    userID,
		Action:    "failed_2fa",
		IPAddress: ip,
		CreatedAt: time.Now(),
	})
}

// ----------- ساختارهای داده -----------
type Verify2FAInput struct {
	Code string `json:"code" binding:"required,len=6"`
}

type Login2FAInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Code     string `json:"code,omitempty"`
}

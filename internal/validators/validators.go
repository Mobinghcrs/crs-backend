// internal/validators/validators.go
package validators

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	validate          *validator.Validate
	db                *sql.DB
	logger            *logrus.Logger
	translations      map[string]map[string]string
	sanitizeRegex     = regexp.MustCompile(`[<>"'%();\\]`)
	phoneRegex        = regexp.MustCompile(`^\+98\d{10}$`)
	usernameRegex     = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{5,20}$`)
	persianAlphaRegex = regexp.MustCompile("^[\u0600-\u06FF\\s]+$")
)

// InitValidators مقداردهی اولیه سیستم اعتبارسنجی
func InitValidators(database *sql.DB, log *logrus.Logger) {
	db = database
	logger = log

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		
		// ثبت اعتبارسنج‌های سفارشی
		_ = validate.RegisterValidation("strong_password", validateStrongPassword)
		_ = validate.RegisterValidation("persian_alpha", validatePersianAlpha)
		_ = validate.RegisterValidation("unique", validateUnique)
		_ = validate.RegisterValidation("secure_string", validateSecureString)
		_ = validate.RegisterValidation("phone", validatePhone)
		_ = validate.RegisterValidation("username", validateUsername)
		_ = validate.RegisterValidation("future_date", validateFutureDate)
	}

	initTranslations()
}

// ValidateRequest اعتبارسنجی درخواست با مدیریت خطاهای پیشرفته
func ValidateRequest(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		HandleValidationErrors(c, err)
		return false
	}
	return true
}

// HandleValidationErrors مدیریت یکپارچه خطاهای اعتبارسنجی
func HandleValidationErrors(c *gin.Context, err error) {
	var errors []ValidationError

	if verr, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verr {
			errors = append(errors, ValidationError{
				Field:   fe.Field(),
				Message: getTranslatedMessage(fe.Tag(), fe.Param(), c.GetHeader("Accept-Language")),
				Code:    "VALIDATION_ERROR",
			})
		}
	} else {
		errors = append(errors, ValidationError{
			Message: "Invalid request structure",
			Code:    "INVALID_REQUEST",
		})
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"errors":  errors,
		"success": false,
	})
}

// ----------- اعتبارسنج‌های سفارشی -----------

// اعتبارسنجی قدرت رمز عبور
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial && len(password) >= 12
}

// اعتبارسنجی رشته فارسی
func validatePersianAlpha(fl validator.FieldLevel) bool {
	return persianAlphaRegex.MatchString(fl.Field().String())
}

// اعتبارسنجی یکتایی در دیتابیس
func validateUnique(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ";")
	if len(params) != 2 {
		return false
	}

	table := params[0]
	column := params[1]

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", table, column)
	var count int
	err := db.QueryRowContext(context.Background(), query, fl.Field().String()).Scan(&count)
	if err != nil {
		logger.Errorf("خطا در اعتبارسنجی یکتایی: %v", err)
		return false
	}

	return count == 0
}

// اعتبارسنجی امنیت رشته
func validateSecureString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !sanitizeRegex.MatchString(value)
}

// اعتبارسنجی شماره تلفن ایرانی
func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// اعتبارسنجی نام کاربری
func validateUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

// اعتبارسنجی تاریخ آینده
func validateFutureDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

// ----------- توابع سودمند عمومی -----------

// ValidatePasswordStrength بررسی قدرت رمز عبور به صورت مستقیم
func ValidatePasswordStrength(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial && len(password) >= 12
}

// ValidateNationalID اعتبارسنجی کد ملی ایران
func ValidateNationalID(nid string) bool {
	if len(nid) != 10 {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		num := int(nid[i] - '0')
		sum += num * (10 - i)
	}

	rem := sum % 11
	checkDigit := int(nid[9] - '0')

	return (rem < 2 && checkDigit == rem) || (rem >= 2 && checkDigit == (11-rem))
}

// ----------- ترجمه پیام‌های خطا -----------

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func initTranslations() {
	translations = map[string]map[string]string{
		"en": {
			"required":         "Field is required",
			"email":            "Invalid email format",
			"strong_password":  "Password must contain at least 12 characters with uppercase, lowercase, number and special character",
			"unique":           "This value already exists",
			"secure_string":    "Invalid characters detected",
			"phone":            "Invalid Iranian phone number format. Example: +989123456789",
			"username":         "Username must be 5-20 characters and can contain letters, numbers, _ - .",
			"future_date":      "Date must be in the future",
			"persian_alpha":    "Must contain only Persian characters",
		},
		"fa": {
			"required":         "این فیلد اجباری است",
			"email":            "فرمت ایمیل نامعتبر است",
			"strong_password":  "رمز عبور باید حداقل ۱۲ کاراکتر با حروف بزرگ، کوچک، عدد و کاراکتر ویژه باشد",
			"unique":           "این مقدار قبلاً ثبت شده است",
			"secure_string":    "کاراکترهای غیرمجاز وجود دارد",
			"phone":            "فرمت شماره تلفن نامعتبر. مثال: +989123456789",
			"username":         "نام کاربری باید ۵ تا ۲۰ کاراکتر و شامل حروف، اعداد، _ - . باشد",
			"future_date":      "تاریخ باید در آینده باشد",
			"persian_alpha":    "فقط حروف فارسی مجاز است",
		},
	}
}

func getTranslatedMessage(tag, param, lang string) string {
	lang = normalizeLang(lang)
	msgKey := translations[lang][tag]
	
	p := message.NewPrinter(language.Make(lang))
	
	switch tag {
	case "min":
		return p.Sprintf("حداقل طول مجاز: %s کاراکتر", param)
	case "max":
		return p.Sprintf("حداکثر طول مجاز: %s کاراکتر", param)
	default:
		return msgKey
	}
}

func normalizeLang(lang string) string {
	if strings.HasPrefix(lang, "fa") {
		return "fa"
	}
	return "en"
}

// ----------- توابع سودمند عمومی -----------

func SanitizeInput(input string) string {
	return sanitizeRegex.ReplaceAllString(input, "")
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ----------- Middleware اعتبارسنجی خودکار -----------
func ValidationMiddleware(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			HandleValidationErrors(c, err)
			return
		}
		c.Set("validatedRequest", model)
		c.Next()
	}
}

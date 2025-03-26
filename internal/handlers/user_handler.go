package handlers

import (
    "crs-backend/internal/database"
    "crs-backend/internal/models"
    "crs-backend/internal/utils"
    "crs-backend/internal/repositories"
    "net/http"
    "strconv"
    "errors"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    DB *gorm.DB
    UserRepo repositories.IUserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{
        DB: db,
        UserRepo: repositories.NewUserRepository(db),
    }
}
func Register(c *gin.Context) {
    type RegisterInput struct {
        Username string `json:"username" binding:"required,min=3"`
        Password string `json:"password" binding:"required,min=8"`
        FullName string `json:"full_name" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
    }

    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "خطای اعتبارسنجی: " + err.Error()})
        return
    }

    var existing models.User
    if result := database.DB.Where("email = ? OR username = ?", input.Email, input.Username).First(&existing); result.Error == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "ایمیل یا نام کاربری تکراری است"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در پردازش رمز عبور"})
        return
    }

    user := models.User{
        Username:     input.Username,
        PasswordHash: string(hashedPassword),
        FullName:     input.FullName,
        Email:        input.Email,
        Role:         "user",
    }

    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ذخیره کاربر"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "ثبت‌نام موفق",
        "user_id": strconv.FormatUint(uint64(user.ID), 10),
    })
}

func Login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=8"`
    }

    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های نامعتبر"})
        return
    }

    var user models.User
    if err := database.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "اطلاعات ورودی نامعتبر"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "اطلاعات ورودی نامعتبر"})
        return
    }

    token, err := utils.GenerateJWT(user.ID, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در تولید توکن"})
        return
    }

    if errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "کاربری با این ایمیل وجود ندارد"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "role":     user.Role,
        },
    })
}

// GetUsers - دریافت لیست کاربران
func (h *UserHandler) GetUsers(c *gin.Context) {
    users, err := h.UserRepo.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت اطلاعات کاربران"})
        return
    }
    c.JSON(http.StatusOK, users)
}

// GetUser - دریافت کاربر بر اساس ID
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.UserRepo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "کاربر یافت نشد"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// DeleteUser - حذف کاربر
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    if err := h.UserRepo.Delete(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف کاربر"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "کاربر با موفقیت حذف شد"})
}

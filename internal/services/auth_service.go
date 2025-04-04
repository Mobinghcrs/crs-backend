package services

import (
	"booking-system/internal/models"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	SignUp(user *models.User) error
	Login(email string, password string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: db}
}

func (s *authService) SignUp(user *models.User) error {
	var existingUser models.User
	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("ایمیل قبلاً ثبت شده است")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), 
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("خطا در پردازش رمز عبور")
	}
	user.Password = string(hashedPassword)

	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.db.Create(user).Error; err != nil {
		return fmt.Errorf("خطا در ایجاد کاربر: %v", err)
	}

	return nil
}

func (s *authService) Login(email string, password string) (*models.User, error) {
	var user models.User

	// جستجوی کاربر با ایمیل (Case-Insensitive)
	if err := s.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("کاربری با این ایمیل وجود ندارد")
		}
		return nil, fmt.Errorf("خطا در جستجوی کاربر: %v", err)
	}

	// دیباگ: چاپ هش و رمز ورودی
	fmt.Printf("\nرمز ذخیرهشده: %s\nرمز ورودی: %s\n", user.Password, password)

	// مقایسه رمز عبور
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("رمز عبور نادرست است")
	}

	user.Password = ""
	return &user, nil
}

func (s *authService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Select("id", "email", "role", "created_at", "updated_at").Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت کاربران: %v", err)
	}
	return users, nil
}
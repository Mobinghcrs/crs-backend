package vault

import (
	"fmt"
)

// VaultService اینترفیس سرویس Vault
type VaultService interface {
	GetSecret(key string) string
	Close()
}

// vaultClient پیاده‌سازی ساده از VaultService
type vaultClient struct {
	data map[string]string
}

// NewVaultClient ایجاد کلاینت Vault جدید
func NewVaultClient(addr string, opts ...func(*vaultClient)) (VaultService, error) {
	v := &vaultClient{
		data: map[string]string{
			"DB_PASSWORD":     "example-db-password",
			"REDIS_PASSWORD":  "example-redis-password",
			"JWT_SECRET":      "example-jwt-secret",
			"JWT_REFRESH_SECRET": "example-refresh-secret",
		},
	}
	for _, opt := range opts {
		opt(v)
	}
	fmt.Println("Vault client initialized successfully!")
	return v, nil
}

// WithToken تابع اختیاری برای تنظیم توکن در کلاینت Vault
func WithToken(token string) func(*vaultClient) {
	return func(v *vaultClient) {
		fmt.Println("Vault token set:", token)
	}
}

// GetSecret دریافت مقدار متغیر از Vault
func (v *vaultClient) GetSecret(key string) string {
	if value, exists := v.data[key]; exists {
		return value
	}
	return ""
}

// Close بستن ارتباط با Vault
func (v *vaultClient) Close() {
	fmt.Println("Vault client closed")
}

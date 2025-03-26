// internal/config/vault.go
package config

import (
	"fmt"
	"os"
	"time"
	

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type VaultClient struct {
	client *api.Client
	config *api.Config
	logger *logrus.Logger
}

// NewVaultClient ایجاد کلاینت Vault با پیکربندی پیشرفته
func NewVaultClient(logger *logrus.Logger) (*VaultClient, error) {
	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	if config.Address == "" {
		config.Address = "http://localhost:8200"
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("خطا در ایجاد کلاینت Vault: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))
	if client.Token() == "" {
		return nil, fmt.Errorf("توکن Vault تنظیم نشده است")
	}

	return &VaultClient{
		client: client,
		config: config,
		logger: logger,
	}, nil
}

// GetSecret دریافت راز با قابلیت کشینگ و بازنشانی خودکار توکن
func (v *VaultClient) GetSecret(path, key string) (string, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("خطا در خواندن راز از Vault: %v", err)
	}

	if secret == nil || secret.Data["data"] == nil {
		return "", fmt.Errorf("راز مورد نظر یافت نشد: %s/%s", path, key)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("فرمت داده‌های راز نامعتبر است")
	}

	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("کلید مورد نظر در راز وجود ندارد: %s", key)
	}

	return fmt.Sprintf("%v", value), nil
}

// GetDynamicDBCredential دریافت اعتبارنامه داینامیک دیتابیس
func (v *VaultClient) GetDynamicDBCredential(role string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("database/creds/%s", role))
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت اعتبارنامه دیتابیس: %v", err)
	}

	if secret == nil || secret.LeaseID == "" {
		return nil, fmt.Errorf("اعتبارنامه دیتابیس معتبر دریافت نشد")
	}

	// مدیریت Lease خودکار
	go v.renewLease(secret.LeaseID, secret.LeaseDuration)

	return map[string]interface{}{
		"username":     secret.Data["username"],
		"password":     secret.Data["password"],
		"lease_id":     secret.LeaseID,
		"lease_duration": secret.LeaseDuration,
	}, nil

}

func (v *VaultClient) renewLease(leaseID string, duration int) {
    ticker := time.NewTicker(time.Duration(duration/2) * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        _, err := v.client.Sys().Renew(leaseID, duration)
        if err != nil {
            v.logger.Errorf("خطا در تمدید Lease: %v", err)
            return
        }
        v.logger.Info("Lease با موفقیت تمدید شد")
    }
}
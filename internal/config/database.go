package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger" // اضافه کردن این خط برای لاگر
)

func ConnectDB() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres password=mobin2005G dbname=booking port=5432 sslmode=disable"
    // اطمینان از صحت پرانتزها و ویرگول‌ها
    db, err := gorm.Open(
        postgres.Open(dsn),
        &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info), // بدون ویرگول اضافی
        },
    )

    if err != nil {
        return nil, err
    }

    return db, nil // بدون ویرگول اضافی
}
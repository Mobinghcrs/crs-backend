package models

import "time"

type IPWhitelist struct {
    ID          uint      `gorm:"primaryKey"`
    IP          string    `gorm:"type:varchar(45);not null"`  // IPv4/IPv6
    CIDR        int       `gorm:"type:int;default:32"`        // 32 برای IPv4، 128 برای IPv6
    Description string    `gorm:"type:text"`
    Active      bool      `gorm:"default:true"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

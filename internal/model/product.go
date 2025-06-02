package model

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	Stock       int64     `gorm:"not null"` // 当前库存
	StartTime   time.Time `gorm:"not null"` // 秒杀开始时间
	EndTime     time.Time `gorm:"not null"` // 秒杀结束时间
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

package model

import "time"

type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Status    string  `gorm:"type:varchar(20);default:'created'"` // created, success, failed
	CreatedAt time.Time

	// 外键关联（可选）
	User    User
	Product Product
}

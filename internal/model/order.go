package model

import "time"

type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Status    string  `gorm:"type:varchar(20);default:'created'"` // created, success, failed
	CreatedAt time.Time
<<<<<<< HEAD
}

type OrderRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Timestamp int64  `json:"timestamp"` // 抢购时间戳
=======

	// 外键关联（可选）
	User    User
	Product Product
>>>>>>> c50d88e66b527ce124cdb371815bc2fd6947aba9
}

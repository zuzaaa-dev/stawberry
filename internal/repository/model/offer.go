package model

import (
	"time"
)

type Offer struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    uint
	ProductID uint
	StoreID   uint
	Price     float64
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Product   Product `gorm:"foreignKey:ProductID"`
	Store     Store   `gorm:"foreignKey:StoreID"`
}

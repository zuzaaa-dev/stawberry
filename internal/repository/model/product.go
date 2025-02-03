package model

import "time"

type Product struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	StoreID     uint
	Name        string
	Description string
	Price       float64
	Category    string
	InStock     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Store       Store `gorm:"foreignKey:StoreID"`
}

type UpdateProduct struct {
	StoreID     *uint    `gorm:"column:store_id"`
	Name        *string  `gorm:"column:name"`
	Description *string  `gorm:"column:description"`
	Price       *float64 `gorm:"column:price"`
	Category    *string  `gorm:"column:category"`
	InStock     *bool    `gorm:"column:in_stock"`
}

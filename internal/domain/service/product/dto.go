package product

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	StoreID     uint      `json:"store_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	InStock     bool      `json:"in_stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProduct struct {
	StoreID     *uint    `json:"store_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	InStock     *bool    `json:"in_stock,omitempty"`
}

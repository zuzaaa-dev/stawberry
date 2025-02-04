package model

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
)

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

func ConvertProductFromSvc(p product.Product) Product {
	return Product{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		InStock:     p.InStock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ConvertProductToEntity(p Product) entity.Product {
	return entity.Product{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		InStock:     p.InStock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ConvertUpdateProductFromSvc(up product.UpdateProduct) UpdateProduct {
	return UpdateProduct{
		StoreID:     up.StoreID,
		Name:        up.Name,
		Description: up.Description,
		Price:       up.Price,
		Category:    up.Category,
		InStock:     up.InStock,
	}
}

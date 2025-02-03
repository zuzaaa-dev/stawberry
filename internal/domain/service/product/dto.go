package product

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
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

func (p *Product) ConvertToRepo() model.Product {
	return model.Product{
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

type UpdateProduct struct {
	StoreID     *uint    `json:"store_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	InStock     *bool    `json:"in_stock,omitempty"`
}

func (p *UpdateProduct) ConvertToRepo() model.UpdateProduct {
	return model.UpdateProduct{
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		InStock:     p.InStock,
	}
}

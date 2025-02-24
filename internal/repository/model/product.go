package model

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
)

type Product struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	CategoryId  int `gorm:"foreignKey:CategoryId"`
	Description string
}

type UpdateProduct struct {
	Name        string `gorm:"column:name"`
	CategoryId  int    `gorm:"column:category_id"`
	Description string `gorm:"column:description"`
}

func ConvertProductFromSvc(p product.Product) Product {
	return Product{
		ID:          p.ID,
		Name:        p.Name,
		CategoryId:  p.CategoryId,
		Description: p.Description,
	}
}

func ConvertProductToEntity(p Product) entity.Product {
	return entity.Product{
		ID:          p.ID,
		Name:        p.Name,
		CategoryId:  p.CategoryId,
		Description: p.Description,
	}
}

func ConvertUpdateProductFromSvc(up product.UpdateProduct) UpdateProduct {
	return UpdateProduct{
		Name:        up.Name,
		CategoryId:  up.CategoryId,
		Description: up.Description,
	}
}

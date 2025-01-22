package repository

import "marketplace/models"

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
	GetProducts(offset, limit int) ([]models.Product, int, error)
	GetStoreProducts(id string, offset, limit int) ([]models.Product, int, error)
	UpdateProduct(id string, update *models.ProductUpdate) error
}

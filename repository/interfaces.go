package repository

import "marketplace/models"

type OfferRepository interface {
	CreateOffer(offer models.Offer) (uint, error)
	GetUserOffers(userID uint, limit, offset int) ([]models.Offer, int64, error)
	GetOffer(offerID uint) (*models.Offer, error)
	UpdateOfferStatus(offerID uint, status string) (*models.Offer, error)
	DeleteOffer(offerID uint) (*models.Offer, error)

type ProductRepositoryInf interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
	GetProducts(offset, limit int) ([]models.Product, int, error)
	GetStoreProducts(id string, offset, limit int) ([]models.Product, int, error)
	UpdateProduct(id string, update *models.ProductUpdate) error
}

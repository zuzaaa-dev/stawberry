package repository

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
)

type OfferRepository interface {
	CreateOffer(offer model.Offer) (uint, error)
	GetUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error)
	GetOffer(offerID uint) (entity.Offer, error)
	UpdateOfferStatus(offerID uint, status string) (entity.Offer, error)
	DeleteOffer(offerID uint) (entity.Offer, error)
}

type ProductRepository interface {
	CreateProduct(product model.Product) (uint, error)
	GetProductByID(id string) (entity.Product, error)
	GetProducts(offset, limit int) ([]entity.Product, int, error)
	GetStoreProducts(id string, offset, limit int) ([]entity.Product, int, error)
	UpdateProduct(id string, update model.UpdateProduct) error
}

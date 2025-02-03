package service

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
)

type OfferService interface {
	CreateOffer(offer offer.Offer) (uint, error)
	GetUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error)
	GetOffer(offerID uint) (entity.Offer, error)
	UpdateOfferStatus(offerID uint, status string) (entity.Offer, error)
	DeleteOffer(offerID uint) (entity.Offer, error)
}

type ProductService interface {
	CreateProduct(product product.Product) (uint, error)
	GetProductByID(id string) (entity.Product, error)
	GetProducts(offset, limit int) ([]entity.Product, int, error)
	GetStoreProducts(id string, offset, limit int) ([]entity.Product, int, error)
	UpdateProduct(id string, updateProduct product.UpdateProduct) error
}

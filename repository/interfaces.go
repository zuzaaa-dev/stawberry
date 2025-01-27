package repository

import "marketplace/models"

type OfferRepository interface {
	CreateOffer(offer models.Offer) (uint, error)
	GetUserOffers(userID uint, limit, offset int) ([]models.Offer, int64, error)
	GetOffer(offerID uint) (*models.Offer, error)
	UpdateOfferStatus(offerID uint, status string) (*models.Offer, error)
	DeleteOffer(offerID uint) (*models.Offer, error)
}

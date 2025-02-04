package offer

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type Repository interface {
	InsertOffer(offer Offer) (uint, error)
	GetOfferByID(offerID uint) (entity.Offer, error)
	SelectUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error)
	UpdateOfferStatus(offerID uint, status string) (entity.Offer, error)
	DeleteOffer(offerID uint) (entity.Offer, error)
}

type offerService struct {
	offerRepository Repository
}

func NewOfferService(offerRepository Repository) *offerService {
	return &offerService{offerRepository: offerRepository}
}

func (os *offerService) CreateOffer(offer Offer) (uint, error) {
	return os.offerRepository.InsertOffer(offer)
}

func (os *offerService) GetOffer(offerID uint) (entity.Offer, error) {
	return os.offerRepository.GetOfferByID(offerID)
}

func (os *offerService) GetUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error) {
	return os.offerRepository.SelectUserOffers(userID, limit, offset)
}

func (os *offerService) UpdateOfferStatus(offerID uint, status string) (entity.Offer, error) {
	return os.offerRepository.UpdateOfferStatus(offerID, status)
}

func (os *offerService) DeleteOffer(offerID uint) (entity.Offer, error) {
	return os.offerRepository.DeleteOffer(offerID)
}

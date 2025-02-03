package offer

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository"
)

type offerService struct {
	offerRepository repository.OfferRepository
}

func NewOfferService(offerRepository repository.OfferRepository) *offerService {
	return &offerService{offerRepository: offerRepository}
}

func (os *offerService) CreateOffer(offer Offer) (uint, error) {
	return os.offerRepository.CreateOffer(offer.ConvertToRepo())
}

func (os *offerService) GetUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error) {
	return os.offerRepository.GetUserOffers(userID, limit, offset)
}

func (os *offerService) GetOffer(offerID uint) (entity.Offer, error) {
	return os.offerRepository.GetOffer(offerID)
}

func (os *offerService) UpdateOfferStatus(offerID uint, status string) (entity.Offer, error) {
	return os.offerRepository.UpdateOfferStatus(offerID, status)
}

func (os *offerService) DeleteOffer(offerID uint) (entity.Offer, error) {
	return os.offerRepository.DeleteOffer(offerID)
}

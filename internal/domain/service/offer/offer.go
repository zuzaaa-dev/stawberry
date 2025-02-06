package offer

import (
	"context"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type Repository interface {
	InsertOffer(ctx context.Context, offer Offer) (uint, error)
	GetOfferByID(ctx context.Context, offerID uint) (entity.Offer, error)
	SelectUserOffers(ctx context.Context, userID uint, limit, offset int) ([]entity.Offer, int64, error)
	UpdateOfferStatus(ctx context.Context, offerID uint, status string) (entity.Offer, error)
	DeleteOffer(ctx context.Context, offerID uint) (entity.Offer, error)
}

type offerService struct {
	offerRepository Repository
}

func NewOfferService(offerRepository Repository) *offerService {
	return &offerService{offerRepository: offerRepository}
}

func (os *offerService) CreateOffer(ctx context.Context, offer Offer) (uint, error) {
	return os.offerRepository.InsertOffer(ctx, offer)
}

func (os *offerService) GetOffer(ctx context.Context, offerID uint) (entity.Offer, error) {
	return os.offerRepository.GetOfferByID(ctx, offerID)
}

func (os *offerService) GetUserOffers(ctx context.Context, userID uint, limit, offset int) ([]entity.Offer, int64, error) {
	return os.offerRepository.SelectUserOffers(ctx, userID, limit, offset)
}

func (os *offerService) UpdateOfferStatus(ctx context.Context, offerID uint, status string) (entity.Offer, error) {
	return os.offerRepository.UpdateOfferStatus(ctx, offerID, status)
}

func (os *offerService) DeleteOffer(ctx context.Context, offerID uint) (entity.Offer, error) {
	return os.offerRepository.DeleteOffer(ctx, offerID)
}

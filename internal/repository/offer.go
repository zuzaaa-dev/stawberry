package repository

import (
	"context"
	"errors"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"gorm.io/gorm"
)

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *offerRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) InsertOffer(ctx context.Context, offer offer.Offer) (uint, error) {
	offerModel := model.ConvertOfferFromSvc(offer)
	if err := r.db.Create(&offerModel).Error; err != nil {
		if isDuplicateError(err) {
			return 0, &apperror.OfferError{
				Code:    apperror.DuplicateError,
				Message: "offer with this id already exists",
				Err:     err,
			}
		}
		return 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "offer to create product",
			Err:     err,
		}
	}

	return offer.ID, nil
}

func (r *offerRepository) GetOfferByID(ctx context.Context, offerID uint) (entity.Offer, error) {
	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get offer",
			Err:     err,
		}
	}

	return offer, nil
}

func (r *offerRepository) SelectUserOffers(ctx context.Context, userID uint, limit, offset int) ([]entity.Offer, int64, error) {
	var total int64
	if err := r.db.Model(&model.Offer{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to count user offers",
			Err:     err,
		}
	}

	var offers []entity.Offer
	if err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&offers).Error; err != nil {
		return nil, 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user offers",
			Err:     err,
		}
	}

	return offers, total, nil
}

func (r *offerRepository) UpdateOfferStatus(ctx context.Context, offerID uint, status string) (entity.Offer, error) {
	tx := r.db.Model(&model.Offer{}).Where("id = ?", offerID).Update("status", status)
	if tx.Error != nil {
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to update offer status",
			Err:     tx.Error,
		}
	}
	if tx.RowsAffected == 0 {
		return entity.Offer{}, apperror.ErrOfferNotFound
	}

	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get updated offer",
			Err:     err,
		}
	}

	return offer, nil
}

func (r *offerRepository) DeleteOffer(ctx context.Context, offerID uint) (entity.Offer, error) {
	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get offer for deletion",
			Err:     err,
		}
	}

	if err := r.db.Delete(&offer).Error; err != nil {
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to delete offer",
			Err:     err,
		}
	}

	return offer, nil
}

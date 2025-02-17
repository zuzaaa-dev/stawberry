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

func (r *offerRepository) InsertOffer(
	ctx context.Context,
	offer offer.Offer,
) (uint, error) {
	offerModel := model.ConvertOfferFromSvc(offer)

	result := r.db.WithContext(ctx).Create(&offerModel)
	if result.Error != nil {
		if isDuplicateError(result.Error) {
			return 0, &apperror.OfferError{
				Code:    apperror.DuplicateError,
				Message: "offer with this id already exists",
				Err:     result.Error,
			}
		}
		return 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "offer to create product",
			Err:     result.Error,
		}
	}

	return offer.ID, nil
}

func (r *offerRepository) GetOfferByID(
	ctx context.Context,
	offerID uint,
) (entity.Offer, error) {
	var offer entity.Offer

	result := r.db.WithContext(ctx).
		Where("id = ?", offerID).
		First(&offer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get offer",
			Err:     result.Error,
		}
	}

	return offer, nil
}

func (r *offerRepository) SelectUserOffers(
	ctx context.Context,
	userID uint,
	limit, offset int,
) ([]entity.Offer, int64, error) {
	var total int64

	result := r.db.WithContext(ctx).
		Model(&model.Offer{}).
		Where("user_id = ?", userID).
		Count(&total)

	if result.Error != nil {
		return nil, 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to count user offers",
			Err:     result.Error,
		}
	}

	var offers []entity.Offer
	result = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Find(&offers)

	if result.Error != nil {
		return nil, 0, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user offers",
			Err:     result.Error,
		}
	}

	return offers, total, nil
}

func (r *offerRepository) UpdateOfferStatus(
	ctx context.Context,
	offerID uint,
	status string,
) (entity.Offer, error) {
	tx := r.db.WithContext(ctx).
		Model(&model.Offer{}).
		Where("id = ?", offerID).
		Update("status", status)

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
	result := r.db.WithContext(ctx).
		Where("id = ?", offerID).
		First(&offer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get updated offer",
			Err:     result.Error,
		}
	}

	return offer, nil
}

func (r *offerRepository) DeleteOffer(
	ctx context.Context,
	offerID uint,
) (entity.Offer, error) {
	var offer entity.Offer

	result := r.db.WithContext(ctx).
		Where("id = ?", offerID).
		First(&offer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.Offer{}, apperror.ErrOfferNotFound
		}
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to get offer for deletion",
			Err:     result.Error,
		}
	}

	if result := r.db.WithContext(ctx).Delete(&offer); result.Error != nil {
		return entity.Offer{}, &apperror.OfferError{
			Code:    apperror.DatabaseError,
			Message: "failed to delete offer",
			Err:     result.Error,
		}
	}

	return offer, nil
}

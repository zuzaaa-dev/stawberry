package repository

import (
	"errors"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"

	"gorm.io/gorm"
)

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *offerRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) CreateOffer(offer model.Offer) (uint, error) {
	if err := r.db.Create(&offer).Error; err != nil {
		return 0, err
	}

	return offer.ID, nil
}

func (r *offerRepository) GetUserOffers(userID uint, limit, offset int) ([]entity.Offer, int64, error) {
	var total int64
	if err := r.db.Model(&model.Offer{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var offers []entity.Offer
	if err := r.db.Offset(offset).Limit(limit).Find(&offers).Error; err != nil {
		return nil, 0, err
	}

	return offers, total, nil
}

func (r *offerRepository) GetOffer(offerID uint) (entity.Offer, error) {
	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return entity.Offer{}, err
	}

	return offer, nil
}

func (r *offerRepository) UpdateOfferStatus(offerID uint, status string) (entity.Offer, error) {
	result := r.db.Model(&model.Offer{}).Where("id = ?", offerID).Update("status", status)
	if result.Error != nil {
		return entity.Offer{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.Offer{}, errors.New("Offer not found") // изменю хэндлинг ошибок после апрува хэндлинга Артема
	}

	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return entity.Offer{}, err
	}

	return offer, nil
}

func (r *offerRepository) DeleteOffer(offerID uint) (entity.Offer, error) {
	var offer entity.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return entity.Offer{}, err
	}

	if err := r.db.Delete(&offer).Error; err != nil {
		return entity.Offer{}, err
	}

	return entity.Offer{}, nil
}

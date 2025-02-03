package repository

import (
	"errors"
	"marketplace/models"

	"gorm.io/gorm"
)

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *offerRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) CreateOffer(offer models.Offer) (uint, error) {
	if err := r.db.Create(&offer).Error; err != nil {
		return 0, err
	}

	return offer.ID, nil
}

func (r *offerRepository) GetUserOffers(userID uint, limit, offset int) ([]models.Offer, int64, error) {
	var total int64
	if err := r.db.Model(&models.Offer{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var offers []models.Offer
	if err := r.db.Offset(offset).Limit(limit).Find(&offers).Error; err != nil {
		return nil, 0, err
	}

	return offers, total, nil
}

func (r *offerRepository) GetOffer(offerID uint) (*models.Offer, error) {
	var offer models.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepository) UpdateOfferStatus(offerID uint, status string) (*models.Offer, error) {
	result := r.db.Model(&models.Offer{}).Where("id = ?", offerID).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("Offer not found") //изменю хэндлинг ошибок после апрува хэндлинга Артема
	}

	var offer models.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepository) DeleteOffer(offerID uint) (*models.Offer, error) {
	var offer models.Offer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return nil, err
	}

	if err := r.db.Delete(&offer).Error; err != nil {
		return nil, err
	}

	return &offer, nil
}

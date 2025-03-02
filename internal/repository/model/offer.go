package model

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
)

type Offer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OfferPrice float64   `gorm:"type:decimal(10,2);not null" json:"offer_price"`
	Status     string    `gorm:"size:50;not null" json:"status"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID     *uint     `gorm:"index" json:"user_id,omitempty"`
	ProductID  *uint     `gorm:"index" json:"product_id,omitempty"`
}

func ConvertOfferFromSvc(offer offer.Offer) Offer {
	return Offer{
		ID:         offer.ID,
		OfferPrice: offer.OfferPrice,
		UserID:     offer.UserID,
		ProductID:  offer.ProductID,
		Status:     offer.Status,
		CreatedAt:  offer.CreatedAt,
		UpdatedAt:  offer.UpdatedAt,
	}
}

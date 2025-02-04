package model

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
	"time"
)

type Offer struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    uint
	ProductID uint
	StoreID   uint
	Price     float64
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Product   Product `gorm:"foreignKey:ProductID"`
	Store     Store   `gorm:"foreignKey:StoreID"`
}

func ConvertOfferFromSvc(offer offer.Offer) Offer {
	return Offer{
		ID:        offer.ID,
		UserID:    offer.UserID,
		ProductID: offer.ProductID,
		StoreID:   offer.StoreID,
		Price:     offer.Price,
		Status:    offer.Status,
		ExpiresAt: offer.ExpiresAt,
		CreatedAt: offer.CreatedAt,
		UpdatedAt: offer.UpdatedAt,
	}
}

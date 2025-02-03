package dto

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
)

type PostOfferReq struct {
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	StoreID   uint      `json:"store_id"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type PostOfferResp struct {
	ID uint `json:"id"`
}

func (po *PostOfferReq) ConvertToSvc() offer.Offer {
	return offer.Offer{
		UserID:    po.UserID,
		ProductID: po.ProductID,
		StoreID:   po.StoreID,
		Price:     po.Price,
		Status:    po.Status,
		ExpiresAt: po.ExpiresAt,
	}
}

type PatchOfferStatusReq struct {
	Status string `json:"status" binding:"required"`
}

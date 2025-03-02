package dto

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/offer"
)

type PostOfferReq struct {
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
}

type PostOfferResp struct {
	ID uint `json:"id"`
}

func (po *PostOfferReq) ConvertToSvc() offer.Offer {
	return offer.Offer{
		UserID:     &po.UserID,
		ProductID:  &po.ProductID,
		OfferPrice: po.Price,
		Status:     po.Status,
	}
}

type PatchOfferStatusReq struct {
	Status string `json:"status" binding:"required"`
}

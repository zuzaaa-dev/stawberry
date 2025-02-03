package offer

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/store"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
)

type Offer struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id"`
	ProductID uint            `json:"product_id"`
	StoreID   uint            `json:"store_id"`
	Price     float64         `json:"price"`
	Status    string          `json:"status"`
	ExpiresAt time.Time       `json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Product   product.Product `json:"product" gorm:"foreignKey:ProductID"`
	Store     store.Store     `json:"store" gorm:"foreignKey:StoreID"`
}

func (o *Offer) ConvertToRepo() model.Offer {
	return model.Offer{
		ID:        o.ID,
		UserID:    o.UserID,
		ProductID: o.ProductID,
		StoreID:   o.StoreID,
		Price:     o.Price,
		Status:    o.Status,
		ExpiresAt: o.ExpiresAt,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

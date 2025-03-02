package offer

import (
	"time"
)

type Offer struct {
	ID         uint      `json:"id"`
	OfferPrice float64   `json:"offer_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     *uint     `json:"user_id,omitempty"`
	ProductID  *uint     `json:"product_id,omitempty"`
}

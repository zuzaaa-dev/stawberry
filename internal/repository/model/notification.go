package model

import "time"

type Notification struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	OfferID   uint
	Message   string
	Read      bool
	CreatedAt time.Time
}

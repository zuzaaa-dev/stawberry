package model

import "time"

type Notification struct {
	ID      uint      `gorm:"column:id;primary_key"`
	UserID  uint      `gorm:"column:user_id"`
	Message string    `gorm:"column:message"`
	SentAt  time.Time `gorm:"column:sent_at"`
}

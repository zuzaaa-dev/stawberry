package entity

import "time"

type Notification struct {
	ID      uint      `json:"id"`
	UserID  uint      `json:"user_id"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
}

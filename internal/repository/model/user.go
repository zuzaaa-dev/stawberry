package model

import "time"

type User struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Email         string `gorm:"unique"`
	Password      string
	Notifications []Notification
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	UserID    uint      `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	// role в теории, если понадобится
}

type RefreshToken struct {
	UUID        uuid.UUID  `json:"uuid"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
	Fingerprint string     `json:"fingerprint"`
	UserID      uint       `json:"user_id"`
}

func (rt RefreshToken) IsValid() bool {
	now := time.Now()
	if rt.RevokedAt != nil && rt.RevokedAt.Before(now) {
		return false
	}
	return rt.ExpiresAt.After(now)
}

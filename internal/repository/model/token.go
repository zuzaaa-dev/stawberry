package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type RefreshToken struct {
	UUID        uuid.UUID  `gorm:"column:uuid"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	ExpiresAt   time.Time  `gorm:"column:expires_at"`
	RevokedAt   *time.Time `gorm:"column:revoked_at"`
	Fingerprint string     `gorm:"column:fingerprint"`
	UserID      uint       `gorm:"column:user_id"`
}

func ConvertTokenFromEntity(t entity.RefreshToken) RefreshToken {
	return RefreshToken{
		UUID:        t.UUID,
		CreatedAt:   t.CreatedAt,
		ExpiresAt:   t.ExpiresAt,
		RevokedAt:   t.RevokedAt,
		Fingerprint: t.Fingerprint,
		UserID:      t.UserID,
	}
}

func ConvertTokenToEntity(t RefreshToken) entity.RefreshToken {
	return entity.RefreshToken{
		UUID:        t.UUID,
		CreatedAt:   t.CreatedAt,
		ExpiresAt:   t.ExpiresAt,
		RevokedAt:   t.RevokedAt,
		Fingerprint: t.Fingerprint,
		UserID:      t.UserID,
	}
}

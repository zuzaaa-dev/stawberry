package repository

import (
	"context"
	"errors"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) InsertToken(ctx context.Context, token entity.RefreshToken) error {
	tokenModel := model.ConvertTokenFromEntity(token)
	if err := r.db.WithContext(ctx).Create(tokenModel).Error; err != nil {
		if isDuplicateError(err) {
			return &apperror.ProductError{
				Code:    apperror.DuplicateError,
				Message: "token with this uuid already exists",
				Err:     err,
			}
		}
		return &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to create token",
			Err:     err,
		}
	}

	return nil
}

func (r *tokenRepository) GetActivesTokenByUserID(ctx context.Context, userID uint) ([]entity.RefreshToken, error) {
	var tokensModel []model.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&tokensModel).Error; err != nil {
		return nil, &apperror.TokenError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user tokens",
			Err:     err,
		}
	}

	tokens := make([]entity.RefreshToken, 0, len(tokensModel))
	for _, token := range tokensModel {
		tokens = append(tokens, model.ConvertTokenToEntity(token))
	}

	return tokens, nil
}

func (r *tokenRepository) RevokeActivesByUserID(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		return &apperror.TokenError{
			Code:    apperror.DatabaseError,
			Message: "failed to revoke user tokens",
			Err:     result.Error,
		}
	}

	if result.RowsAffected == 0 {
		return apperror.ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepository) GetByUUID(ctx context.Context, uuid string) (entity.RefreshToken, error) {
	var tokenModel model.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&tokenModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.RefreshToken{}, apperror.ErrInvalidToken
		}
		return entity.RefreshToken{}, &apperror.TokenError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch token by uuid",
			Err:     err,
		}
	}

	return model.ConvertTokenToEntity(tokenModel), nil
}

func (r *tokenRepository) Update(ctx context.Context, refresh entity.RefreshToken) (entity.RefreshToken, error) {
	refreshModel := model.ConvertTokenFromEntity(refresh)

	result := r.db.WithContext(ctx).
		Model(&refreshModel).
		Where("uuid = ?", refresh.UUID).
		Updates(refreshModel)

	if result.Error != nil {
		return entity.RefreshToken{}, &apperror.TokenError{
			Code:    apperror.DatabaseError,
			Message: "failed to update refresh token",
			Err:     result.Error,
		}
	}

	if result.RowsAffected == 0 {
		return entity.RefreshToken{}, apperror.ErrTokenNotFound
	}

	return model.ConvertTokenToEntity(refreshModel), nil
}

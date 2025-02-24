package user

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/pkg/security"
)

const maxUsers = 5

type Repository interface {
	InsertUser(ctx context.Context, user User) (uint, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
	UpdateUser(ctx context.Context, user User) error
}

type TokenService interface {
	GenerateTokens(ctx context.Context, fingerprint string, userID uint) (string, entity.RefreshToken, error)
	InsertToken(ctx context.Context, token entity.RefreshToken) error
	GetActivesTokenByUserID(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
	RevokeActivesByUserID(ctx context.Context, userID uint) error
	GetByUUID(ctx context.Context, uuid string) (entity.RefreshToken, error)
	Update(ctx context.Context, refresh entity.RefreshToken) (entity.RefreshToken, error)
}

type userService struct {
	userRepository Repository
	tokenService   TokenService
}

func NewUserService(userRepo Repository, tokenService TokenService) *userService {
	return &userService{userRepository: userRepo, tokenService: tokenService}
}

// CreateUser создает пользователя, хэшируя его пароль, используя HashArgon2id
// генерирует access токен и uuid refresh uuid.
func (us *userService) CreateUser(
	ctx context.Context,
	user User,
	fingerprint string,
) (string, string, error) {
	hash, err := security.HashArgon2id(user.Password)
	if err != nil {
		appError := apperror.ErrFailedToGeneratePassword
		appError.Err = fmt.Errorf("failed to generate password %w, password = %s", err, user.Password)
		return "", "", appError
	}
	user.Password = hash

	id, err := us.userRepository.InsertUser(ctx, user)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := us.tokenService.GenerateTokens(ctx, fingerprint, id)
	if err != nil {
		return "", "", err
	}

	if err = us.tokenService.InsertToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.UUID.String(), nil
}

// Authenticate аутентифицирует пользователя по email и паролю, создавая новые токены.
func (us *userService) Authenticate(
	ctx context.Context,
	email,
	password,
	fingerprint string,
) (string, string, error) {
	user, err := us.userRepository.GetUser(ctx, email)
	if err != nil {
		return "", "", apperror.ErrUserNotFound
	}

	compared, err := security.ComparePasswordAndArgon2id(password, user.Password)
	if err != nil {
		return "", "", err
	}

	if !compared {
		return "", "", apperror.ErrInvalidPassword
	}

	// проверяет количество токенов у пользователя
	userActiveTokens, err := us.tokenService.GetActivesTokenByUserID(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	if len(userActiveTokens) >= maxUsers {
		if err := us.tokenService.RevokeActivesByUserID(ctx, user.ID); err != nil {
			return "", "", err
		}
	}

	accessToken, refreshToken, err := us.tokenService.GenerateTokens(ctx, fingerprint, user.ID)
	if err != nil {
		return "", "", err
	}

	if err = us.tokenService.InsertToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.UUID.String(), nil
}

// Refresh обновляет пару токенов аутентификации.
func (us *userService) Refresh(
	ctx context.Context,
	refreshToken,
	fingerprint string,
) (string, string, error) {
	refresh, err := us.tokenService.GetByUUID(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if !refresh.IsValid() {
		return "", "", apperror.ErrInvalidToken
	}

	if refresh.Fingerprint != fingerprint {
		return "", "", apperror.ErrInvalidFingerprint
	}

	now := time.Now()
	refresh.RevokedAt = &now

	refresh, err = us.tokenService.Update(ctx, refresh)
	if err != nil {
		return "", "", err
	}

	user, err := us.userRepository.GetUserByID(ctx, refresh.UserID)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := us.tokenService.GenerateTokens(ctx, fingerprint, user.ID)
	if err != nil {
		return "", "", err
	}

	err = us.tokenService.InsertToken(ctx, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh.UUID.String(), nil
}

func (us *userService) Logout(
	ctx context.Context,
	refreshToken,
	fingerprint string,
) error {
	refresh, err := us.tokenService.GetByUUID(ctx, refreshToken)
	if err != nil {
		return apperror.ErrInvalidToken
	}

	if !refresh.IsValid() {
		return apperror.ErrInvalidToken
	}

	if refresh.Fingerprint != fingerprint {
		return apperror.ErrInvalidFingerprint
	}

	now := time.Now()
	refresh.RevokedAt = &now

	_, err = us.tokenService.Update(ctx, refresh)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

func (us *userService) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	return us.userRepository.GetUserByID(ctx, id)
}

func (us *userService) UpdateUser(ctx context.Context, id uint, updateUser UpdateUser) error {
	panic("implement me")
}

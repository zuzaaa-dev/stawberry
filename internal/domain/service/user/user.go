package user

import (
	"context"
	"fmt"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/pkg/security"
)

type Repository interface {
	InsertUser(ctx context.Context, user User) (uint, error)
	GetUser(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
	UpdateUser(ctx context.Context, user User) error
}

type TokenService interface {
	GenerateTokens(ctx context.Context, userID uint) (string, entity.RefreshToken, error)
	InsertToken(ctx context.Context, token string) error
}

type userService struct {
	userRepository Repository
	tokenService   TokenService
}

func NewUserService(userRepo Repository) *userService {
	return &userService{userRepository: userRepo}
}

func (us *userService) CreateUser(ctx context.Context, user User) (string, string, error) {
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

	accessToken, refreshToken, err := us.tokenService.GenerateTokens(ctx, id)
	if err != nil {
		return "", "", err
	}

	if err = us.tokenService.InsertToken(ctx, refreshToken.UUID.String()); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.UUID.String(), nil
}

func (us *userService) Authenticate(ctx context.Context, email, password string) (string, error) {
	panic("implement me")
}

func (us *userService) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	panic("implement me")
}

func (us *userService) UpdateUser(ctx context.Context, id string, updateUser UpdateUser) error {
	panic("implement me")
}

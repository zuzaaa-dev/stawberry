package repository

import (
	"context"
	"errors"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/user"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// InsertUser вставляет пользователя в БД
func (r *userRepository) InsertUser(
	ctx context.Context,
	user user.User,
) (uint, error) {
	userModel := model.ConvertUserFromSvc(user)
	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		if isDuplicateError(err) {
			return 0, &apperror.UserError{
				Code:    apperror.DuplicateError,
				Message: "user with this email already exists",
				Err:     err,
			}
		}
		return 0, &apperror.UserError{
			Code:    apperror.DatabaseError,
			Message: "failed to create user",
			Err:     err,
		}
	}

	return userModel.ID, nil
}

// GetUser получает пользователя по почте
func (r *userRepository) GetUser(
	ctx context.Context,
	email string,
) (entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, apperror.ErrUserNotFound
		}
		return entity.User{}, &apperror.UserError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user by email",
			Err:     err,
		}
	}

	return model.ConvertUserToEntity(userModel), nil
}

// GetUserByID получает пользователя по айди
func (r *userRepository) GetUserByID(
	ctx context.Context,
	id uint,
) (entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, apperror.ErrUserNotFound
		}
		return entity.User{}, &apperror.UserError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user by ID",
			Err:     err,
		}
	}

	return model.ConvertUserToEntity(userModel), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user user.User) error {
	panic("implement me")
}

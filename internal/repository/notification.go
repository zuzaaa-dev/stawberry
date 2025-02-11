package repository

import (
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *notificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) SelectUserNotifications(id string, offset, limit int) ([]entity.Notification, int, error) {
	var total int64
	if err := r.db.Model(&model.Notification{}).Where("user_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, &apperror.NotificationError{
			Code:    apperror.DatabaseError,
			Message: "failed to count user notifications",
			Err:     err,
		}
	}

	var notifications []entity.Notification
	if err := r.db.Where("user_id = ?", id).Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, &apperror.NotificationError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch user not notifications",
			Err:     err,
		}
	}

	return notifications, int(total), nil
}

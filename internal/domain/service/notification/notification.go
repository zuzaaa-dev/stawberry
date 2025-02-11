package notification

import "github.com/zuzaaa-dev/stawberry/internal/domain/entity"

type Repository interface {
	SelectUserNotifications(id string, offset, limit int) ([]entity.Notification, int, error)
}

type notificationService struct {
	notificationRepository Repository
}

func NewNotificationService(notificationRepository Repository) *notificationService {
	return &notificationService{notificationRepository}
}

func (ns *notificationService) GetNotification(id string, offset int, limit int) ([]entity.Notification, int, error) {
	return ns.notificationRepository.SelectUserNotifications(id, offset, limit)
}

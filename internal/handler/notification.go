package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type NotificationService interface {
	GetNotification(id string, offset int, limit int) ([]entity.Notification, int, error)
}

type notificationHandler struct {
	offerService NotificationService
}

func NewNotificationHandler(notificationService NotificationService) *notificationHandler {
	return &notificationHandler{offerService: notificationService}
}

// GetNotification обработчик уведомлений
// получает все уведомления (одобрение или неодобрение заявки) авторизированного пользователя
func (h *notificationHandler) GetNotification(c *gin.Context) {
	// сделать user, ok := c.Get("user")
	// сейчас в контекст кладется userID, в будущем структура user
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID from context"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	offset := (page - 1) * limit

	// исправить вызов h.offerService.GetNotification()
	// сейчас в контекст кладется userID, в будущем структура user
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	notifications, total, err := h.offerService.GetNotification(strconv.FormatUint(uint64(uid), 10), offset, limit)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}

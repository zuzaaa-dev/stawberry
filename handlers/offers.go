package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"marketplace/models"
	"marketplace/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OfferHandlers struct {
	repo repository.OfferRepository
}

func NewOfferHandler(repo repository.OfferRepository) *OfferHandlers {
	return &OfferHandlers{repo: repo}
}

type UpdateOfferStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *OfferHandlers) CreateOffer(c *gin.Context) {
	userID, _ := c.Get("userID")

	var offer models.Offer
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer.UserID = userID.(uint)
	offer.Status = "pending"
	offer.ExpiresAt = time.Now().Add(24 * time.Hour)

	var err error
	if offer.ID, err = h.repo.CreateOffer(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create notification for store
	// notification := models.Notification{
	// 	UserID:  offer.StoreID, // Store notification
	// 	OfferID: offer.ID,
	// 	Message: fmt.Sprintf("New offer received for product %d", offer.ProductID),
	// }
	// h.notifyRepo.Create(&notification)

	c.JSON(http.StatusCreated, offer)
}

func (h *OfferHandlers) GetUserOffers(c *gin.Context) {
	userID, _ := c.Get("userID")

	//надо вынести в pkg/utils/paginators.go или что то подобное
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

	var (
		offers []models.Offer
		total  int64
	)
	if offers, total, err = h.repo.GetUserOffers(userID.(uint), limit, offset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": offers,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}

func (h *OfferHandlers) GetOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nondigit offer id"})
		return
	}

	var offer *models.Offer
	if offer, err = h.repo.GetOffer(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": offer,
	})
}

func (h *OfferHandlers) UpdateOfferStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nondigit offer id"})
		return
	}

	//предположу что статус будет лежать в теле запроса
	var req UpdateOfferStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var offer *models.Offer
	if offer, err = h.repo.UpdateOfferStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create notification for store
	// notification := models.Notification{
	// 	UserID:  offer.StoreID, // Store notification
	// 	OfferID: offer.ID,
	// 	Message: fmt.Sprintf("Offer %d has changed status to %s", offer.ID, offer.Status),
	// }
	// h.notifyRepo.Create(&notification)

	c.JSON(http.StatusCreated, offer)
}

func (h *OfferHandlers) CancelOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nondigit offer id"})
		return
	}

	var offer *models.Offer
	if offer, err = h.repo.DeleteOffer(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create notification for store
	// notification := models.Notification{
	// 	UserID:  offer.StoreID, // Store notification
	// 	OfferID: offer.ID,
	// 	Message: fmt.Sprintf("Offer %d canceled", offer.ID),
	// }
	// h.notifyRepo.Create(&notification)

	c.JSON(http.StatusCreated, offer)
}

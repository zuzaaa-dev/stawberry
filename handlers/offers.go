package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"marketplace/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOffer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var offer models.Offer
		if err := c.ShouldBindJSON(&offer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offer.UserID = userID.(uint)
		offer.Status = "pending"
		offer.ExpiresAt = time.Now().Add(24 * time.Hour)

		if err := db.Create(&offer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create notification for store
		notification := models.Notification{
			UserID:  offer.StoreID, // Store notification
			OfferID: offer.ID,
			Message: fmt.Sprintf("New offer received for product %d", offer.ProductID),
		}
		db.Create(&notification)

		c.JSON(http.StatusCreated, offer)
	}
}

func GetUserOffers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		//TODO:надо вынести в pkg/utils/paginators.go или что то подобное
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

		var total int64
		if err := db.Model(&models.Offer{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var offers []models.Offer
		if err := db.Offset(offset).Limit(limit).Find(&offers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
}

func GetOffer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nondigit offer id"})
			return
		}

		var offer models.Offer
		if err := db.Where("id = ?", id).First(&offer).Error; err != nil {
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
}

type UpdateOfferStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func UpdateOfferStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		result := db.Model(&models.Offer{}).Where("id = ?", id).Update("status", req.Status)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
			return
		}

		var offer models.Offer
		if err := db.Where("id = ?", id).First(&offer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create notification for store
		notification := models.Notification{
			UserID:  offer.StoreID, // Store notification
			OfferID: offer.ID,
			Message: fmt.Sprintf("Offer %d has changed status to %s", offer.ID, offer.Status),
		}
		db.Create(&notification)

		c.JSON(http.StatusCreated, offer)
	}
}

func CancelOffer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nondigit offer id"})
			return
		}

		var offer models.Offer
		if err := db.Where("id = ?", id).First(&offer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Delete(&offer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create notification for store
		notification := models.Notification{
			UserID:  offer.StoreID, // Store notification
			OfferID: offer.ID,
			Message: fmt.Sprintf("Offer %d canceled", offer.ID),
		}
		db.Create(&notification)

		c.JSON(http.StatusCreated, offer)
	}
}

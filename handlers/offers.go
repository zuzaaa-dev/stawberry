package handlers

import (
	"fmt"
	"net/http"
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

package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service"
	"github.com/zuzaaa-dev/stawberry/internal/handler/dto"
	"gorm.io/gorm"
)

// Здесь надо по-новому ошибки обработать как в product
type OfferHandlers struct {
	offerService service.OfferService
}

func NewOfferHandler(offerService service.OfferService) *OfferHandlers {
	return &OfferHandlers{offerService: offerService}
}

func (h *OfferHandlers) PostOffer(c *gin.Context) {
	userID, _ := c.Get("userID")

	var offer dto.PostOfferReq
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer.UserID = userID.(uint)
	offer.Status = "pending"
	offer.ExpiresAt = time.Now().Add(24 * time.Hour)

	var response dto.PostOfferResp
	var err error
	if response.ID, err = h.offerService.CreateOffer(offer.ConvertToSvc()); err != nil {
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
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
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

	offers, total, err := h.offerService.GetUserOffers(userID.(uint), offset, limit)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid non digit offer id"})
		return
	}

	offer, err := h.offerService.GetOffer(uint(id))
	if err != nil {
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

	// предположу что статус будет лежать в теле запроса
	var req dto.PatchOfferStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	offer, err := h.offerService.UpdateOfferStatus(uint(id), req.Status)
	if err != nil {
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

	offer, err := h.offerService.DeleteOffer(uint(id))
	if err != nil {
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

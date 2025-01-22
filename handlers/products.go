package handlers

import (
	"errors"
	"marketplace/models/apperrors"
	"marketplace/repository"
	"math"
	"net/http"
	"strconv"

	"marketplace/models"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// AddProduct handles the creation of a new product in the database.
// It expects a JSON payload with product details
// and responds with the created product.
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct retrieves a specific product by its ID from the database.
// If the product is not found, it responds with a 404 error.
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.repo.GetProductByID(id)
	if err != nil {
		var productErr *apperrors.ProductError
		if errors.As(err, &productErr) {
			status := http.StatusInternalServerError

			switch productErr.Code {
			case apperrors.NotFound:
				status = http.StatusNotFound
			case apperrors.DatabaseError:
				status = http.StatusInternalServerError
			}

			c.JSON(status, gin.H{
				"code":    productErr.Code,
				"message": productErr.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    apperrors.InternalError,
			"message": "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts fetches a paginated list of products from the database.
// It supports query parameters for pagination: `page` and `limit`.
func (h *ProductHandler) GetProducts(c *gin.Context) {
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

	// Get recordings with pagination
	products, total, err := h.repo.GetProducts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}

// GetStoreProducts fetches a paginated list of products for
// a specific store from the database. It filters products
// by `store_id` and supports query parameters for
// pagination: `page` and `limit`.
func (h *ProductHandler) GetStoreProducts(c *gin.Context) {
	id := c.Param("id")

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

	products, total, err := h.repo.GetStoreProducts(id, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"meta": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  totalPages,
		},
	})
}

// UpdateProduct updates specific fields of a product identified by its ID.
// It expects a JSON payload with the fields to be updated
// and responds with the update status.
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var update models.ProductUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update"})
		return
	}

	err := h.repo.UpdateProduct(id, &update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func handleError(c *gin.Context, err error) {
	var productErr *apperrors.ProductError
	if errors.As(err, &productErr) {
		status := http.StatusInternalServerError

		switch productErr.Code {
		case apperrors.NotFound:
			status = http.StatusNotFound
		case apperrors.DuplicateError:
			status = http.StatusConflict
		case apperrors.DatabaseError:
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{
			"code":    productErr.Code,
			"message": productErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apperrors.InternalError,
		"message": "An unexpected error occurred",
	})
}

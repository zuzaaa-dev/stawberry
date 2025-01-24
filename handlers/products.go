package handlers

import (
	"errors"
	objectstorage "marketplace/s3"
	"math"
	"net/http"
	"strconv"

	"marketplace/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddProduct handles the creation of a new product in the database.
// It expects a JSON payload with product details
// and responds with the created product.
func AddProduct(db *gorm.DB, s3 *objectstorage.BucketBasics) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

// GetProduct retrieves a specific product by its ID from the database.
// If the product is not found, it responds with a 404 error.
func GetProduct(db *gorm.DB, s3 *objectstorage.BucketBasics) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var product models.Product

		if err := db.Where("id = ?", id).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

// GetProducts fetches a paginated list of products from the database.
// It supports query parameters for pagination: `page` and `limit`.
func GetProducts(db *gorm.DB, s3 *objectstorage.BucketBasics) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		if err := db.Model(&models.Product{}).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get recordings with pagination
		var products []models.Product
		if err := db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
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
}

// GetStoreProducts fetches a paginated list of products for
// a specific store from the database. It filters products
// by `store_id` and supports query parameters for
// pagination: `page` and `limit`.
func GetStoreProducts(db *gorm.DB, s3 *objectstorage.BucketBasics) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var total int64
		if err := db.Model(&models.Product{}).Where("store_id = ?", id).Count(&total).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Products not found in this store"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var products []models.Product
		if err := db.Where("store_id = ?", id).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
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
}

// UpdateProduct updates specific fields of a product identified by its ID.
// It expects a JSON payload with the fields to be updated
// and responds with the update status.
func UpdateProduct(db *gorm.DB, s3 *objectstorage.BucketBasics) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result := db.Model(&models.Product{}).Where("id = ?", id).Updates(input)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

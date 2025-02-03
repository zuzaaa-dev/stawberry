package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/gin-gonic/gin"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service"
	"github.com/zuzaaa-dev/stawberry/internal/handler/dto"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return ProductHandler{productService: productService}
}

func (h *ProductHandler) PostProduct(c *gin.Context) {
	var postProductReq dto.PostProductReq

	if err := c.ShouldBindJSON(&postProductReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid product data",
			"details": err.Error(),
		})
		return
	}

	var response dto.PostProductResp
	var err error
	if response.ID, err = h.productService.CreateProduct(postProductReq.ConvertToSvc()); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid limit value (should be between 1 and 100)",
		})
		return
	}

	offset := (page - 1) * limit

	products, total, err := h.productService.GetProducts(offset, limit)
	if err != nil {
		handleError(c, err)
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

func (h *ProductHandler) GetStoreProducts(c *gin.Context) {
	id := c.Param("id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid limit value (should be between 1 and 100)",
		})
		return
	}

	offset := (page - 1) * limit

	products, total, err := h.productService.GetStoreProducts(id, offset, limit)
	if err != nil {
		handleError(c, err)
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

func (h *ProductHandler) PatchProduct(c *gin.Context) {
	id := c.Param("id")

	var update dto.PatchProductReq
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid update data",
			"details": err.Error(),
		})
		return
	}

	if err := h.productService.UpdateProduct(id, update.ConvertToSvc()); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

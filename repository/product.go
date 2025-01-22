package repository

import (
	"errors"
	"strings"

	"marketplace/models"
	"marketplace/models/apperrors"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		if isDuplicateError(err) {
			return &apperrors.ProductError{
				Code:    apperrors.DuplicateError,
				Message: "product with this id already exists",
				Err:     err,
			}
		}
		return &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to create product",
			Err:     err,
		}
	}

	return nil
}

func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to fetch product",
			Err:     err,
		}
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(offset, limit int) ([]models.Product, int, error) {
	var total int64
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to count products",
			Err:     err,
		}
	}

	var products []models.Product
	if err := r.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to fetch products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *ProductRepository) GetStoreProducts(id string, offset, limit int) ([]models.Product, int, error) {
	var exists bool
	if err := r.db.Model(&models.Store{}).Where("id = ?", id).Select("1").Scan(&exists).Error; err != nil {
		return nil, 0, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to check store existence",
			Err:     err,
		}
	}
	if !exists {
		return nil, 0, apperrors.ErrStoreNotFound
	}

	var total int64
	if err := r.db.Model(&models.Product{}).Where("store_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to count store products",
			Err:     err,
		}
	}

	var products []models.Product
	if err := r.db.Where("store_id = ?", id).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to fetch store products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *ProductRepository) UpdateProduct(id string, update *models.ProductUpdate) error {
	tx := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(update)

	if tx.Error != nil {
		if isDuplicateError(tx.Error) {
			return &apperrors.ProductError{
				Code:    apperrors.DuplicateError,
				Message: "product with these details already exists",
				Err:     tx.Error,
			}
		}
		return &apperrors.ProductError{
			Code:    apperrors.DatabaseError,
			Message: "failed to update product",
			Err:     tx.Error,
		}
	}

	if tx.RowsAffected == 0 {
		return apperrors.ErrProductNotFound
	}

	return nil
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate") ||
		strings.Contains(err.Error(), "unique violation")
}

package repository

import (
	"errors"
	"strings"

	error2 "github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product model.Product) (uint, error) {
	if err := r.db.Create(product).Error; err != nil {
		if isDuplicateError(err) {
			return 0, &error2.ProductError{
				Code:    error2.DuplicateError,
				Message: "product with this id already exists",
				Err:     err,
			}
		}
		return 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to create product",
			Err:     err,
		}
	}

	return product.ID, nil
}

func (r *productRepository) GetProductByID(id string) (entity.Product, error) {
	var product entity.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, error2.ErrProductNotFound
		}
		return entity.Product{}, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to fetch product",
			Err:     err,
		}
	}

	return product, nil
}

func (r *productRepository) GetProducts(offset, limit int) ([]entity.Product, int, error) {
	var total int64
	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to count products",
			Err:     err,
		}
	}

	var products []entity.Product
	if err := r.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to fetch products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *productRepository) GetStoreProducts(id string, offset, limit int) ([]entity.Product, int, error) {
	var exists bool
	if err := r.db.Model(&model.Store{}).Where("id = ?", id).Select("1").Scan(&exists).Error; err != nil {
		return nil, 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to check store existence",
			Err:     err,
		}
	}
	if !exists {
		return nil, 0, error2.ErrStoreNotFound
	}

	var total int64
	if err := r.db.Model(&model.Product{}).Where("store_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to count store products",
			Err:     err,
		}
	}

	var products []entity.Product
	if err := r.db.Where("store_id = ?", id).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to fetch store products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *productRepository) UpdateProduct(id string, update model.UpdateProduct) error {
	tx := r.db.Model(&model.Product{}).Where("id = ?", id).Updates(update)

	if tx.Error != nil {
		if isDuplicateError(tx.Error) {
			return &error2.ProductError{
				Code:    error2.DuplicateError,
				Message: "product with these details already exists",
				Err:     tx.Error,
			}
		}
		return &error2.ProductError{
			Code:    error2.DatabaseError,
			Message: "failed to update product",
			Err:     tx.Error,
		}
	}

	if tx.RowsAffected == 0 {
		return error2.ErrProductNotFound
	}

	return nil
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate") ||
		strings.Contains(err.Error(), "unique violation")
}

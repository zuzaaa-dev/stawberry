package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) InsertProduct(
	ctx context.Context,
	product product.Product,
) (uint, error) {
	productModel := model.ConvertProductFromSvc(product)
	if err := r.db.WithContext(ctx).Create(productModel).Error; err != nil {
		if isDuplicateError(err) {
			return 0, &apperror.ProductError{
				Code:    apperror.DuplicateError,
				Message: "product with this id already exists",
				Err:     err,
			}
		}
		return 0, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to create product",
			Err:     err,
		}
	}

	return productModel.ID, nil
}

func (r *productRepository) GetProductByID(
	ctx context.Context,
	id string,
) (entity.Product, error) {
	var productModel model.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&productModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, apperror.ErrProductNotFound
		}
		return entity.Product{}, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch product",
			Err:     err,
		}
	}

	return model.ConvertProductToEntity(productModel), nil
}

func (r *productRepository) SelectProducts(
	ctx context.Context,
	offset,
	limit int,
) ([]entity.Product, int, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to count products",
			Err:     err,
		}
	}

	var products []entity.Product
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *productRepository) SelectStoreProducts(
	ctx context.Context,
	id string, offset, limit int,
) ([]entity.Product, int, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("store_id = ?", id).
		Count(&total).Error; err != nil {
		return nil, 0, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to count store products",
			Err:     err,
		}
	}

	var products []entity.Product
	if err := r.db.WithContext(ctx).
		Where("store_id = ?", id).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to fetch store products",
			Err:     err,
		}
	}

	return products, int(total), nil
}

func (r *productRepository) UpdateProduct(
	ctx context.Context,
	id string,
	update product.UpdateProduct,
) error {
	updateModel := model.ConvertUpdateProductFromSvc(update)
	tx := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Updates(updateModel)
	if tx.Error != nil {
		if isDuplicateError(tx.Error) {
			return &apperror.ProductError{
				Code:    apperror.DuplicateError,
				Message: "product with these details already exists",
				Err:     tx.Error,
			}
		}
		return &apperror.ProductError{
			Code:    apperror.DatabaseError,
			Message: "failed to update product",
			Err:     tx.Error,
		}
	}

	if tx.RowsAffected == 0 {
		return apperror.ErrProductNotFound
	}

	return nil
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate") ||
		strings.Contains(err.Error(), "unique violation")
}

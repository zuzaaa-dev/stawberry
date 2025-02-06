package product

import (
	"context"

	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type Repository interface {
	InsertProduct(ctx context.Context, product Product) (uint, error)
	GetProductByID(ctx context.Context, id string) (entity.Product, error)
	SelectProducts(ctx context.Context, offset, limit int) ([]entity.Product, int, error)
	SelectStoreProducts(ctx context.Context, id string, offset, limit int) ([]entity.Product, int, error)
	UpdateProduct(ctx context.Context, id string, update UpdateProduct) error
}

type productService struct {
	productRepository Repository
}

func NewProductService(productRepo Repository) *productService {
	return &productService{productRepository: productRepo}
}

func (ps *productService) CreateProduct(ctx context.Context, product Product) (uint, error) {
	return ps.productRepository.InsertProduct(ctx, product)
}

func (ps *productService) GetProductByID(ctx context.Context, id string) (entity.Product, error) {
	return ps.productRepository.GetProductByID(ctx, id)
}

func (ps *productService) GetProducts(ctx context.Context, offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.SelectProducts(ctx, offset, limit)
}

func (ps *productService) GetStoreProducts(ctx context.Context, id string, offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.SelectStoreProducts(ctx, id, offset, limit)
}

func (ps *productService) UpdateProduct(ctx context.Context, id string, updateProduct UpdateProduct) error {
	return ps.productRepository.UpdateProduct(ctx, id, updateProduct)
}

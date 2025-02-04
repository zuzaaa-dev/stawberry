package product

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

type Repository interface {
	InsertProduct(product Product) (uint, error)
	GetProductByID(id string) (entity.Product, error)
	SelectProducts(offset, limit int) ([]entity.Product, int, error)
	SelectStoreProducts(id string, offset, limit int) ([]entity.Product, int, error)
	UpdateProduct(id string, update UpdateProduct) error
}

type productService struct {
	productRepository Repository
}

func NewProductService(productRepo Repository) *productService {
	return &productService{productRepository: productRepo}
}

func (ps *productService) CreateProduct(product Product) (uint, error) {
	return ps.productRepository.InsertProduct(product)
}

func (ps *productService) GetProductByID(id string) (entity.Product, error) {
	return ps.productRepository.GetProductByID(id)
}

func (ps *productService) GetProducts(offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.SelectProducts(offset, limit)
}

func (ps *productService) GetStoreProducts(id string, offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.SelectStoreProducts(id, offset, limit)
}

func (ps *productService) UpdateProduct(id string, updateProduct UpdateProduct) error {
	return ps.productRepository.UpdateProduct(id, updateProduct)
}

package product

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/repository"
)

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *productService {
	return &productService{productRepository: productRepo}
}

func (ps *productService) CreateProduct(product Product) (uint, error) {
	return ps.productRepository.CreateProduct(product.ConvertToRepo())
}

func (ps *productService) GetProductByID(id string) (entity.Product, error) {
	return ps.productRepository.GetProductByID(id)
}

func (ps *productService) GetProducts(offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.GetProducts(offset, limit)
}

func (ps *productService) GetStoreProducts(id string, offset, limit int) ([]entity.Product, int, error) {
	return ps.productRepository.GetStoreProducts(id, offset, limit)
}

func (ps *productService) UpdateProduct(id string, updateProduct UpdateProduct) error {
	return ps.productRepository.UpdateProduct(id, updateProduct.ConvertToRepo())
}

package store

import (
	"time"

	"github.com/zuzaaa-dev/stawberry/internal/domain/service/product"
	"github.com/zuzaaa-dev/stawberry/internal/repository/model"
)

// Store здесь еще надо думать тому, кто возьмется Store реализовывать
type Store struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Products    []product.Product `json:"products,omitempty" gorm:"foreignKey:StoreID"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (s *Store) ConvertToRepo() model.Store {
	products := make([]model.Product, len(s.Products))
	for _, p := range s.Products {
		products = append(products, p.ConvertToRepo())
	}
	return model.Store{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Products:    products,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

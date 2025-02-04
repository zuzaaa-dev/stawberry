package store

import (
	"time"
)

// Store здесь еще надо думать тому, кто возьмется Store реализовывать
type Store struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//Products    []product.Product `json:"products,omitempty" gorm:"foreignKey:StoreID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

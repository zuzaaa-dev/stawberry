package dto

import "github.com/zuzaaa-dev/stawberry/internal/domain/service/product"

type PostProductReq struct {
	StoreID     uint    `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
}

type PostProductResp struct {
	ID uint `json:"id"`
}

func (pp *PostProductReq) ConvertToSvc() product.Product {
	return product.Product{
		StoreID:     pp.StoreID,
		Name:        pp.Name,
		Description: pp.Description,
		Price:       pp.Price,
		Category:    pp.Category,
		InStock:     pp.InStock,
	}
}

type PatchProductReq struct {
	StoreID     *uint    `json:"store_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	InStock     *bool    `json:"in_stock,omitempty"`
}

func (pp *PatchProductReq) ConvertToSvc() product.UpdateProduct {
	return product.UpdateProduct{
		StoreID:     pp.StoreID,
		Name:        pp.Name,
		Description: pp.Description,
		Price:       pp.Price,
		Category:    pp.Category,
		InStock:     pp.InStock,
	}
}

package dto

import "github.com/zuzaaa-dev/stawberry/internal/domain/service/product"

type PostProductReq struct {
	Name        string `json:"name"`
	CategoryId  int    `json:"category_id"`
	Description string `json:"description"`
}

type PostProductResp struct {
	ID uint `json:"id"`
}

func (pp *PostProductReq) ConvertToSvc() product.Product {
	return product.Product{
		Name:        pp.Name,
		CategoryId:  pp.CategoryId,
		Description: pp.Description,
	}
}

type PatchProductReq struct {
	Name        *string `json:"name,omitempty"`
	CategoryId  *int    `json:"category_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (pp *PatchProductReq) ConvertToSvc() product.UpdateProduct {
	return product.UpdateProduct{
		Name:        *pp.Name,
		CategoryId:  *pp.CategoryId,
		Description: *pp.Description,
	}
}

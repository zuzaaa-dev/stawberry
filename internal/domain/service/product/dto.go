package product

type Product struct {
	ID          uint   `json:"id" gorm:"primary_key;"`
	Name        string `json:"name"`
	CategoryId  int    `json:"category_id"`
	Description string `json:"description"`
}

type UpdateProduct struct {
	Name        string `json:"name,omitempty"`
	CategoryId  int    `json:"category_id,omitempty"`
	Description string `json:"description,omitempty"`
}

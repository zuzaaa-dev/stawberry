package entity

type Product struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CategoryId  int    `json:"category_id"`
	Description string `json:"description"`
}

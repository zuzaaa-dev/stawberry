package entity

type Category struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Lft      uint   `json:"lft"`
	Rgt      uint   `json:"rgt"`
	ParentId uint   `json:"parent_id"`
}

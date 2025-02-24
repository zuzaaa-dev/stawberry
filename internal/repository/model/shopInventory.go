package model

type ShopInventory struct {
	ProductID   uint `gorm:"not null;primaryKey;index" json:"product_id"`
	ShopID      uint `gorm:"not null;primaryKey;index" json:"shop_id"`
	IsAvailable bool `gorm:"not null;default:false" json:"is_available"`
}

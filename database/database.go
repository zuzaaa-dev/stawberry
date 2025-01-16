package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"marketplace/config"
	"marketplace/models"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Product{},
		&models.Offer{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

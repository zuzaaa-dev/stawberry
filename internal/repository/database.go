package repository

import (
	"log"

	"github.com/zuzaaa-dev/stawberry/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

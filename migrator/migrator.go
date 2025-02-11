package migrator

import (
	"log"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// RunMigrations applies database migrations using *gorm.DB.
func RunMigrations(gormDB *gorm.DB, migrationsDir string) {
	// Get the underlying *sql.DB from *gorm.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}

	// Set the database dialect (PostgreSQL)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	// Apply migrations
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully!")
}

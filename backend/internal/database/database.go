package database

import (
	"log"

	"backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dbURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Critical: Failed to connect to PostgreSQL database: %v", err)
	}

	log.Println("Database connection established successfully")

	// Trigger GORM AutoMigrate schema sync
	err = db.AutoMigrate(&models.Wheel{}, &models.WheelEntry{}, &models.SpinHistory{})
	if err != nil {
		log.Fatalf("Critical: Failed to run database auto-migration: %v", err)
	}

	log.Println("Database schema auto-migrations verified successfully")
	DB = db
	return db
}

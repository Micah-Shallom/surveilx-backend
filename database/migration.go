package database

import (
	"log"
	"survielx-backend/models"
)

func MigrateDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.AccessExitPoint{},
		&models.Profile{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

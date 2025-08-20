package database

import (
	"log"
	"survielx-backend/models"
)

func MigrateDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.AccessExitPoint{},
		&models.Vehicle{},
		&models.VehicleActivity{},
		&models.GuestVehicleActivity{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

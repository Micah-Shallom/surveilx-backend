package main

import (
	"log"
	"survielx-backend/database"
	"survielx-backend/middleware"
	"survielx-backend/models"
	"survielx-backend/routers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.VehicleLog{},
		&models.Watchlist{},
		&models.GuestLog{},
		&models.AccessExitPoint{},
		&models.Profile{},
	)

	r := routers.SetupRouter()
	r.Use(middleware.CORSMiddleware())

	r.Run()
}
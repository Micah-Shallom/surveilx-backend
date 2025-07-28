package main

import (
	"boilerplate/database"
	"boilerplate/models"
	"boilerplate/routers"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.User{})

	r := routers.SetupRouter()
	r.Run()
}

package main

import (
	"fmt"
	"log"
	"os"
	"survielx-backend/database"
	"survielx-backend/routers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()
	database.MigrateDatabase()

	r := routers.SetupRouter()

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))) // Use PORT from .env file
}

package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("POSTGRES_DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("[error] failed to initialize database, got error %v\n", err)
		panic("Failed to connect to database!")
	}

	DB = database
}

package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		panic("Unable to load environment variables")
	}

	dsn := os.Getenv("POSTGRES_DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("[error] failed to initialize database, got error %v\n", err)
		panic("Failed to connect to database!")
	}

	DB = database
}

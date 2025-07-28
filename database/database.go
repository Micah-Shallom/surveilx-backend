package database

import (
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
		panic("Failed to connect to database!")
	}

	DB = database
}

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.Vehicle{}, &models.VehicleLog{}, &models.Admin{})
}

func DropTables() {
	DB.Migrator().DropTable(&models.User{}, &models.Vehicle{}, &models.VehicleLog{}, &models.Admin{})
}

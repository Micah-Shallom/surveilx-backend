package main

import (
	"log"
	"os"
	"testing"
	"survielx-backend/database"
	"survielx-backend/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Set up the test database
	database.Connect()
	database.Migrate()

	// Run the tests
	exitVal := m.Run()

	// Clean up the test database
	database.DropTables()

	os.Exit(exitVal)
}

func setupRouter() *gin.Engine {
	return routers.SetupRouter()
}

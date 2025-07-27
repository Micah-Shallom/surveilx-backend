package main

import (
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/routers"
)

func main() {
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.User{})

	r := routers.SetupRouter()
	r.Run()
}

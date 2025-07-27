package main

import (
	"boilerplate/database"
	"boilerplate/models"
	"boilerplate/routers"
)

func main() {
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.User{})

	r := routers.SetupRouter()
	r.Run()
}

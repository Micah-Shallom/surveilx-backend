package services

import (
	"boilerplate/database"
	"boilerplate/models"
)

func GetUsers(users *[]models.User) error {
	result := database.DB.Find(users)
	return result.Error
}

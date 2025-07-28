package services

import (
	"survielx-backend/database"
	"survielx-backend/models"
)

func CreateUser(user *models.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func GetUsers(users *[]models.User) error {
	result := database.DB.Find(users)
	return result.Error
}
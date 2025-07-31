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

func GetUserByID(userID string) (models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

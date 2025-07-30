package services

import (
	"survielx-backend/database"
	"survielx-backend/models"
)

func GetProfile(userID string) (models.Profile, error) {
	var profile models.Profile
	if err := database.DB.Where("id = ?", userID).First(&profile).Error; err != nil {
		return profile, err
	}
	return profile, nil
}

func UpdateProfile(profile *models.Profile) error {
	result := database.DB.Save(profile)
	return result.Error
}

func DeleteProfile(userID string) error {
	result := database.DB.Delete(&models.Profile{}, "id = ?", userID)
	return result.Error
}

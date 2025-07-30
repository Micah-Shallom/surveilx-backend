package services

import (
	"fmt"
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

func UpdateProfile(profile *models.Profile, req models.UpdateProfileInput) error {

	profileUpdates := models.Profile{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Phone:       req.Phone,
		AvatarURL:   req.AvatarURL,
		Department:  req.Department,
	}

	result, err := models.UpdateFields(database.DB, profile, profileUpdates, "id = ?", profile.ID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return result.Error
}

func DeleteProfile(userID string) error {
	result := database.DB.Delete(&models.Profile{}, "id = ?", userID)
	return result.Error
}

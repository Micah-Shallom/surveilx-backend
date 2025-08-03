package services

import (
	"errors"
	"net/http"
	"survielx-backend/models"

	"gorm.io/gorm"
)

func UpdateUserProfile(db *gorm.DB, userID string, req *models.UpdateUserProfileInput) (int, error) {
	var (
		user    models.User
		profile models.Profile
	)

	exists := models.CheckExists(db, &user, "id  = ?", userID)
	if !exists {
		return http.StatusNotFound, errors.New("user not found")
	}

	exists = models.CheckExists(db, &profile, "user_id = ?", userID)
	if !exists {
		return http.StatusNotFound, errors.New("profile not found")
	}

	profileUpdates := models.Profile{
		FullName:    req.FullName,
		UserName:    req.UserName,
		Phone:       req.Phone,
		AvatarURL:   req.AvatarURL,
		DisplayName: req.DisplayName,
	}

	result, err := models.UpdateFields(db, &profile, profileUpdates, "user_id = ?", userID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if result.RowsAffected == 0 {
		return http.StatusBadRequest, errors.New("no fields updated")
	}

	return http.StatusOK, nil
}

func GetUserProfile(db *gorm.DB, user_id string) (models.Profile, int, error) {
	var (
		profile models.Profile
		user    models.User
	)

	exists := models.CheckExists(db, &user, "id  = ?", user_id)
	if !exists {
		return models.Profile{}, http.StatusNotFound, errors.New("user not found")
	}

	err := profile.GetUserProfile(db, user_id)
	if err != nil {
		return models.Profile{}, http.StatusBadRequest, err
	}

	return profile, http.StatusOK, nil
}

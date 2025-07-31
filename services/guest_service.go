package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"
)

func RegisterGuest(guest *models.Guest) (*models.Guest, int, error) {
	db := database.DB

	checkExists := models.CheckExists(db, &models.Guest{}, "plate_number = ?", guest.PlateNumber)
	if checkExists {
		return nil, http.StatusConflict, fmt.Errorf("guest with plate number %s already exists", guest.PlateNumber)
	}

	if err := db.Create(guest).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var fullGuest models.Guest
	if err := db.
		Preload("User").
		First(&fullGuest, "id = ?", guest.ID).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &fullGuest, http.StatusCreated, nil
}

func GetGuests(from, to time.Time) (*[]models.Guest, int, error) {
	var guests []models.Guest
	if err := database.DB.Where("created_at BETWEEN ? AND ?", from, to).Find(&guests).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &guests, http.StatusOK, nil
}

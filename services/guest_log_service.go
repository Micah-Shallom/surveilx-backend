package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"
)

func CreateGuestLog(logInput *models.LogGuestInput) (*models.GuestLog, int, error) {
	db := database.DB

	if !logInput.IsEntry {
		// For exits, check if the last log was an entry
		var lastLog models.GuestLog
		if err := db.Where("plate_number = ?", logInput.PlateNumber).Order("timestamp desc").First(&lastLog).Error; err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("no entry record found for this vehicle")
		}

		if !lastLog.IsEntry {
			return nil, http.StatusBadRequest, fmt.Errorf("last log for this vehicle was not an entry")
		}
	}

	log := models.GuestLog{
		PlateNumber: logInput.PlateNumber,
		IsEntry:     logInput.IsEntry,
		Timestamp:   time.Now(),
	}

	if err := db.Create(&log).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &log, http.StatusCreated, nil
}

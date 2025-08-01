package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"

	"gorm.io/gorm"
)

func AddToWatchlist(db *gorm.DB, req models.LogGuestInput) (*models.GuestWatchList, int, error) {

	watchlist := models.GuestWatchList{
		PlateNumber:  req.PlateNumber,
		Type:         req.Type,
		IsEntry:      req.IsEntry,
		RegisteredBy: req.RegisteredBy,
		Timestamp:    time.Now(),
		CreatedAt:    time.Now(),
	}

	checkExists := models.CheckExists(db, &models.GuestWatchList{}, "plate_number = ?", watchlist.PlateNumber)
	if checkExists {
		return nil, http.StatusConflict, fmt.Errorf("watchlist with plate number %s already exists", watchlist.PlateNumber)
	}

	if err := db.Create(watchlist).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var fullWatchlist models.GuestWatchList
	if err := db.
		Preload("User").
		First(&fullWatchlist, "id = ?", watchlist.ID).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &fullWatchlist, http.StatusCreated, nil
}

func GetWatchlist(from, to time.Time) (*[]models.GuestWatchList, int, error) {
	var watchlist []models.GuestWatchList
	if err := database.DB.Where("created_at BETWEEN ? AND ?", from, to).Find(&watchlist).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &watchlist, http.StatusOK, nil
}

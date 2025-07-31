package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"
)

func AddToWatchlist(watchlist *models.Watchlist) (*models.Watchlist, int, error) {
	db := database.DB

	checkExists := models.CheckExists(db, &models.Watchlist{}, "plate_number = ?", watchlist.PlateNumber)
	if checkExists {
		return nil, http.StatusConflict, fmt.Errorf("watchlist with plate number %s already exists", watchlist.PlateNumber)
	}

	if err := db.Create(watchlist).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var fullWatchlist models.Watchlist
	if err := db.
		Preload("User").
		First(&fullWatchlist, "id = ?", watchlist.ID).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &fullWatchlist, http.StatusCreated, nil
}

func GetWatchlist(from, to time.Time) (*[]models.Watchlist, int, error) {
	var watchlist []models.Watchlist
	if err := database.DB.Where("created_at BETWEEN ? AND ?", from, to).Find(&watchlist).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &watchlist, http.StatusOK, nil
}

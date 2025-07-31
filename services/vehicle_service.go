package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"
)

func RegisterVehicle(vehicle *models.Vehicle) (*models.Vehicle, int, error) {
	db := database.DB

	checkExists := models.CheckExists(db, &models.Vehicle{}, "plate_number = ?", vehicle.PlateNumber)
	if checkExists {
		return nil, http.StatusConflict, fmt.Errorf("vehicle with plate number %s already exists", vehicle.PlateNumber)
	}

	if err := db.Create(vehicle).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var fullVehicle models.Vehicle
	if err := db.
		Preload("User").
		First(&fullVehicle, "id = ?", vehicle.ID).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &fullVehicle, http.StatusCreated, nil
}

func GetVehicleByPlateNumber(plateNumber string) (*models.Vehicle, int, error) {
	var vehicle models.Vehicle
	if err := database.DB.Preload("User").Where("plate_number = ?", plateNumber).First(&vehicle).Error; err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to fetch plate number: %v", err)
	}
	return &vehicle, http.StatusOK, nil
}

func LogVehicle(log *models.VehicleLog) (*models.VehicleLog, int, error) {
	if err := database.DB.Create(log).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return log, http.StatusCreated, nil
}

func GetVehicleLogs(userId string) (*[]models.VehicleLog, int, error) {
	var logs []models.VehicleLog
	if err := database.DB.Where("user_id = ?", userId).Find(&logs).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &logs, http.StatusOK, nil
}

func CreateVehicleLog(vehicle *models.Vehicle, isEntry bool, entryPointID, exitPointID string) (*models.VehicleLog, int, error) {
	log := models.VehicleLog{
		VehicleID:    vehicle.ID,
		UserID:       vehicle.UserID,
		Timestamp:    time.Now(),
		IsEntry:      isEntry,
	}

	if entryPointID != "" {
		log.EntryPointID = &entryPointID
	}

	if exitPointID != "" {
		log.ExitPointID = &exitPointID
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &log, http.StatusCreated, nil
}

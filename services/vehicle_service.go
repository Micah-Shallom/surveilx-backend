package services

import (
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"
)

func RegisterVehicle(vehicle *models.Vehicle) (*models.Vehicle, int, error) {
	if err := database.DB.Create(vehicle).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return vehicle, http.StatusCreated, nil
}

func GetVehicleByPlateNumber(plateNumber string) (*models.Vehicle, int, error) {
	var vehicle models.Vehicle
	if err := database.DB.Where("plate_number = ?", plateNumber).First(&vehicle).Error; err != nil {
		return nil, http.StatusNotFound, err
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

func CreateVehicleLog(vehicle *models.Vehicle, isEntry bool) (*models.VehicleLog, int, error) {
	log := models.VehicleLog{
		VehicleID: vehicle.ID,
		UserID:    vehicle.UserID,
		Timestamp: time.Now(),
		IsEntry:   isEntry,
	}
	if err := database.DB.Create(&log).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &log, http.StatusCreated, nil
}

package services

import (
	"fmt"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"time"

	"gorm.io/gorm"
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

func CreateVehicleLog(vehicle *models.Vehicle, req models.LogVehicleInput) (*models.VehicleLog, int, error) {
	var (
		isEntry      = req.IsEntry
		entryPointID = req.EntryPointID
		exitPointID  = req.ExitPointID	
	)

	if err := validateVehicleEntryExit(database.DB, vehicle.ID, isEntry); err != nil {
		return nil, http.StatusBadRequest, err
	}

	log := models.VehicleLog{
		VehicleID: vehicle.ID,
		UserID:    vehicle.UserID,
		Timestamp: time.Now(),
		IsEntry:   isEntry,
		Type:      req.Type,
		CreatedAt: time.Now(),
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

func GetVehicleStatusByPlateNumber(plateNumber string) (string, error) {
	vehicle, _, err := GetVehicleByPlateNumber(plateNumber)
	if err != nil {
		return "", err
	}

	return GetVehicleStatus(vehicle.ID)
}

func GetVehicleStatus(vehicleID string) (string, error) {
	db := database.DB
	var lastLog models.VehicleLog

	err := db.Where("vehicle_id = ?", vehicleID).
		Order("timestamp desc").
		First(&lastLog).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "outside", nil // Never entered
		}
		return "", fmt.Errorf("database error: %v", err)
	}

	if lastLog.IsEntry {
		return "inside", nil
	}
	return "outside", nil
}

func GetVehicleLogHistory(vehicleID string, limit int) ([]models.VehicleLog, error) {
	db := database.DB
	var logs []models.VehicleLog

	query := db.Where("vehicle_id = ?", vehicleID).Order("timestamp desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get vehicle log history: %v", err)
	}

	return logs, nil
}

func validateVehicleEntryExit(db *gorm.DB, vehicleID string, isEntry bool) error {
	var lastLog models.VehicleLog

	err := db.Where("vehicle_id = ?", vehicleID).
		Order("timestamp desc").
		First(&lastLog).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if !isEntry {
				return fmt.Errorf("vehicle must enter before it can exit")
			}
			return nil
		}
		return fmt.Errorf("database error while checking vehicle logs: %v", err)
	}

	if isEntry {
		if lastLog.IsEntry {
			return fmt.Errorf("vehicle is already inside - cannot enter again without exiting first")
		}
	} else {
		if !lastLog.IsEntry {
			return fmt.Errorf("vehicle is already outside - cannot exit without entering first")
		}
	}

	return nil
}

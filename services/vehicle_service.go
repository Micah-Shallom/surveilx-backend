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

func LogVehicleActivity(db *gorm.DB, req models.LogVehicleActivityInput, loggedByUserID string) (*models.VehicleActivityResponse, int, error) {

	activity := models.VehicleActivity{
		PlateNumber: req.PlateNumber,
		VisitorType: req.VisitorType,
		IsEntry:     req.IsEntry,
		VehicleType: req.VehicleType,
		Timestamp:   time.Now(),
	}

	// Handle entry/exit points
	if req.EntryPointID != "" {
		activity.EntryPointID = &req.EntryPointID
	}
	if req.ExitPointID != "" {
		activity.ExitPointID = &req.ExitPointID
	}

	// Handle based on visitor type
	switch req.VisitorType {
	case models.VisitorTypeRegistered:
		// Get registered vehicle
		vehicle, _, err := GetVehicleByPlateNumber(req.PlateNumber)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("registered vehicle not found: %v", err)
		}

		activity.VehicleID = &vehicle.ID
		activity.UserID = &vehicle.UserID

		// Validate entry/exit logic for registered vehicles
		if err := validateVehicleEntryExit(db, vehicle.ID, req.IsEntry); err != nil {
			return nil, http.StatusBadRequest, err
		}

	case models.VisitorTypeGuest:
		// For guests, record who logged them
		activity.RegisteredBy = &loggedByUserID

		// Optional: Validate guest entry/exit logic if needed
		if err := validateGuestEntryExit(db, req.PlateNumber, req.IsEntry); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	// Create the activity record
	if err := db.Create(&activity).Error; err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create activity log: %v", err)
	}

	// Return with preloaded relationships
	return getVehicleActivityResponse(db, activity.ID)
}

func getVehicleActivityResponse(db *gorm.DB, activityID string) (*models.VehicleActivityResponse, int, error) {
	var activity models.VehicleActivity

	err := db.Preload("Vehicle").
		Preload("User").
		Preload("EntryPoint").
		Preload("ExitPoint").
		Preload("RegisteredByUser").
		First(&activity, "id = ?", activityID).Error

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get activity: %v", err)
	}

	response := convertToActivityResponse(activity)
	return &response, http.StatusOK, nil
}

func GetVehicleLogs(userId string) (*[]models.VehicleActivity, int, error) {
	var logs []models.VehicleActivity
	if err := database.DB.Where("user_id = ?", userId).Find(&logs).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &logs, http.StatusOK, nil
}

func CreateVehicleLog(vehicle *models.Vehicle, req models.LogVehicleInput) (*models.VehicleActivity, int, error) {
	var (
		isEntry      = req.IsEntry
		entryPointID = req.EntryPointID
		exitPointID  = req.ExitPointID
	)

	if err := validateVehicleEntryExit(database.DB, vehicle.ID, isEntry); err != nil {
		return nil, http.StatusBadRequest, err
	}

	log := models.VehicleActivity{
		VehicleID:    &vehicle.ID,
		UserID:       &vehicle.UserID,
		Timestamp:    time.Now(),
		IsEntry:      isEntry,
		VehicleType:  vehicle.Type,
		PlateNumber:  vehicle.PlateNumber,
		RegisteredBy: &vehicle.UserID, // Assuming the vehicle is registered by the user
		CreatedAt:    time.Now(),
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
	var lastLog models.VehicleActivity

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

func GetVehicleLogHistory(vehicleID string, limit int) ([]models.VehicleActivity, error) {
	db := database.DB
	var logs []models.VehicleActivity

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
	var lastLog models.VehicleActivity

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

func GetUserVehicles(db *gorm.DB, userID string) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	if err := db.Where("user_id = ?", userID).Find(&vehicles).Error; err != nil {
		return nil, fmt.Errorf("failed to get user vehicles: %v", err)
	}

	return vehicles, nil
}

func GetVehicleActivities(db *gorm.DB, plateNumber string, limit int, visitorType *models.VisitorType) ([]models.VehicleActivityResponse, error) {
	var activities []models.VehicleActivity

	query := db.Preload("Vehicle").
		Preload("User").
		Preload("EntryPoint").
		Preload("ExitPoint").
		Preload("RegisteredByUser")

	if plateNumber != "" {
		query = query.Where("plate_number = ?", plateNumber)
	}

	if visitorType != nil {
		query = query.Where("visitor_type = ?", *visitorType)
	}

	query = query.Order("timestamp desc")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("failed to get vehicle activities: %v", err)
	}

	// Convert to response format
	responses := make([]models.VehicleActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = convertToActivityResponse(activity)
	}

	return responses, nil
}

func convertToActivityResponse(activity models.VehicleActivity) models.VehicleActivityResponse {
	return models.VehicleActivityResponse{
		ID:           activity.ID,
		PlateNumber:  activity.PlateNumber,
		VisitorType:  activity.VisitorType,
		IsEntry:      activity.IsEntry,
		VehicleType:  activity.VehicleType,
		Timestamp:    activity.Timestamp,
		Vehicle:      activity.Vehicle,
		User:         activity.User,
		EntryPoint:   activity.EntryPoint,
		ExitPoint:    activity.ExitPoint,
		RegisteredBy: activity.RegisteredByUser,
	}
}

func validateGuestEntryExit(db *gorm.DB, plateNumber string, isEntry bool) error {
	var lastActivity models.VehicleActivity

	err := db.Where("plate_number = ? AND visitor_type = ?", plateNumber, models.VisitorTypeGuest).
		Order("timestamp desc").
		First(&lastActivity).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if !isEntry {
				return fmt.Errorf("guest vehicle must enter before it can exit")
			}
			return nil
		}
		return fmt.Errorf("database error while checking guest vehicle logs: %v", err)
	}

	if isEntry {
		if lastActivity.IsEntry {
			return fmt.Errorf("guest vehicle is already inside - cannot enter again without exiting first")
		}
	} else {
		if !lastActivity.IsEntry {
			return fmt.Errorf("guest vehicle is already outside - cannot exit without entering first")
		}
	}

	return nil
}

func GetAllVehicleActivities(from, to time.Time, visitorType *models.VisitorType) ([]models.VehicleActivityResponse, error) {
	db := database.DB
	var activities []models.VehicleActivity

	query := db.Preload("Vehicle").
		Preload("User").
		Preload("EntryPoint").
		Preload("ExitPoint").
		Preload("RegisteredByUser").
		Where("timestamp BETWEEN ? AND ?", from, to)

	if visitorType != nil {
		query = query.Where("visitor_type = ?", *visitorType)
	}

	query = query.Order("timestamp desc")

	if err := query.Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("failed to get activities: %v", err)
	}

	responses := make([]models.VehicleActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = convertToActivityResponse(activity)
	}

	return responses, nil
}

func GenerateActivitySummary(activities []models.VehicleActivityResponse) map[string]any {
	summary := map[string]any{
		"total_activities": len(activities),
		"by_visitor_type": map[string]int{
			"registered": 0,
			"guest":      0,
		},
		"by_action": map[string]int{
			"entries": 0,
			"exits":   0,
		},
		"by_vehicle_type": map[string]int{
			"car":  0,
			"bike": 0,
			"bus":  0,
		},
		"by_time_period": map[string]int{
			"morning":   0, // 06:00 - 12:00
			"afternoon": 0, // 12:00 - 18:00
			"evening":   0, // 18:00 - 24:00
			"night":     0, // 00:00 - 06:00
		},
		"unique_vehicles":     make(map[string]bool),
		"entry_exit_balance":  0, // entries - exits
		"busiest_entry_point": "",
		"busiest_exit_point":  "",
	}

	entryPoints := make(map[string]int)
	exitPoints := make(map[string]int)
	hourlyData := make(map[int]int)

	for _, activity := range activities {
		// Count by visitor type
		summary["by_visitor_type"].(map[string]int)[string(activity.VisitorType)]++

		// Count by action and calculate balance
		if activity.IsEntry {
			summary["by_action"].(map[string]int)["entries"]++
			summary["entry_exit_balance"] = summary["entry_exit_balance"].(int) + 1

			// Track entry points
			if activity.EntryPoint != nil {
				entryPoints[activity.EntryPoint.Name]++
			}
		} else {
			summary["by_action"].(map[string]int)["exits"]++
			summary["entry_exit_balance"] = summary["entry_exit_balance"].(int) - 1

			// Track exit points
			if activity.ExitPoint != nil {
				exitPoints[activity.ExitPoint.Name]++
			}
		}

		// Count by vehicle type
		if _, exists := summary["by_vehicle_type"].(map[string]int)[activity.VehicleType]; exists {
			summary["by_vehicle_type"].(map[string]int)[activity.VehicleType]++
		}

		// Track unique vehicles
		summary["unique_vehicles"].(map[string]bool)[activity.PlateNumber] = true

		// Time period analysis
		hour := activity.Timestamp.Hour()
		hourlyData[hour]++

		var period string
		switch {
		case hour >= 6 && hour < 12:
			period = "morning"
		case hour >= 12 && hour < 18:
			period = "afternoon"
		case hour >= 18 && hour < 24:
			period = "evening"
		default:
			period = "night"
		}
		summary["by_time_period"].(map[string]int)[period]++
	}

	// Find busiest entry and exit points
	var busiestEntry, busiestExit string
	var maxEntryCount, maxExitCount int

	for point, count := range entryPoints {
		if count > maxEntryCount {
			maxEntryCount = count
			busiestEntry = point
		}
	}

	for point, count := range exitPoints {
		if count > maxExitCount {
			maxExitCount = count
			busiestExit = point
		}
	}

	summary["busiest_entry_point"] = busiestEntry
	summary["busiest_exit_point"] = busiestExit

	// Convert unique vehicles to count
	uniqueCount := len(summary["unique_vehicles"].(map[string]bool))
	summary["unique_vehicles"] = uniqueCount

	// Add peak hour information

	// Add calculated metrics
	summary["average_activities_per_vehicle"] = 0.0
	if uniqueCount > 0 {
		summary["average_activities_per_vehicle"] = float64(len(activities)) / float64(uniqueCount)
	}

	return summary
}

package services

import (
	"errors"
	"fmt"
	"math"
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
		return nil, http.StatusBadRequest, err
	}

	var fullVehicle models.Vehicle
	if err := db.
		First(&fullVehicle, "id = ?", vehicle.ID).Error; err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &fullVehicle, http.StatusCreated, nil
}

func DeRegisterVehicle(db *gorm.DB, vehicle_id string) (int, error) {
	var vehicle models.Vehicle

	exists := models.CheckExists(db, &vehicle, "id = ?", vehicle_id)
	if !exists {
		return http.StatusNotFound, errors.New("vehicle does not exist")
	}

	err := vehicle.DeRegister(db)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func GetVehicleByPlateNumber(plateNumber string) (*models.Vehicle, int, error) {
	var vehicle models.Vehicle
	if err := database.DB.Where("plate_number = ?", plateNumber).First(&vehicle).Error; err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to fetch plate number: %v", err)
	}
	return &vehicle, http.StatusOK, nil
}

func LogVehicleActivity(db *gorm.DB, req models.LogVehicleActivityInput) (int, error) {

	activity := models.VehicleActivity{
		PlateNumber: req.PlateNumber,
		VisitorType: models.VisitorTypeRegistered,
		IsEntry:     req.IsEntry,
		Timestamp:   time.Now(),
	}

	if (req.EntryPointID != "" && req.ExitPointID != "") || (req.EntryPointID == "" && req.ExitPointID == "") {
		return http.StatusBadRequest, fmt.Errorf("either both entry and exit points must be provided or neither")
	}

	if req.IsEntry {
		if req.EntryPointID == "" {
			return http.StatusBadRequest, fmt.Errorf("entry point ID is required for entry activity")
		}
	} else {
		if req.ExitPointID == "" {
			return http.StatusBadRequest, fmt.Errorf("exit point ID is required for exit activity")
		}
	}

	if req.EntryPointID != "" {
		exist := models.CheckExists(db, &models.AccessExitPoint{}, "id = ?", req.EntryPointID)
		if !exist {
			return http.StatusNotFound, fmt.Errorf("entry point with ID %s not found", req.EntryPointID)
		}
		activity.EntryPointID = &req.EntryPointID
	}
	if req.ExitPointID != "" {
		exist := models.CheckExists(db, &models.AccessExitPoint{}, "id = ?", req.ExitPointID)
		if !exist {
			return http.StatusNotFound, fmt.Errorf("exit point with ID %s not found", req.ExitPointID)
		}
		activity.ExitPointID = &req.ExitPointID
	}

	vehicle, _, err := GetVehicleByPlateNumber(req.PlateNumber)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("registered vehicle not found: %v", err)
	}

	activity.VehicleID = &vehicle.ID
	activity.VehicleType = vehicle.Type
	activity.Model = vehicle.Model

	if err := validateVehicleEntryExit(db, vehicle.ID, req.IsEntry); err != nil {
		return http.StatusBadRequest, err
	}

	if err := db.Create(&activity).Error; err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create activity log: %v", err)
	}

	return http.StatusOK, nil
}

func getVehicleActivityResponse(db *gorm.DB, activityID string) (*models.VehicleActivityResponse, int, error) {
	var activity models.VehicleActivity

	err := db.
		Preload("EntryPoint").
		Preload("ExitPoint").
		First(&activity, "id = ?", activityID).Error

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get activity: %v", err)
	}

	response := convertToActivityResponse(activity)
	return &response, http.StatusOK, nil
}

func GetVehicleLogs(userId string) (*[]models.VehicleActivity, int, error) {
	var logs []models.VehicleActivity
	if err := database.DB.Where("user_id = ?", userId).Find(&logs).Error; err != nil {
		return nil, http.StatusBadRequest, err
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
		VehicleID:   &vehicle.ID,
		Timestamp:   time.Now(),
		IsEntry:     isEntry,
		VehicleType: vehicle.Type,
		PlateNumber: vehicle.PlateNumber,
		CreatedAt:   time.Now(),
	}

	if entryPointID != "" {
		log.EntryPointID = &entryPointID
	}

	if exitPointID != "" {
		log.ExitPointID = &exitPointID
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &log, http.StatusCreated, nil
}

// the model backend calls this to find out if the vehicle exists either as a registered user or guest user
func IdentifyVehicle(plateNumber string) (models.VehicleIdentity, int, error) {
	vehicle, statuscode, err := GetVehicleByPlateNumber(plateNumber)
	if err != nil {
		// if not record found...security personnel should log this vehicle entry as a guest entry
		return models.VehicleIdentity{}, statuscode, err
	}

	return GetVehicleStatus(vehicle.ID)
}

func GetVehicleStatus(vehicleID string) (models.VehicleIdentity, int, error) {
	db := database.DB
	var (
		lastLog         models.VehicleActivity
		vehicle         models.Vehicle
		vehicleIdentity models.VehicleIdentity
	)

	vehicleIdentity.Status = "outside"
	exists := models.CheckExists(db, &vehicle, "id = ?", vehicleID)
	if exists {
		vehicleIdentity.IsRegistered = true
	}

	err := db.Where("vehicle_id = ?", vehicleID).
		Order("timestamp desc").
		First(&lastLog).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return vehicleIdentity, http.StatusNotFound, fmt.Errorf("no logs found for vehicle")
		}
		return vehicleIdentity, http.StatusBadRequest, fmt.Errorf("database error while checking vehicle logs: %v", err)
	}

	if lastLog.IsEntry {
		vehicleIdentity.Status = "inside"
		return vehicleIdentity, http.StatusOK, nil
	}
	return vehicleIdentity, http.StatusOK, nil
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
			//send notification to security personnel
			return fmt.Errorf("vehicle is already inside - cannot enter again without exiting first")
		}
	} else {
		if !lastLog.IsEntry {
			return fmt.Errorf("vehicle is already outside - cannot exit without entering first")
		}
	}

	return nil
}

func GetUserVehicles(db *gorm.DB, userID string) ([]models.Vehicle, int, error) {
	var vehicles []models.Vehicle

	if err := db.Where("user_id = ?", userID).Find(&vehicles).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get user vehicles: %v", err)
	}

	return vehicles, http.StatusOK, nil
}

func GetVehicleActivities(db *gorm.DB, vehicle_id string) ([]models.VehicleActivityResponse, int, error) {
	responses := []models.VehicleActivityResponse{}

	exists := models.CheckExists(db, &models.Vehicle{}, "id = ?", vehicle_id)
	if !exists {
		return nil, http.StatusNotFound, fmt.Errorf("vehicle with ID %s not found", vehicle_id)
	}

	query := db.
		Model(&models.VehicleActivity{}).
		Select(`
            vehicle_activities.id,
            vehicle_activities.plate_number,
            vehicle_activities.visitor_type,
            vehicle_activities.is_entry,
            vehicle_activities.vehicle_type,
            vehicle_activities.timestamp,
			vehicle_activities.model
        `).
		Joins("LEFT JOIN vehicles ON vehicle_activities.vehicle_id = vehicles.id").
		Joins("LEFT JOIN access_exit_points AS entry_points ON vehicle_activities.entry_point_id = entry_points.id").
		Joins("LEFT JOIN access_exit_points AS exit_points ON vehicle_activities.exit_point_id = exit_points.id").
		Where("vehicle_activities.vehicle_id = ?", vehicle_id).
		Order("vehicle_activities.timestamp desc")

	if err := query.Scan(&responses).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get vehicle activities: %v", err)
	}

	return responses, http.StatusOK, nil
}

func GetGuestVehicleActivitiesByPlateNumber(db *gorm.DB, plateNumber string) ([]models.VehicleActivityResponse, int, error) {
	var activities []models.VehicleActivity

	query := db.
		Model(&models.VehicleActivity{}).
		Select(`
			vehicle_activities.id,
			vehicle_activities.plate_number,
			vehicle_activities.visitor_type,
			vehicle_activities.is_entry,
			vehicle_activities.vehicle_type,
			vehicle_activities.timestamp
		`).
		Joins("LEFT JOIN vehicles ON vehicle_activities.vehicle_id = vehicles.id").
		Joins("LEFT JOIN access_exit_points AS entry_points ON vehicle_activities.entry_point_id = entry_points.id").
		Joins("LEFT JOIN access_exit_points AS exit_points ON vehicle_activities.exit_point_id = exit_points.id").
		Where("vehicle_activities.plate_number = ? AND vehicle_activities.visitor_type = ?", plateNumber, models.VisitorTypeGuest).
		Order("vehicle_activities.timestamp desc")

	if err := query.Find(&activities).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get guest vehicle activities: %v", err)
	}

	responses := make([]models.VehicleActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = convertToActivityResponse(activity)
	}

	return responses, http.StatusOK, nil
}

func convertToActivityResponse(activity models.VehicleActivity) models.VehicleActivityResponse {
	return models.VehicleActivityResponse{
		ID:          activity.ID,
		PlateNumber: activity.PlateNumber,
		VisitorType: activity.VisitorType,
		IsEntry:     activity.IsEntry,
		VehicleType: activity.VehicleType,
		Timestamp:   activity.Timestamp,
		Model:       activity.Model,
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
		Preload("EntryPoint").
		Preload("ExitPoint").
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
			// if activity.EntryPoint != nil {
			// 	entryPoints[activity.EntryPoint.Name]++
			// }
		} else {
			summary["by_action"].(map[string]int)["exits"]++
			summary["entry_exit_balance"] = summary["entry_exit_balance"].(int) - 1

			// Track exit points
			// if activity.ExitPoint != nil {
			// 	exitPoints[activity.ExitPoint.Name]++
			// }
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

func LogGuestVehicleActivity(db *gorm.DB, req models.LogVehicleActivityInput) (int, error) {
	activity := models.GuestVehicleActivity{
		PlateNumber: req.PlateNumber,
		IsEntry:     req.IsEntry,
		Timestamp:   time.Now(),
	}

	// Validate entry/exit points
	if (req.EntryPointID != "" && req.ExitPointID != "") || (req.EntryPointID == "" && req.ExitPointID == "") {
		return http.StatusBadRequest, fmt.Errorf("either entry point or exit point must be provided, not both or neither")
	}

	if req.IsEntry {
		if req.EntryPointID == "" {
			return http.StatusBadRequest, fmt.Errorf("entry point ID is required for entry activity")
		}
	} else {
		if req.ExitPointID == "" {
			return http.StatusBadRequest, fmt.Errorf("exit point ID is required for exit activity")
		}
	}

	if req.EntryPointID != "" {
		exist := models.CheckExists(db, &models.AccessExitPoint{}, "id = ?", req.EntryPointID)
		if !exist {
			return http.StatusNotFound, fmt.Errorf("entry point with ID %s not found", req.EntryPointID)
		}
		activity.EntryPointID = &req.EntryPointID
	}
	if req.ExitPointID != "" {
		exist := models.CheckExists(db, &models.AccessExitPoint{}, "id = ?", req.ExitPointID)
		if !exist {
			return http.StatusNotFound, fmt.Errorf("exit point with ID %s not found", req.ExitPointID)
		}
		activity.ExitPointID = &req.ExitPointID
	}

	// Validate guest vehicle entry/exit sequence
	// if err := validateGuestEntryExit(db, req.PlateNumber, req.IsEntry); err != nil {
	// 	return http.StatusBadRequest, err
	// }

	if err := db.Create(&activity).Error; err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to create guest activity log: %v", err)
	}

	return http.StatusOK, nil
}

func FetchRegisteredVehiclesLogs(db *gorm.DB, pagination models.Pagination, filters models.VehicleFilters) (*models.PaginatedVehicleResponse, int, error) {
	var vehicles []models.Vehicle
	var count int64

	query := db.Model(&models.Vehicle{})

	if filters.PlateNumber != "" {
		query = query.Where("plate_number ILIKE ?", "%"+filters.PlateNumber+"%")
	}
	if filters.Model != "" {
		query = query.Where("model ILIKE ?", "%"+filters.Model+"%")
	}
	if filters.Color != "" {
		query = query.Where("color ILIKE ?", "%"+filters.Color+"%")
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to count vehicles: %v", err)
	}

	offset := (pagination.Page - 1) * pagination.Limit
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	if err := query.Offset(offset).Limit(pagination.Limit).Order("created_at desc").Find(&vehicles).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to fetch vehicles: %v", err)
	}

	paginationResponse := models.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       len(vehicles),
		TotalPagesCount: totalPages,
	}

	response := &models.PaginatedVehicleResponse{
		Data:       vehicles,
		Pagination: paginationResponse,
	}

	return response, http.StatusOK, nil
}

func FetchGuestVehiclesLogs(db *gorm.DB, pagination models.Pagination, plateNumber string) (*models.PaginatedVehicleResponse, int, error) {
	var activities []models.GuestVehicleActivity
	var count int64

	query := db.Model(&models.GuestVehicleActivity{})

	if plateNumber != "" {
		query = query.Where("plate_number ILIKE ?", "%"+plateNumber+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to count guest activities: %v", err)
	}

	offset := (pagination.Page - 1) * pagination.Limit
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	if err := db.Model(&models.GuestVehicleActivity{}).
		Scopes(func(d *gorm.DB) *gorm.DB {
			if plateNumber != "" {
				return d.Where("plate_number ILIKE ?", "%"+plateNumber+"%")
			}
			return d
		}).
		Offset(offset).
		Limit(pagination.Limit).
		Order("timestamp desc").
		Find(&activities).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to fetch guest activities: %v", err)
	}

	paginationResponse := models.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       len(activities),
		TotalPagesCount: totalPages,
	}

	response := &models.PaginatedVehicleResponse{
		Data:       activities,
		Pagination: paginationResponse,
	}

	return response, http.StatusOK, nil
}

func GetVehicleOwnerProfile(db *gorm.DB, vehicleID string) (*models.VehicleOwnerProfileResponse, int, error) {
	var vehicle models.Vehicle
	var user models.User
	var profile models.Profile

	if err := db.Where("id = ?", vehicleID).First(&vehicle).Error; err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("vehicle not found: %v", err)
	}

	if err := db.Where("id = ?", vehicle.UserID).First(&user).Error; err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("vehicle owner not found: %v", err)
	}

	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		profile = models.Profile{
			FullName: user.Name,
		}
	}

	activities, _, err := GetVehicleActivities(db, vehicleID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get vehicle activities: %v", err)
	}

	ownerInfo := models.VehicleOwnerInfo{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    profile.Phone,
		UserName: profile.UserName,
		FullName: profile.FullName,
	}

	response := &models.VehicleOwnerProfileResponse{
		Vehicle:    vehicle,
		Owner:      ownerInfo,
		Activities: activities,
	}

	return response, http.StatusOK, nil
}

func GenerateActivityReport(db *gorm.DB, from, to time.Time, visitorType *models.VisitorType, pagination models.Pagination) (*models.ActivityReportResponse, int, error) {
	var activities []models.VehicleActivity
	var count int64

	query := db.Model(&models.VehicleActivity{}).Where("timestamp BETWEEN ? AND ?", from, to)

	if visitorType != nil {
		query = query.Where("visitor_type = ?", *visitorType)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to count activities: %v", err)
	}

	offset := (pagination.Page - 1) * pagination.Limit
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	if err := query.
		Offset(offset).
		Limit(pagination.Limit).
		Order("timestamp desc").
		Find(&activities).Error; err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to fetch activities: %v", err)
	}

	activityResponses := make([]models.VehicleActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = convertToActivityResponse(activity)
	}

	allActivities, err := GetAllVehicleActivities(from, to, visitorType)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to get all activities for summary: %v", err)
	}

	summary := GenerateActivitySummary(allActivities)

	paginationResponse := models.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       len(activityResponses),
		TotalPagesCount: totalPages,
	}

	reportData := models.ActivityReportData{
		Activities: activityResponses,
		Summary:    summary,
		DateRange: map[string]string{
			"from": from.Format("2006-01-02"),
			"to":   to.Format("2006-01-02"),
		},
	}

	response := &models.ActivityReportResponse{
		Data:       reportData,
		Pagination: paginationResponse,
	}

	return response, http.StatusOK, nil
}

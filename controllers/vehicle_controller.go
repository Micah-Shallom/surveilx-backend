package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"
)

func RegisterVehicle(c *gin.Context) {
	var input models.RegisterVehicleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID := c.MustGet("user_id").(string)
	vehicle := models.Vehicle{
		UserID:      userID,
		PlateNumber: input.PlateNumber,
		Model:       input.Model,
		Color:       input.Color,
		Type:        input.Type,
	}

	createdVehicle, code, err := services.RegisterVehicle(&vehicle)
	if err != nil {
		log.Default().Println("Error registering vehicle:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to register vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Vehicle successfully registered", createdVehicle)
	c.JSON(code, rd)
}

func UpdateVehicle(c *gin.Context) {
	vehicle_id := c.Param("vehicle_id")
	var input models.UpdateVehicleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := utility.ValidateUUID(vehicle_id); err != nil {
		log.Default().Println("invalid vehicle ID:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid vehicle id", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	updatedVehicle, code, err := services.UpdateVehicle(database.DB, vehicle_id, input)
	if err != nil {
		log.Default().Println("Error updating user vehicles:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to update user vehicles", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User vehicle updated successfully")
	rd := utility.BuildSuccessResponse(code, "User vehicle updated successfully", updatedVehicle)
	c.JSON(code, rd)
}

func DeRegisterVehicle(c *gin.Context) {
	vehicle_id := c.Param("vehicle_id")

	code, err := services.DeRegisterVehicle(database.DB, vehicle_id)
	if err != nil {
		log.Default().Println("Error deregistering vehicle:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to deregister vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Vehicle successfully deregistered", nil)
	c.JSON(code, rd)
}

func GetUserVehicles(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	vehicles, code, err := services.GetUserVehicles(database.DB, userID)
	if err != nil {
		log.Default().Println("Error fetching user vehicles:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to fetch user vehicles", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User vehicles retrieved successfully for user ID:", userID)
	rd := utility.BuildSuccessResponse(code, "User vehicles retrieved successfully", vehicles)
	c.JSON(code, rd)
}

func LogVehicleActivity(c *gin.Context) {
	var input models.LogVehicleActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	input.VisitorType = models.VisitorTypeRegistered
	code, err := services.LogVehicleActivity(database.DB, input)
	if err != nil {
		log.Default().Println("Error logging vehicle activity:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to log vehicle activity", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	message := "Vehicle activity logged successfully"
	rd := utility.BuildSuccessResponse(code, message, nil)
	c.JSON(code, rd)
}

func SystemLogVehicleActivity(c *gin.Context) {
	var input models.LogVehicleActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := services.LogVehicleActivity(database.DB, input)
	if err != nil {
		log.Default().Println("Error logging vehicle activity:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to log vehicle activity", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	message := "Vehicle activity logged successfully"
	if input.VisitorType == models.VisitorTypeGuest {
		message = "Guest vehicle activity logged successfully"
	}

	rd := utility.BuildSuccessResponse(code, message, nil)
	c.JSON(code, rd)
}

func GetVehicleLogs(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	logs, code, err := services.GetVehicleLogs(userID)
	if err != nil {
		log.Default().Println("Error fetching vehicle logs:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle logs", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("Successfully fetched vehicle logs for user ID:", userID)
	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched vehicle logs", logs)
	c.JSON(http.StatusOK, rd)
}

// GetUserVehicleActivities returns activities for vehicles owned by the authenticated user
func GetVehicleActivities(c *gin.Context) {
	vehicle_id := c.Param("vehicle_id")

	err := utility.ValidateUUID(vehicle_id)
	if err != nil {
		log.Default().Println("Invalid vehicle ID:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid vehicle ID", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	pagination := models.GetPagination(c)

	data, code, err := services.GetVehicleActivities(database.DB, vehicle_id, pagination)
	if err != nil {
		log.Default().Println("Error fetching vehicle activities:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle activities", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User vehicle activities retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "User vehicle activities retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func GetVehiclesActivities(c *gin.Context) {

	userID := c.MustGet("user_id").(string)

	pagination := models.GetPagination(c)

	data, code, err := services.GetVehiclesActivities(database.DB, userID, pagination)
	if err != nil {
		log.Default().Println("Error fetching vehicles activities:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicles activities", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User vehicles activities retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "User vehicles activities retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func GetPendingVehicles(c *gin.Context) {

	userID := c.MustGet("user_id").(string)

	pagination := models.GetPagination(c)

	data, code, err := services.GetPendingVehicles(database.DB, userID, pagination)
	if err != nil {
		log.Default().Println("Error fetching vehicles activities:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get pending vehicles", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User pending vehicles retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "User pending vehicles retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func UpdatePendingVehicle(c *gin.Context) {

	pending_id := c.Param("pending_id")

	var req models.PendingUpdateReq

	err := utility.ValidateUUID(pending_id)
	if err != nil {
		log.Default().Println("Invalid pending ID:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid pending ID", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid req", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(req); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID := c.MustGet("user_id").(string)
	req.UserID = userID
	req.ID = pending_id

	code, err := services.UpdatePendingVehicle(database.DB, req)
	if err != nil {
		log.Default().Println("Error fetching vehicles activities:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get pending vehicles", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("User pending vehicles updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "User pending vehicles updated successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func GetGuestVehicleActivitiesByPlateNumber(c *gin.Context) {
	plateNumber := c.Param("plateNumber")

	if plateNumber == "" {
		log.Default().Println("Plate number is required")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Plate number is required", "Please provide a valid plate number", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	pagination := models.GetPagination(c)

	data, code, err := services.GetGuestVehicleActivitiesByPlateNumber(database.DB, plateNumber, pagination)
	if err != nil {
		log.Default().Println("Error fetching guest vehicle activities:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get guest vehicle activities", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("Guest vehicle activities retrieved successfully for plate number:", plateNumber)
	rd := utility.BuildSuccessResponse(http.StatusOK, "Guest vehicle activities retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func IdentifyVehicle(c *gin.Context) {
	plateNumber := c.Param("plateNumber")

	resp, code, err := services.IdentifyVehicle(plateNumber)
	if err != nil {
		log.Default().Println("Error getting vehicle status:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle status", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	data := models.VehicleIdentity{
		PlateNumber:  plateNumber,
		Status:       resp.Status,
		IsRegistered: resp.IsRegistered,
	}

	log.Default().Println("Vehicle status retrieved successfully for plate number:", plateNumber)
	rd := utility.BuildSuccessResponse(http.StatusOK, "Vehicle status retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

// GetVehicleLogHistory returns detailed log history for a vehicle
func GetVehicleLogHistory(c *gin.Context) {
	plateNumber := c.Param("plateNumber")
	limitStr := c.DefaultQuery("limit", "50")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50 // Default limit
	}

	vehicle, _, err := services.GetVehicleByPlateNumber(plateNumber)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Vehicle not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	logs, err := services.GetVehicleLogHistory(vehicle.ID, limit)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to get vehicle log history", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data := map[string]any{
		"vehicle":    vehicle,
		"logs":       logs,
		"total_logs": len(logs),
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Vehicle log history retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func LogGuestVehicleActivity(c *gin.Context) {
	var input models.LogVehicleActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Default().Println("Validation error:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	input.VisitorType = models.VisitorTypeGuest

	code, err := services.LogGuestVehicleActivity(database.DB, input)
	if err != nil {
		log.Default().Println("Error logging guest vehicle activity:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to log guest vehicle activity", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("Guest vehicle activity logged successfully")
	rd := utility.BuildSuccessResponse(code, "Guest vehicle activity logged successfully", nil)
	c.JSON(code, rd)
}

func FetchRegisteredVehiclesLogs(c *gin.Context) {
	pagination := models.GetPagination(c)

	// Get filter parameters
	plateNumber := c.Query("plate_number")
	model := c.Query("model")
	color := c.Query("color")
	vehicleType := c.Query("type")

	filters := models.VehicleFilters{
		PlateNumber: plateNumber,
		Model:       model,
		Color:       color,
		Type:        vehicleType,
	}

	response, statusCode, err := services.FetchRegisteredVehiclesLogs(database.DB, pagination, filters)
	if err != nil {
		log.Default().Println("Failed to fetch registered vehicles logs:", err)
		rd := utility.BuildErrorResponse(statusCode, "error", "Failed to fetch registered vehicles logs", err.Error(), nil)
		c.JSON(statusCode, rd)
		return
	}

	log.Default().Println("Successfully fetched registered vehicles logs")
	rd := utility.BuildSuccessResponse(statusCode, "Successfully fetched registered vehicles logs", response.Data, response.Pagination)
	c.JSON(statusCode, rd)
}

func FetchGuestVehiclesLogs(c *gin.Context) {
	pagination := models.GetPagination(c)
	plateNumber := c.Query("plate_number")

	response, statusCode, err := services.FetchGuestVehiclesLogs(database.DB, pagination, plateNumber)
	if err != nil {
		log.Default().Println("Failed to fetch guest vehicles logs:", err)
		rd := utility.BuildErrorResponse(statusCode, "error", "Failed to fetch guest vehicles logs", err.Error(), nil)
		c.JSON(statusCode, rd)
		return
	}

	log.Default().Println("Successfully fetched guest vehicles logs")
	rd := utility.BuildSuccessResponse(statusCode, "Successfully fetched guest vehicles logs", response.Data, response.Pagination)
	c.JSON(statusCode, rd)
}

// GetVehicleOwnerProfile returns vehicle owner profile and vehicle activity logs
func GetVehicleOwnerProfile(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	err := utility.ValidateUUID(vehicleID)
	if err != nil {
		log.Default().Println("Invalid vehicle ID:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid vehicle ID", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
	}

	data, code, err := services.GetVehicleOwnerProfile(database.DB, vehicleID, pagination)
	if err != nil {
		log.Default().Println("Error fetching vehicle owner profile:", err)
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle owner profile", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log.Default().Println("Vehicle owner profile retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Vehicle owner profile retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

func GenerateActivityReport(c *gin.Context) {
	pagination := models.GetPagination(c)

	// Get date range parameters
	fromStr := c.Query("from")
	toStr := c.Query("to")
	visitorTypeStr := c.Query("visitor_type")

	var from, to time.Time
	var err error

	// Default to last 30 days if no dates provided
	if fromStr == "" {
		from = time.Now().AddDate(0, 0, -30)
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid from date format. Use YYYY-MM-DD", err.Error(), nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
	}

	if toStr == "" {
		to = time.Now()
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid to date format. Use YYYY-MM-DD", err.Error(), nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
	}

	var visitorType *models.VisitorType
	if visitorTypeStr != "" {
		vt := models.VisitorType(visitorTypeStr)
		if vt != models.VisitorTypeRegistered && vt != models.VisitorTypeGuest {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid visitor type. Use 'registered' or 'guest'", nil, nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
		visitorType = &vt
	}

	response, statusCode, err := services.GenerateActivityReport(database.DB, from, to, visitorType, pagination)
	if err != nil {
		log.Default().Println("Failed to generate activity report:", err)
		rd := utility.BuildErrorResponse(statusCode, "error", "Failed to generate activity report", err.Error(), nil)
		c.JSON(statusCode, rd)
		return
	}

	log.Default().Println("Successfully generated activity report")
	rd := utility.BuildSuccessResponse(statusCode, "Successfully generated activity report", response.Data, response.Pagination)
	c.JSON(statusCode, rd)
}

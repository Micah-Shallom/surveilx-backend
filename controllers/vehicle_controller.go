package controllers

import (
	"log"
	"net/http"
	"strconv"
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"
	"time"

	"github.com/gin-gonic/gin"
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

	// Get the user ID of who is logging this activity
	loggedByUserID := c.MustGet("user_id").(string)

	activity, code, err := services.LogVehicleActivity(database.DB, input, loggedByUserID)
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

	log.Default().Println(message, activity.PlateNumber, "by user ID:", loggedByUserID)
	rd := utility.BuildSuccessResponse(code, message, activity)
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
func GetUserVehicleActivities(c *gin.Context) {
	vehicle_id := c.Param("vehicle_id")

	err := utility.ValidateUUID(vehicle_id)
	if err != nil {
		log.Default().Println("Invalid vehicle ID:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid vehicle ID", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data, code, err := services.GetVehicleActivities(database.DB, vehicle_id)
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

func GetVehicleStatus(c *gin.Context) {
	plateNumber := c.Param("plateNumber")

	status, err := services.GetVehicleStatusByPlateNumber(plateNumber)
	if err != nil {
		log.Default().Println("Error getting vehicle status:", err)
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Failed to get vehicle status", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	data := map[string]any{
		"plate_number": plateNumber,
		"status":       status,
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
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to get vehicle log history", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
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


// GetActivityReport returns activity report for a date range
func GetActivityReport(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	visitorTypeStr := c.Query("visitor_type") // optional filter

	if fromStr == "" || toStr == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Missing date parameters", "Both 'from' and 'to' query parameters are required", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'from' date format", "Use RFC3339 format (e.g., 2024-01-01T00:00:00Z)", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'to' date format", "Use RFC3339 format (e.g., 2024-01-01T23:59:59Z)", nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var visitorType *models.VisitorType
	if visitorTypeStr != "" {
		vt := models.VisitorType(visitorTypeStr)
		if vt == models.VisitorTypeRegistered || vt == models.VisitorTypeGuest {
			visitorType = &vt
		} else {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid visitor type", "visitor_type must be 'registered' or 'guest'", nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
	}

	activities, err := services.GetAllVehicleActivities(from, to, visitorType)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to get activity report", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	// Generate summary statistics
	summary := services.GenerateActivitySummary(activities)

	data := map[string]any{
		"activities": activities,
		"summary":    summary,
		"period": map[string]any{
			"from": from,
			"to":   to,
		},
		"filters": map[string]any{
			"visitor_type": visitorTypeStr,
		},
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Activity report retrieved successfully", data)
	c.JSON(http.StatusOK, rd)
}

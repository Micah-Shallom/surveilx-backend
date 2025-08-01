package controllers

import (
	"net/http"
	"strconv"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func RegisterVehicle(c *gin.Context) {
	var input models.RegisterVehicleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
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
		rd := utility.BuildErrorResponse(code, "error", "Failed to register vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Vehicle successfully registered", createdVehicle)
	c.JSON(code, rd)
}

func LogVehicle(c *gin.Context) {
	var input models.LogVehicleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := validate.Struct(input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vehicle, code, err := services.GetVehicleByPlateNumber(input.PlateNumber)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to find vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	log, code, err := services.CreateVehicleLog(vehicle, input)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to log vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(code, "Vehicle log successfully created", log)
	c.JSON(code, rd)
}

func GetVehicleLogs(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	logs, code, err := services.GetVehicleLogs(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle logs", err.Error(), nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched vehicle logs", logs)
	c.JSON(http.StatusOK, rd)
}

func GetVehicleStatus(c *gin.Context) {
	plateNumber := c.Param("plateNumber")

	status, err := services.GetVehicleStatusByPlateNumber(plateNumber)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Failed to get vehicle status", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	data := map[string]any{
		"plate_number": plateNumber,
		"status":       status,
	}

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

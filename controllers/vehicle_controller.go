package controllers

import (
	"net/http"
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

	userResponse := models.UserResponse{
		ID:        createdVehicle.User.ID,
		Name:      createdVehicle.User.Name,
		Email:     createdVehicle.User.Email,
		Role:      createdVehicle.User.Role,
		CreatedAt: createdVehicle.User.CreatedAt,
	}
	response := models.VehicleResponse{
		ID:          createdVehicle.ID,
		UserID:      createdVehicle.UserID,
		User:        userResponse,
		PlateNumber: createdVehicle.PlateNumber,
		Type:        createdVehicle.Type,
		Model:       createdVehicle.Model,
		Color:       createdVehicle.Color,
		CreatedAt:   createdVehicle.CreatedAt,
	}
	rd := utility.BuildSuccessResponse(code, "Vehicle successfully registered", response)
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

	log, code, err := services.CreateVehicleLog(vehicle, input.IsEntry)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to log vehicle", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	userResponse := models.UserResponse{
		ID:        log.User.ID,
		Name:      log.User.Name,
		Email:     log.User.Email,
		Role:      log.User.Role,
		CreatedAt: log.User.CreatedAt,
	}

	vehicleResponse := models.VehicleResponse{
		ID:          log.Vehicle.ID,
		UserID:      log.Vehicle.UserID,
		PlateNumber: log.Vehicle.PlateNumber,
		Type:        log.Vehicle.Type,
		Model:       log.Vehicle.Model,
		Color:       log.Vehicle.Color,
		CreatedAt:   log.Vehicle.CreatedAt,
	}

	response := models.VehicleLogResponse{
		ID:        log.ID,
		VehicleID: log.VehicleID,
		Vehicle:   vehicleResponse,
		UserID:    log.UserID,
		User:      userResponse,
		Timestamp: log.Timestamp,
		IsEntry:   log.IsEntry,
		CreatedAt: log.CreatedAt,
	}

	rd := utility.BuildSuccessResponse(code, "Vehicle log successfully created", response)
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

	var logResponses []models.VehicleLogResponse
	for _, log := range logs {
		userResponse := models.UserResponse{
			ID:        log.User.ID,
			Name:      log.User.Name,
			Email:     log.User.Email,
			Role:      log.User.Role,
			CreatedAt: log.User.CreatedAt,
		}

		vehicleResponse := models.VehicleResponse{
			ID:          log.Vehicle.ID,
			UserID:      log.Vehicle.UserID,
			PlateNumber: log.Vehicle.PlateNumber,
			Type:        log.Vehicle.Type,
			Model:       log.Vehicle.Model,
			Color:       log.Vehicle.Color,
			CreatedAt:   log.Vehicle.CreatedAt,
		}
		logResponses = append(logResponses, models.VehicleLogResponse{
			ID:        log.ID,
			VehicleID: log.VehicleID,
			Vehicle:   vehicleResponse,
			UserID:    log.UserID,
			User:      userResponse,
			Timestamp: log.Timestamp,
			IsEntry:   log.IsEntry,
			CreatedAt: log.CreatedAt,
		})
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched vehicle logs", logResponses)
	c.JSON(http.StatusOK, rd)
}

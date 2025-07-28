package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func RegisterGuest(c *gin.Context) {
	var input models.RegisterInput
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

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     "guest",
	}

	createdUser, code, err := services.Register(&user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to register guest", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	response := models.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt,
	}

	rd := utility.BuildSuccessResponse(code, "Guest successfully registered", response)
	c.JSON(code, rd)
}

func RegisterGuestVehicle(c *gin.Context) {
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

	userID := c.Param("userID")
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

func GetGuestRegistrations(c *gin.Context) {
	var guests []models.User
	if err := services.GetUsers(&guests); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to get guests", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	var guestResponses []models.UserResponse
	for _, g := range guests {
		if g.Role == "guest" {
			guestResponses = append(guestResponses, models.UserResponse{
				ID:        g.ID,
				Name:      g.Name,
				Email:     g.Email,
				Role:      g.Role,
				CreatedAt: g.CreatedAt,
			})
		}
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched guest registrations", guestResponses)
	c.JSON(http.StatusOK, rd)
}

func GetGuestVehicleLogs(c *gin.Context) {
	logs, code, err := services.GetAllVehicleLogs()
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to get vehicle logs", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched guest vehicle logs", logs)
	c.JSON(http.StatusOK, rd)
}

package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterGuestInput struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	Model       string `json:"model"`
	Color       string `json:"color"`
	Type        string `json:"type" validate:"oneof=bus car bike"`
}

func RegisterGuest(c *gin.Context) {
	var input RegisterGuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID := c.MustGet("user_id").(string)
	guest := models.Guest{
		PlateNumber:  input.PlateNumber,
		Model:        input.Model,
		Color:        input.Color,
		Type:         input.Type,
		RegisteredBy: userID,
	}

	createdGuest, code, err := services.RegisterGuest(&guest)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to register guest", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Guest successfully registered", createdGuest)
	c.JSON(code, rd)
}

func GetGuests(c *gin.Context) {
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")

	if fromStr == "" || toStr == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Missing 'from' or 'to' query parameters", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'from' date format", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'to' date format", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	guests, code, err := services.GetGuests(from, to)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to get guests", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched guests", guests)
	c.JSON(http.StatusOK, rd)
}

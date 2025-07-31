package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func LogGuest(c *gin.Context) {
	var input models.LogGuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	log, code, err := services.CreateGuestLog(&input)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to log guest", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Guest log successfully created", log)
	c.JSON(code, rd)
}

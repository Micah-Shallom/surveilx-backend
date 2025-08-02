package controllers

import (
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func UpdateUserProfile(c *gin.Context) {
	var input models.UpdateUserProfileInput

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
	code, err := services.UpdateUserProfile(database.DB, userID, &input)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to update profile", err.Error(), nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(code, "Profile successfully updated", nil)
	c.JSON(code, rd)
}

func GetUserProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	userProfile, code, err := services.GetUserProfile(database.DB, userID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to fetch user profile", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Profile successfully fetched", userProfile)
	c.JSON(code, rd)
}

package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	profile, err := services.GetProfile(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched profile", profile)
	c.JSON(http.StatusOK, rd)
}

func UpdateProfile(c *gin.Context) {
	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID := c.MustGet("user_id").(string)
	profile, err := services.GetProfile(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	if err := services.UpdateProfile(&profile, input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update user", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully updated profile", profile)
	c.JSON(http.StatusOK, rd)
}

func DeleteProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	if err := services.DeleteProfile(userID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete user", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully deleted profile", nil)
	c.JSON(http.StatusOK, rd)
}

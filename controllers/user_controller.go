package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := services.GetUsers(&users); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to get users", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched users", users)
	c.JSON(http.StatusOK, rd)
}

func GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	user, err := services.GetUserByID(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched profile", user)
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
	user, err := services.GetUserByID(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	user.Name = input.Name

	if err := services.UpdateUser(&user); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update user", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully updated profile", user)
	c.JSON(http.StatusOK, rd)
}

func DeleteProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	if err := services.DeleteUser(userID); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to delete user", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully deleted profile", nil)
	c.JSON(http.StatusOK, rd)
}

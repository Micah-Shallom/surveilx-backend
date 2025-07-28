package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := services.CreateUser(&user); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create user", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "User successfully created", user)
	c.JSON(http.StatusCreated, rd)
}

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

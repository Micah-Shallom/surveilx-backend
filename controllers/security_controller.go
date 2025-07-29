package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func RegisterSecurity(c *gin.Context) {
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
		Role:     "security",
	}

	createdUser, code, err := services.Register(&user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to register security personnel", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	// Create a security record
	security := models.Security{
		UserID: createdUser.ID,
	}
	if err := services.CreateSecurity(&security); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to create security record", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Security personnel successfully registered", createdUser)
	c.JSON(code, rd)
}

func LoginSecurity(c *gin.Context) {
	var input models.LoginInput
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

	user, code, err := services.Login(input.Email, input.Password)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Login failed", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	if user.Role != "security" {
		rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "Forbidden", "User is not a security personnel", nil)
		c.JSON(http.StatusForbidden, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Login successful", user)
	c.JSON(http.StatusOK, rd)
}

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

	rd := utility.BuildSuccessResponse(code, "Guest successfully registered", createdUser)
	c.JSON(code, rd)
}

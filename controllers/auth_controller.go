package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Register(c *gin.Context) {
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

	role := "user"
	if input.Role != "" {
		role = input.Role
	}
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     role,
	}

	createdUser, code, err := services.Register(&user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to register user", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "User successfully registered", createdUser)
	c.JSON(code, rd)
}

func Login(c *gin.Context) {
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

	rd := utility.BuildSuccessResponse(http.StatusOK, "Login successful", user)
	c.JSON(http.StatusOK, rd)
}

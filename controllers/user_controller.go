package controllers

import (
	"log"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {

	users, err := services.GetUsers(database.DB)
	if err != nil {
		log.Default().Println("Error fetching users:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to get users", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	log.Default().Println("Successfully fetched users")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched users", users)
	c.JSON(http.StatusOK, rd)
}

func GetUser(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	user, err := services.GetUserByID(userID)
	if err != nil {
		log.Default().Println("Error fetching user:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to get user", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	log.Default().Println("Successfully fetched user")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched user", user)
	c.JSON(http.StatusOK, rd)
}
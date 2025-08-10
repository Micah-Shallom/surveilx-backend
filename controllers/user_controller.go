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
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to get users", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	log.Default().Println("Successfully fetched users")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched users", users)
	c.JSON(http.StatusOK, rd)
}

package controllers

import (
	"log"
	"net/http"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

type VerifyVehicleInput struct {
	PlateNumber string `json:"plate_number" binding:"required"`
}

func VerifyVehicle(c *gin.Context) {
	var input VerifyVehicleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vehicle, code, err := services.GetVehicleByPlateNumber(input.PlateNumber)
	if err != nil {
		log.Default().Println("Error fetching vehicle by plate number:", err)
		rd := utility.BuildErrorResponse(code, "error", "Vehicle has not been registered", err.Error(), nil)
		c.JSON(code, rd)
		return
	}
	
	log.Default().Println("Vehicle verification successful for plate number:", input.PlateNumber)
	rd := utility.BuildSuccessResponse(http.StatusOK, "Vehicle is registered", vehicle)
	c.JSON(http.StatusOK, rd)
}

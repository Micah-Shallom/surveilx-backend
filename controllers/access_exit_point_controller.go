package controllers

import (
	"log"
	"net/http"
	"survielx-backend/database"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func CreateAccessExitPoint(c *gin.Context) {
	var point models.AccessExitPoint
	if err := c.ShouldBindJSON(&point); err != nil {
		log.Default().Println("Error binding JSON:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := services.CreateAccessExitPoint(database.DB, &point); err != nil {
		log.Default().Println("Error creating access exit point:", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to create access exit point", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	log.Default().Println("Access exit point created successfully:", point.Name)
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Access exit point created successfully", point)
	c.JSON(http.StatusCreated, rd)
}

func GetAccessExitPoints(c *gin.Context) {
	var points []models.AccessExitPoint
	if err := services.GetAccessExitPoints(database.DB, &points); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to get access exit points", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Access exit points retrieved successfully", points)
	c.JSON(http.StatusOK, rd)
}

func GetAccessExitPoint(c *gin.Context) {
	id := c.Param("id")
	var point models.AccessExitPoint
	if err := services.GetAccessExitPoint(database.DB, id, &point); err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Access exit point not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Access exit point retrieved successfully", point)
	c.JSON(http.StatusOK, rd)
}

func DeleteAccessExitPoint(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteAccessExitPoint(database.DB, id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to delete access exit point", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Access exit point deleted successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func UpdateAccessExitPoint(c *gin.Context) {
	id := c.Param("id")
	var point models.AccessExitPoint
	if err := services.GetAccessExitPoint(database.DB, id, &point); err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Access exit point not found", err.Error(), nil)
		c.JSON(http.StatusNotFound, rd)
		return
	}

	if err := c.ShouldBindJSON(&point); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := services.UpdateAccessExitPoint(database.DB, &point); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to update access exit point", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Access exit point updated successfully", point)
	c.JSON(http.StatusOK, rd)
}

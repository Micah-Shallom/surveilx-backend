package controllers

import (
	"net/http"
	"survielx-backend/models"
	"survielx-backend/services"
	"survielx-backend/utility"
	"time"

	"github.com/gin-gonic/gin"
)

type AddToWatchlistInput struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	Model       string `json:"model"`
	Color       string `json:"color"`
	Type        string `json:"type" validate:"oneof=bus car bike"`
}

func AddToWatchlist(c *gin.Context) {
	var input AddToWatchlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid input", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userID := c.MustGet("user_id").(string)
	watchlist := models.Watchlist{
		PlateNumber:  input.PlateNumber,
		Model:        input.Model,
		Color:        input.Color,
		Type:         input.Type,
		RegisteredBy: userID,
	}

	createdWatchlist, code, err := services.AddToWatchlist(&watchlist)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to add to watchlist", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(code, "Successfully added to watchlist", createdWatchlist)
	c.JSON(code, rd)
}

func GetWatchlist(c *gin.Context) {
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")

	if fromStr == "" || toStr == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Missing 'from' or 'to' query parameters", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'from' date format", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid 'to' date format", err.Error(), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	watchlist, code, err := services.GetWatchlist(from, to)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "Failed to get watchlist", err.Error(), nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Successfully fetched watchlist", watchlist)
	c.JSON(http.StatusOK, rd)
}

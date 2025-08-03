package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(r *gin.Engine, ApiVersion string) {
	r.GET(ApiVersion+"/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "SurvielX Backend API is running smoothly.",
		})
	})
}

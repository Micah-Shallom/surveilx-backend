package routers

import (
	"survielx-backend/controllers"

	"github.com/gin-gonic/gin"
)

func VerifyRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	// authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/verify-vehicle", controllers.VerifyVehicle)
	}
}

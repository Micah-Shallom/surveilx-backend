package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VerifyRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	// authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/verify-vehicle", controllers.VerifyVehicle)
	}
}

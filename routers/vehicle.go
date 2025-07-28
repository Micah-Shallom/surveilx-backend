package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VehicleRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/vehicles", controllers.RegisterVehicle)
		authorized.POST("/vehicles/log", controllers.LogVehicle)
		authorized.GET("/vehicles/logs", controllers.GetVehicleLogs)
	}
}

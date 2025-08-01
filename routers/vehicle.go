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
		// Vehicle registration and logging
		authorized.POST("/vehicles", controllers.RegisterVehicle)
		authorized.POST("/vehicles/log", controllers.LogVehicle)
		authorized.GET("/vehicles/logs", controllers.GetVehicleLogs)
		
		// Vehicle status and monitoring
		authorized.GET("/vehicles/status/:plateNumber", controllers.GetVehicleStatus) //returns the current status of a vehicle (inside/outside)
		authorized.GET("/vehicles/history/:plateNumber", controllers.GetVehicleLogHistory) //returns the log history of a vehicle
	}
}
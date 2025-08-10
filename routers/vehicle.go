package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VehicleActivityRoutes(r *gin.Engine, api_version string) {

	// Regular user routes
	activityRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version))
	activityRoutes.Use(middleware.AuthMiddleware())
	{
		// Vehicle management
		activityRoutes.POST("/register", controllers.RegisterVehicle)
		activityRoutes.GET("/fetch_vehicles", controllers.GetUserVehicles)
		activityRoutes.GET("/:vehicle_id/activities", controllers.GetVehicleActivities)
		activityRoutes.GET("/guest/activities/:plateNumber", controllers.GetGuestVehicleActivitiesByPlateNumber)
		activityRoutes.POST("/log-vehicle", controllers.LogVehicleActivity) //keep here temporarily for testing

	}

	securityRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version))
	securityRoutes.Use(middleware.SecurityMiddleware())
	{
		// securityRoutes.POST("/log-vehicle", controllers.LogVehicleActivity)
	}

	unauthRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version))
	unauthRoutes.GET("/identify/:plateNumber", controllers.IdentifyVehicle)
	unauthRoutes.POST("/sys-log-vehicle", controllers.SystemLogVehicleActivity) //the model backend logs vehicle activity without user context

}

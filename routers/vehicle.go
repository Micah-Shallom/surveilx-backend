package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VehicleActivityRoutes(r *gin.Engine, api_version string) {

	// Regular user routes
	activityRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version), middleware.AuthMiddleware())
	{
		// Vehicle management
		activityRoutes.POST("/register", controllers.RegisterVehicle)
		activityRoutes.DELETE("/:vehicle_id/deregister", controllers.DeRegisterVehicle)
		activityRoutes.GET("/fetch_vehicles", controllers.GetUserVehicles)
		activityRoutes.GET("/:vehicle_id/activities", controllers.GetVehicleActivities)
	}

	securityRoutes := r.Group(fmt.Sprintf("%v/sec/vehicles", api_version), middleware.SecurityMiddleware())
	{
		//security personnels should be able to get all registerd and guest vehicles
		securityRoutes.POST("/log-vehicle", controllers.LogVehicleActivity)
		securityRoutes.GET("/fetch_vehicle_logs", controllers.FetchVehiclesLogs)
		activityRoutes.GET("/:vehicle_id/activities", controllers.GetVehicleActivities)


	}

	unauthRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version))
	{
		unauthRoutes.GET("/guest/activities/:plateNumber", controllers.GetGuestVehicleActivitiesByPlateNumber)
		unauthRoutes.GET("/identify/:plateNumber", controllers.IdentifyVehicle)
		unauthRoutes.POST("/sys-log-vehicle", controllers.SystemLogVehicleActivity) //the model backend logs vehicle activity without user context
	}
}

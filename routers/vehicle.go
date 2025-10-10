package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"survielx-backend/controllers"
	"survielx-backend/middleware"
)

func VehicleActivityRoutes(r *gin.Engine, api_version string) {

	activityRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version), middleware.AuthMiddleware())
	{
		activityRoutes.POST("/register", controllers.RegisterVehicle)
		activityRoutes.PATCH("/:vehicle_id", controllers.UpdateVehicle)
		activityRoutes.DELETE("/:vehicle_id/deregister", controllers.DeRegisterVehicle)
		activityRoutes.GET("/fetch_vehicles", controllers.GetUserVehicles)
		activityRoutes.GET("/activities", controllers.GetVehiclesActivities)
		activityRoutes.GET("/pending", controllers.GetPendingVehicles)
		activityRoutes.PUT("/pending/:pending_id", controllers.GetPendingVehicles)
		activityRoutes.GET("/:vehicle_id/activities", controllers.GetVehicleActivities)
	}

	securityRoutes := r.Group(fmt.Sprintf("%v/security", api_version), middleware.AuthMiddleware(), middleware.SecurityMiddleware())
	{
		securityRoutes.POST("/log-vehicle", controllers.LogVehicleActivity)
		securityRoutes.POST("/log-guest-vehicle", controllers.LogGuestVehicleActivity)
		securityRoutes.GET("/vehicle/:vehicle_id/activities", controllers.GetVehicleActivities)
		securityRoutes.GET("/activities/:plateNumber", controllers.GetGuestVehicleActivitiesByPlateNumber)
		securityRoutes.GET("/registered-logs", controllers.FetchRegisteredVehiclesLogs)
		securityRoutes.GET("/guest-logs", controllers.FetchGuestVehiclesLogs)
		securityRoutes.GET("/:vehicle_id/owner-profile", controllers.GetVehicleOwnerProfile)
		securityRoutes.GET("/activity-report", controllers.GenerateActivityReport)
	}

	unauthRoutes := r.Group(fmt.Sprintf("%v/vehicles", api_version))
	{
		unauthRoutes.GET("/identify/:plateNumber", controllers.IdentifyVehicle)
		unauthRoutes.POST("/sys-log-vehicle", controllers.SystemLogVehicleActivity) //the model backend logs vehicle activity without user context
	}
}

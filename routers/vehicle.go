package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VehicleActivityRoutes(r *gin.Engine, api_version string) {
	// Public routes
	// r.POST(fmt.Sprintf("%v/

	// Regular user routes
	user := r.Group(fmt.Sprintf("%v/users", api_version))
	user.Use(middleware.AuthMiddleware())
	{
		// Vehicle management
		user.POST("/vehicles", controllers.RegisterVehicle)
		user.POST("/vehicles/log-entry-exit", controllers.LogVehicleActivity) // For registered vehicles only

		// User's vehicle activities
		user.GET("/vehicles/my-activities", controllers.GetUserVehicleActivities)
		user.GET("/vehicles/status/:plateNumber", controllers.GetVehicleStatus)
	}

	// Security personnel routes
	security_admin := r.Group(fmt.Sprintf("%v/security", api_version))
	security_admin.Use(middleware.AuthMiddleware(), middleware.SecurityMiddleware())
	{

		// Guest vehicle management
		security_admin.POST("/vehicles/log-vehicle", controllers.LogVehicleActivity)

		// Monitoring and reports
		security_admin.GET("/vehicles/activities", controllers.GetVehicleActivities)
		security_admin.GET("/vehicles/status/:plateNumber", controllers.GetVehicleStatus)
		security_admin.GET("/reports/activity", controllers.GetActivityReport)
	}

	// Admin routes (if needed)
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // Add admin middleware if you have one
	{
		// Could add endpoints for managing access points, users, etc.
	}
}

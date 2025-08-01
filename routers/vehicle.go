package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VehicleActivityRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/verify-vehicle", controllers.VerifyVehicle)

	// Regular user routes
	user := r.Group("/user")
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
	security_admin := r.Group("/security")
	security_admin.Use(middleware.AuthMiddleware(), middleware.SecurityMiddleware())
	{
		// Guest vehicle management
		security_admin.POST("/vehicles/log-guest", controllers.LogVehicleActivity)
		security_admin.POST("/vehicles/log-registered", controllers.LogVehicleActivity)
		
		// Monitoring and reports
		security_admin.GET("/vehicles/activities", controllers.GetVehicleActivities)
		security_admin.GET("/vehicles/status/:plateNumber", controllers.GetVehicleStatus)
		security_admin.GET("/reports/activity", controllers.GetActivityReport)
		security_admin.GET("/reports/daily-summary", controllers.GetActivityReport) // Could be a separate endpoint
	}

	// Admin routes (if needed)
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // Add admin middleware if you have one
	{
		// Could add endpoints for managing access points, users, etc.
	}
}

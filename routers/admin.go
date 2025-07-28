package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("/guests", controllers.RegisterGuest)
		admin.POST("/guests/:userID/vehicles", controllers.RegisterGuestVehicle)
		admin.GET("/guests", controllers.GetGuestRegistrations)
		admin.GET("/guests/logs", controllers.GetGuestVehicleLogs)
	}
}

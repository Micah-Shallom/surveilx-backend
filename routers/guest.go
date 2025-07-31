package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func GuestRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/guests", middleware.SecurityMiddleware(), controllers.RegisterGuest)
		authorized.GET("/guests", controllers.GetGuests)
	}
}

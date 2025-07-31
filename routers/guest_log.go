package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func GuestLogRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/guestlogs", middleware.SecurityMiddleware(), controllers.LogGuest)
	}
}

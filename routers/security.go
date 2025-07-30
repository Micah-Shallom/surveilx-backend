package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SecurityRoutes(r *gin.Engine) {
	security := r.Group("/security")
	{
		security.POST("/register", controllers.RegisterSecurity)
		security.POST("/login", controllers.LoginSecurity)
	}

	guest := r.Group("/guest")
	guest.Use(middleware.AuthMiddleware(), middleware.SecurityMiddleware())
	{
		guest.POST("/register", controllers.RegisterGuest)
	}
}

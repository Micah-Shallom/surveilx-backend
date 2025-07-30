package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/profile", controllers.GetProfile)
		authorized.PUT("/profile", controllers.UpdateProfile)
		authorized.DELETE("/profile", controllers.DeleteProfile)
	}
}

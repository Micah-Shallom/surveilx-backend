package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AccessExitPointRoutes(r *gin.Engine) {
	access := r.Group("/access-points")
	access.Use(middleware.AuthMiddleware())
	{
		access.POST("/", controllers.CreateAccessExitPoint)
		access.GET("/", controllers.GetAccessExitPoints)
		access.GET("/:id", controllers.GetAccessExitPoint)
		access.DELETE("/:id", controllers.DeleteAccessExitPoint)
		access.PUT("/:id", middleware.SecurityMiddleware(), controllers.UpdateAccessExitPoint)
	}
}

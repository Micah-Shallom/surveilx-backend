package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AccessExitPointRoutes(r *gin.Engine, api_version string) {
	access := r.Group(fmt.Sprintf("%v/access-exit-points", api_version))
	access.Use(middleware.AuthMiddleware())
	{
		access.POST("/", controllers.CreateAccessExitPoint, middleware.SecurityMiddleware())
		access.GET("/", controllers.GetAccessExitPoints)
		access.GET("/:id", controllers.GetAccessExitPoint)
		access.DELETE("/:id", controllers.DeleteAccessExitPoint, middleware.SecurityMiddleware())
		access.PUT("/:id", middleware.SecurityMiddleware(), controllers.UpdateAccessExitPoint)
	}
}
																
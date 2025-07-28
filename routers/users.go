package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/users", controllers.GetUsers)
	}
}

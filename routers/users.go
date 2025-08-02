package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(r *gin.Engine, api_version string) {
	authorized := r.Group(fmt.Sprintf("%v", api_version))
	authorized.Use(middleware.AuthMiddleware(), middleware.SecurityMiddleware())
	{
		authorized.GET("/users", controllers.GetUsers)
	}
}

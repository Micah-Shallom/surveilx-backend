package routers

import (
	"fmt"
	"survielx-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, api_version string) {
	authRoutes := r.Group(fmt.Sprintf("%v/auth", api_version))
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}
}

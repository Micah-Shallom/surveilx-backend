package routers

import (
	"fmt"
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserProfileRoutes(r *gin.Engine, api_version string) {
	profileRoutes := r.Group(fmt.Sprintf("%v/profile", api_version))
	profileRoutes.Use(middleware.AuthMiddleware())
	{
		// Current user profile operations
		profileRoutes.PUT("/", controllers.UpdateUserProfile)
		profileRoutes.GET("/", controllers.GetUserProfile)
	}

}

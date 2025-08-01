package routers

import (
	"github.com/gin-gonic/gin"
)

func UserProfileRoutes(r *gin.Engine) {
	// profile := r.Group("/profile")
	// profile.Use(middleware.AuthMiddleware())
	// {
	// 	// Current user profile operations
	// 	profile.PUT("/", controllers.UpdateUserProfile)
	// 	profile.GET("/", controllers.GetUserProfile)
	// }

	// // Admin routes for user management
	// adminUsers := r.Group("/admin/users")
	// adminUsers.Use(middleware.AuthMiddleware())
	// {
	// 	adminUsers.GET("/:id", controllers.GetUserByIDAdmin)
	// }
}

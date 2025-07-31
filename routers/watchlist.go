package routers

import (
	"survielx-backend/controllers"
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func WatchlistRoutes(r *gin.Engine) {
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/watchlist", middleware.SecurityMiddleware(), controllers.AddToWatchlist)
		authorized.GET("/watchlist", controllers.GetWatchlist)
	}
}

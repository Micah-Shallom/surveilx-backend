package routers

import (
	"survielx-backend/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(r *gin.Engine) {
	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.GetUsers)
}

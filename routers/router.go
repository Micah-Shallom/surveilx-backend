package routers

import (
	"boilerplate/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.GetUsers)

	return r
}

package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	AuthRoutes(r)
	UsersRoutes(r)

	return r
}

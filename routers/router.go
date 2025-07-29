package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	AuthRoutes(r)
	UsersRoutes(r)
	VehicleRoutes(r)
	SecurityRoutes(r)

	return r
}

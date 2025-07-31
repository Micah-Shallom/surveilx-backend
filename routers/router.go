package routers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	AuthRoutes(r)
	UsersRoutes(r)
	VehicleRoutes(r)
	AccessExitPointRoutes(r)
	VerifyRoutes(r)
	GuestRoutes(r)

	return r
}

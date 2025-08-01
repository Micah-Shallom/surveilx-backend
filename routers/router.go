package routers

import (
	"survielx-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	AuthRoutes(r)
	UsersRoutes(r)
	VehicleActivityRoutes(r)
	AccessExitPointRoutes(r)

	return r
}

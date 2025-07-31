package middleware

import (
	"net/http"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
)

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)
		user, err := services.GetUserByID(userID)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Unauthorized", "Invalid user", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
			return
		}

		if user.Role != "security" {
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "Forbidden", "Security access required", nil)
			c.AbortWithStatusJSON(http.StatusForbidden, rd)
			return
		}

		c.Next()
	}
}
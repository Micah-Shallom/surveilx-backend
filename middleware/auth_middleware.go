package middleware

import (
	"net/http"
	"os"
	"strings"
	"survielx-backend/services"
	"survielx-backend/utility"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Authorization header is required", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Invalid token format", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			if err == jwt.ErrTokenExpired {
				// Attempt to refresh the token
				newToken, err := services.RefreshToken(tokenString)
				if err != nil {
					rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Failed to refresh token", err.Error(), nil)
					c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
					return
				}
				c.Header("X-Refreshed-Token", newToken)
				claims, _ := token.Claims.(jwt.MapClaims)
				c.Set("user_id", claims["sub"])
				c.Next()
				return

			}
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Invalid token", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["sub"])
			c.Next()
		} else {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "Invalid token", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
		}
	}
}

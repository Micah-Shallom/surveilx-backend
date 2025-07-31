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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, "error", "User ID not found in context", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rd)
			return
		}

		user, err := services.GetUserByID(userID.(string))
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "User not found", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusNotFound, rd)
			return
		}

		if user.Role != "security" {
			rd := utility.BuildErrorResponse(http.StatusForbidden, "error", "Forbidden", "You are not authorized to perform this action", nil)
			c.AbortWithStatusJSON(http.StatusForbidden, rd)
			return
		}

		c.Next()
	}
}

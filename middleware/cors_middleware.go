package middleware

import (
	"github.com/gin-gonic/gin"
)

// func CORSMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         origin := c.Request.Header.Get("Origin")
//         if origin != "" {
//             c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
//         }
//         c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//         c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
//         c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

//         if c.Request.Method == "OPTIONS" {
//             c.AbortWithStatus(204)
//             return
//         }
//         c.Next()
//     }
// }

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        allowedOrigins := []string{
            "http://localhost:3000",
            "https://abunrfedusec.com",
            "https://www.abunrfedusec.com",
            "https://abunrf-security-portal.vercel.app/",
            "https://www.abunrf-security-portal.vercel.app/",
        }

        for _, o := range allowedOrigins {
            if o == origin {
                c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
                break
            }
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}


// func CORSMiddlewareProduction(allowedOrigins []string) gin.HandlerFunc {
// 	return gin.HandlerFunc(func(c *gin.Context) {
// 		origin := c.Request.Header.Get("Origin")

// 		// Check if origin is in allowed origins
// 		allowed := false
// 		for _, allowedOrigin := range allowedOrigins {
// 			if origin == allowedOrigin {
// 				allowed = true
// 				break
// 			}
// 		}

// 		if allowed {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 		}

// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	})
// }

package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS enables Cross-Origin Resource Sharing for web dashboard and external API clients
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from web dashboard domain and localhost (development)
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"https://api-krisho.aarcsx.com":       true,
			"https://krisho.aarcsx.com":           true,
			"https://dashboard.krisho.aarcsx.com": true,
			"http://localhost:3000":               true, // Development Next.js
			"http://localhost:8080":               true, // Development API
		}

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Trace-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Trace-ID, Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

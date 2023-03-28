package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func WhitelistIPsMiddleware() gin.HandlerFunc {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get whitelist IPs from environment variable
	whitelistIPs := strings.Split(os.Getenv("WHITELIST_IPS"), ",")

	return func(c *gin.Context) {
		// Get client IP address
		clientIP := c.ClientIP()

		// Check if client IP is in whitelist
		for _, ip := range whitelistIPs {
			if ip == "0.0.0.0" {
				c.Next()
				return
			}
			if ip == clientIP {
				// Client IP is in whitelist, continue processing request
				c.Next()
				return
			}
		}

		// Client IP is not in whitelist, return 403 Forbidden error
		c.AbortWithStatus(http.StatusForbidden)
	}
}

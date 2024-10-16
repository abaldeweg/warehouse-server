package router

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/framework/apikey"
	"github.com/gin-gonic/gin"
)

// ApiKeyMiddleware is a middleware to check for API key authentication.
func ApiKeyMiddleware(c *gin.Context) {
	key := c.GetHeader("X-API-Key")

	if !apikey.IsValidAPIKey(key) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	c.Next()
}

// permissionsMiddleware is a middleware to check for API key permissions.
func PermissionsMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")

		for _, permission := range permissions {
			if !apikey.HasPermission(key, permission) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return

			}
		}

		c.Next()
	}
}

// Engine creates a gin engine with CORS and sets it to release mode.
func Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	return r
}

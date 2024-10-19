package router

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/framework/apikey"
	"github.com/gin-gonic/gin"
)

// Engine creates a gin engine with CORS and sets it to release mode.
func Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	return r
}

// ApiKeyMiddleware is a middleware to check for API key authentication.
func ApiKeyMiddleware(data []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")

		k, _ := apikey.NewAPIKeys(data)

		if !k.IsValidAPIKey(key) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Next()
	}
}

// permissionsMiddleware is a middleware to check for API key permissions.
func PermissionsMiddleware(data []byte, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")

		k, _ := apikey.NewAPIKeys(data)

		for _, permission := range permissions {
			if !k.HasPermission(key, permission) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return

			}
		}

		c.Next()
	}
}

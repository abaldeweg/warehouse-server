package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config represents CORS configuration options.
type Config struct {
	cors.Config
}

// NewCorsConfig creates a new CORS configuration with default values.
func NewCors() *Config {
	return &Config{
		Config: cors.Config{
			AllowOrigins:  []string{"http://127.0.0.1"},
			AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders: []string{"Content-Length"},
		},
	}
}

// SetCorsHeaders sets up CORS middleware with the provided configuration.
func (c *Config) SetCorsHeaders() gin.HandlerFunc {
	return cors.New(c.Config)
}

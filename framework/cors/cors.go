package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// DefaultConfig returns a generic default configuration mapped to localhost.
func DefaultConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")

	allowOrigin := "*"
	if viper.IsSet("CORS_ALLOW_ORIGIN") {
		allowOrigin = viper.GetString("CORS_ALLOW_ORIGIN")
	}
	config.AllowOrigins = []string{allowOrigin}

	return config
}

// SetDefaultCorsHeaders returns the location middleware with default configuration.
func SetDefaultCorsHeaders() gin.HandlerFunc {
	return cors.New(DefaultConfig())
}

// SetCorsHeaders sets up CORS middleware with the provided configuration.
func SetCorsHeaders(c cors.Config) gin.HandlerFunc {
	return cors.New(c)
}

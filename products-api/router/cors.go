package router

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func headers() gin.HandlerFunc {
	origins, exists := os.LookupEnv("CORS_ALLOW_ORIGIN")
	if !exists {
		origins = "*"
	}

	return cors.New(cors.Config{
		AllowOrigins:  strings.Split(origins, ","),
		AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	})
}

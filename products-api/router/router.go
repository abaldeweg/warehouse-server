package router

import (
	"os"

	"github.com/abaldeweg/warehouse-server/products/mock"
	"github.com/abaldeweg/warehouse-server/products/web"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.New()
	r.SetTrustedProxies(nil)

	if os.Getenv("ENV") != "prod" {
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())

	r.Use(headers())

	auth := r.Group("/"+os.Getenv("BASE_PATH")+"/api/v1", checkAuth)

	auth.GET("/products", web.List)
	auth.POST("/products", web.Create)
	auth.PUT("/products/:id", web.Update)
	auth.DELETE("/products/:id", web.Delete)

	if os.Getenv("ENV") != "prod" {
		r.GET("/api/v1/me", mock.Me)
	}

	r.Run(":8080")
}

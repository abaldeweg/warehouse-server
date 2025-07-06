package router

import (
	"github.com/abaldeweg/warehouse-server/framework/router"
	"github.com/abaldeweg/warehouse-server/framework/storage"
	"github.com/abaldeweg/warehouse-server/logs_web/controller"
	"github.com/gin-gonic/gin"
)

// Routes sets up the Gin router.
func Routes() *gin.Engine {
	s := storage.NewStorage()
	s.FileSystem.Path = "./data/auth"
	s.FileSystem.FileName = "api_keys.json"
	k, _ := s.Load()

	r := router.Engine()

	api := r.Group("/apis/logs/1", router.ApiKeyMiddleware(k))
	{
		api.POST("/events", controller.GetEvents)
		api.GET("/count", controller.Count)
	}

	return r
}

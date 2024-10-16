package router

import (
	"net/http"
	"path/filepath"

	"github.com/abaldeweg/warehouse-server/gateway/proxy"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Routes() *gin.Engine {
    gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Any(`/apis/core/1/*path`, func(c *gin.Context) {
		path := c.Param("path")

		safePath := filepath.Join("/", path)

		if err := proxy.Proxy(c, viper.GetString("API_CORE"), safePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
			return
		}
	})

	return r
}

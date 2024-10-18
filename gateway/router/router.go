package router

import (
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"github.com/abaldeweg/warehouse-server/gateway/proxy"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Any(`/apis/core/1/*path`, func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")

		path := c.Param("path")

		uploadCover(c, path)

		safePath := filepath.Join("/", path)

		if err := proxy.Proxy(c, viper.GetString("API_CORE"), safePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
			return
		}
	})

	return r
}

func uploadCover(c *gin.Context, path string) {
	if c.Request.Method == http.MethodPost {
		if match, _ := regexp.MatchString(`/api/book/cover/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`, path); match {
			if auth.Authenticate(c) {
				re := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
				imageUUID := re.FindString(path)

				cover.SaveCover(c, imageUUID)

				c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
				return
			}

			c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
			return
		}
	}
}

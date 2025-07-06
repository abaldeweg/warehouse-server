package router

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/blog/content/article"
	"github.com/abaldeweg/warehouse-server/blog/content/home"
	"github.com/abaldeweg/warehouse-server/framework/router"
	"github.com/abaldeweg/warehouse-server/framework/storage"
	"github.com/gin-gonic/gin"
)

// Routes sets up the Gin router.
func Routes() *gin.Engine {
	s := storage.NewStorage()
	s.FileSystem.Path = "./data/auth"
	s.FileSystem.FileName = "api_keys.json"
	k, _ := s.Load()

	r := router.Engine()

	api := r.Group("/", router.ApiKeyMiddleware(k))
	{
		api.GET("/home", router.PermissionsMiddleware(k, "articles"), func(c *gin.Context) {
			index, err := home.GetHome()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.String(http.StatusOK, index)
		})
		api.GET("/home/new/:days", router.PermissionsMiddleware(k, "articles"), func(c *gin.Context) {
			daysStr := c.Param("days")

			days, err := strconv.Atoi(daysStr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'days' parameter"})
				return
			}

			newCount, err := home.GetNewArticles(days)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"new_articles": newCount})
		})
		api.GET("/article/*path", router.PermissionsMiddleware(k, "articles"), func(c *gin.Context) {
			path := c.Param("path")

			cnt, err := article.GetArticle(path)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.String(http.StatusOK, cnt)
		})
	}

	return r
}

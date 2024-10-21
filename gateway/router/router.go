package router

import (
	"net/http"
	"path/filepath"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"github.com/abaldeweg/warehouse-server/gateway/proxy"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Routes sets up the routes for the gateway.
func Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Next()
	})

	apiCore := r.Group(`/apis/core/1`)
	{
		apiCoreAuthor := apiCore.Group(`/api/author`)
		{
			apiCoreAuthor.GET(`/find`, handleCoreAPI("/apis/core/1/api/author/find"))
			apiCoreAuthor.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/author"))
			apiCoreAuthor.POST(`/new`, handleCoreAPI("/apis/core/1/api/author/new"))
			apiCoreAuthor.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/author"))
			apiCoreAuthor.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/author"))
		}

		apiCoreBook := apiCore.Group(`/api/book`)
		{
			apiCoreBook.GET(`/find`, handleCoreAPI("/apis/core/1/api/book/find"))
			apiCoreBook.DELETE(`/clean`, handleCoreAPI("/apis/core/1/api/book/clean"))
			apiCoreBook.GET(`/stats`, handleCoreAPI("/apis/core/1/api/book/stats"))
			apiCoreBook.PUT(`/inventory/found/:id`, handleCoreAPIWithId("/apis/core/1/api/book/inventory/found"))
			apiCoreBook.PUT(`/inventory/notfound/:id`, handleCoreAPIWithId("/apis/core/1/api/book/inventory/notfound"))
			apiCoreBook.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/book"))
			apiCoreBook.POST(`/new`, handleCoreAPI("/apis/core/1/api/book/new"))
			apiCoreBook.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/book"))
			apiCoreBook.GET(`/cover/:id`, handleCoreAPIWithId("/apis/core/1/api/book/cover"))
			apiCoreBook.POST(`/cover/:id`, handleCover)
			apiCoreBook.DELETE(`/cover/:id`, handleCoreAPIWithId("/apis/core/1/api/book/cover"))
			apiCoreBook.PUT(`/sell/:id`, handleCoreAPIWithId("/apis/core/1/api/book/sell"))
			apiCoreBook.PUT(`/remove/:id`, handleCoreAPIWithId("/apis/core/1/api/book/remove"))
			apiCoreBook.PUT(`/reserve/:id`, handleCoreAPIWithId("/apis/core/1/api/book/reserve"))
			apiCoreBook.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/book"))
		}

		apiCoreBranch := apiCore.Group(`/api/branch`)
		{
			apiCoreBranch.GET(`/`, handleCoreAPI("/apis/core/1/api/branch/"))
			apiCoreBranch.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/branch"))
			apiCoreBranch.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/branch"))
		}

		apiCoreCondition := apiCore.Group(`/api/condition`)
		{
			apiCoreCondition.GET(`/`, handleCoreAPI("/apis/core/1/api/condition/"))
			apiCoreCondition.POST(`/new`, handleCoreAPI("/apis/core/1/api/condition/new"))
			apiCoreCondition.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/condition"))
			apiCoreCondition.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/condition"))
			apiCoreCondition.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/condition"))
		}

		apiCoreDirectory := apiCore.Group(`/api/directory`)
		{
			apiCoreDirectory.GET(`/`, handleCoreAPI("/apis/core/1/api/directory/"))
			apiCoreDirectory.POST(`/cover/:id`, handleCoreAPIWithId("/apis/core/1/api/directory/cover"))
			apiCoreDirectory.POST(`/new`, handleCoreAPI("/apis/core/1/api/directory/new"))
			apiCoreDirectory.POST(`/upload`, handleCoreAPI("/apis/core/1/api/directory/upload"))
			apiCoreDirectory.PUT(`/edit`, handleCoreAPI("/apis/core/1/api/directory/edit"))
			apiCoreDirectory.DELETE(`/`, handleCoreAPI("/apis/core/1/api/directory/"))
		}

		apiCoreFormat := apiCore.Group(`/api/format`)
		{
			apiCoreFormat.GET(`/`, handleCoreAPI("/apis/core/1/api/format/"))
			apiCoreFormat.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/format"))
			apiCoreFormat.POST(`/new`, handleCoreAPI("/apis/core/1/api/format/new"))
			apiCoreFormat.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/format"))
			apiCoreFormat.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/format"))
		}

		apiCoreGenre := apiCore.Group(`/api/genre`)
		{
			apiCoreGenre.GET(`/`, handleCoreAPI("/apis/core/1/api/genre/"))
			apiCoreGenre.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/genre"))
			apiCoreGenre.POST(`/new`, handleCoreAPI("/apis/core/1/api/genre/new"))
			apiCoreGenre.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/genre"))
			apiCoreGenre.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/genre"))
		}

		apiCoreInventory := apiCore.Group(`/api/inventory`)
		{
			apiCoreInventory.GET(`/`, handleCoreAPI("/apis/core/1/api/inventory/"))
			apiCoreInventory.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/inventory"))
			apiCoreInventory.POST(`/new`, handleCoreAPI("/apis/core/1/api/inventory/new"))
			apiCoreInventory.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/inventory"))
			apiCoreInventory.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/inventory"))
		}

		apiCore.GET(`/api/me`, handleCoreAPI("/apis/core/1/api/me"))
		apiCore.PUT(`/api/password`, handleCoreAPI("/apis/core/1/api/password"))

		apiCorePublic := apiCore.Group(`/api/public`)
		{
			apiCorePublicBook := apiCorePublic.Group(`/book`)
			{
				apiCorePublicBook.GET(`/find`, handleCoreAPI("/apis/core/1/api/public/book/find"))
				apiCorePublicBook.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/public/book"))
				apiCorePublicBook.GET(`/recommendation/:id`, handleCoreAPIWithId("/apis/core/1/api/public/book/recommendation"))
				apiCorePublicBook.GET(`/cover/:id`, handleCoreAPIWithId("/apis/core/1/api/public/book/cover"))
			}
			apiCorePublic.GET(`/branch/`, handleCoreAPI("/apis/core/1/api/public/branch/"))
			apiCorePublic.GET(`/branch/show/:id`, handleCoreAPIWithId("/apis/core/1/api/public/branch/show"))
			apiCorePublic.GET(`/genre/:id`, handleCoreAPIWithId("/apis/core/1/api/public/genre"))
			apiCorePublic.POST(`/reservation/new`, handleCoreAPI("/apis/core/1/api/public/reservation/new"))
		}

		apiCoreReservation := apiCore.Group(`/api/reservation`)
		{
			apiCoreReservation.GET(`/list`, handleCoreAPI("/apis/core/1/api/reservation/list"))
			apiCoreReservation.GET(`/status`, handleCoreAPI("/apis/core/1/api/reservation/status"))
			apiCoreReservation.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/reservation"))
			apiCoreReservation.POST(`/new`, handleCoreAPI("/apis/core/1/api/reservation/new"))
			apiCoreReservation.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/reservation"))
			apiCoreReservation.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/reservation"))
		}

		apiCoreTag := apiCore.Group(`/api/tag`)
		{
			apiCoreTag.GET(`/`, handleCoreAPI("/apis/core/1/api/tag/"))
			apiCoreTag.GET(`/:id`, handleCoreAPIWithId("/apis/core/1/api/tag"))
			apiCoreTag.POST(`/new`, handleCoreAPI("/apis/core/1/api/tag/new"))
			apiCoreTag.PUT(`/:id`, handleCoreAPIWithId("/apis/core/1/api/tag"))
			apiCoreTag.DELETE(`/:id`, handleCoreAPIWithId("/apis/core/1/api/tag"))
		}
	}

	return r
}

// handleCoreAPI handles requests to the core API.
func handleCoreAPI(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		safePath := filepath.Join("/", path)

		if err := proxy.Proxy(c, viper.GetString("API_CORE"), safePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
			return
		}
	}
}

// handleCoreAPIWithId handles requests to the core API.
func handleCoreAPIWithId(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		safePath := filepath.Join("/", path, id)

		if err := proxy.Proxy(c, viper.GetString("API_CORE"), safePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
			return
		}
	}
}

// handleCover handles requests to the core API.
func handleCover(c *gin.Context) {
	imageUUID := c.Param("id")

	validate := validator.New()
	err := validate.Var(imageUUID, "uuid")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid UUID"})
		return
	}

	if auth.Authenticate(c) {
		cover.SaveCover(c, imageUUID)

		c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
}

// RoleMiddleware ensures that the user has the specified role before allowing access.
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, ok := c.Get("user")
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
            return
        }

        userRoles := user.(auth.User).Roles

        for _, role := range userRoles {
            if role == requiredRole {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
    }
}


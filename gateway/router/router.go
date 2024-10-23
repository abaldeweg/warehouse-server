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
			apiCoreAuthor.GET(`/find`, handleCoreAPI("/api/author/find"))
			apiCoreAuthor.GET(`/:id`, handleCoreAPIWithId("/api/author"))
			apiCoreAuthor.POST(`/new`, handleCoreAPI("/api/author/new"))
			apiCoreAuthor.PUT(`/:id`, handleCoreAPIWithId("/api/author"))
			apiCoreAuthor.DELETE(`/:id`, handleCoreAPIWithId("/api/author"))
		}

		apiCoreBook := apiCore.Group(`/api/book`)
		{
			apiCoreBook.GET(`/find`, handleCoreAPI("/api/book/find"))
			apiCoreBook.DELETE(`/clean`, handleCoreAPI("/api/book/clean"))
			apiCoreBook.GET(`/stats`, handleCoreAPI("/api/book/stats"))
			apiCoreBook.PUT(`/inventory/found/:id`, handleCoreAPIWithId("/api/book/inventory/found"))
			apiCoreBook.PUT(`/inventory/notfound/:id`, handleCoreAPIWithId("/api/book/inventory/notfound"))
			apiCoreBook.GET(`/:id`, handleCoreAPIWithId("/api/book"))
			apiCoreBook.POST(`/new`, handleCoreAPI("/api/book/new"))
			apiCoreBook.PUT(`/:id`, handleCoreAPIWithId("/api/book"))
			apiCoreBook.GET(`/cover/:id`, handleCoreAPIWithId("/api/book/cover"))
			apiCoreBook.POST(`/cover/:id`, handleCover)
			apiCoreBook.DELETE(`/cover/:id`, handleCoreAPIWithId("/api/book/cover"))
			apiCoreBook.PUT(`/sell/:id`, handleCoreAPIWithId("/api/book/sell"))
			apiCoreBook.PUT(`/remove/:id`, handleCoreAPIWithId("/api/book/remove"))
			apiCoreBook.PUT(`/reserve/:id`, handleCoreAPIWithId("/api/book/reserve"))
			apiCoreBook.DELETE(`/:id`, handleCoreAPIWithId("/api/book"))
		}

		apiCoreBranch := apiCore.Group(`/api/branch`)
		{
			apiCoreBranch.GET(`/`, handleCoreAPI("/api/branch/"))
			apiCoreBranch.GET(`/:id`, handleCoreAPIWithId("/api/branch"))
			apiCoreBranch.PUT(`/:id`, handleCoreAPIWithId("/api/branch"))
		}

		apiCoreCondition := apiCore.Group(`/api/condition`)
		{
			apiCoreCondition.GET(`/`, handleCoreAPI("/api/condition/"))
			apiCoreCondition.POST(`/new`, handleCoreAPI("/api/condition/new"))
			apiCoreCondition.GET(`/:id`, handleCoreAPIWithId("/api/condition"))
			apiCoreCondition.PUT(`/:id`, handleCoreAPIWithId("/api/condition"))
			apiCoreCondition.DELETE(`/:id`, handleCoreAPIWithId("/api/condition"))
		}

		apiCoreDirectory := apiCore.Group(`/api/directory`)
		{
			apiCoreDirectory.GET(`/`, handleCoreAPI("/api/directory/"))
			apiCoreDirectory.POST(`/cover/:id`, handleCoreAPIWithId("/api/directory/cover"))
			apiCoreDirectory.POST(`/new`, handleCoreAPI("/api/directory/new"))
			apiCoreDirectory.POST(`/upload`, handleCoreAPI("/api/directory/upload"))
			apiCoreDirectory.PUT(`/edit`, handleCoreAPI("/api/directory/edit"))
			apiCoreDirectory.DELETE(`/`, handleCoreAPI("/api/directory/"))
		}

		apiCoreFormat := apiCore.Group(`/api/format`)
		{
			apiCoreFormat.GET(`/`, handleCoreAPI("/api/format/"))
			apiCoreFormat.GET(`/:id`, handleCoreAPIWithId("/api/format"))
			apiCoreFormat.POST(`/new`, handleCoreAPI("/api/format/new"))
			apiCoreFormat.PUT(`/:id`, handleCoreAPIWithId("/api/format"))
			apiCoreFormat.DELETE(`/:id`, handleCoreAPIWithId("/api/format"))
		}

		apiCoreGenre := apiCore.Group(`/api/genre`)
		{
			apiCoreGenre.GET(`/`, handleCoreAPI("/api/genre/"))
			apiCoreGenre.GET(`/:id`, handleCoreAPIWithId("/api/genre"))
			apiCoreGenre.POST(`/new`, handleCoreAPI("/api/genre/new"))
			apiCoreGenre.PUT(`/:id`, handleCoreAPIWithId("/api/genre"))
			apiCoreGenre.DELETE(`/:id`, handleCoreAPIWithId("/api/genre"))
		}

		apiCoreInventory := apiCore.Group(`/api/inventory`)
		{
			apiCoreInventory.GET(`/`, handleCoreAPI("/api/inventory/"))
			apiCoreInventory.GET(`/:id`, handleCoreAPIWithId("/api/inventory"))
			apiCoreInventory.POST(`/new`, handleCoreAPI("/api/inventory/new"))
			apiCoreInventory.PUT(`/:id`, handleCoreAPIWithId("/api/inventory"))
			apiCoreInventory.DELETE(`/:id`, handleCoreAPIWithId("/api/inventory"))
		}

		apiCore.GET(`/api/me`, handleCoreAPI("/api/me"))
		apiCore.POST(`/api/login_check`, handleCoreAPI("/api/login_check"))
		apiCore.PUT(`/api/password`, handleCoreAPI("/api/password"))

		apiCorePublic := apiCore.Group(`/api/public`)
		{
			apiCorePublicBook := apiCorePublic.Group(`/book`)
			{
				apiCorePublicBook.GET(`/find`, handleCoreAPI("/api/public/book/find"))
				apiCorePublicBook.GET(`/:id`, handleCoreAPIWithId("/api/public/book"))
				apiCorePublicBook.GET(`/recommendation/:id`, handleCoreAPIWithId("/api/public/book/recommendation"))
				apiCorePublicBook.GET(`/cover/:id`, handleCoreAPIWithId("/api/public/book/cover"))
			}
			apiCorePublic.GET(`/branch/`, handleCoreAPI("/api/public/branch/"))
			apiCorePublic.GET(`/branch/show/:id`, handleCoreAPIWithId("/api/public/branch/show"))
			apiCorePublic.GET(`/genre/:id`, handleCoreAPIWithId("/api/public/genre"))
			apiCorePublic.POST(`/reservation/new`, handleCoreAPI("/api/public/reservation/new"))
		}

		apiCoreReservation := apiCore.Group(`/api/reservation`)
		{
			apiCoreReservation.GET(`/list`, handleCoreAPI("/api/reservation/list"))
			apiCoreReservation.GET(`/status`, handleCoreAPI("/api/reservation/status"))
			apiCoreReservation.GET(`/:id`, handleCoreAPIWithId("/api/reservation"))
			apiCoreReservation.POST(`/new`, handleCoreAPI("/api/reservation/new"))
			apiCoreReservation.PUT(`/:id`, handleCoreAPIWithId("/api/reservation"))
			apiCoreReservation.DELETE(`/:id`, handleCoreAPIWithId("/api/reservation"))
		}

		apiCoreTag := apiCore.Group(`/api/tag`)
		{
			apiCoreTag.GET(`/`, handleCoreAPI("/api/tag/"))
			apiCoreTag.GET(`/:id`, handleCoreAPIWithId("/api/tag"))
			apiCoreTag.POST(`/new`, handleCoreAPI("/api/tag/new"))
			apiCoreTag.PUT(`/:id`, handleCoreAPIWithId("/api/tag"))
			apiCoreTag.DELETE(`/:id`, handleCoreAPIWithId("/api/tag"))
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


package router

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/core/controllers"
	"github.com/abaldeweg/warehouse-server/gateway/core/database"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"github.com/abaldeweg/warehouse-server/gateway/proxy"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var authenticator = auth.Authenticate

// Routes sets up the routes for the gateway.
func Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Next()
	})

	db := database.Connect()

	apiCore := r.Group(`/apis/core/1`)
	{
		apiCoreAuthor := apiCore.Group(`/api/author`)
		{
			apiCoreAuthor.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreAuthor.GET(`/find`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewAuthorController(db)
				ac.GetAuthors(c)
			})

			apiCoreAuthor.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewAuthorController(db)
				ac.GetAuthor(c)
			})
			apiCoreAuthor.POST(`/new`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewAuthorController(db)
				ac.CreateAuthor(c)
			})
			apiCoreAuthor.PUT(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewAuthorController(db)
				ac.UpdateAuthor(c)
			})
			apiCoreAuthor.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ac := controllers.NewAuthorController(db)
				ac.DeleteAuthor(c)
			})
		}

		// @fix: port to new API
		apiCoreBook := apiCore.Group(`/api/book`)
		{
			apiCoreBook.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreBook.GET(`/find`, handleCoreAPI("/api/book/find"))
			// apiCoreBook.DELETE(`/clean`, handleCoreAPI("/api/book/clean"))
			apiCoreBook.DELETE(`/clean`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.CleanBooks(c)
			})
			apiCoreBook.GET(`/stats`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.ShowStats(c)
			})
			// apiCoreBook.PUT(`/inventory/found/:id`, handleCoreAPIWithId("/api/book/inventory/found"))
			apiCoreBook.PUT(`/inventory/found/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.FindInventory(c)
			})
			// apiCoreBook.PUT(`/inventory/notfound/:id`, handleCoreAPIWithId("/api/book/inventory/notfound"))
			apiCoreBook.PUT(`/inventory/notfound/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.NotFoundInventory(c)
			})
			apiCoreBook.GET(`/:id`, handleCoreAPIWithId("/api/book"))
			apiCoreBook.POST(`/new`, handleCoreAPI("/api/book/new"))
			apiCoreBook.PUT(`/:id`, handleCoreAPIWithId("/api/book"))
			apiCoreBook.GET(`/cover/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.ShowCover(c)
			})
			apiCoreBook.POST(`/cover/:id`, handleCover)
			apiCoreBook.DELETE(`/cover/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.DeleteCover(c)
			})
			// apiCoreBook.PUT(`/sell/:id`, handleCoreAPIWithId("/api/book/sell"))
			apiCoreBook.PUT(`/sell/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBookController(db)
				bc.SellBook(c)
			})
			apiCoreBook.PUT(`/remove/:id`, handleCoreAPIWithId("/api/book/remove"))
			apiCoreBook.PUT(`/reserve/:id`, handleCoreAPIWithId("/api/book/reserve"))
			apiCoreBook.DELETE(`/:id`, handleCoreAPIWithId("/api/book"))
		}

		apiCoreBranch := apiCore.Group(`/api/branch`)
		{
			apiCoreBranch.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreBranch.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBranchController(db)
				bc.List(c)
			})
			apiCoreBranch.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewBranchController(db)
				bc.Show(c)
			})
			apiCoreBranch.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), IsOwnBranchMiddleware(), func(c *gin.Context) {
				bc := controllers.NewBranchController(db)
				bc.Update(c)
			})
		}

		apiCoreCondition := apiCore.Group(`/api/condition`)
		{
			apiCoreCondition.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreCondition.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewConditionController(db)
				bc.FindAll(c)
			})
			apiCoreCondition.POST(`/new`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				bc := controllers.NewConditionController(db)
				bc.Create(c)
			})
			apiCoreCondition.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				bc := controllers.NewConditionController(db)
				bc.FindOne(c)
			})
			apiCoreCondition.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				bc := controllers.NewConditionController(db)
				bc.Update(c)
			})
			apiCoreCondition.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				bc := controllers.NewConditionController(db)
				bc.Delete(c)
			})
		}

		apiCoreFormat := apiCore.Group(`/api/format`)
		{
			apiCoreFormat.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreFormat.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				fc := controllers.NewFormatController(db)
				fc.FindAll(c)
			})
			apiCoreFormat.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				fc := controllers.NewFormatController(db)
				fc.FindOne(c)
			})
			apiCoreFormat.POST(`/new`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				fc := controllers.NewFormatController(db)
				fc.Create(c)
			})
			apiCoreFormat.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				fc := controllers.NewFormatController(db)
				fc.Update(c)
			})
			apiCoreFormat.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				fc := controllers.NewFormatController(db)
				fc.Delete(c)
			})
		}

		apiCoreGenre := apiCore.Group(`/api/genre`)
		{
			apiCoreGenre.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreGenre.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewGenreController(db)
				ac.FindAll(c)
			})
			apiCoreGenre.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ac := controllers.NewGenreController(db)
				ac.FindOne(c)
			})
			apiCoreGenre.POST(`/new`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ac := controllers.NewGenreController(db)
				ac.Create(c)
			})
			apiCoreGenre.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ac := controllers.NewGenreController(db)
				ac.Update(c)
			})
			apiCoreGenre.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ac := controllers.NewGenreController(db)
				ac.Delete(c)
			})
		}

		apiCoreInventory := apiCore.Group(`/api/inventory`)
		{
			apiCoreInventory.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreInventory.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ic := controllers.NewInventoryController(db)
				ic.List(c)
			})
			apiCoreInventory.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				ic := controllers.NewInventoryController(db)
				ic.Show(c)
			})
			apiCoreInventory.POST(`/new`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ic := controllers.NewInventoryController(db)
				ic.Create(c)
			})
			apiCoreInventory.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ic := controllers.NewInventoryController(db)
				ic.Update(c)
			})
			apiCoreInventory.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				ic := controllers.NewInventoryController(db)
				ic.Delete(c)
			})
		}

		// @fix: port to new API
		apiCore.GET(`/api/me`, handleCoreAPI("/api/me"))
		apiCore.POST(`/api/login_check`, handleCoreAPI("/api/login_check"))
		apiCore.PUT(`/api/password`, handleCoreAPI("/api/password"))

		// apiCorePublic := apiCore.Group(`/api/public`)
		// {
		// 	apiCorePublicBook := apiCorePublic.Group(`/book`)
		// 	{
		// 		apiCorePublicBook.GET(`/find`, handleCoreAPI("/api/public/book/find"))
		// 		apiCorePublicBook.GET(`/:id`, handleCoreAPIWithId("/api/public/book"))
		// 		apiCorePublicBook.GET(`/recommendation/:id`, handleCoreAPIWithId("/api/public/book/recommendation"))
		// 		apiCorePublicBook.GET(`/cover/:id`, handleCoreAPIWithId("/api/public/book/cover"))
		// 	}
		// 	apiCorePublic.GET(`/branch/`, handleCoreAPI("/api/public/branch/"))
		// 	apiCorePublic.GET(`/branch/show/:id`, handleCoreAPIWithId("/api/public/branch/show"))
		// 	apiCorePublic.GET(`/genre/:id`, handleCoreAPIWithId("/api/public/genre"))
		// 	apiCorePublic.POST(`/reservation/new`, handleCoreAPI("/api/public/reservation/new"))
		// }

		apiCorePublic := apiCore.Group(`/api/public`)
		{
			apiCorePublicBook := apiCorePublic.Group(`/book`)
			apiCorePublicBook.GET(`/find`, handleCoreAPI("/api/public/book/find"))
			{
				apiCorePublicBook.GET(`/:id`, func(c *gin.Context) {
					pbc := controllers.NewPublicBookController(db)
					pbc.Show(c)
				})
				apiCorePublicBook.GET(`/recommendation/:branch`, func(c *gin.Context) {
					pbc := controllers.NewPublicBookController(db)
					pbc.Recommendation(c)
				})
				apiCorePublicBook.GET(`/cover/:image`, func(c *gin.Context) {
					pbc := controllers.NewPublicBookController(db)
					pbc.Image(c)
				})
			}
			apiCorePublic.GET(`/branch/`, func(c *gin.Context) {
				ac := controllers.NewPublicBranchController(db)
				ac.GetBranches(c)
			})
			apiCorePublic.GET(`/branch/show/:id`, func(c *gin.Context) {
				ac := controllers.NewPublicBranchController(db)
				ac.GetBranch(c)
			})
			apiCorePublic.GET(`/genre/:id`, handleCoreAPIWithId("/api/public/genre"))
			apiCorePublic.POST(`/reservation/new`, handleCoreAPI("/api/public/reservation/new"))
		}

		apiCoreReservation := apiCore.Group(`/api/reservation`)
		{
			apiCoreReservation.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreReservation.GET(`/list`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.FindAll(c)
			})
			apiCoreReservation.GET(`/status`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.ReservationStatus(c)
			})
			apiCoreReservation.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.FindOne(c)
			})
			apiCoreReservation.POST(`/new`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.Create(c)
			})
			apiCoreReservation.PUT(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.Update(c)
			})
			apiCoreReservation.DELETE(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				rc := controllers.NewReservationController(db)
				rc.Delete(c)
			})
		}

		// apiCoreReservation := apiCore.Group(`/api/reservation`)
		// {
		// 	apiCoreReservation.GET(`/list`, handleCoreAPI("/api/reservation/list"))
		// 	apiCoreReservation.GET(`/status`, handleCoreAPI("/api/reservation/status"))
		// 	apiCoreReservation.GET(`/:id`, handleCoreAPIWithId("/api/reservation"))
		// 	apiCoreReservation.POST(`/new`, handleCoreAPI("/api/reservation/new"))
		// 	apiCoreReservation.PUT(`/:id`, handleCoreAPIWithId("/api/reservation"))
		// 	apiCoreReservation.DELETE(`/:id`, handleCoreAPIWithId("/api/reservation"))
		// }

		apiCoreTag := apiCore.Group(`/api/tag`)
		{
			apiCoreTag.Use(func(c *gin.Context) {
				if !authenticator(c) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
					return
				}
				c.Next()
			})

			apiCoreTag.GET(`/`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				tc := controllers.NewTagController(db)
				tc.FindAll(c)
			})
			apiCoreTag.GET(`/:id`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				tc := controllers.NewTagController(db)
				tc.FindOne(c)
			})
			apiCoreTag.POST(`/new`, RoleMiddleware("ROLE_USER"), func(c *gin.Context) {
				tc := controllers.NewTagController(db)
				tc.Create(c)
			})
			apiCoreTag.PUT(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				tc := controllers.NewTagController(db)
				tc.Update(c)
			})
			apiCoreTag.DELETE(`/:id`, RoleMiddleware("ROLE_ADMIN"), func(c *gin.Context) {
				tc := controllers.NewTagController(db)
				tc.Delete(c)
			})
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
	imageID := c.Param("id")

	validate := validator.New()
	errUUID := validate.Var(imageID, "uuid")
	errInt := validate.Var(imageID, "numeric")

	if errUUID != nil && errInt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID"})
		return
	}

	if authenticator(c) {
		cover.SaveCover(c, imageID)

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

// IsOwnBranchMiddleware ensures that the user is associated with the branch being accessed.
// It checks if the user's branch ID matches the branch ID provided in the request context.
func IsOwnBranchMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
			return
		}

		branchId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
			return
		}

		if user.(auth.User).Branch.Id == branchId {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
	}
}

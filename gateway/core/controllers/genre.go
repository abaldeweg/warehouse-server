package controllers

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GenreController struct defines the database connection.
type GenreController struct {
	DB        *gorm.DB
	GenreRepo repository.GenreRepository
}

// NewGenreController creates a new instance of GenreController.
func NewGenreController(db *gorm.DB) *GenreController {
	return &GenreController{
		DB:        db,
		GenreRepo: repository.NewGenreRepository(db),
	}
}

// FindAll retrieves all genres.
func (gc *GenreController) FindAll(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	genres, err := gc.GenreRepo.FindAllByBranchID(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// FindOne retrieves a specific genre by ID.
func (gc *GenreController) FindOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	genre, err := gc.GenreRepo.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// Create creates a new genre.
func (gc *GenreController) Create(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	genre.BranchID = uint(user.(auth.User).Branch.Id)

	if !genre.Validate(gc.DB) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := gc.GenreRepo.Create(&genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// Update updates an existing genre.
func (gc *GenreController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	existingGenre, err := gc.GenreRepo.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != existingGenre.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	existingGenre.Name = genre.Name

	if !existingGenre.Validate(gc.DB) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := gc.GenreRepo.Update(uint(id), &existingGenre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update genre"})
		return
	}

	c.JSON(http.StatusOK, existingGenre)
}

// Delete deletes a genre by ID.
func (gc *GenreController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	genre, err := gc.GenreRepo.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != genre.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	if err := gc.GenreRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete genre"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

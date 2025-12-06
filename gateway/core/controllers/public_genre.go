package controllers

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PublicGenreController struct defines the database connection.
type PublicGenreController struct {
	DB         *gorm.DB
	GenreRepo  repository.GenreRepository
	BranchRepo *repository.BranchRepository
}

// NewPublicGenreController creates a new instance of PublicGenreController.
func NewPublicGenreController(db *gorm.DB) *PublicGenreController {
	return &PublicGenreController{
		DB:         db,
		GenreRepo:  repository.NewGenreRepository(db),
		BranchRepo: repository.NewBranchRepository(db),
	}
}

// FindAll retrieves all genres by the given branch ID.
func (gc *PublicGenreController) FindAll(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	branchID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	branch, err := gc.BranchRepo.FindOne(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	if !branch.Public {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access to this branch is restricted"})
		return
	}

	genres, err := gc.GenreRepo.FindAllByBranchID(uint(branchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

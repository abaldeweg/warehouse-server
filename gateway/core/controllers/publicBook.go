package controllers

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicBookController struct defines the database connection.
type PublicBookController struct {
	DB   *gorm.DB
	Repo *repository.PublicBookRepository
}

// NewPublicBookController creates a new instance of PublicBookController.
func NewPublicBookController(db *gorm.DB) *PublicBookController {
	return &PublicBookController{
		DB:   db,
		Repo: repository.NewPublicBookRepository(db),
	}
}

// Show retrieves a public book by its ID.
func (pbc *PublicBookController) Show(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format"})
		return
	}

	var book models.PublicBook
	if err := pbc.DB.Preload("Branch").Preload("Genre").Preload("Condition").Preload("Format").Preload("Author").First(&book, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

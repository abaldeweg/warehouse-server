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

	if book.Sold || book.Removed || book.Reserved {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Recommendation retrieves recommended books for a specific branch.
func (pbc *PublicBookController) Recommendation(c *gin.Context) {
	branchID := c.Param("branch")
	var branch models.Branch
	if err := pbc.DB.First(&branch, "id = ?", branchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Branch not found"})
		return
	}

	if !branch.Public {
		c.JSON(http.StatusOK, gin.H{"books": []models.PublicBook{}, "counter": 0})
		return
	}

	var books []models.PublicBook
	if err := pbc.DB.Preload("Branch").Preload("Genre").Preload("Condition").Preload("Format").Preload("Author").Where("branch_id = ? AND sold = ? AND removed = ? AND reserved = ? AND recommendation = ?", branchID, false, false, false, true).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books, "counter": len(books)})
}

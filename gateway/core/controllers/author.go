package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// AuthorController struct for author controller.
type AuthorController struct {
	repo *repository.AuthorRepository
	v    *validator.Validate
}

// NewAuthorController creates a new author controller.
func NewAuthorController(db *gorm.DB) *AuthorController {
	return &AuthorController{
		repo: repository.NewAuthorRepository(db),
		v:    validator.New(),
	}
}

// GetAuthors
func (ac *AuthorController) GetAuthors(c *gin.Context) {
	term := c.Query("term")
	term = strings.ReplaceAll(term, "%", "")
	term = strings.ReplaceAll(term, "*", "")

	authors, err := ac.repo.FindAllByTerm(term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive authors"})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// GetAuthor
func (ac *AuthorController) GetAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	author, err := ac.repo.FindOneById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// CreateAuthor
func (ac *AuthorController) CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if err := ac.v.Struct(author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Valid"})
		return
	}

	if err := ac.repo.Create(&author); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, author)
}

// UpdateAuthor
func (ac *AuthorController) UpdateAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	author.ID = id

	if err := ac.v.Struct(author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Valid"})
		return
	}

	if err := ac.repo.Update(&author); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// DeleteAuthor
func (ac *AuthorController) DeleteAuthor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	if err := ac.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

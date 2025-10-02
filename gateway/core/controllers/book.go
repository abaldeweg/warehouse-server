package controllers

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
  "net/http"
)

// BookController struct defines the database connection.
type BookController struct {
	DB   *gorm.DB
	Repo *repository.BookRepository
}

// NewBookController creates a new instance of BookController.
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{
		DB:   db,
		Repo: repository.NewBookRepository(db),
	}
}

// ShowCover retrieves the cover images.
func (pbc *BookController) ShowCover(ctx *gin.Context) {
	id := ctx.Param("id")
	bookID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid book id"})
		return
	}

	covers := map[string]string{
		"cover_s": cover.ShowCover("s", bookID),
		"cover_m": cover.ShowCover("m", bookID),
		"cover_l": cover.ShowCover("l", bookID),
	}

	ctx.JSON(http.StatusOK, covers)
}

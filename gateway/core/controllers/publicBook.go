package controllers

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
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

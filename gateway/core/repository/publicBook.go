package repository

import (
	"gorm.io/gorm"
)

// PublicBookRepository defines the methods for interacting with the PublicBook model.
type PublicBookRepository struct {
	DB *gorm.DB
}

// NewPublicBookRepository creates a new instance of PublicBookRepository.
func NewPublicBookRepository(db *gorm.DB) *PublicBookRepository {
	return &PublicBookRepository{DB: db}
}

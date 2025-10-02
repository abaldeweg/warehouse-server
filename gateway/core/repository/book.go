package repository

import (
	"gorm.io/gorm"
)

// BookRepository defines the methods for interacting with the Book model.
type BookRepository struct {
	DB *gorm.DB
}

// NewBookRepository creates a new instance of BookRepository.
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

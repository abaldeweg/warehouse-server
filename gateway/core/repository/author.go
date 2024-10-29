package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// AuthorRepository struct for author repository.
type AuthorRepository struct {
	db *gorm.DB
}

const limit = 100

// NewAuthorRepository creates a new author repository.
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// FindAllByTerm returns all authors by term.
func (r *AuthorRepository) FindAllByTerm(term string) ([]models.Author, error) {
	var authors []models.Author
	result := r.db.Where("firstname LIKE ? OR surname LIKE ? OR CONCAT(firstname, ' ', surname) LIKE ? OR CONCAT(surname, ' ', firstname) LIKE ? OR CONCAT(firstname, ',', surname) LIKE ? OR CONCAT(firstname, ', ', surname) LIKE ?", "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%").Limit(limit).Find(&authors)
	return authors, result.Error
}

// FindOneById returns one author by id.
func (r *AuthorRepository) FindOneById(id uint64) (models.Author, error) {
	var author models.Author
	result := r.db.First(&author, id)
	return author, result.Error
}

// Create an author.
func (r *AuthorRepository) Create(author *models.Author) error {
	result := r.db.Create(author)
	return result.Error
}

// Update an author.
func (r *AuthorRepository) Update(author *models.Author) error {
	result := r.db.Save(author)
	return result.Error
}

// Delete an author.
func (r *AuthorRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Author{}, id)
	return result.Error
}

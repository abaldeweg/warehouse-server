package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// GenreRepository defines the interface for interacting with Genre data.
type GenreRepository interface {
	FindAllByBranchID(uint) ([]models.Genre, error)
	FindOne(id uint) (models.Genre, error)
	Create(genre *models.Genre) error
	Update(id uint, genre *models.Genre) error
	Delete(id uint) error
}

type genreRepository struct {
	db *gorm.DB
}

// NewGenreRepository instantiates a new GenreRepository using GORM.
func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

// FindAllByBranchID retrieves all genres associated with a specific branch ID, ordered alphabetically by name.
func (r *genreRepository) FindAllByBranchID(branchID uint) ([]models.Genre, error) {
	var genres []models.Genre
	if err := r.db.Where("branch_id = ?", branchID).Order("name asc").Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

// FindOne retrieves a specific genre by ID.
func (r *genreRepository) FindOne(id uint) (models.Genre, error) {
	var genre models.Genre
	if err := r.db.First(&genre, id).Error; err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}

// Create creates a new genre.
func (r *genreRepository) Create(genre *models.Genre) error {
	return r.db.Create(genre).Error
}

// Update updates an existing genre.
func (r *genreRepository) Update(id uint, genre *models.Genre) error {
	return r.db.Model(&models.Genre{}).Where("id = ?", id).Updates(genre).Error
}

// Delete deletes a genre by ID.
func (r *genreRepository) Delete(id uint) error {
	return r.db.Delete(&models.Genre{}, id).Error
}

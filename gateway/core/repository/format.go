package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// FormatRepository struct for format repository.
type FormatRepository struct {
	db *gorm.DB
}

// NewFormatRepository creates a new format repository.
func NewFormatRepository(db *gorm.DB) *FormatRepository {
	return &FormatRepository{db: db}
}

// FindAllByBranchID returns all formats for a given branch ID, ordered alphabetically by name.
func (r *FormatRepository) FindAllByBranchID(branchID uint) ([]models.Format, error) {
	var formats []models.Format
	result := r.db.Where("branch_id = ?", branchID).Order("name asc").Find(&formats)
	return formats, result.Error
}

// FindOne returns one format by id and branchID.
func (r *FormatRepository) FindOne(id uint) (models.Format, error) {
	var format models.Format
	result := r.db.Where("id = ?", id).First(&format)
	return format, result.Error
}

// Create creates a new format.
func (r *FormatRepository) Create(format *models.Format) error {
	result := r.db.Create(format)
	return result.Error
}

// Update updates a format.
func (r *FormatRepository) Update(format *models.Format) error {
	result := r.db.Save(format)
	return result.Error
}

// Delete deletes a format.
func (r *FormatRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Format{}, id)
	return result.Error
}

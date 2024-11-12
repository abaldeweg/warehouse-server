package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// TagRepository represents a tag repository.
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository.
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db}
}

// FindAllByBranchID finds all tags by branch ID.
func (r *TagRepository) FindAllByBranchID(branchID uint) ([]models.Tag, error) {
	var tags []models.Tag
	result := r.db.Where("branch_id = ?", branchID).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

// FindOne finds a tag by ID.
func (r *TagRepository) FindOne(id uint) (*models.Tag, error) {
	var tag models.Tag
	result := r.db.First(&tag, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

// Create creates a new tag.
func (r *TagRepository) Create(tag *models.Tag) error {
	result := r.db.Create(tag)
	return result.Error
}

// Update updates a tag.
func (r *TagRepository) Update(tag *models.Tag) error {
	result := r.db.Save(tag)
	return result.Error
}

// Delete deletes a tag by ID.
func (r *TagRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Tag{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

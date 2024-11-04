package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// BranchRepository struct for branch repository.
type BranchRepository struct {
	db *gorm.DB
}

// NewBranchRepository creates a new branch repository.
func NewBranchRepository(db *gorm.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

// FindAll returns all branches.
func (r *BranchRepository) FindAll() ([]models.Branch, error) {
	var branches []models.Branch
	result := r.db.Find(&branches)
	return branches, result.Error
}

// FindOne returns one branch by id.
func (r *BranchRepository) FindOne(id uint) (models.Branch, error) {
	var branch models.Branch
	result := r.db.First(&branch, id)
	return branch, result.Error
}

// Update a branch.
func (r *BranchRepository) Update(branch *models.Branch) error {
	result := r.db.Save(branch)
	return result.Error
}

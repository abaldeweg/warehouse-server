package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// ConditionRepository defines the interface for interacting with Condition data.
type ConditionRepository interface {
	FindAllByBranchID(branchID uint) ([]models.Condition, error)
	FindOneByID(id uint) (models.Condition, error)
	Create(condition *models.Condition) error
	Update(condition *models.Condition) error
	Delete(id uint) error
}

// NewConditionRepository instantiates a new ConditionRepository using GORM.
func NewConditionRepository(db *gorm.DB) ConditionRepository {
	return &conditionRepository{db: db}
}

// conditionRepository implements the ConditionRepository interface.
type conditionRepository struct {
	db *gorm.DB
}

// FindAllByBranchID retrieves all conditions associated with a specific branch ID.
func (r *conditionRepository) FindAllByBranchID(branchID uint) ([]models.Condition, error) {
	var conditions []models.Condition
	if err := r.db.Where("branch_id = ?", branchID).Find(&conditions).Error; err != nil {
		return nil, err
	}
	return conditions, nil
}

// FindOneByID retrieves a single condition by its ID.
func (r *conditionRepository) FindOneByID(id uint) (models.Condition, error) {
	var condition models.Condition
	if err := r.db.First(&condition, id).Error; err != nil {
		return models.Condition{}, err
	}

	return condition, nil
}

// Create inserts a new condition into the database.
func (r *conditionRepository) Create(condition *models.Condition) error {
	return r.db.Create(condition).Error
}

// Update modifies an existing condition in the database.
func (r *conditionRepository) Update(condition *models.Condition) error {
	return r.db.Save(condition).Error
}

// Delete removes a condition from the database by its ID.
func (r *conditionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Condition{}, id).Error
}

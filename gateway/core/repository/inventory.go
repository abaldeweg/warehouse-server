package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"gorm.io/gorm"
)

// InventoryRepository handles the CRUD operations for Inventory.
type InventoryRepository struct {
	DB *gorm.DB
}

// NewInventoryRepository creates a new InventoryRepository.
func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

// FindAll retrieves all inventory items by branch ID from the database.
func (repo *InventoryRepository) FindAllByBranch(branchID uint) ([]models.Inventory, error) {
	var inventories []models.Inventory
	if err := repo.DB.Preload("Branch").Where("branch_id = ?", branchID).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

// FindByID retrieves a single inventory item by ID from the database.
func (repo *InventoryRepository) FindByID(id uint) (*models.Inventory, error) {
	var inventory models.Inventory
	if err := repo.DB.Preload("Branch").First(&inventory, id).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

// Create adds a new inventory item to the database.
func (repo *InventoryRepository) Create(inventory *models.Inventory) error {
	return repo.DB.Create(inventory).Error
}

// Update modifies an existing inventory item in the database.
func (repo *InventoryRepository) Update(inventory *models.Inventory) error {
	return repo.DB.Save(inventory).Error
}

// Delete removes an inventory item by ID from the database.
func (repo *InventoryRepository) Delete(id uint) error {
	return repo.DB.Delete(&models.Inventory{}, id).Error
}

// FindActive retrieves the count of active inventory items by branch value from the database.
func (repo *InventoryRepository) FindActive(branch uint) (int64, error) {
	var count int64
	if err := repo.DB.Model(&models.Inventory{}).Where("branch_id = ? AND ended_at IS NULL", branch).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

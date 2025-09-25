package controllers

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InventoryController handles inventory-related requests.
type InventoryController struct {
	DB   *gorm.DB
	Repo *repository.InventoryRepository
}

// NewInventoryController creates a new InventoryController.
func NewInventoryController(db *gorm.DB) *InventoryController {
	return &InventoryController{DB: db, Repo: repository.NewInventoryRepository(db)}
}

// List lists all inventory items.
func (ctrl *InventoryController) List(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	inventories, err := ctrl.Repo.FindAllByBranch(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	c.JSON(http.StatusOK, inventories)
}

// Show shows a single inventory item by ID.
func (ctrl *InventoryController) Show(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID"})
		return
	}

	inventory, err := ctrl.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	if inventory.BranchID != uint(user.(auth.User).Branch.Id) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// Create creates a new inventory item.
func (ctrl *InventoryController) Create(c *gin.Context) {
	inventory := models.NewInventory()

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	inventory.BranchID = uint(user.(auth.User).Branch.Id)

	activeCount, err := ctrl.Repo.FindActive(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	if activeCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"msg": "Active inventory already exists"})
		return
	}

	if err := ctrl.Repo.Create(inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}

	createdInventory, err := ctrl.Repo.FindByID(inventory.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created inventory"})
		return
	}

	c.JSON(http.StatusCreated, createdInventory)
}

// Update updates an existing inventory item by ID.
func (ctrl *InventoryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID"})
		return
	}

	var jsonBody struct {
		EndedAt float64 `json:"endedAt"`
	}
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON"})
		return
	}

	inventory := models.NewInventory()

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	existingInventory, err := ctrl.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	if existingInventory.BranchID != uint(user.(auth.User).Branch.Id) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	inventory.ID = uint(id)
	inventory.BranchID = uint(user.(auth.User).Branch.Id)

	endedAt := int64(jsonBody.EndedAt)
	inventory.EndedAtTimestamp = &endedAt

	if err := ctrl.Repo.Update(inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	c.JSON(http.StatusOK, inventory)
}

// Delete deletes an inventory item by ID.
func (ctrl *InventoryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID"})
		return
	}

	inventory, err := ctrl.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	if inventory.BranchID != uint(user.(auth.User).Branch.Id) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	if err := ctrl.Repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

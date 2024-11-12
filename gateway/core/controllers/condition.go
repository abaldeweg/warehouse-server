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

// ConditionController handles requests related to book conditions.
type ConditionController struct {
	DB            *gorm.DB
	ConditionRepo repository.ConditionRepository
}

// NewConditionController instantiates a new ConditionController.
func NewConditionController(db *gorm.DB) *ConditionController {
	return &ConditionController{
		DB:            db,
		ConditionRepo: repository.NewConditionRepository(db),
	}
}

// FindAll retrieves all conditions for the authenticated user's branch.
func (cc *ConditionController) FindAll(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	conditions, err := cc.ConditionRepo.FindAllByBranchID(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conditions"})
		return
	}

	c.JSON(http.StatusOK, conditions)
}

// FindOne retrieves a single condition by ID.
func (cc *ConditionController) FindOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	condition, err := cc.ConditionRepo.FindOneByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Condition not found"})
		return
	}

	c.JSON(http.StatusOK, condition)
}

// Create creates a new condition.
func (cc *ConditionController) Create(c *gin.Context) {
	var condition models.Condition
	if err := c.ShouldBindJSON(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	condition.BranchID = uint(user.(auth.User).Branch.Id)

	if !condition.Validate(cc.DB) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := cc.ConditionRepo.Create(&condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create condition"})
		return
	}

	c.JSON(http.StatusCreated, condition)
}

// Update updates an existing condition.
func (cc *ConditionController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var condition models.Condition
	if err := c.ShouldBindJSON(&condition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	condition.ID = uint(id)

	// Retrieve user from context
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	// Check if the user's branch ID matches the condition's branch ID
	if user.(auth.User).Branch.Id == int(condition.BranchID) {
		if !condition.Validate(cc.DB) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
			return
		}

		if err := cc.ConditionRepo.Update(&condition); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update condition"})
			return
		}

		c.JSON(http.StatusOK, condition)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this condition"})
	}
}

// Delete deletes a condition by ID.
func (cc *ConditionController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Retrieve user from context
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	// Fetch the condition to check its branch ID
	condition, err := cc.ConditionRepo.FindOneByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Condition not found"})
		return
	}

	// Check if the user's branch ID matches the condition's branch ID
	if user.(auth.User).Branch.Id == int(condition.BranchID) {
		if err := cc.ConditionRepo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete condition"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this condition"})
	}
}

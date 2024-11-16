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

// FormatController handles requests related to book formats.
type FormatController struct {
	formatRepo *repository.FormatRepository
	db         *gorm.DB
}

// NewFormatController instantiates a new FormatController.
func NewFormatController(db *gorm.DB) *FormatController {
	return &FormatController{
		formatRepo: repository.NewFormatRepository(db),
		db:         db,
	}
}

// FindAll retrieves all formats for the authenticated user's branch.
func (fc *FormatController) FindAll(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	formats, err := fc.formatRepo.FindAllByBranchID(uint(user.(auth.User).Branch.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve formats"})
		return
	}

	c.JSON(http.StatusOK, formats)
}

// FindOne retrieves a specific format by ID for the authenticated user's branch.
func (fc *FormatController) FindOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	format, err := fc.formatRepo.FindOne(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Format not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve format"})
		return
	}

	c.JSON(http.StatusOK, format)
}

// Create creates a new format for the authenticated user's branch.
func (fc *FormatController) Create(c *gin.Context) {
	var format models.Format
	if err := c.ShouldBindJSON(&format); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	format.BranchID = uint(user.(auth.User).Branch.Id)

	if !format.Validate(fc.db) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := fc.formatRepo.Create(&format); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create format"})
		return
	}

	c.JSON(http.StatusCreated, format)
}

// Update updates an existing format for the authenticated user's branch.
func (fc *FormatController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingFormat, err := fc.formatRepo.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Format not found"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != existingFormat.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	var format models.Format
	if err := c.ShouldBindJSON(&format); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format has invalid data"})
		return
	}

	format.ID = existingFormat.ID
	format.BranchID = existingFormat.BranchID

	if !format.Validate(fc.db) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := fc.formatRepo.Update(&format); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update format"})
		return
	}

	c.JSON(http.StatusOK, format)
}

// Delete deletes a format by ID for the authenticated user's branch.
func (fc *FormatController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	format, err := fc.formatRepo.FindOne(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Format not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve format"})
		return
	}

	if uint(user.(auth.User).Branch.Id) != format.BranchID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
		return
	}

	if err := fc.formatRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete format"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

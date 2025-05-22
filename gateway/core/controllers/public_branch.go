package controllers

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// PublicBranchController handles HTTP requests for branches.
type PublicBranchController struct {
	repo *repository.BranchRepository
}

// NewPublicBranchController creates a new branch controller.
func NewPublicBranchController(db *gorm.DB) *PublicBranchController {
	return &PublicBranchController{
		repo: repository.NewBranchRepository(db),
	}
}

// GetBranches retrieves all branches.
func (ctrl *PublicBranchController) GetBranches(c *gin.Context) {
	branches, err := ctrl.repo.FindAllByPublic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches"})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// GetBranch retrieves a branch by its ID.
func (ctrl *PublicBranchController) GetBranch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	branch, err := ctrl.repo.FindOne(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	if !branch.Public {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access to this branch is restricted"})
		return
	}

	c.JSON(http.StatusOK, branch)
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// BranchController struct for branch controller.
type BranchController struct {
	repo *repository.BranchRepository
	v    *validator.Validate
}

// NewBranchController creates a new branch controller.
func NewBranchController(db *gorm.DB) *BranchController {
	return &BranchController{
		repo: repository.NewBranchRepository(db),
		v:    validator.New(),
	}
}

// List returns all branches.
func (c *BranchController) List(ctx *gin.Context) {
	branches, err := c.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No branches found"})
		return
	}

	ctx.JSON(http.StatusOK, branches)
}

// Show returns one branch by id.
func (c *BranchController) Show(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	branch, err := c.repo.FindOne(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	ctx.JSON(http.StatusOK, branch)
}

// Update updates a branch.
func (c *BranchController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var branch models.Branch
	if err := ctx.ShouldBindJSON(&branch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	branch.ID = uint(id)

	if err := c.v.Struct(branch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Branch not valid"})
		return
	}

	if err := c.repo.Update(&branch); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Branch not updated"})
		return
	}

	ctx.JSON(http.StatusOK, branch)
}

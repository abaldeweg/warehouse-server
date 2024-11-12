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

// TagController represents a tag controller.
type TagController struct {
	DB      *gorm.DB
	TagRepo repository.TagRepository
}

// NewTagController creates a new tag controller.
func NewTagController(db *gorm.DB) *TagController {
	return &TagController{
		DB:      db,
		TagRepo: *repository.NewTagRepository(db),
	}
}

// FindAll finds all tags.
func (c *TagController) FindAll(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	tags, err := c.TagRepo.FindAllByBranchID(uint(user.(auth.User).Branch.Id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// FindOne finds a tag by ID.
func (c *TagController) FindOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := c.TagRepo.FindOne(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tag"})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

// Create creates a new tag.
func (c *TagController) Create(ctx *gin.Context) {
	var tag models.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	tag.BranchID = uint(user.(auth.User).Branch.Id)

	if !tag.Validate(c.DB) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	if err := c.TagRepo.Create(&tag); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	ctx.JSON(http.StatusCreated, tag)
}

// Update updates a tag.
func (c *TagController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var tag models.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tag.ID = uint(id)

	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	tag.BranchID = uint(user.(auth.User).Branch.Id)

	if user.(auth.User).Branch.Id == int(tag.BranchID) {
		if !tag.Validate(c.DB) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
			return
		}

		if err := c.TagRepo.Update(&tag); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
			return
		}

		ctx.JSON(http.StatusOK, tag)
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this tag"})
	}
}

// Delete deletes a tag.
func (c *TagController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	tag, err := c.TagRepo.FindOne(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tag"})
		return
	}

	if user.(auth.User).Branch.Id == int(tag.BranchID) {
		if err := c.TagRepo.Delete(uint(id)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this tag"})
	}
}

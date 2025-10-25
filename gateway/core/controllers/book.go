package controllers

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/gateway/auth"
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/core/repository"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookController struct defines the database connection.
type BookController struct {
	DB   *gorm.DB
	Repo *repository.BookRepository
}

// NewBookController creates a new instance of BookController.
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{
		DB:   db,
		Repo: repository.NewBookRepository(db),
	}
}

// CleanBooks removes books that are marked as removed or sold.
func (pbc *BookController) CleanBooks(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	branchId := uint(user.(auth.User).Branch.Id)
	if err := pbc.Repo.DeleteBooksByBranch(branchId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to clean books"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Cleaned up successfully!"})
}

// FindInventory marks books as found in inventory.
func (pbc *BookController) FindInventory(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	branchId := uint(user.(auth.User).Branch.Id)

	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}
	book, err := pbc.Repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		return
	}

	invRepo := repository.NewInventoryRepository(pbc.DB)
	inventory, err := invRepo.FindActiveInventory(branchId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	if inventory == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Active inventory not found"})
		return
	}
	if user.(auth.User).Branch.Id != int(inventory.Branch.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	if book.Inventory != nil && *book.Inventory {
		inventory.Found = inventory.Found - 1
	} else {
		inventory.Found = inventory.Found + 1
	}

	if book.Inventory != nil && *book.Inventory {
		book.Inventory = nil
	} else {
		val := true
		book.Inventory = &val
	}

	if err := invRepo.Update(inventory); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update inventory"})
		return
	}
	if err := pbc.Repo.Update(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update book"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// NotFoundInventory marks books as not found in inventory.
func (pbc *BookController) NotFoundInventory(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	branchId := uint(user.(auth.User).Branch.Id)

	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}
	book, err := pbc.Repo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		return
	}

	invRepo := repository.NewInventoryRepository(pbc.DB)
	inventory, err := invRepo.FindActiveInventory(branchId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	if inventory == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Active inventory not found"})
		return
	}
	if user.(auth.User).Branch.Id != int(inventory.Branch.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	if book.Inventory != nil && *book.Inventory == false {
		inventory.NotFound = inventory.NotFound - 1
	} else {
		inventory.NotFound = inventory.NotFound + 1
	}

	if book.Inventory != nil && *book.Inventory == false {
		book.Inventory = nil
	} else {
		val := false
		book.Inventory = &val
	}

	if err := invRepo.Update(inventory); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update inventory"})
		return
	}
	if err := pbc.Repo.Update(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update book"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// ShowStats retrieves book statistics.
func (pbc *BookController) ShowStats(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	branchId := uint(user.(auth.User).Branch.Id)

	var (
		all       int64
		available int64
		reserved  int64
		sold      int64
		removed   int64
	)

	// total books for branch
	if err := pbc.DB.Model(&models.Book{}).Where("branch_id = ?", branchId).Count(&all).Error; err != nil {
		all = 0
	}

	// available: not sold, not removed, not reserved
	if err := pbc.DB.Model(&models.Book{}).
		Where("branch_id = ? AND sold = ? AND removed = ? AND reserved = ?", branchId, false, false, false).
		Count(&available).Error; err != nil {
		available = 0
	}

	// reserved
	if err := pbc.DB.Model(&models.Book{}).
		Where("branch_id = ? AND reserved = ?", branchId, true).
		Count(&reserved).Error; err != nil {
		reserved = 0
	}

	// sold
	if err := pbc.DB.Model(&models.Book{}).
		Where("branch_id = ? AND sold = ?", branchId, true).
		Count(&sold).Error; err != nil {
		sold = 0
	}

	// removed
	if err := pbc.DB.Model(&models.Book{}).
		Where("branch_id = ? AND removed = ?", branchId, true).
		Count(&removed).Error; err != nil {
		removed = 0
	}

	// storage size
	var size int64 = 0
	if s, err := cover.GetSize(); err == nil {
		size = s
	}

	ctx.JSON(http.StatusOK, gin.H{
		"all":       all,
		"available": available,
		"reserved":  reserved,
		"sold":      sold,
		"removed":   removed,
		"storage":   float64(size) / 1000000.0,
	})
}

// ShowCover retrieves the cover images.
func (pbc *BookController) ShowCover(ctx *gin.Context) {
	id := ctx.Param("id")
	bookID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid book id"})
		return
	}

	covers := map[string]string{
		"cover_s": cover.ShowCover("s", bookID),
		"cover_m": cover.ShowCover("m", bookID),
		"cover_l": cover.ShowCover("l", bookID),
	}

	ctx.JSON(http.StatusOK, covers)
}

// DeleteCover deletes the cover images for a book.
func (pbc *BookController) DeleteCover(ctx *gin.Context) {
	id := ctx.Param("id")
	bookID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	cover.DeleteCover(bookID)
	ctx.JSON(http.StatusOK, gin.H{"message": "cover images deleted"})
}

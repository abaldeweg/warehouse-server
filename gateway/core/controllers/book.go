package controllers

import (
	"net/http"
	"strings"
	"time"

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

// SellBook marks books as sold.
func (pbc *BookController) SellBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
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

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	book.Sold = !book.Sold
	if book.Sold {
		t := time.Now()
		book.SoldOn = &t
	} else {
		book.SoldOn = nil
	}

	book.Reserved = false
	book.ReservedAt = nil
	book.Reservation = nil
	book.ReservationID = nil

	if err := pbc.Repo.Update(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update book"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// RemoveBook marks books as removed.
func (pbc *BookController) RemoveBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
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

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	book.Removed = !book.Removed
	if book.Removed {
		t := time.Now()
		book.RemovedOn = &t
	} else {
		book.RemovedOn = nil
	}

	book.Reserved = false
	book.ReservedAt = nil

	if err := pbc.Repo.Update(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update book"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// ReserveBook marks books as reserved.
func (pbc *BookController) ReserveBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
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

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	if book.Reserved {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Book is already reserved"})
		return
	}

	book.Reserved = true
	t := time.Now()
	book.ReservedAt = &t

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

// ShowBook retrieves a book by its ID.
func (pbc *BookController) ShowBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	book, err := pbc.Repo.FindByIDAndPreload(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		return
	}

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book.
func (pbc *BookController) DeleteBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
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

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	if err := pbc.Repo.Delete(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to delete book"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Book deleted successfully"})
}

// UpdateBook updates a book.
func (pbc *BookController) UpdateBook(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
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

	if book.BranchID == nil || user.(auth.User).Branch.Id != int(*book.BranchID) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Invalid Branch"})
		return
	}

	// Bind into BookUpdate to allow partial updates, then map to existing book
	var bu models.BookUpdate
	if err := ctx.ShouldBindJSON(&bu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Please enter a valid book! \n " + err.Error()})
		return
	}

	// Apply non-nil fields from BookUpdate to the existing book
	if bu.Added != nil {
		book.AddedUnix = *bu.Added
	}
	if bu.Title != nil {
		book.Title = *bu.Title
	}
	if bu.ShortDescription != nil {
		book.ShortDescription = bu.ShortDescription
	}
	if bu.Author != nil {
		fname := ""
		sname := ""
		if *bu.Author != "" {
			authorStr := *bu.Author
			if strings.Contains(authorStr, ",") {
				parts := strings.SplitN(authorStr, ",", 2)
				sname = strings.TrimSpace(parts[0])
				fname = strings.TrimSpace(parts[1])
			} else {
				parts := strings.Fields(authorStr)
				if len(parts) == 1 {
					fname = parts[0]
				} else if len(parts) >= 2 {
					fname = parts[0]
					sname = strings.Join(parts[1:], " ")
				}
			}
		}
		if book.Author == nil {
			book.Author = &models.Author{Firstname: fname, Surname: sname}
		} else {
			if fname != "" {
				book.Author.Firstname = fname
			}
			if sname != "" {
				book.Author.Surname = sname
			}
		}
	}
	if bu.GenreID != nil {
		if bu.GenreID.Val == nil {
			book.GenreID = nil
		} else {
			v := *bu.GenreID.Val
			book.GenreID = &v
		}
	}
	if bu.Price != nil {
		book.Price = *bu.Price
	}
	if bu.Sold != nil {
		book.Sold = *bu.Sold
	}
	if bu.Removed != nil {
		book.Removed = *bu.Removed
	}
	if bu.Reserved != nil {
		book.Reserved = *bu.Reserved
	}
	if bu.ReleaseYear != nil {
		book.ReleaseYear = *bu.ReleaseYear
	}
	if bu.CondID != nil {
		if bu.CondID.Val == nil {
			book.ConditionID = nil
		} else {
			v := *bu.CondID.Val
			book.ConditionID = &v
		}
	}
	if bu.Tags != nil {
		var tags []*models.Tag
		for _, t := range bu.Tags {
			if t == nil {
				continue
			}
			id := uint(*t)
			var tag models.Tag
			if err := pbc.DB.First(&tag, "id = ?", id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Tag not found", "tag_id": id})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch tag"})
				return
			}
			tags = append(tags, &tag)
		}
		book.Tags = tags
	}
	if bu.Recommendation != nil {
		book.Recommendation = *bu.Recommendation
	}
	if bu.FormatID != nil {
		if bu.FormatID.Val == nil {
			book.FormatID = 0
		} else {
			book.FormatID = *bu.FormatID.Val
		}
	}
	if bu.Subtitle != nil {
		book.Subtitle = bu.Subtitle
	}
	if bu.Duplicate != nil {
		book.Duplicate = *bu.Duplicate
	}

	existing, err := pbc.Repo.FindDuplicate(book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Error"})
		return
	}
	if existing != nil && existing.ID != book.ID {
		ctx.JSON(http.StatusConflict, gin.H{"msg": "Book not saved, because it exists already!"})
		return
	}

	// sold
	if book.Sold && book.SoldOn == nil {
		t := time.Now()
		book.SoldOn = &t
	}

	// revert sold
	if !book.Sold && book.SoldOn != nil {
		book.SoldOn = nil
	}

	// removed
	if book.Removed && book.RemovedOn == nil {
		t := time.Now()
		book.RemovedOn = &t
	}

	// revert removed
	if !book.Removed && book.RemovedOn != nil {
		book.RemovedOn = nil
	}

	// reserved
	if book.Reserved && book.ReservedAt == nil {
		t := time.Now()
		book.ReservedAt = &t
	}

	// revert reserved
	if !book.Reserved && book.ReservedAt != nil {
		book.ReservedAt = nil
		book.Reservation = nil
		book.ReservationID = nil
	}

	if err := pbc.Repo.Update(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update book"})
		return
	}

	updatedBook, err := pbc.Repo.FindByIDAndPreload(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "Book not found"})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

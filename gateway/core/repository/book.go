package repository

import (
	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/cover"
	"gorm.io/gorm"
)

// BookRepository defines the methods for interacting with the Book model.
type BookRepository struct {
	DB *gorm.DB
}

// NewBookRepository creates a new instance of BookRepository.
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

// DeleteBooksByBranch finds all books for the given branch that are marked
// as sold or removed. Then, deletes their cover files and removes them from the DB.
func (r *BookRepository) DeleteBooksByBranch(branchID uint) error {
	var books []models.Book

	if err := r.DB.Where("branch_id = ? AND (sold = ? OR removed = ?)", branchID, true, true).Find(&books).Error; err != nil {
		return err
	}

	tx := r.DB.Begin()
	for _, b := range books {
		cover.DeleteCover(b.ID)

		if err := tx.Delete(&b).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// FindByID retrieves a book by UUID.
func (r *BookRepository) FindByID(id interface{}) (*models.Book, error) {
	var book models.Book
	if err := r.DB.First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// Update saves the provided book.
func (r *BookRepository) Update(book *models.Book) error {
	return r.DB.Save(book).Error
}

// FindByID retrieves a book by UUID.
func (r *BookRepository) FindByIDAndPreload(id interface{}) (*models.Book, error) {
	var book models.Book
	if err := r.DB.Preload("Branch").Preload("Author").Preload("Genre").Preload("Condition").Preload("Format").Preload("Reservation").Preload("Tags").First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

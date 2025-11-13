package repository

import (
	"time"

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

// KEEP_REMOVED_DAYS is the default number of days to keep removed/sold books
// before permanently deleting them.
const KEEP_REMOVED_DAYS = 28

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
	tx := r.DB.Begin()

	if err := tx.Omit("Tags").Save(book).Error; err != nil {
		tx.Rollback()
		return err
	}

	if book.Tags != nil {
		if err := tx.Exec("DELETE FROM book_tag WHERE book_id = ?", book.ID).Error; err != nil {
			tx.Rollback()
			return err
		}

		if len(book.Tags) > 0 {
			vals := make([]map[string]any, 0, len(book.Tags))
			for _, t := range book.Tags {
				vals = append(vals, map[string]any{"book_id": book.ID, "tag_id": t.ID})
			}

			if err := tx.Table("book_tag").Create(vals).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// FindDuplicate searches for an existing book that would be considered a duplicate
func (r *BookRepository) FindDuplicate(b *models.Book) (*models.Book, error) {
	var existing models.Book
	query := r.DB.Where("title = ?", b.Title).
		Where("sold = ?", b.Sold).
		Where("removed = ?", b.Removed).
		Where("release_year = ?", b.ReleaseYear).
		Where("format_id = ?", b.FormatID)

	if b.BranchID != nil {
		query = query.Where("branch_id = ?", *b.BranchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}

	if b.AuthorID != nil {
		query = query.Where("author_id = ?", *b.AuthorID)
	} else {
		query = query.Where("author_id IS NULL")
	}

	if b.GenreID != nil {
		query = query.Where("genre_id = ?", *b.GenreID)
	} else {
		query = query.Where("genre_id IS NULL")
	}

	if err := query.First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &existing, nil
}

// FindByID retrieves a book by UUID.
func (r *BookRepository) FindByIDAndPreload(id interface{}) (*models.Book, error) {
	var book models.Book
	if err := r.DB.Preload("Branch").Preload("Author").Preload("Genre").Preload("Condition").Preload("Format").Preload("Reservation").Preload("Tags").First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// Delete removes the given book from the database and deletes its cover files.
func (r *BookRepository) Delete(book *models.Book) error {
	tx := r.DB.Begin()

	cover.DeleteCover(book.ID)

	if err := tx.Delete(book).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeleteBooks deletes books whose `sold_on` or `removed_on` timestamp is
// older than `clearLimit` days. If `clearLimit` <= 0 the default
// `KEEP_REMOVED_DAYS` is used.
func (r *BookRepository) DeleteBooks(clearLimit int) error {
	if clearLimit <= 0 {
		clearLimit = KEEP_REMOVED_DAYS
	}

	cutoff := time.Now().AddDate(0, 0, -clearLimit)

	var books []models.Book
	if err := r.DB.Where("sold_on <= ? OR removed_on <= ?", cutoff, cutoff).Find(&books).Error; err != nil {
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

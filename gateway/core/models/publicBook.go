package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicBook represents a public book entity.
type PublicBook struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Currency         string    `json:"currency"`
	Title            string    `json:"title" binding:"required"`
	AuthorFirstname  string    `json:"authorFirstname"`
	AuthorSurname    string    `json:"authorSurname"`
	BranchID         *uint     `json:"branch_id"`
	ShortDescription *string   `json:"shortDescription"`
	Genre            *string   `json:"genre"`
	BranchName       *string   `json:"branchName"`
	BranchOrdering   *string   `json:"branchOrdering"`
	Price            float64   `json:"price"`
	ReleaseYear      int       `json:"releaseYear"`
	Cond             *string   `json:"cond"`
	FormatName       *string   `json:"format_name"`
	Subtitle         *string   `json:"subtitle"`
}

// BeforeCreate is a GORM hook that generates a UUID before creating a new record.
func (book *PublicBook) BeforeCreate(tx *gorm.DB) (err error) {
	book.ID = uuid.New()
	return
}

// TableName overrides the default table name for PublicBook model.
func (PublicBook) TableName() string {
	return "book"
}

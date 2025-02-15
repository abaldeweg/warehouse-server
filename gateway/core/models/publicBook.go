package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicBook represents a public book entity.
type PublicBook struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;->"`
	Currency         string     `json:"currency" gorm:"-"`
	Title            string     `json:"title" binding:"required"`
	Subtitle         string     `json:"subtitle" validate:"max=255"`
	AuthorID         uint       `json:"author_id" gorm:"index"`
	Author           Author     `json:"-" gorm:"foreignKey:AuthorID"`
	AuthorFirstname  string     `json:"authorFirstname" gorm:"-"`
	AuthorSurname    string     `json:"authorSurname" gorm:"-"`
	BranchID         *int       `json:"branch_id"`
	Branch           Branch     `json:"branch" gorm:"foreignKey:BranchID"`
	ShortDescription *string    `json:"shortDescription"`
	GenreID          uint       `json:"genre_id" gorm:"index"`
	Genre            Genre      `json:"-" gorm:"foreignKey:GenreID"`
	GenreName        string     `json:"genre" gorm:"-"`
	BranchName       string     `json:"branchName" gorm:"-"`
	BranchOrdering   string     `json:"branchOrdering" gorm:"-"`
	Price            float32    `json:"price" gorm:"default:0.00"`
	ReleaseYear      int        `json:"releaseYear"`
	Condition        *Condition `json:"-" gorm:"foreignKey:ConditionID"`
	ConditionID      uint       `json:"condition_id"`
	Cond             string     `json:"cond" binding:"required" gorm:"-"`
	Format           Format     `json:"-" gorm:"foreignKey:FormatID"`
	FormatID         uint       `json:"format_id" gorm:""`
	FormatName       string     `json:"format_name" gorm:"-"`
	BranchCart       bool       `json:"branchCart" gorm:"-"`
}

// TableName overrides the default table name for PublicBook model.
func (PublicBook) TableName() string {
	return "book"
}

// BeforeCreate is a GORM hook that generates a UUID before creating a new record.
func (book *PublicBook) BeforeCreate(tx *gorm.DB) (err error) {
	book.ID = uuid.New()
	return
}

// AfterFind is a GORM hook that populates the PublicBook struct with related entities.
func (book *PublicBook) AfterFind(tx *gorm.DB) (err error) {
	book.Currency = book.Branch.Currency
	book.AuthorFirstname = book.Author.Firstname
	book.AuthorSurname = book.Author.Surname
	book.GenreName = book.Genre.Name
	book.BranchName = book.Branch.Name
	book.BranchOrdering = book.Branch.Ordering
	book.Cond = book.Condition.Name
	book.FormatName = book.Format.Name
	book.BranchCart = book.Branch.Cart
	return
}

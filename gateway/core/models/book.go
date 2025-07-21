package models

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a book entity.
type Book struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	BranchID         uint         `json:"branch_id"`
	Branch           *Branch      `json:"branch" gorm:"foreignKey:BranchID"`
	Added            time.Time    `json:"added"`
	Title            string       `json:"title" gorm:"type:varchar(255);not null" validate:"required"`
	ShortDescription string       `json:"short_description"`
	AuthorID         uint         `json:"author_id" gorm:"index"`
	Author           Author       `json:"-" gorm:"foreignKey:AuthorID"`
	GenreID          uint         `json:"genre_id" gorm:"index"`
	Genre            *Genre       `json:"genre" gorm:"foreignKey:GenreID"`
	Price            float32      `json:"price" gorm:"default:0.00"`
	Sold             bool         `json:"sold" gorm:"default:false"`
	SoldOn           time.Time    `json:"sold_on,omitempty"`
	Removed          bool         `json:"removed" gorm:"default:false"`
	RemovedOn        time.Time    `json:"removed_on,omitempty"`
	Reserved         bool         `json:"reserved" gorm:"default:false"`
	ReservedAt       time.Time    `json:"reserved_at,omitempty"`
	ReleaseYear      int          `json:"published_date" validate:"release_year"`
	Condition        *Condition   `json:"condition" gorm:"foreignKey:ConditionID"`
	ConditionID      uint         `json:"condition_id"`
	Tags             []*Tag       `json:"tags" gorm:"many2many:book_tag;"`
	Recommendation   bool         `json:"recommendations" gorm:"foreignKey:BookID"`
	Inventory        bool         `json:"inventory" gorm:"default:false"`
	Format           *Format      `json:"format" gorm:"foreignKey:FormatID"`
	FormatID         uint         `json:"format_id" gorm:"not null"`
	Subtitle         string       `json:"subtitle" validate:"max=255"`
	Duplicate        bool         `json:"duplicate" gorm:"default:false"`
	ReservationID    uuid.UUID    `json:"reservation_id"`
	Reservation      *Reservation `json:"reservation" gorm:"foreignKey:ReservationID"`
}

// TableName overrides the default table name for Book model.
func (Book) TableName() string {
	return "book"
}

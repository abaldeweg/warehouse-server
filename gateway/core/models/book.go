package models

import (
	"time"
)

// Book represents a book entity.
type Book struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement;->"`
	BranchID         uint       `json:"branch_id"`
	Branch           *Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	Added            time.Time  `json:"added"`
	ShortDescription string     `json:"summary"`
	Author           []*Author  `json:"authors" gorm:"many2many:book_author;"`
	GenreID          uint       `json:"genre_id" gorm:"index"`
	Genre            *Genre     `json:"genre" gorm:"foreignKey:GenreID"`
	Price            float32    `json:"price" gorm:"default:0.00"`
	Sold             bool       `gorm:"default:false"`
	SoldOn           time.Time  `json:"sold_on,omitempty"`
	Removed          bool       `gorm:"default:false"`
	RemovedOn        time.Time  `json:"removed_on,omitempty"`
	Reserved         bool       `gorm:"default:false"`
	ReservedAt       time.Time  `json:"reserved_at,omitempty"`
	ReleaseYear      time.Time  `json:"published_date"`
	Condition        *Condition `json:"condition" gorm:"foreignKey:ConditionID"`
	ConditionID      uint       `json:"condition_id"`
	Tags             []*Tag     `json:"tags" gorm:"many2many:book_tag;"`
	// Reservation      []*Reservation `json:"reservations" gorm:"foreignKey:BookID"`
	Recommendation bool `json:"recommendations" gorm:"foreignKey:BookID"`
	Inventory      bool `gorm:"default:false"`
	// Format           *Format        `json:"format" gorm:"foreignKey:FormatID"`
	FormatID  uint   `json:"format_id"`
	Subtitle  string `json:"subtitle" validate:"max=255"`
	Duplicate bool   `gorm:"default:false"`
}

// TableName overrides the default table name for Book model.
func (Book) TableName() string {
	return "book"
}

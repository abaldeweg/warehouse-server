package models

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a book entity.
type Book struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	BranchID         *uint        `json:"branch_id" gorm:"default:null"`
	Branch           *Branch      `json:"branch" gorm:"foreignKey:BranchID"`
	Added            time.Time    `json:"added" gorm:"column:added;autoCreateTime"`
	Title            string       `json:"title" gorm:"type:varchar(255)" validate:"required,min=1,max=255"`
	ShortDescription *string      `json:"shortDescription" gorm:"default:null"`
	AuthorID         *uint        `json:"author_id" gorm:"default:null"`
	Author           *Author      `json:"author" gorm:"foreignKey:AuthorID"`
	GenreID          *uint        `json:"genre_id" gorm:"default:null"`
	Genre            *Genre       `json:"genre" gorm:"foreignKey:GenreID"`
	Price            float64      `json:"price" gorm:"type:decimal(10,2);default:0.00" validate:"gte=0.00"`
	Sold             bool         `json:"sold" gorm:"default:false"`
	SoldOn           *time.Time   `json:"soldOn,omitempty" gorm:"default:null"`
	Removed          bool         `json:"removed" gorm:"default:false"`
	RemovedOn        *time.Time   `json:"removedOn,omitempty" gorm:"default:null"`
	Reserved         bool         `json:"reserved" gorm:"default:false"`
	ReservedAt       *time.Time   `json:"reservedAt,omitempty" gorm:"default:null"`
	ReleaseYear      int          `json:"releaseYear" gorm:"type:int;column:release_year" validate:"gte=1000,lte=9999"`
	Condition        *Condition   `json:"condition" gorm:"foreignKey:ConditionID"`
	ConditionID      *uint        `json:"cond_id" gorm:"default:null"`
	Tags             []*Tag       `json:"tags" gorm:"many2many:book_tag;default:null"`
	Recommendation   bool         `json:"recommendations" gorm:"foreignKey:BookID;default:false"`
	Inventory        *bool        `json:"inventory" gorm:"default:null"`
	Format           *Format      `json:"format" gorm:"foreignKey:FormatID"`
	FormatID         uint         `json:"format_id" gorm:"not null"`
	Subtitle         *string      `json:"subtitle" gorm:"default:null" validate:"max=255"`
	Duplicate        bool         `json:"duplicate" gorm:"default:false"`
	ReservationID    *uuid.UUID   `json:"reservation_id" gorm:"default:null"`
	Reservation      *Reservation `json:"reservation" gorm:"foreignKey:ReservationID"`
}

// TableName overrides the default table name for Book model.
func (Book) TableName() string {
	return "book"
}

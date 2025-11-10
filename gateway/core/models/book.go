package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Book represents a book entity.
type Book struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	BranchID         *uint        `json:"branch_id" gorm:"default:null"`
	Branch           *Branch      `json:"branch" gorm:"foreignKey:BranchID"`
	Added            time.Time    `json:"-" gorm:"column:added;autoCreateTime"`
	Title            string       `json:"title" gorm:"type:varchar(255)" validate:"required,min=1,max=255"`
	ShortDescription *string      `json:"shortDescription" gorm:"default:null"`
	AuthorID         *uint        `json:"author_id" gorm:"default:null"`
	Author           *Author      `json:"author" gorm:"foreignKey:AuthorID"`
	GenreID          *uint        `json:"genre_id" gorm:"default:null"`
	Genre            *Genre       `json:"genre" gorm:"foreignKey:GenreID"`
	Price            float64      `json:"price" gorm:"type:decimal(10,2);default:0.00" validate:"gte=0.00"`
	Sold             bool         `json:"sold" gorm:"default:false"`
	SoldOn           *time.Time   `json:"-" gorm:"default:null"`
	Removed          bool         `json:"removed" gorm:"default:false"`
	RemovedOn        *time.Time   `json:"-" gorm:"default:null"`
	Reserved         bool         `json:"reserved" gorm:"default:false"`
	ReservedAt       *time.Time   `json:"-" gorm:"default:null"`
	ReleaseYear      int          `json:"releaseYear" gorm:"type:int;column:release_year" validate:"gte=1000,lte=9999"`
	Condition        *Condition   `json:"condition" gorm:"foreignKey:ConditionID"`
	ConditionID      *uint        `json:"cond_id" gorm:"column:cond_id;default:null"`
	Tags             []*Tag       `json:"tags" gorm:"many2many:book_tag;default:null"`
	Recommendation   bool         `json:"recommendation" gorm:"foreignKey:BookID;default:false"`
	Inventory        *bool        `json:"inventory" gorm:"default:null"`
	Format           *Format      `json:"format" gorm:"foreignKey:FormatID"`
	FormatID         uint         `json:"format_id" gorm:"not null"`
	Subtitle         *string      `json:"subtitle" gorm:"default:null" validate:"max=255"`
	Duplicate        bool         `json:"duplicate" gorm:"default:false"`
	ReservationID    *uuid.UUID   `json:"reservation_id" gorm:"default:null"`
	Reservation      *Reservation `json:"reservation" gorm:"foreignKey:ReservationID"`
	AddedUnix        int64        `json:"added" gorm:"-"`
	SoldOnUnix       *int64       `json:"soldOn,omitempty" gorm:"-"`
	RemovedOnUnix    *int64       `json:"removedOn,omitempty" gorm:"-"`
	ReservedAtUnix   *int64       `json:"reservedAt,omitempty" gorm:"-"`
}

// BookUpdate represents an update payload.
type BookUpdate struct {
	Added            *int64   `json:"added,omitempty"`
	Title            *string  `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	ShortDescription *string  `json:"shortDescription,omitempty"`
	Author           *string  `json:"author,omitempty"`
	GenreID          *uint    `json:"genre,omitempty"`
	Price            *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	Sold             *bool    `json:"sold,omitempty"`
	Removed          *bool    `json:"removed,omitempty"`
	Reserved         *bool    `json:"reserved,omitempty"`
	ReleaseYear      *int     `json:"releaseYear,omitempty" validate:"omitempty,gte=1000,lte=9999"`
	CondID           *uint    `json:"cond,omitempty"`
	Tags             []*int64 `json:"tags,omitempty"`
	Recommendation   *bool    `json:"recommendation,omitempty"`
	FormatID         *uint    `json:"format,omitempty"`
	Subtitle         *string  `json:"subtitle,omitempty" validate:"omitempty,max=255"`
	Duplicate        *bool    `json:"duplicate,omitempty"`
}

// TableName overrides the default table name for Book model.
func (Book) TableName() string {
	return "book"
}

// AfterFind is a GORM hook that populates AddedUnix after loading from DB.
func (b *Book) AfterFind(tx *gorm.DB) (err error) {
	b.AddedUnix = b.Added.Unix()
	if b.SoldOn != nil {
		v := b.SoldOn.Unix()
		b.SoldOnUnix = &v
	}
	if b.RemovedOn != nil {
		v := b.RemovedOn.Unix()
		b.RemovedOnUnix = &v
	}
	if b.ReservedAt != nil {
		v := b.ReservedAt.Unix()
		b.ReservedAtUnix = &v
	}
	return nil
}

// BeforeSave will copy AddedUnix into Added if the unix value is set.
func (b *Book) BeforeSave(tx *gorm.DB) (err error) {
	if b.AddedUnix != 0 {
		b.Added = time.Unix(b.AddedUnix, 0)
	}
	if b.SoldOnUnix != nil {
		t := time.Unix(*b.SoldOnUnix, 0)
		b.SoldOn = &t
	}
	if b.RemovedOnUnix != nil {
		t := time.Unix(*b.RemovedOnUnix, 0)
		b.RemovedOn = &t
	}
	if b.ReservedAtUnix != nil {
		t := time.Unix(*b.ReservedAtUnix, 0)
		b.ReservedAt = &t
	}
	return nil
}

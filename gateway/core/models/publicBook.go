package models

import (
	"gorm.io/gorm"
)

// PublicBook represents a public book entity.
type PublicBook struct {
	gorm.Model
}

// TableName overrides the default table name for PublicBook model.
func (PublicBook) TableName() string {
	return "book"
}

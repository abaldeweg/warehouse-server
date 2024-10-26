package models

import (
	"github.com/go-playground/validator/v10"
)

// Author represents an author entity with ID, Firstname, and Surname.
type Author struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Firstname string `json:"firstname" gorm:"size:255" validate:"required,min=1,max=255"`
	Surname   string `json:"surname" gorm:"size:255" validate:"required,min=1,max=255"`
}

// Validate validates the Author struct based on defined validation tags.
func (a *Author) Validate(v *validator.Validate) error {
	return v.Struct(a)
}

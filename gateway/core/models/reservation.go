package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Reservation represents a reservation.
type Reservation struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	BranchID   uint      `json:"branch_id" gorm:"index"`
	Branch     Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	CreatedAt  time.Time `json:"created_at"`
	Notes      string    `json:"notes"`
	Books      []*Book   `json:"books" gorm:"foreignKey:ReservationID"`
	Salutation string    `json:"salutation" validate:"required,oneof=m f d"`
	Firstname  string    `json:"firstname" validate:"required,max=255"`
	Surname    string    `json:"surname" validate:"required,max=255"`
	Mail       string    `json:"mail" validate:"required,email,max=255"`
	Phone      string    `json:"phone" validate:"max=255"`
	Open       bool      `json:"open" gorm:"default:true"`
}

// ReservationForm represents a form for creating or updating a reservation.
type ReservationForm struct {
	CreatedAt  time.Time `json:"created_at"`
	Notes      string    `json:"notes"`
	Books      string    `json:"books"`
	Salutation string    `json:"salutation" validate:"required,oneof=m f d"`
	Firstname  string    `json:"firstname" validate:"required,max=255"`
	Surname    string    `json:"surname" validate:"required,max=255"`
	Mail       string    `json:"mail" validate:"required,email,max=255"`
	Phone      string    `json:"phone" validate:"max=255"`
	Open       bool      `json:"open"`
}

// TableName overrides the default table name for the Reservation model.
func (Reservation) TableName() string {
	return "reservation"
}

// Validate validates the Reservation model.
func (r *Reservation) Validate(db *gorm.DB) bool {
	validate := validator.New()
	return validate.StructExcept(r, "Branch") == nil
}

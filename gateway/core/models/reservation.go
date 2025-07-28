package models

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Reservation represents a reservation.
type Reservation struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	BranchID   uint      `json:"branch_id" gorm:"index"`
	Branch     Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	CreatedAt  time.Time `json:"createdAt"`
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
	CreatedAt  int64  `json:"created_at"`
	Notes      string `json:"notes"`
	Books      string `json:"books"`
	Salutation string `json:"salutation" validate:"required,oneof=m f d"`
	Firstname  string `json:"firstname" validate:"required,max=255"`
	Surname    string `json:"surname" validate:"required,max=255"`
	Mail       string `json:"mail" validate:"required,email,max=255"`
	Phone      string `json:"phone" validate:"max=255"`
	Open       bool   `json:"open"`
}

// ReservationUpdateForm represents a form for creating or updating a reservation.
type ReservationUpdateForm struct {
	CreatedAt  int64  `json:"createdAt"`
	Notes      string `json:"notes"`
	Salutation string `json:"salutation" validate:"required,oneof=m f d"`
	Firstname  string `json:"firstname" validate:"required,max=255"`
	Surname    string `json:"surname" validate:"required,max=255"`
	Mail       string `json:"mail" validate:"required,email,max=255"`
	Phone      string `json:"phone" validate:"max=255"`
	Open       bool   `json:"open"`
}

// ToTime converts the CreatedAt unix timestamp to time.Time
func (f *ReservationForm) ToTime() time.Time {
	return time.Unix(f.CreatedAt, 0)
}

// ToTime converts the CreatedAt unix timestamp to time.Time
func (f *ReservationUpdateForm) ToTime() time.Time {
	return time.Unix(f.CreatedAt, 0)
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

// MarshalJSON customizes the JSON output for Reservation.
func (r Reservation) MarshalJSON() ([]byte, error) {
	type Alias Reservation
	return json.Marshal(&struct {
		CreatedAt int64 `json:"createdAt"`
		*Alias
	}{
		CreatedAt: r.CreatedAt.Unix(),
		Alias:     (*Alias)(&r),
	})
}

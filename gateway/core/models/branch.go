package models

import (
	"github.com/go-playground/validator/v10"
)

// Branch represents a branch entity with various attributes.
type Branch struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Name      string  `json:"name" validate:"required,max=255"`
	Steps     float32 `json:"steps" validate:"gte=0.00" gorm:"default:0.00"`
	Currency  string  `json:"currency" validate:"required,oneof=EUR USD" gorm:"default:'EUR'"`
	Ordering  string  `json:"ordering" gorm:"type:text"`
	Public    bool    `json:"public" gorm:"default:false"`
	Pricelist string  `json:"pricelist" gorm:"type:text"`
	Cart      bool    `json:"cart" gorm:"default:false"`
	Content   string  `json:"content" gorm:"type:text"`
}

// TableName returns the branch table name.
func (b *Branch) TableName() string {
	return "branch"
}

// Validate validates the Branch struct based on defined validation tags.
func (b *Branch) Validate(v *validator.Validate) error {
	return v.Struct(b)
}

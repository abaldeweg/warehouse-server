package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Format represents a book format entity.
type Format struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	BranchID uint   `json:"branch_id" gorm:"index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`
	Books    []Book `json:"-" gorm:"foreignKey:FormatID"`
}

// TableName overrides the default table name for Format model.
func (Format) TableName() string {
	return "format"
}

// Validate validates the Format model based on defined rules.
func (f *Format) Validate(db *gorm.DB) bool {
	validate := validator.New()
	if err := validate.StructExcept(f, "Branch", "Books"); err != nil {
		return false
	}

	return true
}

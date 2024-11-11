package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Condition model represents a book's condition.
type Condition struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Name     string `json:"name" gorm:"type:varchar(255);not null;unique" validate:"required,min=1,max=255"`
	BranchID uint   `json:"branch_id" gorm:"index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`
}

// TableName overrides the default table name for Condition model.
func (Condition) TableName() string {
	return "cond"
}

// Validate validates the Condition model based on defined rules.
func (c *Condition) Validate(db *gorm.DB) bool {
	validate := validator.New()
	if err := validate.StructExcept(c, "Branch"); err != nil {
		return false
	}

	return true
}

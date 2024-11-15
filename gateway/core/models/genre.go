package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Genre represents a genre entity.
type Genre struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Name     string `json:"name" gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	BranchID uint   `json:"branch_id" gorm:"index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`
}

// TableName overrides the default table name for Genre model.
func (Genre) TableName() string {
	return "genre"
}

// Validate validates the Genre model based on defined rules.
func (g *Genre) Validate(db *gorm.DB) bool {
	validate := validator.New()
	if err := validate.StructExcept(g, "Branch", "BranchID"); err != nil {
		return false
	}

	return true
}

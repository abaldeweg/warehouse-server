package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Tag represents a tag model.
type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement;->"`
	Name string `json:"name" validate:"required,max=255"`

	BranchID uint   `json:"branch_id" gorm:"index"`
	Branch   Branch `json:"branch" gorm:"foreignKey:BranchID"`

	Books []*Book `json:"books" gorm:"many2many:book_tag;"`
}

// TableName overrides the default table name.
func (Tag) TableName() string {
	return "tag"
}

// Validate validates the Tag model.
func (t *Tag) Validate(db *gorm.DB) bool {
	validate := validator.New()
	if err := validate.StructExcept(t, "Branch", "Books"); err != nil {
		return false
	}

	return true
}

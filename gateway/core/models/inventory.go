package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Inventory represents the inventory model.
type Inventory struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement;->"`
	BranchID           uint       `json:"branch_id" gorm:"index"`
	Branch             Branch     `json:"branch" gorm:"foreignKey:BranchID"`
	StartedAt          time.Time  `json:"-"`
	StartedAtTimestamp int64      `json:"startedAt" gorm:"-"`
	EndedAt            *time.Time `json:"-"`
	EndedAtTimestamp   *int64     `json:"endedAt" gorm:"-"`
	Found              int        `json:"found" validate:"min=0"`
	NotFound           int        `json:"notFound" validate:"min=0"`
}

// TableName sets the table name for the Inventory model.
func (Inventory) TableName() string {
	return "inventory"
}

// Validate validates the Inventory model.
func (i *Inventory) Validate(db *gorm.DB) bool {
	validate := validator.New()
	if err := validate.StructExcept(i, "Branch"); err != nil {
		return false
	}

	return true
}

// BeforeSave sets StartedAt and EndedAt before saving.
func (i *Inventory) BeforeSave(tx *gorm.DB) (err error) {
	i.StartedAt = time.Unix(i.StartedAtTimestamp, 0)

	if i.EndedAtTimestamp == nil {
		i.EndedAt = nil
	} else {
		endedAt := time.Unix(*i.EndedAtTimestamp, 0)
		i.EndedAt = &endedAt
	}
	return
}

// AfterFind sets StartedAtTimestamp and EndedAtTimestamp after loading.
func (i *Inventory) AfterFind(tx *gorm.DB) (err error) {
	i.StartedAtTimestamp = i.StartedAt.Unix()

	if i.EndedAt != nil {
		endedAtTimestamp := i.EndedAt.Unix()
		i.EndedAtTimestamp = &endedAtTimestamp
	} else {
		i.EndedAtTimestamp = nil
	}
	return
}

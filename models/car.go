package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	UserID      uint           `json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"` // URLs
}

package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model `json:"-"`
	UserID      uint     `json:"user_id"`
	Title       string   `gorm:"not null" json:"title"`
	Description string   `json:"description"`
	Tags        []string `gorm:"type:text[]" json:"tags"`
	Images      []string `gorm:"type:text[]" json:"images"` // URLs or paths
}

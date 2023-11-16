package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResNews struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	UserID      int    `json:"user_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
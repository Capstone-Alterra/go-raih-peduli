package news

import (
	"raihpeduli/features/auth"
	"time"

	"gorm.io/gorm"
)

type News struct {
	ID          int    `gorm:"type:int(11)"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	Photo       string `gorm:"type:varchar(255)"`
	UserID      int    `gorm:"type:int(11)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User auth.User
}

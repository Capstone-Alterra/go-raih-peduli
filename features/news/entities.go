package news

import (
	"raihpeduli/features/auth"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model

	ID          int    `gorm:"type:int(11)"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	Photo       string `gorm:"type:varchar(255)"`
	UserID      int    `gorm:"type:int(11)"`

	User auth.User
}

package _blueprint

import (
	"gorm.io/gorm"
)

type Placeholder struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
	Name string `gorm:"type:varchar(255)"`
}


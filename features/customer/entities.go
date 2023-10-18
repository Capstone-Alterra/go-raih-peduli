package customer

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
	Name string `gorm:"type:varchar(255)"`
}


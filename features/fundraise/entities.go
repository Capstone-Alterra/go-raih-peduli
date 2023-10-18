package fundraise

import (
	"gorm.io/gorm"
)

type Fundraise struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
	Name string `gorm:"type:varchar(255)"`
}


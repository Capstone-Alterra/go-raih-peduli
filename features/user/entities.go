package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
}
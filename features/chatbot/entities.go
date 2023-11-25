package chatbot

import (
	"gorm.io/gorm"
)

type Chatbot struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
	Name string `gorm:"type:varchar(255)"`
}


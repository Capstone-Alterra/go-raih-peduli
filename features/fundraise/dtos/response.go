package dtos

import (
	"gorm.io/gorm"
	"time"
)

type ResFundraise struct {
	gorm.Model
	Target     string    `json:"target" form:"target"`
	User_id    int       `json:"user_id" form:"user_id"`
	Start_date time.Time `json:"start_date" form:"start_date"`
	End_date   time.Time `json:"end_date" form:"end_date"`
}

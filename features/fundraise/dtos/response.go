package dtos

import (
	"time"
)

type ResFundraise struct {
	Title	   string    `json:"title" form:"title"`
	
	Target     string    `json:"target" form:"target"`
	User_id    int       `json:"user_id" form:"user_id"`
	Start_date time.Time `json:"start_date" form:"start_date"`
	End_date   time.Time `json:"end_date" form:"end_date"`
}

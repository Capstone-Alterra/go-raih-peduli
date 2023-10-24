package dtos

import "time"

type InputFundraise struct {
	Target     string    `json:"target" form:"target" validate:"required"`
	User_id    int       `json:"user_id" form:"user_id" validate:"required"`
	Start_date time.Time `json:"start_date" form:"start_date" validate:"required"`
	End_date   time.Time `json:"end_date" form:"end_date" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

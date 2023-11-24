package dtos

import "mime/multipart"

type InputFundraise struct {
	Title       string         `json:"title" form:"title" validate:"required"`
	Description string         `json:"description" form:"description" validate:"required"`
	Photo       multipart.File `json:"photo" form:"photo"`
	Target      int32          `json:"target" form:"target" validate:"required"`
	StartDate   string         `json:"start_date" form:"start_date" validate:"required"`
	EndDate     string         `json:"end_date" form:"end_date" validate:"required"`
}

type InputFundraiseStatus struct {
	Status string `json:"status" form:"status" validate:"oneof=pending accepted rejected"`
}

type Pagination struct {
	Page int `query:"page"`
	PageSize int `query:"page_size"`
}

type SearchAndFilter struct {
	Title string `query:"title"`
	MinTarget int32 `query:"min_target"`
	MaxTarget int32 `query:"max_target"`
}

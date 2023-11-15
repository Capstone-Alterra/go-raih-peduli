package dtos

import "mime/multipart"

type InputNews struct {
	Title       string         `json:"title" form:"title" validate:"required"`
	Description string         `json:"description" form:"description" validate:"required"`
	Photo       multipart.File `json:"photo" form:"photo"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

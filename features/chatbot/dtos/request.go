package dtos

type InputMessage struct {
	Message string `json:"message" form:"message" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
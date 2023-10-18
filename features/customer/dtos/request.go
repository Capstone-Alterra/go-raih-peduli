package dtos

type InputCustomer struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
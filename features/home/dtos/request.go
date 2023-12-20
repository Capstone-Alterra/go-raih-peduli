package dtos

type InputHome struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type Pagination struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

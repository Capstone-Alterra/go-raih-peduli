package dtos

type InputBookmarkPost struct {
	PostID   int    `form:"post_id" json:"post_id" validate:"required"`
	PostType string `form:"post_type" json:"post_type" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
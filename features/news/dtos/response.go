package dtos

type ResNews struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	UserID      int    `json:"user_id"`
}

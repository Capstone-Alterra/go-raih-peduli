package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResFundraise struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Target      int32     `json:"target"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	UserID      int       `json:"user_id"`
	BookmarkID  *string	  `json:"bookmark_id"`	

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type PaginationResponse struct {
	TotalData    int64 `json:"total_data"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage int   `json:"previous_page"`
	NextPage     int   `json:"next_page"`
	TotalPage    int   `json:"total_page"`
}

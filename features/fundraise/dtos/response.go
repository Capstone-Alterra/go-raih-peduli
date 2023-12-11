package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResFundraise struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Photo          string    `json:"photo"`
	Target         int32     `json:"target"`
	FundAcquired   int32     `json:"fund_acquired"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	RejectedReason string    `json:"rejected_reason"`
	UserID         int       `json:"user_id"`

	BookmarkID *string `json:"bookmark_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type FundraiseDetails struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Photo          string    `json:"photo"`
	Target         int32     `json:"target"`
	FundAcquired   int32     `json:"fund_acquired"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	RejectedReason string    `json:"rejected_reason"`
	UserID         int       `json:"user_id"`
	UserFullname   string	 `json:"user_fullname"`
	UserPhoto	   string 	 `json:"user_photo"`

	BookmarkID *string `json:"bookmark_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

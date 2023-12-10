package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResFundraisesHistory struct {
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
	PostType  string         `json:"post_type"`
}

type ResVolunteersVacancyHistory struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	SkillsRequired      []string  `json:"skills_required"`
	NumberOfVacancies   int       `json:"number_of_vacancies"`
	ApplicationDeadline time.Time `json:"application_deadline"`
	ContactEmail        string    `json:"contact_email"`
	Province            string    `json:"province"`
	City                string    `json:"city"`
	SubDistrict         string    `json:"sub_district"`
	DetailLocation      string    `json:"detail_location"`
	Photo               string    `json:"photo"`
	Status              string    `json:"status"`
	TotalRegistrar      int       `json:"total_registrar"`

	BookmarkID *string `json:"bookmark_id"`

	RejectedReason string         `json:"rejected_reason,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
	PostType       string         `json:"post_type"`
}

type ResRegistrantVacancyHistory struct {
	ID          int      `json:"id"`
	Email       string   `json:"email"`
	Fullname    string   `json:"fullname"`
	Address     string   `json:"address"`
	PhoneNumber string   `json:"phone_number"`
	Gender      string   `json:"gender"`
	Nik         string   `json:"nik"`
	Skills      []string `json:"skills_required"`
	Resume      string   `json:"resume"`
	Reason      string   `json:"reason"`
	Photo       string   `json:"photo"`
	Status      string   `json:"status"`
	PostType    string   `json:"post_type"`
}

type ResTransactionHistory struct {
	ID             int    `json:"transaction_id"`
	UserID         int    `json:"user_id"`
	Email          string `json:"email"`
	Fullname       string `json:"fullname"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	FundraiseID    int    `json:"fundraise_id"`
	FudraiseTitle  string `json:"fundraise_title"`
	FundraisePhoto string `json:"fundraise_photo"`
	Amount         int    `json:"amount"`
	PaymentType    string `json:"payment_type"`
	VirtualAccount string `json:"virtual_account"`
	UrlCallback    string `json:"url_callback"`
	PaidAt         string `json:"paid_at"`
	Status         string `json:"status"`
	PostType       string `json:"post_type"`
}

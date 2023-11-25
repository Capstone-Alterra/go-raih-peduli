package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResVacancy struct {
	ID                  int            `json:"id"`
	UserID              int            `json:"user_id"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	SkillsRequired      []string       `json:"skills_requred"`
	NumberOfVacancies   int            `json:"number_of_vacancies"`
	ApplicationDeadline time.Time      `json:"application_deadline"`
	ContactEmail        string         `json:"contact_email"`
	Province            string         `json:"province"`
	City                string         `json:"city"`
	SubDistrict         string         `json:"sub_district"`
	DetailLocation      string         `json:"detail_location"`
	Photo               string         `json:"photo"`
	Status              string         `json:"status"`
	TotalRegistrar      int            `json:"total_registrar"`
	BookmarkID          *string        `json:"bookmark_id"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
}

type ResRegistrantVacancy struct {
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	Nik      string `json:"nik"`
	Resume   string `json:"resume"`
	Photo    string `json:"photo"`
	Status   string `json:"status"`
}

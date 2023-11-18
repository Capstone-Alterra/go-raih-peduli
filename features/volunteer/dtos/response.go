package dtos

import (
	"time"

	"gorm.io/gorm"
)

type ResVolunteer struct {
	ID                  int            `json:"id"`
	UserID              int            `json:"user_id"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	SkillsRequired      string         `json:"skills_requred"`
	NumberOfVacancies   int            `json:"number_of_vacancies"`
	ApplicationDeadline time.Time      `json:"application_deadline"`
	ContactEmail        string         `json:"contact_email"`
	Province            string         `json:"province"`
	City                string         `json:"city"`
	SubDistrict         string         `json:"sub_district"`
	Photo               string         `json:"photo"`
	Status              string         `json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
}

type PaginationResponse struct {
	TotalData    int64 `json:"total_data"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage int   `json:"previous_page"`
	NextPage     int   `json:"next_page"`
	TotalPage    int   `json:"total_page"`
}

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
	SkillsRequired      []string       `json:"skills_required"`
	NumberOfVacancies   int            `json:"number_of_vacancies"`
	ApplicationDeadline time.Time      `json:"application_deadline"`
	ContactEmail        string         `json:"contact_email"`
	Province            string         `json:"province"`
	City                string         `json:"city"`
	SubDistrict         string         `json:"sub_district"`
	DetailLocation      string         `json:"detail_location"`
	Photo               string         `json:"photo"`
	Status              string         `json:"status"`
	TotalRegistrar      int            `json:"total_registrants"`
	BookmarkID          string        `json:"bookmark_id"`
	RejectedReason      string         `json:"rejected_reason,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
}

type ResRegistrantVacancy struct {
	ID       		int    `json:"id"`
	Email 	 		string `json:"email"`
	Fullname 		string `json:"fullname"`
	Address  		string `json:"address"`
	PhoneNumber 	string `json:"phone_number"`
	Gender    		string `json:"gender"`
	Nik      		string `json:"nik"`
	Skills   		[]string `json:"skills_required"`
	Resume   		string `json:"resume"`
	Reason 			string `json:"reason"`
	Photo    		string `json:"photo"`
	Status   		string `json:"status"`
}

type Skill struct {
	ID		int `json:"id"`
	Name 	string `json:"name"`
}

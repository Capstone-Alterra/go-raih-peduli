package dtos

import (
	"mime/multipart"
	"time"
)

type InputVacancy struct {
	Title               string         `form:"title" validate:"required" json:"title"`
	Description         string         `form:"description" validate:"required" json:"description"`
	SkillsRequired      []string       `form:"skills_required" validate:"required" json:"skills_required"`
	NumberOfVacancies   int            `form:"number_of_vacancies" validate:"required" json:"number_of_vacancies"`
	ApplicationDeadline time.Time      `form:"application_deadline" validate:"required" json:"application_deadline"`
	ContactEmail        string         `form:"contact_email" validate:"required" json:"contact_email"`
	Province            string         `form:"province" validate:"required" json:"province"`
	City                string         `form:"city" validate:"required" json:"city"`
	SubDistrict         string         `form:"sub_district" validate:"required" json:"sub_district"`
	DetailLocation      string         `form:"detail_location" validate:"required" json:"detail_location"`
	Photo               multipart.File `json:"photo" form:"photo"`
	Status              string         `json:"status" form:"status"`
}

type ApplyVacancy struct {
	VolunteerID int    `json:"volunteer_id" form:"volunteer_id" validate:"required"`
	Skills      string `json:"skills" form:"skills" validate:"required"`
	Resume      string `json:"resume" form:"resume" validate:"required"`
	Reason      string `json:"reason" form:"reason" validate:"required"`
}

type StatusRegistrar struct {
	Status string `json:"status" form:"status"`
}

type Pagination struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type SearchAndFilter struct {
	Title          string `query:"title"`
	City           string `query:"city"`
	Skill          string `query:"skill"`
	MinParticipant int32  `query:"min_participant"`
	MaxParticipant int32  `query:"max_participant"`
}

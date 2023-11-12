package dtos

import "time"

type InputVolunteer struct {
	Title               string    `form:"title" validate:"required" json:"tittle"`
	Description         string    `form:"description" validate:"required" json:"description"`
	SkillsRequired      string    `form:"skills_required" validate:"required" json:"skills_required"`
	NumberOfVacancies   int       `form:"number_of_vacancies" validate:"required" json:"number_of_vacancies"`
	ApplicationDeadline time.Time `form:"application_deadline" validate:"required" json:"application_deadline"`
	ContactEmail        string    `form:"contact_email" validate:"required" json:"contact_email"`
	Location            string    `form:"location" validate:"required" json:"location"`
	Photo               string    `form:"photo" validate:"required" json:"photo"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

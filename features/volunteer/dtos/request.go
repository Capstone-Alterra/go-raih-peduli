package dtos

import "time"

type InputVolunteer struct {
	Title               string    `form:"title" validate:"required"`
	Description         string    `form:"description" validate:"required"`
	SkillsRequired      string    `form:"skills_required" validate:"required"`
	NumberOfVacancies   int       `form:"number_of_vacancies" validate:"required"`
	ApplicationDeadline time.Time `form:"application_deadline" validate:"required"`
	ContactEmail        string    `form:"contact_email" validate:"required"`
	Location            string    `form:"location" validate:"required"`
	Photo               string    `form:"photo" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

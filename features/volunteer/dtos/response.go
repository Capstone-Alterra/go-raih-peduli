package dtos

import "time"

type ResVolunteer struct {
	SkillsRequired string `json:"skills_requred"`
	NumberOfVacancies int `json:"number_of_vacancies"`
	ApplicationDeadline time.Duration `json:"application_deadline"`
	ContactEmail string `json:"contact_email"`
	Location string `json:"location"`
}

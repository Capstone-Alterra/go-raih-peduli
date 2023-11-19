package dtos

import "time"

type ResBookmark struct {
	Fundraise []ResFundraise
	News []ResNews
	Vacancy []ResVolunteerVacancy
}

type ResFundraise struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Target      int32     `json:"target"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`

	BookmarkID  string 	  `json:"bookmark_id"`
}

type ResNews struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`

	BookmarkID  string `json:"bookmark_id"`
}

type ResVolunteerVacancy struct {
	ID          		int       `json:"id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	SkillsRequired      string    `json:"skills_requred"`
	NumberOfVacancies   int       `json:"number_of_vacancies"`
	ApplicationDeadline time.Time `json:"application_deadline"`
	ContactEmail        string    `json:"contact_email"`
	Province 			string 	  `json:"province"`
	City 				string 	  `json:"city"`
	SubDistrict 		string 	  `json:"sub_district"`
	Photo               string    `json:"photo"`

	BookmarkID  string 	  `json:"bookmark_id"`
}

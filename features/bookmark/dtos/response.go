package dtos

import (
	"time"
)

type ResOwnerID struct {
	OwnerID int `bson:"owner_id"`
}

type ResBookmark struct {
	Fundraise []ResFundraise `json:"fundraise"`
	News []ResNews `json:"news"`
	Vacancy []ResVolunteerVacancy `json:"vacancy"`
}

type ResFundraise struct {
	BookmarkID  string 	  `json:"bookmark_id" bson:"_id"`

	ID          int       `json:"post_id" bson:"post_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Target      int32     `json:"target"`
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	EndDate     time.Time `json:"end_date" bson:"end_date"`
	Status		string	  `json:"status"`

}

type ResNews struct {
	BookmarkID  string `json:"bookmark_id" bson:"_id"`

	ID          int    `json:"post_id" bson:"post_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

type ResVolunteerVacancy struct {
	BookmarkID  		string 	  `json:"bookmark_id" bson:"_id"`

	ID          		int       `json:"post_id" bson:"post_id"`
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
}

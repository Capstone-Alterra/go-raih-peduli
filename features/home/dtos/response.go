package dtos

type ResHome struct {
	Name string `json:"name"`
}

type ResGetHome struct {
	Fundraise []ResFundraise `json:"fundraise"`
	Volunteer []ResVolunteer `json:"volunteer"`
	News      []ResNews      `json:"news"`
}
type ResFundraise struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Target      int32  `json:"target"`
}
type ResVolunteer struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Photo             string `json:"photo"`
	NumberOfVacancies int    `json:"number_of_vacancies"`
}
type ResNews struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

type ResWebGetHome struct {
	UserAmount      int            `json:"user_amount"`
	FundraiseAmount int            `json:"fundraise_amount"`
	VolunteerAmount int            `json:"volunteer_amount"`
	NewsAmount      int            `json:"news_amount"`
	Fundraise       []ResFundraise `json:"fundraise"`
	Volunteer       []ResVolunteer `json:"volunteer"`
}
type ResWebFundraise struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Target      int32  `json:"target"`
}
type ResWebVolunteer struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Photo             string `json:"photo"`
	NumberOfVacancies int    `json:"number_of_vacancies"`
}

package dtos

type InputFundraiseID struct {
	FundraiseID int `form:"fundraise_id" json:"fundraise_id" validate:"required"`
}

type InputNewsID struct {
	NewsID int `form:"news_id" json:"news_id" validate:"required"`
}

type InputVacancyID struct {
	VacancyID int `form:"vacancy_id" json:"vacancy_id" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
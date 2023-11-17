package bookmark

type NewsBookmark struct {
	ID int `bson:"id"`
	UserID int `bson:"user_id"`
	NewsID int `bson:"news_id"`
}

type FundraiseBookmark struct {
	ID int `bson:"id"`
	UserID int `bson:"user_id"`
	FundraiseID int `bson:"fundraise_id"`
}

type VacancyBookmark struct {
	ID int `bson:"id"`
	UserID int `bson:"user_id"`
	VacancyID int `bson:"vacancy_id"`
}

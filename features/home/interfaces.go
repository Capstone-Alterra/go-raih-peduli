package home

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home/dtos"
	"raihpeduli/features/news"
	"raihpeduli/features/user"
	"raihpeduli/features/volunteer"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	PaginateFundraise(page, size int, likeQuery, notLikeQuery string) []fundraise.Fundraise
	PaginateVolunteer(page, size int, likeQuery, notLikeQuery string) []volunteer.VolunteerVacancies
	PaginateNews(page, size int, likeQuery, notLikeQuery string) []news.News
	CountUser() int
	CountFundraise() int
	CountVolunteer() int
	CountNews() int
	SelectUserByID(userID int) *user.User
	Insert(newHome Home) int64
	SelectByID(homeID int) *Home
	Update(home Home) int64
	DeleteByID(homeID int) int64
}

type Usecase interface {
	FindAll(page, size int, personalization []string) dtos.ResGetHome
	FindAllWeb(page, size int) dtos.ResWebGetHome
	GetPersonalization(userID int) []string
	FindByID(homeID int) *dtos.ResHome
	Create(newHome dtos.InputHome) *dtos.ResHome
	Modify(homeData dtos.InputHome, homeID int) bool
	Remove(homeID int) bool
}

type Handler interface {
	GetMobileLanding() echo.HandlerFunc
	GetWebLanding() echo.HandlerFunc
}

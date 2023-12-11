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
	PaginateFundraise(page, size int, personalization []string) []fundraise.Fundraise
	PaginateVolunteer(page, size int, personalization []string) []volunteer.VolunteerVacancies
	PaginateNews(page, size int, personalization []string) []news.News
	CountUser() int
	CountFundraise() int
	CountVolunteer() int
	CountNews() int
	SelectUserByID(userID int) *user.User
}

type Usecase interface {
	FindAll(page, size int, personalization []string) dtos.ResGetHome
	FindAllWeb(page, size int) dtos.ResWebGetHome
	GetPersonalization(userID int) []string
}

type Handler interface {
	GetMobileLanding() echo.HandlerFunc
	GetWebLanding() echo.HandlerFunc
}

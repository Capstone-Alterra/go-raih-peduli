package home

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home/dtos"
	"raihpeduli/features/news"
	"raihpeduli/features/volunteer"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	PaginateFundraise(page, size int) []fundraise.Fundraise
	PaginateVolunteer(page, size int) []volunteer.VolunteerVacancies
	PaginateNews(page, size int) []news.News
	CountUser() int
	CountFundraise() int
	CountVolunteer() int
	CountNews() int
	Insert(newHome Home) int64
	SelectByID(homeID int) *Home
	Update(home Home) int64
	DeleteByID(homeID int) int64
}

type Usecase interface {
	FindAll(page, size int) dtos.ResGetHome
	FindAllWeb(page, size int) dtos.ResWebGetHome
	FindByID(homeID int) *dtos.ResHome
	Create(newHome dtos.InputHome) *dtos.ResHome
	Modify(homeData dtos.InputHome, homeID int) bool
	Remove(homeID int) bool
}

type Handler interface {
	GetMobileLanding() echo.HandlerFunc
	GetWebLanding() echo.HandlerFunc
}

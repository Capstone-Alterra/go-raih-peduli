package fundraise

import (
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page int, size int, title string) ([]dtos.ResFundraise, error)
	SelectByID(fundraiseID int) (*dtos.ResFundraise, error)
	DeleteByID(fundraiseID int) (int, error)
}

type Usecase interface {
	FindAll(page int, size int, title string) []dtos.ResFundraise
	FindByID(fundraiseID int) *dtos.ResFundraise
	Remove(fundraiseID int) bool
}

type Handler interface {
	GetFundraises() echo.HandlerFunc
	FundraiseDetails() echo.HandlerFunc
	DeleteFundraise() echo.HandlerFunc
}

package fundraise

import (
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) ([]Fundraise, error)
	Insert(newFundraise Fundraise) (int, error)
	SelectByID(fundraiseID int) (*Fundraise, error)
	Update(fundraise Fundraise) (int, error)
	DeleteByID(fundraiseID int) (int, error)
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResFundraise
	FindByID(fundraiseID int) *dtos.ResFundraise
	Create(newFundraise dtos.InputFundraise) (*dtos.ResFundraise, error)
	Modify(fundraiseData dtos.InputFundraise, fundraiseID int) bool
	Remove(fundraiseID int) bool
}

type Handler interface {
	GetFundraises() echo.HandlerFunc
	FundraiseDetails() echo.HandlerFunc
	CreateFundraise() echo.HandlerFunc
	UpdateFundraise() echo.HandlerFunc
	DeleteFundraise() echo.HandlerFunc
}

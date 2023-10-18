package fundraise

import (
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Fundraise
	Insert(newFundraise Fundraise) int64
	SelectByID(fundraiseID int) *Fundraise
	Update(fundraise Fundraise) int64
	DeleteByID(fundraiseID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResFundraise
	FindByID(fundraiseID int) *dtos.ResFundraise
	Create(newFundraise dtos.InputFundraise) *dtos.ResFundraise
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

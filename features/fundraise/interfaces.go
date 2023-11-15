package fundraise

import (
	"mime/multipart"
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page int, size int, title string) ([]Fundraise, error)
	Insert(newFundraise Fundraise) (int, error)
	SelectByID(fundraiseID int) (*Fundraise, error)
	Update(fundraise Fundraise) (int, error)
	DeleteByID(fundraiseID int) (int, error)
	UploadFile(file multipart.File, objectName string) (string, error)
}

type Usecase interface {
	FindAll(page int, size int, title string) []dtos.ResFundraise
	FindByID(fundraiseID int) *dtos.ResFundraise
	Create(newFundraise dtos.InputFundraise, userID int, file multipart.File) (*dtos.ResFundraise, []string, error)
	Modify(fundraiseData dtos.InputFundraise, file multipart.File, oldData dtos.ResFundraise) bool
	Remove(fundraiseID int) bool
}

type Handler interface {
	GetFundraises() echo.HandlerFunc
	FundraiseDetails() echo.HandlerFunc
	CreateFundraise() echo.HandlerFunc
	UpdateFundraise() echo.HandlerFunc
	DeleteFundraise() echo.HandlerFunc
}

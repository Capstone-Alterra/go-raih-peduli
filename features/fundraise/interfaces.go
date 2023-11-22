package fundraise

import (
	"mime/multipart"
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]Fundraise, error)
	Insert(newFundraise Fundraise) (int, error)
	SelectByID(fundraiseID int) (*Fundraise, error)
	Update(fundraise Fundraise) (int, error)
	DeleteByID(fundraiseID int) (int, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	SelectBookmarkedFundraiseID(ownerID int) (map[int]string, error)
	SelectBookmarkByFundraiseAndOwnerID(fundraiseID, ownerID int) (string, error)
}

type Usecase interface {
	FindAll(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter, ownerID int) []dtos.ResFundraise
	FindByID(fundraiseID, ownerID int) *dtos.ResFundraise
	Create(newFundraise dtos.InputFundraise, userID int, file multipart.File) (*dtos.ResFundraise, []string, error)
	Modify(fundraiseData dtos.InputFundraise, file multipart.File, oldData dtos.ResFundraise) ([]string , error)
	ModifyStatus(fundraiseData dtos.InputFundraiseStatus, oldData dtos.ResFundraise) ([]string, error)
	Remove(fundraiseID int) bool
}

type Handler interface {
	GetFundraises() echo.HandlerFunc
	FundraiseDetails() echo.HandlerFunc
	CreateFundraise() echo.HandlerFunc
	UpdateFundraise() echo.HandlerFunc
	UpdateFundraiseStatus() echo.HandlerFunc
	DeleteFundraise() echo.HandlerFunc
}

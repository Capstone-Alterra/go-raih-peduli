package bookmark

import (
	"raihpeduli/features/bookmark/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(size, userID int) (*dtos.ResBookmark, error)
	Insert(document any) (bool, error)
	SelectFundraiseByID(fundraiseID int) (*Fundraise, error)
	SelectNewsByID(newsID int) (*News, error)
	SelectVolunteerByID(volunteerID int) (*VolunteerVacancy, error)
	DeleteByID(bookmarkID int) (int, error)
}

type Usecase interface {
	FindAll(size, userID int) *dtos.ResBookmark
	// FindByID(bookmarkID int) *dtos.ResBookmark
	SetBookmark(input dtos.InputBookmarkPost, userID int) (bool, error)
	UnsetBookmark(bookmarkID int) bool
}

type Handler interface {
	GetBookmarksByUserID() echo.HandlerFunc
	BookmarkAPost() echo.HandlerFunc
	UnBookmarkAPost() echo.HandlerFunc
}

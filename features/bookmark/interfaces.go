package bookmark

import (
	"raihpeduli/features/bookmark/dtos"
	"raihpeduli/features/fundraise"
	fundraiseDtos "raihpeduli/features/fundraise/dtos"

	"raihpeduli/features/news"
	newsDtos "raihpeduli/features/news/dtos"

	"raihpeduli/features/volunteer"
	volunteerDtos "raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Bookmark
	Insert(newBookmark Bookmark) int64
	SelectFundraiseByID(fundraiseID int) (*fundraise.Fundraise, error)
	SelectNewsByID(newsID int) (*news.News, error)
	SelectVolunteerByID(volunteerID int) (*volunteer.VolunteerVacancies, error)
	DeleteByID(bookmarkID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResBookmark
	FindByID(bookmarkID int) *dtos.ResBookmark
	FindFundraiseByID(fundraiseID int) *fundraiseDtos.ResFundraise
	FindNewsByID(newsID int) *newsDtos.ResNews
	FindVacancyByID(volunteerID int) *volunteerDtos.ResVolunteer
	SetBookmark(postID int, userID int, postType string) *dtos.ResBookmark
	UnsetBookmark(postID int, userID int, postType string) bool
}

type Handler interface {
	GetBookmarksByUserID() echo.HandlerFunc
	BookmarkAFundraise() echo.HandlerFunc
	BookmarkANews() echo.HandlerFunc
	BookmarkAVacancy() echo.HandlerFunc
	UnBookmarkAPost() echo.HandlerFunc
}

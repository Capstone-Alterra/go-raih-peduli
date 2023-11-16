package bookmark

import (
	"raihpeduli/features/bookmark/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Bookmark
	Insert(newBookmark Bookmark) int64
	SelectByID(bookmarkID int) *Bookmark
	DeleteByID(bookmarkID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResBookmark
	FindByID(bookmarkID int) *dtos.ResBookmark
	SetBookmark(newBookmark dtos.InputBookmark) *dtos.ResBookmark
	UnsetBookmark(bookmarkID int) bool
}

type Handler interface {
	GetBookmarksByUserID() echo.HandlerFunc
	BookmarkAPost() echo.HandlerFunc
	UnBookmarkAPost() echo.HandlerFunc
}

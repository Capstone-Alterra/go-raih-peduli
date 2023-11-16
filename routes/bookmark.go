package routes

import (
	"raihpeduli/features/bookmark"

	"github.com/labstack/echo/v4"
)

func Bookmarks(e *echo.Echo, handler bookmark.Handler) {
	bookmarks := e.Group("/bookmarks")

	bookmarks.GET("", handler.GetBookmarksByUserID())
	bookmarks.POST("", handler.BookmarkAPost())
	
	bookmarks.DELETE("/:id", handler.UnBookmarkAPost())
}
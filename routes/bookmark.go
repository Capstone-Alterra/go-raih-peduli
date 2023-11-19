package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/bookmark"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Bookmarks(e *echo.Echo, handler bookmark.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	bookmarks := e.Group("/bookmarks")

	bookmarks.GET("", handler.GetBookmarksByUserID(), m.AuthorizeJWT(jwt, 0, config.SECRET))
	bookmarks.POST("", handler.BookmarkAPost(), m.AuthorizeJWT(jwt, 0, config.SECRET))
	
	bookmarks.DELETE("/:id", handler.UnBookmarkAPost(), m.AuthorizeJWT(jwt, 0, config.SECRET))
}
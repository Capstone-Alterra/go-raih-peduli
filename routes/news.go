package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/news"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func News(e *echo.Echo, handler news.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	mobileNews := e.Group("/mobile/news")

	mobileNews.GET("", handler.GetNews(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	mobileNews.GET("/:id", handler.NewsDetails(), m.AuthorizeJWT(jwt, -1, config.SECRET))

	news := e.Group("/news")

	news.GET("", handler.GetNews(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	news.POST("", handler.CreateNews(), m.AuthorizeJWT(jwt, 2, config.SECRET))

	news.GET("/:id", handler.NewsDetails(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	news.PUT("/:id", handler.UpdateNews(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	news.DELETE("/:id", handler.DeleteNews(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}

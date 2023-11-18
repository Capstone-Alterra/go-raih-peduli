package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/news"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func News(e *echo.Echo, handler news.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	news := e.Group("/news")
	
	news.GET("", handler.GetNews())
	news.POST("", handler.CreateNews(), m.AuthorizeJWT(jwt, 1, config.SECRET))

	news.GET("/:id", handler.NewsDetails())
	news.PUT("/:id", handler.UpdateNews())
	news.DELETE("/:id", handler.DeleteNews())
}
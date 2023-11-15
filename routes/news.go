package routes

import (
	"raihpeduli/features/news"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func News(e *echo.Echo, handler news.Handler, jwt helpers.JWTInterface) {
	news := e.Group("/news")
	news.GET("", handler.GetNews())
	news.POST("", handler.CreateNews(), m.AuthorizeJWT(jwt, 1))
}

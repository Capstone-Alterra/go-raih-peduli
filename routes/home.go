package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/home"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Homes(e *echo.Echo, handler home.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	homes := e.Group("/home")

	homes.GET("/mobile", handler.GetMobileLanding(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	homes.GET("/web", handler.GetWebLanding(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}

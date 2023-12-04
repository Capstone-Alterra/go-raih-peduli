package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/history"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func History(e *echo.Echo, handler history.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	history := e.Group("/history")

	history.GET("", handler.GetAllHistory(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	history.GET("/fundraises", handler.GetHistoryFundraiseCreatedByUser(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	history.GET("/volunteer-vacancies", handler.GetHistoryVolunteerVacanciesCreatedByUser(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	history.GET("/volunteer-vacancies/registered", handler.GetHistoryVolunteerVacanciesRegisterByUser(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	history.GET("/donations", handler.GetHistoryUserTransaction(), m.AuthorizeJWT(jwt, 1, config.SECRET))
}

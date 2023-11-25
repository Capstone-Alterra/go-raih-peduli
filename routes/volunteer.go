package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Volunteers(e *echo.Echo, handler volunteer.Handler, jwt helpers.JWTInterface, cfg config.ProgramConfig) {
	mobile := e.Group("/mobile/volunteer-vacancies")
	mobile.GET("", handler.GetVacancies("mobile"), m.AuthorizeJWT(jwt, -1, cfg.SECRET))
	mobile.GET("/:id", handler.VacancyDetails(), m.AuthorizeJWT(jwt, -1, cfg.SECRET))
	mobile.POST("", handler.CreateVacancy(), m.AuthorizeJWT(jwt, 1, cfg.SECRET))
	mobile.POST("/register", handler.ApplyVacancy(), m.AuthorizeJWT(jwt, 1, cfg.SECRET))

	volunteers := e.Group("/volunteer-vacancies")

	volunteers.POST("", handler.CreateVacancy(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.GET("", handler.GetVacancies(""), m.AuthorizeJWT(jwt, -1, cfg.SECRET))
	volunteers.GET("/:id", handler.VacancyDetails(), m.AuthorizeJWT(jwt, -1, cfg.SECRET))
	volunteers.PUT("/:id", handler.UpdateVacancy(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.PATCH("/:id", handler.UpdateStatusVacancy(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.DELETE("/:id", handler.DeleteVacancy(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))

	volunteers.GET("/:vacancy_id/registrants", handler.GetVolunteersByVacancyID(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.GET("/:vacancy_id/registrants/:volunteer_id", handler.GetVolunteer(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.PATCH("/:vacancy_id/registrants/:volunteer_id", handler.UpdateStatusRegistrar(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
}

package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Volunteers(e *echo.Echo, handler volunteer.Handler, jwt helpers.JWTInterface, cfg config.ProgramConfig) {
	volunteers := e.Group("/volunteer-vacancies")

	volunteers.POST("", handler.CreateVacancy(), m.AuthorizeJWT(jwt, 0, cfg.SECRET))
	volunteers.GET("", handler.GetVacancies())
	volunteers.GET("/:id", handler.VacancyDetails())
	volunteers.PUT("/:id", handler.UpdateVacancy())
	volunteers.DELETE("/:id", handler.DeleteVacancy())

	volunteers.POST("/register", handler.ApplyVacancy(), m.AuthorizeJWT(jwt, 1, cfg.SECRET))
	volunteers.GET("/:vacancy_id/registrants", handler.GetVolunteersByVacancyID())
	volunteers.GET("/:vacancy_id/registrants/:volunteer_id", handler.GetVolunteer())
	volunteers.PATCH("/update-status-registrar/:id", handler.UpdateStatusRegistrar(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
}

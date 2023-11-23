package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Volunteers(e *echo.Echo, handler volunteer.Handler, jwt helpers.JWTInterface, cfg config.ProgramConfig) {
	volunteers := e.Group("/volunteers")

	volunteers.GET("", handler.GetVolunteers())
	volunteers.POST("", handler.CreateVolunteer(), m.AuthorizeJWT(jwt, 0, cfg.SECRET))
	volunteers.POST("/register", handler.ApplyVacancies(), m.AuthorizeJWT(jwt, 1, cfg.SECRET))
	volunteers.GET("/registrar/:vacancy_id", handler.GetVolunteerByVacancyID())

	volunteers.GET("/:id", handler.VolunteerDetails())
	volunteers.PUT("/:id", handler.UpdateVolunteer())
	volunteers.PATCH("/:id", handler.UpdateVolunteer())
	volunteers.PATCH("/update-status-registrar/:id", handler.UpdateStatusRegistrar(), m.AuthorizeJWT(jwt, 2, cfg.SECRET))
	volunteers.DELETE("/:id", handler.DeleteVolunteer())
}

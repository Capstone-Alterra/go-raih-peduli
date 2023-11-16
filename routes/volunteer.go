package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Volunteers(e *echo.Echo, handler volunteer.Handler, jwt helpers.JWTInterface ,cfg config.ProgramConfig) {
	volunteers := e.Group("/volunteers")

	volunteers.GET("", handler.GetVolunteers())
	volunteers.POST("", handler.CreateVolunteer(), m.AuthorizeJWT(jwt, 0, cfg.SECRET))
	
	volunteers.GET("/:id", handler.VolunteerDetails())
	volunteers.PUT("/:id", handler.UpdateVolunteer())
	volunteers.DELETE("/:id", handler.DeleteVolunteer())
}
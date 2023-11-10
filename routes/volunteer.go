package routes

import (
	"raihpeduli/features/volunteer"

	"github.com/labstack/echo/v4"
)

func Volunteers(e *echo.Echo, handler volunteer.Handler) {
	volunteers := e.Group("/volunteers")

	volunteers.GET("", handler.GetVolunteers())
	volunteers.POST("", handler.CreateVolunteer())
	
	volunteers.GET("/:id", handler.VolunteerDetails())
	volunteers.PUT("/:id", handler.UpdateVolunteer())
	volunteers.DELETE("/:id", handler.DeleteVolunteer())
}
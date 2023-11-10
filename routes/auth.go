package routes

import (
	"raihpeduli/features/auth"

	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo, handler auth.Handler) {
	admins := e.Group("/auth")
	admins.POST("", handler.LoginCustomer())

}

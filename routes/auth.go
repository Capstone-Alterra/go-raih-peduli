package routes

import (
	"raihpeduli/features/auth"

	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo, handler auth.Handler) {
	users := e.Group("/auth")
	users.POST("/login", handler.LoginCustomer())
	users.POST("/register", handler.RegisterUser())
}

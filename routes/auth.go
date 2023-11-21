package routes

import (
	"raihpeduli/features/auth"

	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo, handler auth.Handler) {
	users := e.Group("/auth")
	users.POST("/login", handler.Login())
	users.POST("/register", handler.RegisterUser())
	users.POST("/resend-otp", handler.ResendOTP())
}

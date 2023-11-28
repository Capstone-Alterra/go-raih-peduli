package routes

import (
	"raihpeduli/features/auth"

	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo, handler auth.Handler) {
	auth := e.Group("/auth")
	auth.POST("/login", handler.Login())
	auth.POST("/register", handler.RegisterUser())
	auth.POST("/resend-otp", handler.ResendOTP())
	auth.POST("/refresh-jwt", handler.RefreshJWT())
}

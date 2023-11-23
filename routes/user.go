package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/user"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Users(e *echo.Echo, handler user.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	users := e.Group("/users")

	users.GET("", handler.GetUsers())
	users.POST("", handler.CreateUser())
	users.POST("/verify", handler.VerifyEmail())
	users.POST("/forget-password", handler.ForgetPassword())
	users.POST("/verify-otp", handler.VerifyOTP())
	users.POST("/reset-password", handler.ResetPassword(), m.AuthorizeJWT(jwt, 1, config.OTP_SECRET))

	users.GET("/:id", handler.UserDetails())
	users.PUT("", handler.UpdateUser(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	users.PATCH("", handler.UpdateProfilePicture(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	users.DELETE("/:id", handler.DeleteUser())
}

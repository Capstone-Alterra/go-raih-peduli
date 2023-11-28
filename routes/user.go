package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/user"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Users(e *echo.Echo, handler user.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	mobile := e.Group("/mobile/users")
	mobile.GET("/my-profile", handler.MyProfile(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	mobile.GET("/:id", handler.UserDetails(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	mobile.POST("/verify", handler.VerifyEmail(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	mobile.POST("/forget-password", handler.ForgetPassword(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	mobile.POST("/verify-otp", handler.VerifyOTP(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	mobile.POST("/reset-password", handler.ResetPassword(), m.AuthorizeJWT(jwt, 1, config.OTP_SECRET))
	mobile.PUT("", handler.UpdateUser(), m.AuthorizeJWT(jwt, 1, config.SECRET))
	mobile.PATCH("", handler.UpdateProfilePicture(), m.AuthorizeJWT(jwt, 1, config.SECRET))

	users := e.Group("/users")
	users.GET("", handler.GetUsers(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	users.POST("", handler.CreateUser(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	users.POST("/verify", handler.VerifyEmail(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	users.POST("/forget-password", handler.ForgetPassword(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	users.POST("/verify-otp", handler.VerifyOTP(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	users.POST("/reset-password", handler.ResetPassword(), m.AuthorizeJWT(jwt, 2, config.OTP_SECRET))
	users.GET("/:id", handler.UserDetails(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	users.PUT("/:id", handler.UpdateUser(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	users.DELETE("/:id", handler.DeleteUser(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}

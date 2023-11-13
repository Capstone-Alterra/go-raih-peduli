package routes

import (
	"raihpeduli/features/user"

	"github.com/labstack/echo/v4"
)

func Users(e *echo.Echo, handler user.Handler) {
	users := e.Group("/users")

	users.GET("", handler.GetUsers())
	users.POST("", handler.CreateUser())
	users.POST("/verify", handler.VerifyEmail())

	users.GET("/:id", handler.UserDetails())
	users.PUT("/:id", handler.UpdateUser())
	users.DELETE("/:id", handler.DeleteUser())
}

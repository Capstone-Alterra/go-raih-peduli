package routes

import (
	"raihpeduli/features/admin"

	"github.com/labstack/echo/v4"
)

func Admins(e *echo.Echo, handler admin.Handler) {
	admins := e.Group("/admins")

	admins.GET("", handler.GetAdmins())
	admins.POST("", handler.CreateAdmin())
	admins.POST("/login", handler.LoginAdmin())
	
	admins.GET("/:id", handler.AdminDetails())
	admins.PUT("/:id", handler.UpdateAdmin())
	admins.DELETE("/:id", handler.DeleteAdmin())
}
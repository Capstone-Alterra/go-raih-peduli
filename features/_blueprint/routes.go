package routes

import (
	"raihpeduli/features/_blueprint"

	"github.com/labstack/echo/v4"
)

func Placeholders(e *echo.Echo, handler _blueprint.Handler) {
	placeholders := e.Group("/placeholders")

	placeholders.GET("", handler.GetPlaceholders())
	placeholders.POST("", handler.CreatePlaceholder())
	
	placeholders.GET("/:id", handler.PlaceholderDetails())
	placeholders.PUT("/:id", handler.UpdatePlaceholder())
	placeholders.DELETE("/:id", handler.DeletePlaceholder())
}
package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RunServer()  {
	e := echo.New()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "Hello World")
	})

	e.Start(":8000")
}
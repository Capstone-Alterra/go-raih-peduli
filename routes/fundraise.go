package routes

import (
	"raihpeduli/features/fundraise"

	"github.com/labstack/echo/v4"
)

func Fundraises(e *echo.Echo, handler fundraise.Handler, jwtMiddleware echo.MiddlewareFunc) {
	fundraises := e.Group("/fundraises")

	fundraises.Use(jwtMiddleware)

	fundraises.GET("", handler.GetFundraises())
	fundraises.POST("", handler.CreateFundraise())

	fundraises.GET("/:id", handler.FundraiseDetails())
	fundraises.PUT("/:id", handler.UpdateFundraise())
	fundraises.DELETE("/:id", handler.DeleteFundraise())
}

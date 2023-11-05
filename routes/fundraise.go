package routes

import (
	"raihpeduli/features/fundraise"

	"github.com/labstack/echo/v4"
)

func Fundraises(e *echo.Echo, handler fundraise.Handler) {
	fundraises := e.Group("/fundraises")

	fundraises.GET("", handler.GetFundraises())
	
	fundraises.GET("/:id", handler.FundraiseDetails())
	fundraises.DELETE("/:id", handler.DeleteFundraise())
}
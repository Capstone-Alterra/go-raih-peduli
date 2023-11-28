package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Fundraises(e *echo.Echo, handler fundraise.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	mobile := e.Group("/mobile/fundraises")
	
	mobile.GET("", handler.GetFundraises("mobile"), m.AuthorizeJWT(jwt, -1, config.SECRET))
	mobile.POST("", handler.CreateFundraise(), m.AuthorizeJWT(jwt, 0, config.SECRET))

	mobile.GET("/:id", handler.FundraiseDetails(), m.AuthorizeJWT(jwt, -1, config.SECRET))

	fundraises := e.Group("/fundraises")
	
	fundraises.GET("", handler.GetFundraises(""), m.AuthorizeJWT(jwt, -1, config.SECRET))
	fundraises.POST("", handler.CreateFundraise(), m.AuthorizeJWT(jwt, 0, config.SECRET))

	fundraises.GET("/:id", handler.FundraiseDetails(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	fundraises.PUT("/:id", handler.UpdateFundraise(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	fundraises.PATCH("/:id", handler.UpdateFundraiseStatus(), m.AuthorizeJWT(jwt, 2, config.SECRET))
	fundraises.DELETE("/:id", handler.DeleteFundraise(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}

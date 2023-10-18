package main

import (
	"raihpeduli/config"
	"raihpeduli/features"
	"raihpeduli/helpers"
	"raihpeduli/middlewares"
	"raihpeduli/routes"

	"github.com/labstack/echo/v4"
)

var (
	fundraiseHandler = features.FundraiseHandler()
)

func main() {

	e := echo.New()
	var config = config.InitConfig()

	jwtInterface := helpers.New(config.Secret, config.RefreshSecret)
	jwtMiddleware := middlewares.AuthorizeJWT(jwtInterface)

	routes.Fundraises(e, fundraiseHandler, jwtMiddleware)

	e.Start(":8000")
}

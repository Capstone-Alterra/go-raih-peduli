package main

import (
	"raihpeduli/features"
	"raihpeduli/routes"

	"github.com/labstack/echo/v4"
)

var (
	userHandler = features.UserHandler()
	authHandler = features.AuthHandler()
)

func main() {

	e := echo.New()
	// var config = config.InitConfig()

	// jwtInterface := helpers.New(config.Secret, config.RefreshSecret)
	// jwtMiddleware := middlewares.AuthorizeJWT(jwtInterface)

	routes.Users(e, userHandler)
	routes.Auth(e, authHandler)

	e.Start(":8000")
}

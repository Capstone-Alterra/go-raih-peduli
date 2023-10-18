package main

import (
	"raihpeduli/features"
	"raihpeduli/routes"

	"github.com/labstack/echo/v4"
)

var (
	fundraiseHandler = features.FundraiseHandler()
)

func main() {

	e := echo.New()

	routes.Fundraises(e, fundraiseHandler)

	e.Start(":8000")
}
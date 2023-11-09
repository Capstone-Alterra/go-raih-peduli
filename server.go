package main

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/routes"
	"raihpeduli/utils"

	"raihpeduli/features/fundraise"
	fh "raihpeduli/features/fundraise/handler"
	fr "raihpeduli/features/fundraise/repository"
	fu "raihpeduli/features/fundraise/usecase"

	"github.com/labstack/echo/v4"
)


func main() {
	cfg := config.LoadServerConfig()
	e := echo.New()

	routes.Fundraises(e, FundraiseHandler())

	e.Start(fmt.Sprintf(":%d", cfg.SERVER_PORT))
}

func FundraiseHandler() fundraise.Handler {
	db := utils.InitDB()
	repo := fr.New(db)
	uc := fu.New(repo)
	return fh.New(uc)
}

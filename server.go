package main

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/helpers"
	"raihpeduli/routes"
	"raihpeduli/utils"

	"raihpeduli/features/auth"
	ah "raihpeduli/features/auth/handler"
	ar "raihpeduli/features/auth/repository"
	au "raihpeduli/features/auth/usecase"

	"raihpeduli/features/user"
	uh "raihpeduli/features/user/handler"
	ur "raihpeduli/features/user/repository"
	uu "raihpeduli/features/user/usecase"

	"raihpeduli/features/fundraise"
	fh "raihpeduli/features/fundraise/handler"
	fr "raihpeduli/features/fundraise/repository"
	fu "raihpeduli/features/fundraise/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	jwtService := helpers.New(cfg.Secret, cfg.RefreshSecret)

	// routes.Auth(e, AuthHandler())
	// routes.Users(e, UserHandler())
	routes.Fundraises(e, FundraiseHandler(), jwtService)

	e.Start(fmt.Sprintf(":%s", cfg.SERVER_PORT))
}

func FundraiseHandler() fundraise.Handler {
	db := utils.InitDB()
	repo := fr.New(db)
	uc := fu.New(repo)
	return fh.New(uc)
}

func UserHandler() user.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := ur.New(db, redis)
	uc := uu.New(repo, jwt, hash)
	return uh.New(uc)
}

func AuthHandler() auth.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := ar.New(db, redis)
	uc := au.New(repo, jwt, hash)
	return ah.New(uc)
}
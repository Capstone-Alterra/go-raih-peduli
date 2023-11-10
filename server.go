package main

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/helpers"
	"raihpeduli/routes"
	"raihpeduli/utils"

	"raihpeduli/features/admin"
	adh "raihpeduli/features/admin/handler"
	adr "raihpeduli/features/admin/repository"
	adu "raihpeduli/features/admin/usecase"

	"raihpeduli/features/auth"
	ah "raihpeduli/features/auth/handler"
	ar "raihpeduli/features/auth/repository"
	au "raihpeduli/features/auth/usecase"

	"raihpeduli/features/customer"
	ch "raihpeduli/features/customer/handler"
	cr "raihpeduli/features/customer/repository"
	cu "raihpeduli/features/customer/usecase"

	"raihpeduli/features/fundraise"
	fh "raihpeduli/features/fundraise/handler"
	fr "raihpeduli/features/fundraise/repository"
	fu "raihpeduli/features/fundraise/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()

	routes.Admins(e, AdminHandler())
	routes.Customers(e, CustomerHandler())
	routes.Auth(e, AuthHandler())
	routes.Fundraises(e, FundraiseHandler())

	e.Start(fmt.Sprintf(":%d", cfg.SERVER_PORT))
}

func FundraiseHandler() fundraise.Handler {
	db := utils.InitDB()
	repo := fr.New(db)
	uc := fu.New(repo)
	return fh.New(uc)
}

func AdminHandler() admin.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := adr.New(db)
	uc := adu.New(repo, jwt, hash)
	return adh.New(uc)
}

func CustomerHandler() customer.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := cr.New(db)
	uc := cu.New(repo, jwt, hash)
	return ch.New(uc)
}

func AuthHandler() auth.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := ar.New(db)
	uc := au.New(repo, jwt, hash)
	return ah.New(uc)
}
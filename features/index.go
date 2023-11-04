package features

import (
	"raihpeduli/config"
	"raihpeduli/features/admin"
	"raihpeduli/features/auth"
	"raihpeduli/features/customer"

	adminHandler "raihpeduli/features/admin/handler"
	adminRepo "raihpeduli/features/admin/repository"
	adminUsecase "raihpeduli/features/admin/usecase"

	customerHandler "raihpeduli/features/customer/handler"
	customerRepo "raihpeduli/features/customer/repository"
	customerUsecase "raihpeduli/features/customer/usecase"

	authHandler "raihpeduli/features/auth/handler"
	authRepo "raihpeduli/features/auth/repository"
	authUsecase "raihpeduli/features/auth/usecase"

	"raihpeduli/helpers"
	"raihpeduli/utils"
)

func AdminHandler() admin.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := adminRepo.New(db)
	uc := adminUsecase.New(repo, jwt, hash)
	return adminHandler.New(uc)
}

func CustomerHandler() customer.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := customerRepo.New(db)
	uc := customerUsecase.New(repo, jwt, hash)
	return customerHandler.New(uc)
}

func AuthHandler() auth.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()

	repo := authRepo.New(db)
	uc := authUsecase.New(repo, jwt, hash)
	return authHandler.New(uc)
}

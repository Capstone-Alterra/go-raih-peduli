package features

import (
	"raihpeduli/config"
	"raihpeduli/features/admin"
	"raihpeduli/features/customer"

	adminHandler "raihpeduli/features/admin/handler"
	adminRepo "raihpeduli/features/admin/repository"
	adminUsecase "raihpeduli/features/admin/usecase"

	customerHandler "raihpeduli/features/customer/handler"
	customerRepo "raihpeduli/features/customer/repository"
	customerUsecase "raihpeduli/features/customer/usecase"

	"raihpeduli/helpers"
	"raihpeduli/utils"
)

func AdminHandler() admin.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)

	repo := adminRepo.New(db)
	uc := adminUsecase.New(repo, jwt)
	return adminHandler.New(uc)
}

func CustomerHandler() customer.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)

	repo := customerRepo.New(db)
	uc := customerUsecase.New(repo, jwt)
	return customerHandler.New(uc)
}

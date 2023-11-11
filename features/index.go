package features

import (
	"raihpeduli/config"
	"raihpeduli/features/auth"
	"raihpeduli/features/user"

	userHandler "raihpeduli/features/user/handler"
	userRepo "raihpeduli/features/user/repository"
	userUsecase "raihpeduli/features/user/usecase"

	authHandler "raihpeduli/features/auth/handler"
	authRepo "raihpeduli/features/auth/repository"
	authUsecase "raihpeduli/features/auth/usecase"

	"raihpeduli/helpers"
	"raihpeduli/utils"
)

func UserHandler() user.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := userRepo.New(db, redis)
	uc := userUsecase.New(repo, jwt, hash)
	return userHandler.New(uc)
}

func AuthHandler() auth.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(config.Secret, config.RefreshSecret)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := authRepo.New(db, redis)
	uc := authUsecase.New(repo, jwt, hash)
	return authHandler.New(uc)
}

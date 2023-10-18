package features

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/handler"
	"raihpeduli/features/fundraise/repository"
	"raihpeduli/features/fundraise/usecase"
	"raihpeduli/utils"
)

func FundraiseHandler() fundraise.Handler {
	db := utils.InitDB()
	repo := repository.New(db)
	uc := usecase.New(repo)
	return handler.New(uc)
}
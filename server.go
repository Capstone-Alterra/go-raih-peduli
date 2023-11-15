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

	"raihpeduli/features/volunteer"
	vh "raihpeduli/features/volunteer/handler"
	vr "raihpeduli/features/volunteer/repository"
	vu "raihpeduli/features/volunteer/usecase"

	"raihpeduli/features/news"
	nh "raihpeduli/features/news/handler"
	nr "raihpeduli/features/news/repository"
	nu "raihpeduli/features/news/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	jwtService := helpers.New(*cfg)

	routes.Auth(e, AuthHandler())
	routes.Users(e, UserHandler(), jwtService, *cfg)
	routes.Fundraises(e, FundraiseHandler(), jwtService, *cfg)
	routes.Volunteers(e, VolunteerHandler())
	routes.News(e, NewsHandler(), jwtService)

	e.Start(fmt.Sprintf(":%s", cfg.SERVER_PORT))
}

func FundraiseHandler() fundraise.Handler {
	config := config.LoadCloudStorageConfig()

	db := utils.InitDB()
	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "fundraises/")
	repo := fr.New(db, clStorage)
	uc := fu.New(repo)
	return fh.New(uc)
}

func UserHandler() user.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(*config)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := ur.New(db, redis)
	uc := uu.New(repo, jwt, hash)
	return uh.New(uc)
}

func AuthHandler() auth.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.New(*config)
	hash := helpers.NewHash()
	redis := utils.ConnectRedis()

	repo := ar.New(db, redis)
	uc := au.New(repo, jwt, hash)
	return ah.New(uc)
}

func VolunteerHandler() volunteer.Handler {

	db := utils.InitDB()

	repo := vr.New(db)
	uc := vu.New(repo)
	return vh.New(uc)
}
<<<<<<< Updated upstream
=======

func NewsHandler() news.Handler {
	db := utils.InitDB()
	config := config.LoadCloudStorageConfig()

	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "news/")
	repo := nr.New(db, clStorage)
	uc := nu.New(repo)
	return nh.New(uc)
}
>>>>>>> Stashed changes

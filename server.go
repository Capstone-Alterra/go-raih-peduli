package main

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/helpers"
	"raihpeduli/middlewares"
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

	"raihpeduli/features/transaction"
	th "raihpeduli/features/transaction/handler"
	tr "raihpeduli/features/transaction/repository"
	tu "raihpeduli/features/transaction/usecase"

	"raihpeduli/features/bookmark"
	bh "raihpeduli/features/bookmark/handler"
	br "raihpeduli/features/bookmark/repository"
	bu "raihpeduli/features/bookmark/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	jwtService := helpers.NewJWT(*cfg)

	routes.Auth(e, AuthHandler())
	routes.Users(e, UserHandler(), jwtService, *cfg)
	routes.Fundraises(e, FundraiseHandler(), jwtService, *cfg)
	routes.Volunteers(e, VolunteerHandler(), jwtService, *cfg)
	routes.News(e, NewsHandler(), jwtService, *cfg)
	routes.Transactions(e, TransactionHandler(), jwtService, *cfg)
	routes.Bookmarks(e, BookmarkHandler(), jwtService, *cfg)

	middlewares.LogMiddlewares(e)

	e.Start(fmt.Sprintf(":%s", cfg.SERVER_PORT))
}

func FundraiseHandler() fundraise.Handler {
	config := config.LoadCloudStorageConfig()
	validation := helpers.NewValidationRequest()

	db := utils.InitDB()
	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "fundraises/")
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")

	repo := fr.New(db, clStorage, collection)
	uc := fu.New(repo, validation)
	return fh.New(uc)
}

func UserHandler() user.Handler {
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.NewJWT(*config)
	hash := helpers.NewHash()
	generator := helpers.NewGenerator()
	redis := utils.ConnectRedis()

	repo := ur.New(db, redis)
	uc := uu.New(repo, jwt, hash, generator)
	return uh.New(uc)
}

func AuthHandler() auth.Handler {
	smtpConfig := config.LoadSMTPConfig()
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.NewJWT(*config)
	hash := helpers.NewHash()
	generator := helpers.NewGenerator()
	validation := helpers.NewValidationRequest()
	converter := helpers.NewConverter()
	redis := utils.ConnectRedis()

	repo := ar.New(db, redis, smtpConfig)
	uc := au.New(repo, jwt, hash, generator, validation, converter)
	return ah.New(uc)
}

func VolunteerHandler() volunteer.Handler {
	config := config.LoadCloudStorageConfig()
	validation := helpers.NewValidationRequest()

	db := utils.InitDB()

	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "fundraises/")
	repo := vr.New(db, clStorage)
	uc := vu.New(repo, validation)
	return vh.New(uc)
}

func NewsHandler() news.Handler {
	db := utils.InitDB()
	config := config.LoadCloudStorageConfig()

	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "news/")
	repo := nr.New(db, clStorage)
	uc := nu.New(repo)
	return nh.New(uc)
}

func TransactionHandler() transaction.Handler {
	db := utils.InitDB()
	repo := tr.New(db)
	coreAPIClient := utils.MidtransCoreAPIClient()

	generator := helpers.NewGenerator()
	midtrans := helpers.NewMidtransRequest()
	tc := tu.New(repo, generator, midtrans, coreAPIClient)
	return th.New(tc)
}

func BookmarkHandler() bookmark.Handler {
	db := utils.InitDB()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")

	repo := br.New(db, collection)
	uc := bu.New(repo)
	return bh.New(uc)
}

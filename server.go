package main

import (
	"fmt"
	"raihpeduli/config"
	"raihpeduli/features/home"
	"raihpeduli/helpers"
	"raihpeduli/middlewares"
	"raihpeduli/routes"
	"raihpeduli/utils"

	"raihpeduli/features/auth"
	ah "raihpeduli/features/auth/handler"
	ar "raihpeduli/features/auth/repository"
	au "raihpeduli/features/auth/usecase"
	"raihpeduli/features/history"

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

	"raihpeduli/features/chatbot"
	ch "raihpeduli/features/chatbot/handler"
	cr "raihpeduli/features/chatbot/repository"
	cu "raihpeduli/features/chatbot/usecase"

	hh "raihpeduli/features/history/handler"
	hr "raihpeduli/features/history/repository"
	hu "raihpeduli/features/history/usecase"

	hoh "raihpeduli/features/home/handler"
	hor "raihpeduli/features/home/repository"
	hou "raihpeduli/features/home/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	jwtService := helpers.NewJWT(*cfg)

	middlewares.LogMiddlewares(e)

	routes.Auth(e, AuthHandler())
	routes.Users(e, UserHandler(), jwtService, *cfg)
	routes.Fundraises(e, FundraiseHandler(), jwtService, *cfg)
	routes.Volunteers(e, VolunteerHandler(), jwtService, *cfg)
	routes.News(e, NewsHandler(), jwtService, *cfg)
	routes.Transactions(e, TransactionHandler(), jwtService, *cfg)
	routes.Bookmarks(e, BookmarkHandler(), jwtService, *cfg)
	routes.Chatbots(e, ChatbotHandler(cfg), jwtService, *cfg)
	routes.History(e, HistoryHandler(), jwtService, *cfg)
	routes.Homes(e, HomeHandler(), jwtService, *cfg)

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

func HomeHandler() home.Handler {
	validation := helpers.NewValidationRequest()

	db := utils.InitDB()

	repo := hor.New(db)
	uc := hou.New(repo, validation)
	return hoh.New(uc)
}

func UserHandler() user.Handler {
	cloud := config.LoadCloudStorageConfig()
	config := config.InitConfig()

	db := utils.InitDB()
	jwt := helpers.NewJWT(*config)
	hash := helpers.NewHash()
	generator := helpers.NewGenerator()
	validation := helpers.NewValidationRequest()
	redis := utils.ConnectRedis()
	clStorage := helpers.NewCloudStorage(cloud.CLOUD_PROJECT_ID, cloud.CLOUD_BUCKET_NAME, "users/")

	repo := ur.New(db, redis, clStorage)
	uc := uu.New(repo, jwt, hash, generator, validation)
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
	redis := utils.ConnectRedis()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("devices")

	repo := ar.New(db, redis, smtpConfig, collection)
	uc := au.New(repo, jwt, hash, generator, validation)
	return ah.New(uc)
}

func VolunteerHandler() volunteer.Handler {
	config := config.LoadCloudStorageConfig()
	validation := helpers.NewValidationRequest()

	db := utils.InitDB()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")

	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "vacancies/")

	repo := vr.New(db, clStorage, collection)
	uc := vu.New(repo, validation)
	return vh.New(uc)
}

func NewsHandler() news.Handler {
	db := utils.InitDB()
	config := config.LoadCloudStorageConfig()
	validation := helpers.NewValidationRequest()

	clStorage := helpers.NewCloudStorage(config.CLOUD_PROJECT_ID, config.CLOUD_BUCKET_NAME, "news/")
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")

	repo := nr.New(db, clStorage, collection)
	uc := nu.New(repo, validation)
	return nh.New(uc)
}

func TransactionHandler() transaction.Handler {
	config.LoadFirebaseConfig()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("devices")
	smtpConfig := config.LoadSMTPConfig()
	db := utils.InitDB()
	repo := tr.New(db, smtpConfig, collection)
	coreAPIClient := utils.MidtransCoreAPIClient()
	validation := helpers.NewValidationRequest()
	nfService := helpers.NewNotificationService()

	generator := helpers.NewGenerator()
	midtrans := helpers.NewMidtransRequest()

	tc := tu.New(repo, generator, midtrans, coreAPIClient, validation, nfService)
	return th.New(tc)
}

func BookmarkHandler() bookmark.Handler {
	db := utils.InitDB()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")
	validation := helpers.NewValidationRequest()

	repo := br.New(db, collection)
	uc := bu.New(repo, validation)
	return bh.New(uc)
}

func ChatbotHandler(cfg *config.ProgramConfig) chatbot.Handler {
	db := utils.InitDB()
	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("chatbot_histories")

	validation := helpers.NewValidationRequest()
	openAI := helpers.NewOpenAI(cfg.OPENAI_KEY)

	repo := cr.New(db, collection)
	uc := cu.New(repo, validation, openAI)
	return ch.New(uc)
}

func HistoryHandler() history.Handler {
	db := utils.InitDB()

	mongoDB := utils.ConnectMongo()
	collection := mongoDB.Collection("bookmarks")

	repo := hr.New(db, collection)
	uc := hu.New(repo)
	return hh.New(uc)
}

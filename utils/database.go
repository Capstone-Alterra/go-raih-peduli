package utils

import (
	"context"
	"raihpeduli/config"
	"raihpeduli/features/auth"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/news"
	"raihpeduli/features/transaction"
	"raihpeduli/features/volunteer"

	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	config := config.LoadDBConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(fundraise.Fundraise{}, &auth.User{}, &volunteer.VolunteerVacancies{}, news.News{}, transaction.Transaction{}, &volunteer.VolunteerRelations{}, &volunteer.Skill{})
}

func ConnectRedis() *redis.Client {
	config := config.LoadRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.REDIS_HOST, config.REDIS_PORT),
		DB:   1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}

	logrus.Info("Connection established")

	return client
}

func ConnectMongo() *mongo.Database {
	config := config.LoadMongoConfig()

	clientOptions := options.Client()
	clientOptions.ApplyURI(config.MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil
	}

	return client.Database(config.MONGO_DB_NAME)
}

func FirebaseInit() *messaging.Client {
	// config := config.LoadFirebaseConfig()

	// Use the path to your service account credential json file
	opt := option.WithCredentialsFile("firebase_key.json")

	// Create a new firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Print("Failed Connect Firbase", err)
		return nil
	}

	// Get the FCM object
	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		logrus.Print("Failed Connect Firbase", err)
		return nil
	}

	return fcmClient
}

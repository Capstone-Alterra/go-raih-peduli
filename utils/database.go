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

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	config := config.LoadDBConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(fundraise.Fundraise{}, &auth.User{}, &volunteer.VolunteerVacancies{}, news.News{}, transaction.Transaction{}, &volunteer.VolunteerRelations{}, )
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
    clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.MONGO_HOST, config.MONGO_PORT))
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil
    }

    return client.Database(config.MONGO_DB_NAME)
}

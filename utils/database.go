package utils

import (
	"raihpeduli/config"
	"raihpeduli/features/auth"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/volunteer"

	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
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
	db.AutoMigrate(fundraise.Fundraise{}, &auth.User{}, &volunteer.VolunteerVacancies{})
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

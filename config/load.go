package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitConfig() *ProgramConfig {
	var res = new(ProgramConfig)
	res = loadConfig()

	if res == nil {
		logrus.Fatal("Config : Cannot start program, failed to load configuration")
		return nil
	}

	return res
}

type DatabaseConfig struct {
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
	DB_NAME string
}

type RedisConfig struct {
	REDIS_HOST string
	REDIS_PORT string
}

type CloudStorageConfig struct {
	GOOGLE_APPLICATION_CREDENTIALS string
	CLOUD_PROJECT_ID               string
	CLOUD_BUCKET_NAME              string
}

type ProgramConfig struct {
	SECRET         string
	REFRESH_SECRET string
	SERVER_PORT    string
	OTP_SECRET     string
}

type SMTPConfig struct {
	SMTP_USER string
	SMTP_PASS string
	SMTP_PORT string
}

func LoadDBConfig() *DatabaseConfig {
	var res = new(DatabaseConfig)

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Error("Config : Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("DB_USER"); found {
		res.DB_USER = val
	}

	if val, found := os.LookupEnv("DB_PASS"); found {
		res.DB_PASS = val
	}

	if val, found := os.LookupEnv("DB_HOST"); found {
		res.DB_HOST = val
	}

	if val, found := os.LookupEnv("DB_PORT"); found {
		res.DB_PORT = val
	}

	if val, found := os.LookupEnv("DB_NAME"); found {
		res.DB_NAME = val
	}

	return res
}

func LoadRedisConfig() *RedisConfig {
	var res = new(RedisConfig)

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Error("Config : Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("REDIS_HOST"); found {
		res.REDIS_HOST = val
	}

	if val, found := os.LookupEnv("REDIS_PORT"); found {
		res.REDIS_PORT = val
	}

	return res
}

func LoadCloudStorageConfig() *CloudStorageConfig {
	var res = new(CloudStorageConfig)

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Error("Config : Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); found {
		res.GOOGLE_APPLICATION_CREDENTIALS = val
	}

	if val, found := os.LookupEnv("CLOUD_PROJECT_ID"); found {
		res.CLOUD_PROJECT_ID = val
	}

	if val, found := os.LookupEnv("CLOUD_BUCKET_NAME"); found {
		res.CLOUD_BUCKET_NAME = val
	}

	return res
}

func LoadSMTPConfig() *SMTPConfig {
	var res = new(SMTPConfig)

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Error("Config : Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("SMTP_USER"); found {
		res.SMTP_USER = val
	}

	if val, found := os.LookupEnv("SMTP_PASS"); found {
		res.SMTP_PASS = val
	}

	if val, found := os.LookupEnv("SMTP_PORT"); found {
		res.SMTP_PORT = val
	}

	return res
}

func loadConfig() *ProgramConfig {
	var res = new(ProgramConfig)

	err := godotenv.Load(".env")

	if err != nil {
		logrus.Error("Config : Cannot load config file,", err.Error())
		return nil
	}

	if val, found := os.LookupEnv("SECRET"); found {
		res.SECRET = val
	}

	if val, found := os.LookupEnv("REFSECRET"); found {
		res.REFRESH_SECRET = val
	}

	if val, found := os.LookupEnv("SERVER_PORT"); found {
		res.SERVER_PORT = val
	}

	if val, found := os.LookupEnv("OTP_SECRET"); found {
		res.OTP_SECRET = val
	}

	return res
}

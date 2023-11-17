package config

import (
	"io"
	"os"

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

type MidtransConfig struct {
	MT_SERVER_KEY string
	MT_CLIENT_KEY string
}

type ProgramConfig struct {
	SECRET         string
	REFRESH_SECRET string
	SERVER_PORT    string
	OTP_SECRET     string
}

func LoadDBConfig() *DatabaseConfig {
	var res = new(DatabaseConfig)

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

	if val, found := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); found {
		gcredentials, _ := os.LookupEnv("APPLICATION_DEFAULT_CREDENTIALS")

		file, err := os.Create("credentials.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.WriteString(file, gcredentials)
		if err != nil {
			panic(err)
		}

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

func LoadMidtransConfig() *MidtransConfig {
	var res = new(MidtransConfig)

	if val, found := os.LookupEnv("MT_SERVER_KEY"); found {
		res.MT_SERVER_KEY = val
	}

	if val, found := os.LookupEnv("MT_CLIENT_KEY"); found {
		res.MT_CLIENT_KEY = val
	}

	return res
}

func loadConfig() *ProgramConfig {
	var res = new(ProgramConfig)

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

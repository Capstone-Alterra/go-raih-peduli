package config

import (
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func InitConfig() *ProgramConfig {
	godotenv.Load()

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

type MongoConfig struct {
	MONGO_URI     string
	MONGO_DB_NAME string
}

type CloudStorageConfig struct {
	GOOGLE_APPLICATION_CREDENTIALS string
	CLOUD_PROJECT_ID               string
	CLOUD_BUCKET_NAME              string
}

type FirebaseConfig struct {
	FIREBASE_API_KEY string
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
	OPENAI_KEY     string
}

type SMTPConfig struct {
	SMTP_USER string
	SMTP_PASS string
	SMTP_PORT string
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

func LoadMongoConfig() *MongoConfig {
	var res = new(MongoConfig)

	if val, found := os.LookupEnv("MONGO_URI"); found {
		res.MONGO_URI = val
	}

	if val, found := os.LookupEnv("MONGO_DB_NAME"); found {
		res.MONGO_DB_NAME = val
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

func LoadFirebaseConfig() *FirebaseConfig {
	var res = new(FirebaseConfig)

	// if val, found := os.LookupEnv("FIREBASE_API_KEY"); found {
	// 	gcredentials, _ := os.LookupEnv("FIREBASE_API_KEY")

	// 	file, err := os.Create("firebase_key.json")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()

	// 	_, err = io.WriteString(file, gcredentials)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	res.FIREBASE_API_KEY = val
	// }

	res.FIREBASE_API_KEY = "firebase_key.json"
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

func LoadSMTPConfig() *SMTPConfig {
	var res = new(SMTPConfig)

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

	if val, found := os.LookupEnv("OPENAI_KEY"); found {
		res.OPENAI_KEY = val
	}

	return res
}

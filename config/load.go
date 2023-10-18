package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DB_USER string
	DB_PASS string
	DB_HOST	string
	DB_PORT	int
	DB_NAME string
}

type ServerConfig struct {
	SERVER_PORT string
}

func LoadDBConfig() DatabaseConfig {
	godotenv.Load(".env")

	DB_PORT, err := strconv.Atoi(os.Getenv("DB_PORT"))
	
	if err != nil {
		panic(err)
	}

	return DatabaseConfig {
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: DB_PORT,
		DB_NAME: os.Getenv("DB_NAME"),
	}
}
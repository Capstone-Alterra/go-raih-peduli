package utils

import (
	"raihpeduli/config"
	"raihpeduli/features/content"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/user"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	config := config.LoadDBConfig()
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(user.User{})
	db.AutoMigrate(fundraise.Fundraise{}, content.Content{})
}
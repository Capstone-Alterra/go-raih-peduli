package utils

import (
	"raihpeduli/config"
	"raihpeduli/features/admin"
	"raihpeduli/features/auth"
	"raihpeduli/features/customer"
	"raihpeduli/features/fundraise"

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
	db.AutoMigrate(fundraise.Fundraise{}, &auth.User{}, &admin.Admin{}, &customer.Customer{})
}

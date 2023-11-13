package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int    `gorm:"primaryKey;type:int(11)"`
	RoleID         int    `gorm:"type:int(1);default:1"`
	IsVerified     bool   `gorm:"default:false"`
	Email          string `gorm:"type:varchar(255);not null"`
	Password       string `gorm:"type:varchar(255);not null"`
	ProfilePicture string `gorm:"varchar(255)"`
	Fullname       string `gorm:"varchar(255)"`
	Gender         string `gorm:"varchar(255)"`
	Address        string `gorm:"varchar(255)"`
	PhoneNumber    string `gorm:"varchar(255)"`
	Nik            string `gorm:"varchar(255)"`
	Status         string `gorm:"varchar(255);default:1"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

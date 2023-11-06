package customer

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID          int    `gorm:"primaryKey;type:int(11)"`
	UserID      int    `gorm:"type:int(11)"`
	Fullname    string `gorm:"type:varchar(255);not null"`
	Address     string `gorm:"type:varchar(255);not null"`
	PhoneNumber string `gorm:"type:varchar(13);not null"`
	Gender      string `gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID             int    `gorm:"primaryKey;type:int(11)"`
	RoleID         int    `gorm:"type:int(1);default:1"`
	Email          string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password       string `gorm:"type:varchar(255);not null"`
	ProfilePicture string `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

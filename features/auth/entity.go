package auth

import (
	"raihpeduli/features/admin"
	"raihpeduli/features/customer"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int               `gorm:"primaryKey;type:int(11)"`
	Admin          admin.Admin       `gorm:"foreignKey:UserID"`
	Custumer       customer.Customer `gorm:"foreignKey:UserID"`
	OTP            customer.OTP      `gorm:"foreignKey:UserID"`
	RoleID         int               `gorm:"type:int(1);default:1"`
	Verified       bool              `gorm:"default:false"`
	Email          string            `gorm:"type:varchar(255);not null"`
	Password       string            `gorm:"type:varchar(255);not null"`
	ProfilePicture string            `gorm:"varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

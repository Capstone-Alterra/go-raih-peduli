package user

import (
	"raihpeduli/features/volunteer"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              int    `gorm:"primaryKey;type:int(11)"`
	RoleID          int    `gorm:"type:int(1);default:1"`
	IsVerified      bool   `gorm:"default:false"`
	Email           string `gorm:"type:varchar(255);not null"`
	Password        string `gorm:"type:varchar(255);not null"`
	ProfilePicture  string `gorm:"type:varchar(255)"`
	Fullname        string `gorm:"type:varchar(100);not null"`
	Gender          string `gorm:"type:varchar(20);not null"`
	Address         string `gorm:"type:varchar(200)"`
	PhoneNumber     string `gorm:"type:varchar(20)"`
	Nik             string `gorm:"type:varchar(17)"`
	Status          string `gorm:"type:int(1);default:1"`
	Personalization string `gorm:"varchar(255)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	VolunteerVacancies     []volunteer.VolunteerVacancies `gorm:"foreignKey:UserID;references:ID"`
	VolunteerRelationships []volunteer.VolunteerRelations `gorm:"foreignKey:UserID;references:ID"`
}

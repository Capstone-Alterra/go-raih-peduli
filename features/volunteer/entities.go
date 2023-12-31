package volunteer

import (
	"time"

	"gorm.io/gorm"
)

type VolunteerVacancies struct {
	ID                  int    `gorm:"type:int(11); primaryKey"`
	UserID              int    `gorm:"type:int(11)"`
	Title               string `gorm:"type:varchar(255)"`
	Description         string `gorm:"type:varchar(255)"`
	SkillsRequired      string `gorm:"type:varchar(255)"`
	NumberOfVacancies   int    `gorm:"type:int(20)"`
	ApplicationDeadline time.Time
	ContactEmail        string `gorm:"type:varchar(30)"`
	Province            string `gorm:"type:varchar(255)"`
	City                string `gorm:"type:varchar(255)"`
	SubDistrict         string `gorm:"type:varchar(255)"`
	DetailLocation      string `gorm:"type:varchar(255)"`
	Photo               string `gorm:"type:varchar(255)"`
	Status              string `gorm:"type:enum('pending','accepted','rejected'); default: 'pending'"`
	RejectedReason      string `gorm:"type:varchar(255)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`

	VolunteerRelationships []VolunteerRelations `gorm:"foreignKey:VolunteerID;references:ID"`
}

type Volunteer struct {
	ID          int    `gorm:"type:int(11)"`
	Email       string `gorm:"type:varchar(255)"`
	Fullname    string `gorm:"type:varchar(255)"`
	Address     string `gorm:"type:varchar(255)"`
	PhoneNumber string `gorm:"type:varchar(255)"`
	Gender      string `gorm:"type:varchar(255)"`
	Skills      string `gorm:"type:varchar(255)"`
	Nik         string `gorm:"type:varchar(255)"`
	Resume      string `gorm:"type:varchar(255)"`
	Reason      string
	Photo       string `gorm:"type:varchar(255)"`
	Status      string `gorm:"type:enum('pending','accepted','rejected')"`
}

type VolunteerRelations struct {
	ID             int    `gorm:"type:int(11); primaryKey"`
	UserID         int    `gorm:"type:int(11)"`
	VolunteerID    int    `gorm:"type:int(11)"`
	Skills         string `gorm:"type:varchar(255)"`
	Reason         string
	Resume         string `gorm:"type:varchar(255)"`
	Photo          string `gorm:"type:varchar(255)"`
	Status         string `gorm:"type:enum('pending','accepted','rejected'); default: 'pending'"`
	RejectedReason string `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Skill struct {
	ID   int    `gorm:"type:int(11)"`
	Name string `gorm:"type:varchar(255)"`
}

type NotificationToken struct {
	DeviceToken string `bson:"device_token" json:"device_token"`
}

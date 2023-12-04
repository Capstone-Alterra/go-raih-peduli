package history

import (
	"raihpeduli/features/auth"
	"time"

	"gorm.io/gorm"
)

type Fundraise struct {
	ID             int       `gorm:"type:int(11)"`
	Title          string    `gorm:"type:varchar(255)"`
	Description    string    `gorm:"text"`
	Photo          string    `gorm:"type:varchar(255)"`
	Target         int32     `gorm:"type:int(11)"`
	StartDate      time.Time `gorm:"type:DATETIME"`
	EndDate        time.Time `gorm:"type:DATETIME"`
	Status         string    `gorm:"type:enum('pending','accepted','rejected')"`
	RejectedReason string    `gorm:"type:varchar(255)"`
	UserID         int       `gorm:"type:int(11)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User auth.User
}

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
}
type Transaction struct {
	gorm.Model

	ID             int    `gorm:"type:int(11)"`
	FundraiseID    int    `gorm:"type:int(11)"`
	UserID         int    `gorm:"type:int(11)"`
	PaymentType    string `gorm:"type:varchar(50)"`
	Amount         int    `gorm:"type:int(11)"`
	Status         string `gorm:"type:varchar(10)"`
	PaidAt         string `gorm:"type:varchar(100)"`
	VirtualAccount string `gorm:"type:varchar(100)"`
	UrlCallback    string `gorm:"type:varchar(250)"`
	CreatedAt      time.Time

	User      auth.User
	Fundraise Fundraise
}

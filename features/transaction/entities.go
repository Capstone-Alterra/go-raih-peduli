package transaction

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/user"

	"gorm.io/gorm"
)

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
	ValidUntil     string `gorm:"type:varchar(250)"`

	User      user.User
	Fundraise fundraise.Fundraise
}

type Status struct {
	Transaction string
	Order       string
}

type NotificationToken struct {
	DeviceToken  string          `bson:"device_token" json:"device_token"`
}

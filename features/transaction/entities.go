package transaction

import (
	"time"

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
	CreatedAt      time.Time
}

type Status struct {
	Transaction string
	Order       string
}

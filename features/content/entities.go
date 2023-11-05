package content

import (
	"raihpeduli/features/fundraise"

	"gorm.io/gorm"
)

type Content struct {
	gorm.Model

	ID int `gorm:"type:int(11)"`
	Title string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`
	Photo string `gorm:"type:varchar(255)"`
	FundraiseID int `gorm:"type:int(11)"`

	Fundraise fundraise.Fundraise
}
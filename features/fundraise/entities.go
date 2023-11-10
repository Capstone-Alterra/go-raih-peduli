package fundraise

import (
	"raihpeduli/features/auth"
	"time"

	"gorm.io/gorm"
)

type Fundraise struct {
	gorm.Model

	ID		  	int		  `gorm:"type:int(11)"`
	Title 	  	string 	  `gorm:"type:varchar(255)"`
	Description string	  `gorm:"text"`
	Photo	  	string 	  `gorm:"type:varchar(255)"`
	Target    	int32     `gorm:"type:int(11)"`
	StartDate 	time.Time `gorm:"type:DATETIME(3)"`
	EndDate   	time.Time `gorm:"type:DATETIME(3)"`
	Status		string	  `gorm:"type:enum('pending','live','closed')"`
	UserID    	int       `gorm:"type:int(11)"`

	User auth.User
}

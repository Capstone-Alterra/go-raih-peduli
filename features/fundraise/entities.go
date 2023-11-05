package fundraise

import (
	"raihpeduli/features/user"
	"time"

	"gorm.io/gorm"
)

type Fundraise struct {
	gorm.Model

	ID		  int		`gorm:"type:int(11)"`
	Target    string    `gorm:"type:varchar(255)"`
	StartDate time.Time `gorm:"type:DATETIME(3)"`
	EndDate   time.Time `gorm:"type:DATETIME(3)"`
	UserID    int       `gorm:"type:int(11)"`

	User user.User
}

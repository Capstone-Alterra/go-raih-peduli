package fundraise

import (
	"gorm.io/gorm"
	"time"
)

type Fundraise struct {
	gorm.Model
	User_id    int       `gorm:"type:int(11)"`
	Target     string    `gorm:"type:varchar(255)"`
	Start_date time.Time `gorm:"type:DATETIME(3)"`
	End_date   time.Time `gorm:"type:DATETIME(3)"`
}

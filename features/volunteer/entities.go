package volunteer

import (
	"time"

	"gorm.io/gorm"
)

type VolunteerVacancies struct {
	gorm.Model

	ID                  int    `gorm:"type:int(11); primaryKey"`
	SkillsRequired      string `gorm:"type:varchar(255)"`
	NumberOfVacancies   int    `gorm:"type:int(20)"`
	ApplicationDeadline time.Time
	ContactEmail        string `gorm:"type:varchar(30)"`
	Location            string `gorm:"type:varchar(255)"`
}

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
	Location            string `gorm:"type:varchar(255)"`
	Photo               string `gorm:"type:varchar(255)"`
	Status              string `gorm:"type:varchar(255)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

package fundraise

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
	RejectedReason string	 `gorm:"type:varchar(255)"`
	UserID         int       `gorm:"type:int(11)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User auth.User
}

func (m *Fundraise) BeforeCreate(tx *gorm.DB) error {
	wibLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	m.CreatedAt = time.Now().In(wibLocation)
	return nil
}

func (m *Fundraise) BeforeUpdate(tx *gorm.DB) error {
	wibLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	m.UpdatedAt = time.Now().In(wibLocation)
	return nil
}
package bookmark

import (
	"raihpeduli/features/auth"
	"time"

	"gorm.io/gorm"
)

type NewsBookmark struct {
	PostID 		int    `bson:"post_id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Photo       string `bson:"photo"`
	
	PostType	string `bson:"post_type"`
	OwnerID 	int    `bson:"owner_id"`
}

type FundraiseBookmark struct {
	PostID 		int 	  `bson:"post_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Photo       string    `bson:"photo"`
	Target      int32     `bson:"target"`
	StartDate   time.Time `bson:"start_date"`
	EndDate     time.Time `bson:"end_date"`
	Status      string    `bson:"status"`

	PostType	string    `bson:"post_type"`
	OwnerID     int 	  `bson:"owner_id"`
}

type VacancyBookmark struct {
	PostID 				int 	  `bson:"post_id"`
	Title               string    `bson:"title"`
	Description         string    `bson:"description"`
	SkillsRequired      string    `bson:"skills_requred"`
	NumberOfVacancies   int       `bson:"number_of_vacancies"`
	ApplicationDeadline time.Time `bson:"application_deadline"`
	ContactEmail        string    `bson:"contact_email"`
	Province 			string 	  `bson:"province"`
	City 				string 	  `bson:"city"`
	SubDistrict 		string 	  `bson:"sub_district"`
	Photo               string    `bson:"photo"`
	
	PostType			string    `bson:"post_type"`
	OwnerID    			int 	  `bson:"owner_id"`
}

type PostBookmark struct {
	BookmarkID int
}

type Fundraise struct {
	ID          int       `gorm:"type:int(11)"`
	Title       string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"text"`
	Photo       string    `gorm:"type:varchar(255)"`
	Target      int32     `gorm:"type:int(11)"`
	StartDate   time.Time `gorm:"type:DATETIME"`
	EndDate     time.Time `gorm:"type:DATETIME"`
	Status      string    `gorm:"type:enum('pending','live','closed')"`
	UserID      int       `gorm:"type:int(11)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User auth.User
}

type News struct {
	ID          int    `gorm:"type:int(11)"`
	Title       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	Photo       string `gorm:"type:varchar(255)"`
	UserID      int    `gorm:"type:int(11)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User auth.User
}

type VolunteerVacancy struct {
	ID                  int    `gorm:"type:int(11); primaryKey"`
	UserID      		int    `gorm:"type:int(11)"`
	Title               string `gorm:"type:varchar(255)"`
	Description         string `gorm:"type:varchar(255)"`
	SkillsRequired      string `gorm:"type:varchar(255)"`
	NumberOfVacancies   int    `gorm:"type:int(20)"`
	ApplicationDeadline time.Time
	ContactEmail        string `gorm:"type:varchar(30)"`
	Province			string `gorm:"type:varchar(255)"`
	City				string `gorm:"type:varchar(255)"`
	SubDistrict			string `gorm:"type:varchar(255)"`
	DetailLocation      string `gorm:"type:varchar(255)"`
	Photo               string `gorm:"type:varchar(255)"`
	Status              string `gorm:"type:enum('pending','accepted','rejected')"`
	
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

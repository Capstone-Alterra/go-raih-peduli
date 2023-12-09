package usecase

import (
	"raihpeduli/features/user"
	"raihpeduli/features/user/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entities = []user.User{
		{
			ID: 1,
			RoleID: 1,
			IsVerified: true,
			Email: "john@example.com",
			Password: "random123",
			ProfilePicture: "google.com",
			Fullname: "John Doe",
			Gender: "Male",
			Address         string `gorm:"type:varchar(200)"`
			PhoneNumber     string `gorm:"type:varchar(20)"`
			Nik             string `gorm:"type:varchar(17)"`
			Status          string `gorm:"type:int(1);default:1"`
			Personalization *string `gorm:"varchar(255)"`
			CreatedAt       time.Time
			UpdatedAt       time.Time
			DeletedAt       gorm.DeletedAt `gorm:"index"`

		}
	}

	var page = 1
	var pageSize = 10

	t.Run("Success", func(t *testing.T) {
		
	})

}
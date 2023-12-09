package usecase

import (
	"errors"
	"mime/multipart"
	"raihpeduli/features/user"
	"raihpeduli/features/user/dtos"
	"raihpeduli/features/user/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
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
			Address: "planet bumi",
			PhoneNumber: "080000000000",
			Nik: "21000000000000",
			Status: "1",
			Personalization: nil,
		},
	}

	var page = 1
	var pageSize = 10

	t.Run("Success", func(t *testing.T) {
		repository.On("Paginate", page, pageSize).Return(entities).Once()
		repository.On("GetTotalData").Return(int64(1)).Once()

		res, total := service.FindAll(page, pageSize)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})

}

func TestFindByID(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entity = user.User{
		ID: 1,
		RoleID: 1,
		IsVerified: true,
		Email: "john@example.com",
		Password: "random123",
		ProfilePicture: "google.com",
		Fullname: "John Doe",
		Gender: "Male",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		Nik: "21000000000000",
		Status: "1",
		Personalization: nil,
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()

		res := service.FindByID(userID)
		assert.Equal(t, res.ID, entity.ID)
		repository.AssertExpectations(t) 
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(nil).Once()

		res := service.FindByID(userID)
		assert.Nil(t, res)
		repository.AssertExpectations(t) 
	})
}

func TestCreate(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.InputUser{
		Fullname: "John Doe",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		Gender: "Male",
		Email: "john@example.com",
		Password: "random123",
	}

	var hashPassword = "iowjdoij1o2i3j12"

	var entity = user.User{
		IsVerified: false,
		Email: "john@example.com",
		Password: hashPassword,
		Fullname: "John Doe",
		Gender: "Male",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		Personalization: nil,
	}

	var otp = "123982"

	var errMap = []string{"email is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByEmail", input.Email).Return(nil, errors.New("user not found")).Once()
		hash.On("HashPassword", input.Password).Return(hashPassword).Once()
		repository.On("InsertUser", &entity).Return(&entity, nil).Once()
		generator.On("GenerateRandomOTP").Return(otp).Once()
		repository.On("SendOTPByEmail", input.Email, otp).Return(nil).Once()

		res, errMap, err := service.Create(input)
		assert.Equal(t, res.Email, entity.Email)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error When Send OTP", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByEmail", input.Email).Return(nil, errors.New("user not found")).Once()
		hash.On("HashPassword", input.Password).Return(hashPassword).Once()
		repository.On("InsertUser", &entity).Return(&entity, nil).Once()
		generator.On("GenerateRandomOTP").Return(otp).Once()
		repository.On("SendOTPByEmail", input.Email, otp).Return(errors.New("error send otp")).Once()

		res, errMap, err := service.Create(input)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error When Insert", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByEmail", input.Email).Return(nil, errors.New("user not found")).Once()
		hash.On("HashPassword", input.Password).Return(hashPassword).Once()
		repository.On("InsertUser", &entity).Return(nil, errors.New("error when insert")).Once()

		res, errMap, err := service.Create(input)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error Already Exist", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()

		res, errMap, err := service.Create(input)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error Validate Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		res, errMap, err := service.Create(input)
		assert.Nil(t, res)
		assert.NotNil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})
}

func TestModify(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.InputUpdate{
		Fullname: "John Doe",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		Gender: "Male",
		Email: "john@example.com",
		Nik: "1231241231232111",

	}

	var file multipart.File

	var oldData = dtos.ResUser{
		ID: 1,
		RoleID: 1,
		Email: "john@example.com",
		Fullname: "John Doe",
		Gender: "Male",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		ProfilePicture: "google.com",
	}

	var entity = user.User{
		ID: 1,
		Email: "john@example.com",
		ProfilePicture: "google.com",
		Fullname: "John Doe",
		Gender: "Male",
		Address: "planet bumi",
		PhoneNumber: "080000000000",
		Nik: "1231241231232111",
	}

	var errMap = []string{"email is required"}
	
	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.True(t, update)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.False(t, update)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return("", errors.New("error when upload")).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.False(t, update)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error NIK Less Than 16 Digits", func(t *testing.T) {
		input.Nik = "210000000"
		validation.On("ValidateRequest", input).Return(nil).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.False(t, update)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error NIK Not a Digit Of Number", func(t *testing.T) {
		input.Nik = "asdaasdas"
		validation.On("ValidateRequest", input).Return(nil).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.False(t, update)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		update, errMap := service.Modify(input, file, oldData)
		assert.False(t, update)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
}
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
			ID:              1,
			RoleID:          1,
			IsVerified:      true,
			Email:           "john@example.com",
			Password:        "random123",
			ProfilePicture:  "google.com",
			Fullname:        "John Doe",
			Gender:          "Male",
			Address:         "planet bumi",
			PhoneNumber:     "080000000000",
			Nik:             "21000000000000",
			Status:          "1",
			Personalization: nil,
		},
	}

	var emptySearchAndFilter = dtos.SearchAndFilter{
		Page:     1,
		PageSize: 10,
		Name:     "",
	}

	var searchAndFilter = dtos.SearchAndFilter{
		Page:     1,
		PageSize: 10,
		Name:     "Nibras",
	}

	t.Run("Success Witout Searching", func(t *testing.T) {
		repository.On("Paginate", emptySearchAndFilter).Return(entities).Once()
		repository.On("GetTotalData").Return(int64(1)).Once()

		res, total := service.FindAll(emptySearchAndFilter)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})

	t.Run("Success With Searching", func(t *testing.T) {
		repository.On("Paginate", searchAndFilter).Return(entities).Once()
		repository.On("GetTotalDataByName", searchAndFilter.Name).Return(int64(1)).Once()

		res, total := service.FindAll(searchAndFilter)
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
		ID:              1,
		RoleID:          1,
		IsVerified:      true,
		Email:           "john@example.com",
		Password:        "random123",
		ProfilePicture:  "google.com",
		Fullname:        "John Doe",
		Gender:          "Male",
		Address:         "planet bumi",
		PhoneNumber:     "080000000000",
		Nik:             "21000000000000",
		Status:          "1",
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
		Fullname:    "John Doe",
		Address:     "planet bumi",
		PhoneNumber: "080000000000",
		Gender:      "Male",
		Email:       "john@example.com",
		Password:    "random123",
	}

	var hashPassword = "iowjdoij1o2i3j12"

	var entity = user.User{
		IsVerified:      false,
		Email:           "john@example.com",
		Password:        hashPassword,
		Fullname:        "John Doe",
		Gender:          "Male",
		Address:         "planet bumi",
		PhoneNumber:     "080000000000",
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
		repository.On("SendOTPByEmail", input.Fullname, input.Email, otp, "1").Return(nil).Once()

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
		repository.On("SendOTPByEmail", input.Fullname, input.Email, otp, "1").Return(errors.New("error send otp")).Once()

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
		Fullname:    "John Doe",
		Address:     "planet bumi",
		PhoneNumber: "080000000000",
		Gender:      "Male",
		Email:       "john@example.com",
		Nik:         "1231241231232111",
	}

	var file multipart.File

	var oldData = dtos.ResUser{
		ID:             1,
		RoleID:         1,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var entity = user.User{
		ID:             1,
		Email:          "john@example.com",
		ProfilePicture: "google.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		Nik:            "1231241231232111",
	}

	var errMap = []string{"email is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return("", errors.New("error when upload")).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error NIK Less Than 16 Digits", func(t *testing.T) {
		input.Nik = "210000000"
		validation.On("ValidateRequest", input).Return(nil).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
	})

	t.Run("Failed : Error NIK Not a Digit Of Number", func(t *testing.T) {
		input.Nik = "asdaasdas"
		validation.On("ValidateRequest", input).Return(nil).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
	})
}

func TestModifyProfilePicture(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var file multipart.File

	var input = dtos.InputUpdateProfilePicture{
		ProfilePicture: file,
	}

	var oldData = dtos.ResUser{
		ID:          1,
		RoleID:      1,
		Email:       "john@example.com",
		Fullname:    "John Doe",
		Gender:      "Male",
		Address:     "planet bumi",
		PhoneNumber: "080000000000",
	}

	var entity = user.User{
		ID: 1,
	}

	var errMap = []string{"photo is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		err, errMap := service.ModifyProfilePicture(input, oldData)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return(oldData.ProfilePicture, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		err, errMap := service.ModifyProfilePicture(input, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("UploadFile", file, oldData.ProfilePicture).Return("", errors.New("error when upload")).Once()

		err, errMap := service.ModifyProfilePicture(input, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		err, errMap := service.ModifyProfilePicture(input, oldData)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		validation.AssertExpectations(t)
	})
}

func TestRemove(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entity = user.User{
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()
		repository.On("DeleteFile", entity.ProfilePicture).Return(nil).Once()
		repository.On("DeleteByID", userID).Return(int64(1)).Once()

		err := service.Remove(userID)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Delete User", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()
		repository.On("DeleteFile", entity.ProfilePicture).Return(nil).Once()
		repository.On("DeleteByID", userID).Return(int64(0)).Once()

		err := service.Remove(userID)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Delete File", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()
		repository.On("DeleteFile", entity.ProfilePicture).Return(errors.New("error when delete")).Once()

		err := service.Remove(userID)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(nil).Once()

		err := service.Remove(userID)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestValidateVerification(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entity = user.User{
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var verificationKey = "123192"

	var email = "john@example.com"

	t.Run("Success", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return(email).Once()
		repository.On("SelectByEmail", email).Return(&entity, nil).Once()
		entity.IsVerified = true
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		verified := service.ValidateVerification(verificationKey)
		assert.True(t, verified)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return(email).Once()
		repository.On("SelectByEmail", email).Return(&entity, nil).Once()
		entity.IsVerified = true
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		verified := service.ValidateVerification(verificationKey)
		assert.False(t, verified)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return(email).Once()
		repository.On("SelectByEmail", email).Return(nil, errors.New("user not found")).Once()

		verified := service.ValidateVerification(verificationKey)
		assert.False(t, verified)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Verification", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return("").Once()

		verified := service.ValidateVerification(verificationKey)
		assert.False(t, verified)
		repository.AssertExpectations(t)
	})
}

func TestForgetPassword(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.ForgetPassword{
		Email: "john@example.com",
	}

	var entity = user.User{
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var otp = "123124"

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()
		generator.On("GenerateRandomOTP").Return(otp).Once()
		repository.On("SendOTPByEmail", entity.Fullname, entity.Email, otp, "2").Return(nil).Once()

		err := service.ForgetPassword(input)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error When Send OTP", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()
		generator.On("GenerateRandomOTP").Return(otp).Once()
		repository.On("SendOTPByEmail", entity.Fullname, entity.Email, otp, "2").Return(errors.New("error when send otp")).Once()

		err := service.ForgetPassword(input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		err := service.ForgetPassword(input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(nil, errors.New("user not found")).Once()

		err := service.ForgetPassword(input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestVerifyOTP(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entity = user.User{
		ID:             1,
		RoleID:         1,
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var userID = "1"
	var roleID = "1"

	var verificationKey = "123124"
	var email = "john@example.com"
	var token = "124233"

	t.Run("Success", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return(email).Once()
		repository.On("SelectByEmail", email).Return(&entity, nil).Once()
		jwt.On("GenerateTokenResetPassword", userID, roleID).Return(token).Once()

		res := service.VerifyOTP(verificationKey)
		assert.Equal(t, res, token)
		repository.AssertExpectations(t)
		jwt.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return(email).Once()
		repository.On("SelectByEmail", email).Return(nil, errors.New("user not found")).Once()

		res := service.VerifyOTP(verificationKey)
		assert.Empty(t, res)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Verification", func(t *testing.T) {
		repository.On("ValidateVerification", verificationKey).Return("").Once()

		res := service.VerifyOTP(verificationKey)
		assert.Empty(t, res)
		repository.AssertExpectations(t)
	})
}

func TestResetPassword(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.ResetPassword{
		Email:    "john@example.com",
		Password: "apo12opkawd",
	}

	var entity = user.User{
		ID:             1,
		RoleID:         1,
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var hashed = "awdapodss"

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()
		hash.On("HashPassword", input.Password).Return(hashed).Once()
		entity.Password = hashed
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		err := service.ResetPassword(input)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(&entity, nil).Once()
		hash.On("HashPassword", input.Password).Return(hashed).Once()
		entity.Password = hashed
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		err := service.ResetPassword(input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("SelectByEmail", input.Email).Return(nil, errors.New("user not found")).Once()

		err := service.ResetPassword(input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestMyProfile(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var entity = user.User{
		ID:             1,
		RoleID:         1,
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()

		res := service.MyProfile(userID)
		assert.Equal(t, res.Email, entity.Email)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(nil).Once()

		res := service.MyProfile(userID)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})
}

func TestCheckPassword(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.CheckPassword{
		OldPassword: "random123",
	}

	var entity = user.User{
		ID:             1,
		RoleID:         1,
		IsVerified:     false,
		Email:          "john@example.com",
		Fullname:       "John Doe",
		Gender:         "Male",
		Address:        "planet bumi",
		Password:       "random123hashed",
		PhoneNumber:    "080000000000",
		ProfilePicture: "google.com",
	}

	var userID = 1

	var errMap = []string{"old password is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByID", userID).Return(&entity).Once()
		hash.On("CompareHash", input.OldPassword, entity.Password).Return(true).Once()

		errMap, err := service.CheckPassword(input, userID)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Password Is Not Match", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByID", userID).Return(&entity).Once()
		hash.On("CompareHash", input.OldPassword, entity.Password).Return(false).Once()

		errMap, err := service.CheckPassword(input, userID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)

		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByID", userID).Return(nil).Once()

		errMap, err := service.CheckPassword(input, userID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		errMap, err := service.CheckPassword(input, userID)
		assert.NotNil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
	})
}

func TestChangePassword(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.ChangePassword{
		NewPassword: "random123",
	}

	var userID = 1
	var hashed = "awpojdapwo1op23k"

	var entity = user.User{
		ID:       1,
		Password: hashed,
	}

	var errMap = []string{"new password is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		hash.On("HashPassword", input.NewPassword).Return(hashed).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		errMap, err := service.ChangePassword(input, userID)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		hash.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		hash.On("HashPassword", input.NewPassword).Return(hashed).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		errMap, err := service.ChangePassword(input, userID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		hash.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Requets", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		errMap, err := service.ChangePassword(input, userID)
		assert.NotNil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
	})
}

func TestAddPersonalization(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var jwt = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, jwt, hash, generator, validation)

	var input = dtos.InputPersonalization{
		Personalization: []string{"koki"},
	}

	var personalization = "koki"

	var entity = user.User{
		ID:              1,
		RoleID:          1,
		IsVerified:      false,
		Email:           "john@example.com",
		Fullname:        "John Doe",
		Gender:          "Male",
		Address:         "planet bumi",
		Password:        "random123hashed",
		PhoneNumber:     "080000000000",
		ProfilePicture:  "google.com",
		Personalization: &personalization,
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()
		repository.On("UpdateUser", entity).Return(int64(1)).Once()

		err := service.AddPersonalization(userID, input)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Update", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(&entity).Once()
		repository.On("UpdateUser", entity).Return(int64(0)).Once()

		err := service.AddPersonalization(userID, input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : User Not Found", func(t *testing.T) {
		repository.On("SelectByID", userID).Return(nil).Once()

		err := service.AddPersonalization(userID, input)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

package usecase

import (
	"errors"
	"raihpeduli/features/auth"
	"raihpeduli/features/auth/dtos"
	"raihpeduli/features/auth/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mashingan/smapping"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	var model = mocks.NewRepository(t)
	var JWT = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validator = helperMocks.NewValidationInterface(t)
	var service = New(model, JWT, hash, generator, validator)

	var newData = dtos.InputUser{
		Fullname:    "Bagus Ario Yudanto",
		Address:     "Kp. Cikempong",
		PhoneNumber: "081234567890",
		Gender:      "Male",
		Email:       "bagus@gmail.com",
		Password:    "bagus123",
	}

	var OTP = "123"

	var errValidation = []string{
		"fullname required",
		"address required",
		"phone_number required",
		"gender required",
		"email required",
		"password required",
	}

	t.Run("Success", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(nil).Once()

		var newUser auth.User
		smapping.FillStruct(&newUser, smapping.MapFields(newData))

		model.On("SelectByEmail", newUser.Email).Return(nil, nil).Once()
		hash.On("HashPassword", newUser.Password).Return("randomhash").Once()
		newUser.Password = "randomhash"
		newUser.ProfilePicture = "https://storage.googleapis.com//users/default"

		model.On("Register", &newUser).Return(&newUser, nil).Once()
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("InsertVerification", newUser.Email, OTP).Return(nil).Once()
		model.On("SendOTPByEmail", newUser.Fullname, newUser.Email, OTP, "1").Return(nil).Once()

		var resUser dtos.ResUser
		smapping.FillStruct(&resUser, smapping.MapFields(newUser))

		result, errMap, err := service.Register(newData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		assert.Equal(t, newUser.Fullname, result.Fullname)
		assert.Equal(t, newUser.Email, result.Email)
		validator.AssertExpectations(t)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Validator error", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(errValidation).Once()

		result, errMap, err := service.Register(newData)
		assert.Nil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, errValidation, errMap)
		validator.AssertExpectations(t)
	})

	t.Run("User already exists", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(nil).Once()
		var newUser auth.User
		smapping.FillStruct(&newUser, smapping.MapFields(newData))
		model.On("SelectByEmail", newUser.Email).Return(&newUser, errors.New("user already exists")).Once()

		result, errMap, err := service.Register(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "User already exists")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		validator.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("Register error", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(nil).Once()

		var newUser auth.User
		smapping.FillStruct(&newUser, smapping.MapFields(newData))

		model.On("SelectByEmail", newUser.Email).Return(nil, nil).Once()
		hash.On("HashPassword", newUser.Password).Return("randomhash").Once()
		newUser.Password = "randomhash"
		newUser.ProfilePicture = "https://storage.googleapis.com//users/default"
		model.On("Register", &newUser).Return(nil, errors.New("failed to register")).Once()

		result, errMap, err := service.Register(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to register")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		validator.AssertExpectations(t)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Insert Verification Failed", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(nil).Once()

		var newUser auth.User
		smapping.FillStruct(&newUser, smapping.MapFields(newData))

		model.On("SelectByEmail", newUser.Email).Return(nil, nil).Once()
		hash.On("HashPassword", newUser.Password).Return("randomhash").Once()
		newUser.Password = "randomhash"
		newUser.ProfilePicture = "https://storage.googleapis.com//users/default"

		model.On("Register", &newUser).Return(&newUser, nil).Once()
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("InsertVerification", newUser.Email, OTP).Return(errors.New("insert verification error")).Once()
		model.On("SendOTPByEmail", newUser.Fullname, newUser.Email, OTP, "1").Return(nil).Once()

		var resUser dtos.ResUser
		smapping.FillStruct(&resUser, smapping.MapFields(newUser))

		result, errMap, err := service.Register(newData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		assert.Equal(t, newUser.Fullname, result.Fullname)
		assert.Equal(t, newUser.Email, result.Email)
		validator.AssertExpectations(t)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})

	t.Run("Send OTP by email failed", func(t *testing.T) {
		validator.On("ValidateRequest", newData).Return(nil).Once()

		var newUser auth.User
		smapping.FillStruct(&newUser, smapping.MapFields(newData))

		model.On("SelectByEmail", newUser.Email).Return(nil, nil).Once()
		hash.On("HashPassword", newUser.Password).Return("randomhash").Once()
		newUser.Password = "randomhash"
		newUser.ProfilePicture = "https://storage.googleapis.com//users/default"

		model.On("Register", &newUser).Return(&newUser, nil).Once()
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("InsertVerification", newUser.Email, OTP).Return(nil).Once()
		model.On("SendOTPByEmail", newUser.Fullname, newUser.Email, OTP, "1").Return(errors.New("send OTP by email failed")).Once()

		var resUser dtos.ResUser
		smapping.FillStruct(&resUser, smapping.MapFields(newUser))

		result, errMap, err := service.Register(newData)
		assert.Error(t, err)
		assert.EqualError(t, err, "send OTP by email failed")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		validator.AssertExpectations(t)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		generator.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	var model = mocks.NewRepository(t)
	var JWT = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validator = helperMocks.NewValidationInterface(t)
	var service = New(model, JWT, hash, generator, validator)

	var loginData = dtos.RequestLogin{
		Email:     "bagus@gmail.com",
		Password:  "bagus123",
		FCMTokens: "fcm_token",
	}

	var userData = auth.User{
		ID:       1,
		RoleID:   1,
		Fullname: "Bagus Ario Yudanto",
		Email:    "bagus@gmail.com",
		Password: "bagus123",
	}

	var token = map[string]any{
		"access_token":  "random_access_token",
		"refresh_token": "random_refresh_token",
	}

	var errValidation = []string{
		"email required",
		"password required",
	}

	t.Run("Success", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(nil).Once()
		model.On("Login", loginData.Email).Return(&userData, nil).Once()
		hash.On("CompareHash", loginData.Password, userData.Password).Return(true).Once()

		userID := strconv.Itoa(userData.ID)
		roleID := strconv.Itoa(userData.RoleID)
		JWT.On("GenerateJWT", userID, roleID).Return(token).Once()
		model.On("InsertToken", userData.ID, loginData.FCMTokens).Return(nil).Once()

		result, errMap, err := service.Login(loginData)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		assert.Equal(t, loginData.Email, result.Email)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		JWT.AssertExpectations(t)
	})

	t.Run("Error validation", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(errValidation).Once()

		result, errMap, err := service.Login(loginData)
		assert.Equal(t, errValidation, errMap)
		assert.Nil(t, err)
		assert.Nil(t, result)
		model.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(nil).Once()
		model.On("Login", loginData.Email).Return(nil, errors.New("user not found")).Once()

		result, errMap, err := service.Login(loginData)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		model.AssertExpectations(t)
	})

	t.Run("Password not match", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(nil).Once()
		model.On("Login", loginData.Email).Return(&userData, nil).Once()
		hash.On("CompareHash", loginData.Password, userData.Password).Return(false).Once()

		result, errMap, err := service.Login(loginData)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid password")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Generate JWT failed", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(nil).Once()
		model.On("Login", loginData.Email).Return(&userData, nil).Once()
		hash.On("CompareHash", loginData.Password, userData.Password).Return(true).Once()

		userID := strconv.Itoa(userData.ID)
		roleID := strconv.Itoa(userData.RoleID)
		JWT.On("GenerateJWT", userID, roleID).Return(nil).Once()

		result, errMap, err := service.Login(loginData)
		assert.Error(t, err)
		assert.EqualError(t, err, "generate token failed")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		JWT.AssertExpectations(t)
	})

	t.Run("Insert Token Failed", func(t *testing.T) {
		validator.On("ValidateRequest", loginData).Return(nil).Once()
		model.On("Login", loginData.Email).Return(&userData, nil).Once()
		hash.On("CompareHash", loginData.Password, userData.Password).Return(true).Once()

		userID := strconv.Itoa(userData.ID)
		roleID := strconv.Itoa(userData.RoleID)
		JWT.On("GenerateJWT", userID, roleID).Return(token).Once()
		model.On("InsertToken", userData.ID, loginData.FCMTokens).Return(errors.New("insert token error")).Once()

		result, errMap, err := service.Login(loginData)
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "insert token error")
		model.AssertExpectations(t)
		hash.AssertExpectations(t)
		JWT.AssertExpectations(t)
	})
}

func TestResendOTP(t *testing.T) {
	var model = mocks.NewRepository(t)
	var JWT = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validator = helperMocks.NewValidationInterface(t)
	var service = New(model, JWT, hash, generator, validator)

	var OTP = "123456"

	var user = auth.User{
		ID:         1,
		RoleID:     1,
		IsVerified: true,
		Email:      "bagus@gmail.com",
		Fullname:   "Bagus Ario Yudanto",
	}

	t.Run("Success", func(t *testing.T) {
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("SelectByEmail", user.Email).Return(&user, nil).Once()
		model.On("InsertVerification", user.Email, OTP).Return(nil).Once()
		model.On("SendOTPByEmail", user.Fullname, user.Email, OTP, "1").Return(nil).Once()

		result := service.ResendOTP(user.Email)
		assert.Equal(t, true, result)
		generator.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("SelectByEmail", user.Email).Return(nil, errors.New("user not found")).Once()

		result := service.ResendOTP(user.Email)
		assert.Equal(t, false, result)
		generator.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("Insert Verification Failed", func(t *testing.T) {
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("SelectByEmail", user.Email).Return(&user, nil).Once()
		model.On("InsertVerification", user.Email, OTP).Return(errors.New("insert verification error")).Once()
		model.On("SendOTPByEmail", user.Fullname, user.Email, OTP, "1").Return(nil).Once()

		result := service.ResendOTP(user.Email)
		assert.Equal(t, true, result)
		generator.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("Send OTP by Email Failed", func(t *testing.T) {
		generator.On("GenerateRandomOTP").Return(OTP).Once()
		model.On("SelectByEmail", user.Email).Return(&user, nil).Once()
		model.On("InsertVerification", user.Email, OTP).Return(nil).Once()
		model.On("SendOTPByEmail", user.Fullname, user.Email, OTP, "1").Return(errors.New("send OTP by email error")).Once()

		result := service.ResendOTP(user.Email)
		assert.Equal(t, false, result)
		generator.AssertExpectations(t)
		model.AssertExpectations(t)
	})
}

func TestRefreshJWT(t *testing.T) {
	var model = mocks.NewRepository(t)
	var JWT = helperMocks.NewJWTInterface(t)
	var hash = helperMocks.NewHashInterface(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var validator = helperMocks.NewValidationInterface(t)
	var service = New(model, JWT, hash, generator, validator)

	var inputToken = dtos.RefreshJWT{
		RefreshToken: "random_refresh_token",
	}

	var parsedToken *jwt.Token

	var token = map[string]any{
		"access_token":  "random_access_token",
		"refresh_token": "random_refresh_token",
	}

	t.Run("Success", func(t *testing.T) {
		JWT.On("ValidateToken", inputToken.RefreshToken, "").Return(parsedToken, nil).Once()
		JWT.On("RefereshJWT", parsedToken).Return(token).Once()

		result, err := service.RefreshJWT(inputToken)
		assert.Nil(t, err)
		assert.Equal(t, token["access_token"], result.AccessToken)
		assert.Equal(t, token["refresh_token"], result.RefreshToken)
		JWT.AssertExpectations(t)
	})

	t.Run("Token not valid", func(t *testing.T) {
		JWT.On("ValidateToken", inputToken.RefreshToken, "").Return(nil, errors.New("token not valid")).Once()

		result, err := service.RefreshJWT(inputToken)
		assert.Error(t, err)
		assert.EqualError(t, err, "validate token failed")
		assert.Nil(t, result)
		JWT.AssertExpectations(t)
	})

	t.Run("Refresh JWT failed", func(t *testing.T) {
		JWT.On("ValidateToken", inputToken.RefreshToken, "").Return(parsedToken, nil).Once()
		JWT.On("RefereshJWT", parsedToken).Return(nil).Once()

		result, err := service.RefreshJWT(inputToken)
		assert.Error(t, err)
		assert.EqualError(t, err, "refresh jwt failed")
		assert.Nil(t, result)
		JWT.AssertExpectations(t)
	})
}

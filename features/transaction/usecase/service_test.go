package usecase

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/transaction"
	"raihpeduli/features/transaction/dtos"
	"raihpeduli/features/transaction/mocks"
	"raihpeduli/features/user"
	helperMocks "raihpeduli/helpers/mocks"
	"strconv"
	"testing"
	"time"

	"errors"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	var transactionData = []transaction.Transaction{
		{
			ID:             1,
			FundraiseID:    1,
			UserID:         1,
			PaymentType:    "6",
			Amount:         100000,
			Status:         "5",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             2,
			FundraiseID:    4,
			UserID:         1,
			PaymentType:    "4",
			Amount:         100000,
			Status:         "2",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             3,
			FundraiseID:    2,
			UserID:         5,
			PaymentType:    "5",
			Amount:         100000,
			Status:         "3",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             4,
			FundraiseID:    6,
			UserID:         8,
			PaymentType:    "7",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             5,
			FundraiseID:    7,
			UserID:         12,
			PaymentType:    "8",
			Amount:         100000,
			Status:         "1",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             6,
			FundraiseID:    4,
			UserID:         3,
			PaymentType:    "10",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             7,
			FundraiseID:    8,
			UserID:         2,
			PaymentType:    "11",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             8,
			FundraiseID:    9,
			UserID:         3,
			PaymentType:    "3",
			Amount:         100000,
			Status:         "3",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
	}

	var userTransactions = []transaction.Transaction{
		{
			ID:             1,
			FundraiseID:    9,
			UserID:         1,
			PaymentType:    "11",
			Amount:         100000,
			Status:         "5",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             8,
			FundraiseID:    2,
			UserID:         1,
			PaymentType:    "11",
			Amount:         100000,
			Status:         "2",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
	}

	t.Run("Success Find All Admin", func(t *testing.T) {
		model.On("Paginate", 1, 10, "").Return(transactionData).Once()
		model.On("GetTotalData", "").Return(int64(8)).Once()

		transactions, totalData := service.FindAll(1, 10, 2, 1, "")
		assert.Equal(t, len(transactionData), len(transactions))
		assert.Equal(t, int64(8), totalData)
		model.AssertExpectations(t)
	})

	t.Run("Success Searching by Name", func(t *testing.T) {
		model.On("Paginate", 1, 10, "Seseorang").Return(userTransactions).Once()
		model.On("GetTotalData", "Seseorang").Return(int64(2)).Once()

		transactions, totalData := service.FindAll(1, 10, 2, 1, "Seseorang")
		assert.Equal(t, len(userTransactions), len(transactions))
		assert.Equal(t, int64(2), totalData)
		model.AssertExpectations(t)
	})

	t.Run("Success Find All User", func(t *testing.T) {
		model.On("PaginateUser", 1, 10, 1).Return(userTransactions).Once()
		model.On("GetTotalDataByUser", 1, "Seseorang").Return(int64(2)).Once()

		transactions, totalData := service.FindAll(1, 10, 1, 1, "Seseorang")
		assert.Equal(t, len(userTransactions), len(transactions))
		assert.Equal(t, int64(2), totalData)
		model.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	var tests = []transaction.Transaction{
		{
			ID:             1,
			FundraiseID:    1,
			UserID:         1,
			PaymentType:    "6",
			Amount:         100000,
			Status:         "5",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             2,
			FundraiseID:    4,
			UserID:         1,
			PaymentType:    "4",
			Amount:         100000,
			Status:         "2",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             3,
			FundraiseID:    2,
			UserID:         5,
			PaymentType:    "5",
			Amount:         100000,
			Status:         "3",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             4,
			FundraiseID:    6,
			UserID:         8,
			PaymentType:    "7",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             5,
			FundraiseID:    7,
			UserID:         12,
			PaymentType:    "8",
			Amount:         100000,
			Status:         "1",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             6,
			FundraiseID:    4,
			UserID:         3,
			PaymentType:    "10",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             7,
			FundraiseID:    8,
			UserID:         2,
			PaymentType:    "11",
			Amount:         100000,
			Status:         "4",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
		{
			ID:             8,
			FundraiseID:    9,
			UserID:         3,
			PaymentType:    "3",
			Amount:         100000,
			Status:         "3",
			VirtualAccount: "678883558856",
			UrlCallback:    "http://test.com",
		},
	}

	for i, test := range tests {
		numberString := strconv.Itoa(i)
		t.Run("Success Transaction ID = "+numberString, func(t *testing.T) {
			model.On("SelectByID", test.ID).Return(&test).Once()

			transaction := service.FindByID(test.ID)
			assert.Equal(t, test.ID, transaction.ID)
			model.AssertExpectations(t)
		})
	}

	t.Run("Not Found", func(t *testing.T) {
		model.On("SelectByID", 1).Return(nil).Once()

		transaction := service.FindByID(1)
		assert.Nil(t, transaction)
		model.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	var user = user.User{
		ID:       1,
		RoleID:   1,
		Email:    "user@example.com",
		Fullname: "Seseorang",
	}

	var errValidation = []string{
		"fundraise_id required",
		"payment_type required",
		"amount required",
	}

	t.Run("Error Validation", func(t *testing.T) {
		var newTransaction dtos.InputTransaction
		validation.On("ValidateRequest", newTransaction).Return(errValidation).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, err)
		assert.Equal(t, errValidation, errMap)
		validation.AssertExpectations(t)
	})

	t.Run("Error Minimum Amount", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      1000,
			PaymentType: "4",
		}
		validation.On("ValidateRequest", newTransaction).Return(nil).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.EqualError(t, err, "Minimum donation ammount is Rp. 10.000")
		validation.AssertExpectations(t)
	})

	t.Run("Get Data Fundraise Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}
		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(0), errors.New("Error get data fundraise")).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.EqualError(t, err, "Error get data fundraise")
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("Fundraise Not Found", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}
		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(0), nil).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.EqualError(t, err, "Fundraise not found")
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
	})

	t.Run("Success Transaction via Bank Transfer", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "4"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(trx.ID)).Once()
		trx.VirtualAccount = "678883558856"
		model.On("Update", trx).Return(int64(1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Equal(t, user.ID, transaction.UserID)
		assert.Equal(t, user.Fullname, transaction.Fullname)
		assert.Equal(t, newTransaction.FundraiseID, transaction.FundraiseID)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, "Bank Transfer", transaction.PaymentType)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
		generator.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Create Transaction Via Bank Transfer Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "4"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("", "", errors.New("create transaction error")).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "create transaction error")
	})

	t.Run("Insert Data Via Bank Transfer Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "4"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Update Data Via Bank Transfer Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "4",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "4"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(1)).Once()
		trx.VirtualAccount = "678883558856"
		model.On("Update", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Success Transaction via Gopay", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "10",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		newTransaction.PaymentType = "10"
		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "10"
		mtRequest.On("CreateTransactionGopay", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", "2023-12-12 22:09:17", nil).Once()
		trx.ValidUntil = "2023-12-12 22:09:17"
		model.On("Insert", trx).Return(int64(trx.ID)).Once()
		trx.UrlCallback = "https://midtrans.com/callback"
		model.On("Update", trx).Return(int64(1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Equal(t, user.ID, transaction.UserID)
		assert.Equal(t, user.Fullname, transaction.Fullname)
		assert.Equal(t, newTransaction.FundraiseID, transaction.FundraiseID)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, "Gopay", transaction.PaymentType)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
		generator.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Create Transaction Via Gopay Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "10",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "10"
		mtRequest.On("CreateTransactionGopay", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("", "", errors.New("create transaction error")).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "create transaction error")
	})

	t.Run("Insert Data Via Gopay Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "10",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "10"
		mtRequest.On("CreateTransactionGopay", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Update Data Via Gopay Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "10",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "10"
		mtRequest.On("CreateTransactionGopay", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(1)).Once()
		trx.UrlCallback = "https://midtrans.com/callback"
		model.On("Update", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Success Transaction via Qris", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "11",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "11"
		mtRequest.On("CreateTransactionQris", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", "2023-12-12 22:09:17", nil).Once()
		trx.ValidUntil = "2023-12-12 22:09:17"
		model.On("Insert", trx).Return(int64(trx.ID)).Once()
		trx.UrlCallback = "https://midtrans.com/callback"
		model.On("Update", trx).Return(int64(1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Equal(t, user.ID, transaction.UserID)
		assert.Equal(t, user.Fullname, transaction.Fullname)
		assert.Equal(t, newTransaction.FundraiseID, transaction.FundraiseID)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, "Qris", transaction.PaymentType)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
		generator.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Create Transaction Via Qris Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "11",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "11"
		mtRequest.On("CreateTransactionQris", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("", "", errors.New("create transaction error")).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "create transaction error")
	})

	t.Run("Insert Data Via Qris Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "11",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "11"
		mtRequest.On("CreateTransactionQris", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Update Data Via Qris Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "11",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "11"
		mtRequest.On("CreateTransactionQris", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("https://midtrans.com/callback", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(1)).Once()
		trx.UrlCallback = "https://midtrans.com/callback"
		model.On("Update", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Success Transaction Default", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "1",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "1"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(trx.ID)).Once()
		model.On("Update", trx).Return(int64(1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Equal(t, user.ID, transaction.UserID)
		assert.Equal(t, user.Fullname, transaction.Fullname)
		assert.Equal(t, newTransaction.FundraiseID, transaction.FundraiseID)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, "Bank Transfer", transaction.PaymentType)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		validation.AssertExpectations(t)
		model.AssertExpectations(t)
		generator.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Create Transaction Default Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "1",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "1"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("", "", errors.New("create transaction error")).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "create transaction error")
	})

	t.Run("Insert Data Default Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "1",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "1"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})

	t.Run("Update Data Default Failed", func(t *testing.T) {
		var newTransaction = dtos.InputTransaction{
			FundraiseID: 1,
			Amount:      100000,
			PaymentType: "1",
		}

		var trx = transaction.Transaction{
			ID:          876565,
			FundraiseID: 1,
			UserID:      1,
			Amount:      100000,
			Status:      "2",
			ValidUntil:  "2023-12-12 22:09:17",
		}

		validation.On("ValidateRequest", newTransaction).Return(nil).Once()
		model.On("CountByID", newTransaction.FundraiseID).Return(int64(1), nil).Once()
		model.On("SelectUserByID", user.ID).Return(&user).Once()
		generator.On("GenerateRandomID").Return(876565).Once()
		trx.PaymentType = "1"
		mtRequest.On("CreateTransactionBank", strconv.Itoa(trx.ID), trx.PaymentType, int64(trx.Amount)).Return("678883558856", trx.ValidUntil, nil).Once()
		model.On("Insert", trx).Return(int64(1)).Once()
		model.On("Update", trx).Return(int64(-1)).Once()

		transaction, err, errMap := service.Create(user.ID, newTransaction)
		assert.Nil(t, transaction)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
	})
}

func TestModify(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	var transactionData = dtos.InputTransaction{
		FundraiseID: 1,
		PaymentType: "10",
		Amount:      100000,
	}

	var newTransaction = transaction.Transaction{
		ID:          876565,
		FundraiseID: 1,
		PaymentType: "10",
		Amount:      100000,
	}

	t.Run("Success", func(t *testing.T) {
		model.On("Update", newTransaction).Return(int64(1)).Once()

		transaction := service.Modify(transactionData, 876565)
		assert.Equal(t, true, transaction)
		model.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		model.On("Update", newTransaction).Return(int64(0)).Once()

		transaction := service.Modify(transactionData, 876565)
		assert.Equal(t, false, transaction)
		model.AssertExpectations(t)
	})
}

func TestRemove(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	t.Run("Success", func(t *testing.T) {
		model.On("DeleteByID", 876565).Return(int64(1)).Once()

		transaction := service.Remove(876565)
		assert.Equal(t, true, transaction)
		model.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		model.On("DeleteByID", 876565).Return(int64(0))

		transaction := service.Remove(876565)
		assert.Equal(t, false, transaction)
		model.AssertExpectations(t)
	})
}

func TestNotifications(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	var notificationsPayload = map[string]any{
		"order_id": "876565",
	}

	var emptyNotificationsPayload = map[string]any{}

	var wrongTrxID = map[string]any{
		"order_id": "wrongID",
	}

	var trx = transaction.Transaction{
		ID:          876565,
		FundraiseID: 1,
		UserID:      1,
		PaymentType: "4",
		Amount:      100000,
		Status:      "5",
		User: user.User{
			Email: "seseorang@gmail.com",
		},
		Fundraise: fundraise.Fundraise{
			Title: "Bantu operasi katarak",
		},
	}

	var trx1 = transaction.Transaction{
		ID:          876565,
		FundraiseID: 1,
		UserID:      1,
		PaymentType: "4",
		Amount:      100000,
		Status:      "2",
		User: user.User{
			Email: "seseorang@gmail.com",
		},
		Fundraise: fundraise.Fundraise{
			Title: "Bantu operasi katarak",
		},
	}

	var paymentName = "Bank Permata"

	t.Run("Success", func(t *testing.T) {
		model.On("SelectByID", trx.ID).Return(&trx).Once()
		mtRequest.On("CheckTransactionStatus", notificationsPayload["order_id"].(string)).Return("5", nil).Once()
		mtRequest.On("MappingPaymentName", trx.PaymentType).Return(paymentName).Once()
		model.On("GetDeviceToken", trx.UserID).Return("deviceToken").Once()
		strAmount := strconv.Itoa(trx.Amount)
		message := "Terimakasih orang baik, donasi sebesar Rp. " + strAmount + "akan sangat membantu " + trx.Fundraise.Title
		nsRequest.On("SendNotifications", "deviceToken", "Pembayaran Berhasil", message).Return(nil).Once()
		model.On("SendPaymentConfirmation", trx.User.Email, trx.Amount, trx.FundraiseID, paymentName).Return(nil).Once()
		trx.PaidAt = time.Now().Format("2006-01-02 15:04:05")
		model.On("Update", trx).Return(int64(1)).Once()

		err := service.Notifications(notificationsPayload)
		assert.Nil(t, err)
		model.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Invalid notification payload", func(t *testing.T) {
		err := service.Notifications(emptyNotificationsPayload)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid notification payload")
	})

	t.Run("Failed Parse TRX ID", func(t *testing.T) {
		err := service.Notifications(wrongTrxID)

		assert.Error(t, err)
		assert.EqualError(t, err, "Failed Parse TRX ID")
	})

	t.Run("Check Transaction Error", func(t *testing.T) {
		model.On("SelectByID", trx.ID).Return(&trx).Once()
		mtRequest.On("CheckTransactionStatus", notificationsPayload["order_id"].(string)).Return("", errors.New("transaction not found")).Once()

		err := service.Notifications(notificationsPayload)
		assert.Error(t, err)
		assert.EqualError(t, err, "transaction not found")
		mtRequest.AssertExpectations(t)
	})

	t.Run("Send Payment Confirmation Error", func(t *testing.T) {
		model.On("SelectByID", trx.ID).Return(&trx).Once()
		mtRequest.On("CheckTransactionStatus", notificationsPayload["order_id"].(string)).Return("5", nil).Once()
		mtRequest.On("MappingPaymentName", trx.PaymentType).Return(paymentName).Once()
		model.On("GetDeviceToken", trx.UserID).Return("deviceToken").Once()
		strAmount := strconv.Itoa(trx.Amount)
		message := "Terimakasih orang baik, donasi sebesar Rp. " + strAmount + "akan sangat membantu " + trx.Fundraise.Title
		nsRequest.On("SendNotifications", "deviceToken", "Pembayaran Berhasil", message).Return(nil).Once()
		model.On("SendPaymentConfirmation", trx.User.Email, trx.Amount, trx.FundraiseID, paymentName).Return(errors.New("send payment confirmation error")).Once()
		trx.PaidAt = time.Now().Format("2006-01-02 15:04:05")
		model.On("Update", trx).Return(int64(1)).Once()

		err := service.Notifications(notificationsPayload)
		assert.Nil(t, err)
		model.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Update Error", func(t *testing.T) {
		model.On("SelectByID", trx.ID).Return(&trx).Once()
		mtRequest.On("CheckTransactionStatus", notificationsPayload["order_id"].(string)).Return("5", nil).Once()
		mtRequest.On("MappingPaymentName", trx.PaymentType).Return(paymentName).Once()
		model.On("GetDeviceToken", trx.UserID).Return("deviceToken").Once()
		strAmount := strconv.Itoa(trx.Amount)
		message := "Terimakasih orang baik, donasi sebesar Rp. " + strAmount + "akan sangat membantu " + trx.Fundraise.Title
		nsRequest.On("SendNotifications", "deviceToken", "Pembayaran Berhasil", message).Return(nil).Once()
		model.On("SendPaymentConfirmation", trx.User.Email, trx.Amount, trx.FundraiseID, paymentName).Return(nil).Once()
		trx.PaidAt = time.Now().Format("2006-01-02 15:04:05")
		model.On("Update", trx).Return(int64(-1)).Once()

		err := service.Notifications(notificationsPayload)
		assert.Nil(t, err)
		model.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})

	t.Run("Update Error(1)", func(t *testing.T) {
		model.On("SelectByID", trx1.ID).Return(&trx1).Once()
		mtRequest.On("CheckTransactionStatus", notificationsPayload["order_id"].(string)).Return("2", nil).Once()
		mtRequest.On("MappingPaymentName", trx1.PaymentType).Return(paymentName).Once()
		model.On("Update", trx1).Return(int64(-1)).Once()

		err := service.Notifications(notificationsPayload)
		assert.Nil(t, err)
		model.AssertExpectations(t)
		mtRequest.AssertExpectations(t)
	})
}

func TestSendPaymentConfirmation(t *testing.T) {
	var model = mocks.NewRepository(t)
	var generator = helperMocks.NewGeneratorInterface(t)
	var mtRequest = helperMocks.NewMidtransInterface(t)
	var coreAPIClient = coreapi.Client{}
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(model, generator, mtRequest, coreAPIClient, validation, nsRequest)

	t.Run("Success", func(t *testing.T) {
		message := "Terimakasih orang baik, donasimu membantu palestina"
		nsRequest.On("SendNotifications", "", "Donasi sebesar Rp. 10.000 Berhasil", message).Return(nil).Once()

		err := service.SendPaymentConfirmation()
		assert.Nil(t, err)
		nsRequest.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		message := "Terimakasih orang baik, donasimu membantu palestina"
		nsRequest.On("SendNotifications", "", "Donasi sebesar Rp. 10.000 Berhasil", message).Return(errors.New("failed send notification")).Once()

		err := service.SendPaymentConfirmation()
		assert.NotNil(t, err)
		assert.EqualError(t, err, "failed send notification")
		nsRequest.AssertExpectations(t)
	})
}

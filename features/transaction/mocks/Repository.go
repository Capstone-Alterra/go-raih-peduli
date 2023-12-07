// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	transaction "raihpeduli/features/transaction"

	mock "github.com/stretchr/testify/mock"

	user "raihpeduli/features/user"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CountByID provides a mock function with given fields: fundraiseID
func (_m *Repository) CountByID(fundraiseID int) (int64, error) {
	ret := _m.Called(fundraiseID)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (int64, error)); ok {
		return rf(fundraiseID)
	}
	if rf, ok := ret.Get(0).(func(int) int64); ok {
		r0 = rf(fundraiseID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(fundraiseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: transactionID
func (_m *Repository) DeleteByID(transactionID int) int64 {
	ret := _m.Called(transactionID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int) int64); ok {
		r0 = rf(transactionID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetTotalData provides a mock function with given fields: keyword
func (_m *Repository) GetTotalData(keyword string) int64 {
	ret := _m.Called(keyword)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(keyword)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetTotalDataByUser provides a mock function with given fields: userID, keyword
func (_m *Repository) GetTotalDataByUser(userID int, keyword string) int64 {
	ret := _m.Called(userID, keyword)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int, string) int64); ok {
		r0 = rf(userID, keyword)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Insert provides a mock function with given fields: newTransaction
func (_m *Repository) Insert(newTransaction transaction.Transaction) int64 {
	ret := _m.Called(newTransaction)

	var r0 int64
	if rf, ok := ret.Get(0).(func(transaction.Transaction) int64); ok {
		r0 = rf(newTransaction)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Paginate provides a mock function with given fields: page, size, keyword
func (_m *Repository) Paginate(page int, size int, keyword string) []transaction.Transaction {
	ret := _m.Called(page, size, keyword)

	var r0 []transaction.Transaction
	if rf, ok := ret.Get(0).(func(int, int, string) []transaction.Transaction); ok {
		r0 = rf(page, size, keyword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]transaction.Transaction)
		}
	}

	return r0
}

// PaginateUser provides a mock function with given fields: page, size, userID
func (_m *Repository) PaginateUser(page int, size int, userID int) []transaction.Transaction {
	ret := _m.Called(page, size, userID)

	var r0 []transaction.Transaction
	if rf, ok := ret.Get(0).(func(int, int, int) []transaction.Transaction); ok {
		r0 = rf(page, size, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]transaction.Transaction)
		}
	}

	return r0
}

// SelectByID provides a mock function with given fields: transactionID
func (_m *Repository) SelectByID(transactionID int) *transaction.Transaction {
	ret := _m.Called(transactionID)

	var r0 *transaction.Transaction
	if rf, ok := ret.Get(0).(func(int) *transaction.Transaction); ok {
		r0 = rf(transactionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Transaction)
		}
	}

	return r0
}

// SelectUserByID provides a mock function with given fields: userID
func (_m *Repository) SelectUserByID(userID int) *user.User {
	ret := _m.Called(userID)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(int) *user.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	return r0
}

// SendPaymentConfirmation provides a mock function with given fields: email, amount, idFundraise, paymentType
func (_m *Repository) SendPaymentConfirmation(email string, amount int, idFundraise int, paymentType string) error {
	ret := _m.Called(email, amount, idFundraise, paymentType)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int, int, string) error); ok {
		r0 = rf(email, amount, idFundraise, paymentType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 transaction.Transaction) int64 {
	ret := _m.Called(_a0)

	var r0 int64
	if rf, ok := ret.Get(0).(func(transaction.Transaction) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

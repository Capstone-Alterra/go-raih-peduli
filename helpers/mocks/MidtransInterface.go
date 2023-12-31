// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	coreapi "github.com/midtrans/midtrans-go/coreapi"

	mock "github.com/stretchr/testify/mock"

	transaction "raihpeduli/features/transaction"
)

// MidtransInterface is an autogenerated mock type for the MidtransInterface type
type MidtransInterface struct {
	mock.Mock
}

// CheckTransactionStatus provides a mock function with given fields: IDTransaction
func (_m *MidtransInterface) CheckTransactionStatus(IDTransaction string) (string, error) {
	ret := _m.Called(IDTransaction)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(IDTransaction)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(IDTransaction)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(IDTransaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTransactionBank provides a mock function with given fields: IDTransaction, PaymentType, Amount
func (_m *MidtransInterface) CreateTransactionBank(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
	ret := _m.Called(IDTransaction, PaymentType, Amount)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int64) (string, string, error)); ok {
		return rf(IDTransaction, PaymentType, Amount)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64) string); ok {
		r0 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, int64) string); ok {
		r1 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string, int64) error); ok {
		r2 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateTransactionGopay provides a mock function with given fields: IDTransaction, PaymentType, Amount
func (_m *MidtransInterface) CreateTransactionGopay(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
	ret := _m.Called(IDTransaction, PaymentType, Amount)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int64) (string, string, error)); ok {
		return rf(IDTransaction, PaymentType, Amount)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64) string); ok {
		r0 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, int64) string); ok {
		r1 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string, int64) error); ok {
		r2 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateTransactionQris provides a mock function with given fields: IDTransaction, PaymentType, Amount
func (_m *MidtransInterface) CreateTransactionQris(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
	ret := _m.Called(IDTransaction, PaymentType, Amount)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, int64) (string, string, error)); ok {
		return rf(IDTransaction, PaymentType, Amount)
	}
	if rf, ok := ret.Get(0).(func(string, string, int64) string); ok {
		r0 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, int64) string); ok {
		r1 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string, int64) error); ok {
		r2 = rf(IDTransaction, PaymentType, Amount)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MappingPaymentName provides a mock function with given fields: paymentType
func (_m *MidtransInterface) MappingPaymentName(paymentType string) string {
	ret := _m.Called(paymentType)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(paymentType)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TransactionStatus provides a mock function with given fields: transactionStatusResp
func (_m *MidtransInterface) TransactionStatus(transactionStatusResp *coreapi.TransactionStatusResponse) transaction.Status {
	ret := _m.Called(transactionStatusResp)

	var r0 transaction.Status
	if rf, ok := ret.Get(0).(func(*coreapi.TransactionStatusResponse) transaction.Status); ok {
		r0 = rf(transactionStatusResp)
	} else {
		r0 = ret.Get(0).(transaction.Status)
	}

	return r0
}

// NewMidtransInterface creates a new instance of MidtransInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMidtransInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MidtransInterface {
	mock := &MidtransInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

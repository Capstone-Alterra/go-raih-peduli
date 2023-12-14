// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	dtos "raihpeduli/features/auth/dtos"

	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Login provides a mock function with given fields: _a0
func (_m *Usecase) Login(_a0 dtos.RequestLogin) (*dtos.LoginResponse, []string, error) {
	ret := _m.Called(_a0)

	var r0 *dtos.LoginResponse
	var r1 []string
	var r2 error
	if rf, ok := ret.Get(0).(func(dtos.RequestLogin) (*dtos.LoginResponse, []string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(dtos.RequestLogin) *dtos.LoginResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.LoginResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(dtos.RequestLogin) []string); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	if rf, ok := ret.Get(2).(func(dtos.RequestLogin) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RefreshJWT provides a mock function with given fields: jwt
func (_m *Usecase) RefreshJWT(jwt dtos.RefreshJWT) (*dtos.ResJWT, error) {
	ret := _m.Called(jwt)

	var r0 *dtos.ResJWT
	var r1 error
	if rf, ok := ret.Get(0).(func(dtos.RefreshJWT) (*dtos.ResJWT, error)); ok {
		return rf(jwt)
	}
	if rf, ok := ret.Get(0).(func(dtos.RefreshJWT) *dtos.ResJWT); ok {
		r0 = rf(jwt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.ResJWT)
		}
	}

	if rf, ok := ret.Get(1).(func(dtos.RefreshJWT) error); ok {
		r1 = rf(jwt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: newUser
func (_m *Usecase) Register(newUser dtos.InputUser) (*dtos.ResUser, []string, error) {
	ret := _m.Called(newUser)

	var r0 *dtos.ResUser
	var r1 []string
	var r2 error
	if rf, ok := ret.Get(0).(func(dtos.InputUser) (*dtos.ResUser, []string, error)); ok {
		return rf(newUser)
	}
	if rf, ok := ret.Get(0).(func(dtos.InputUser) *dtos.ResUser); ok {
		r0 = rf(newUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.ResUser)
		}
	}

	if rf, ok := ret.Get(1).(func(dtos.InputUser) []string); ok {
		r1 = rf(newUser)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	if rf, ok := ret.Get(2).(func(dtos.InputUser) error); ok {
		r2 = rf(newUser)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ResendOTP provides a mock function with given fields: email
func (_m *Usecase) ResendOTP(email string) bool {
	ret := _m.Called(email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewUsecase creates a new instance of Usecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Usecase {
	mock := &Usecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
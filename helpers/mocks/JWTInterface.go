// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt/v5"
	mock "github.com/stretchr/testify/mock"
)

// JWTInterface is an autogenerated mock type for the JWTInterface type
type JWTInterface struct {
	mock.Mock
}

// ExtractToken provides a mock function with given fields: token
func (_m *JWTInterface) ExtractToken(token *jwt.Token) interface{} {
	ret := _m.Called(token)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*jwt.Token) interface{}); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GenerateJWT provides a mock function with given fields: userID, roleID
func (_m *JWTInterface) GenerateJWT(userID string, roleID string) map[string]interface{} {
	ret := _m.Called(userID, roleID)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string, string) map[string]interface{}); ok {
		r0 = rf(userID, roleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// GenerateToken provides a mock function with given fields: userID, roleID
func (_m *JWTInterface) GenerateToken(userID string, roleID string) string {
	ret := _m.Called(userID, roleID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(userID, roleID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GenerateTokenResetPassword provides a mock function with given fields: userID, roleID
func (_m *JWTInterface) GenerateTokenResetPassword(userID string, roleID string) string {
	ret := _m.Called(userID, roleID)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(userID, roleID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ValidateToken provides a mock function with given fields: token, secret
func (_m *JWTInterface) ValidateToken(token string, secret string) (*jwt.Token, error) {
	ret := _m.Called(token, secret)

	var r0 *jwt.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*jwt.Token, error)); ok {
		return rf(token, secret)
	}
	if rf, ok := ret.Get(0).(func(string, string) *jwt.Token); ok {
		r0 = rf(token, secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(token, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJWTInterface creates a new instance of JWTInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJWTInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *JWTInterface {
	mock := &JWTInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ConverterInterface is an autogenerated mock type for the ConverterInterface type
type ConverterInterface struct {
	mock.Mock
}

// Convert provides a mock function with given fields: target, value
func (_m *ConverterInterface) Convert(target interface{}, value interface{}) error {
	ret := _m.Called(target, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) error); ok {
		r0 = rf(target, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewConverterInterface creates a new instance of ConverterInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConverterInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConverterInterface {
	mock := &ConverterInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

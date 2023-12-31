// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// OpenAIInterface is an autogenerated mock type for the OpenAIInterface type
type OpenAIInterface struct {
	mock.Mock
}

// GetAppInformation provides a mock function with given fields: question, qnaList
func (_m *OpenAIInterface) GetAppInformation(question string, qnaList map[string]string) (string, error) {
	ret := _m.Called(question, qnaList)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) (string, error)); ok {
		return rf(question, qnaList)
	}
	if rf, ok := ret.Get(0).(func(string, map[string]string) string); ok {
		r0 = rf(question, qnaList)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(question, qnaList)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNewsContent provides a mock function with given fields: prompt
func (_m *OpenAIInterface) GetNewsContent(prompt string) (string, error) {
	ret := _m.Called(prompt)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(prompt)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(prompt)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(prompt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOpenAIInterface creates a new instance of OpenAIInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOpenAIInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *OpenAIInterface {
	mock := &OpenAIInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

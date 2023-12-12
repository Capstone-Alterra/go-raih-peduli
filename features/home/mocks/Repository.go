// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	fundraise "raihpeduli/features/fundraise"

	mock "github.com/stretchr/testify/mock"

	news "raihpeduli/features/news"

	user "raihpeduli/features/user"

	volunteer "raihpeduli/features/volunteer"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CountFundraise provides a mock function with given fields:
func (_m *Repository) CountFundraise() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// CountNews provides a mock function with given fields:
func (_m *Repository) CountNews() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// CountUser provides a mock function with given fields:
func (_m *Repository) CountUser() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// CountVolunteer provides a mock function with given fields:
func (_m *Repository) CountVolunteer() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// PaginateFundraise provides a mock function with given fields: page, size, personalization
func (_m *Repository) PaginateFundraise(page int, size int, personalization []string) []fundraise.Fundraise {
	ret := _m.Called(page, size, personalization)

	var r0 []fundraise.Fundraise
	if rf, ok := ret.Get(0).(func(int, int, []string) []fundraise.Fundraise); ok {
		r0 = rf(page, size, personalization)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fundraise.Fundraise)
		}
	}

	return r0
}

// PaginateNews provides a mock function with given fields: page, size, personalization
func (_m *Repository) PaginateNews(page int, size int, personalization []string) []news.News {
	ret := _m.Called(page, size, personalization)

	var r0 []news.News
	if rf, ok := ret.Get(0).(func(int, int, []string) []news.News); ok {
		r0 = rf(page, size, personalization)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]news.News)
		}
	}

	return r0
}

// PaginateVolunteer provides a mock function with given fields: page, size, personalization
func (_m *Repository) PaginateVolunteer(page int, size int, personalization []string) []volunteer.VolunteerVacancies {
	ret := _m.Called(page, size, personalization)

	var r0 []volunteer.VolunteerVacancies
	if rf, ok := ret.Get(0).(func(int, int, []string) []volunteer.VolunteerVacancies); ok {
		r0 = rf(page, size, personalization)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]volunteer.VolunteerVacancies)
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
// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	dtos "raihpeduli/features/volunteer/dtos"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// CheckUser provides a mock function with given fields: userID
func (_m *Usecase) CheckUser(userID int) bool {
	ret := _m.Called(userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CreateVacancy provides a mock function with given fields: newVacancy, UserID, file
func (_m *Usecase) CreateVacancy(newVacancy dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error) {
	ret := _m.Called(newVacancy, UserID, file)

	var r0 *dtos.ResVacancy
	var r1 []string
	var r2 error
	if rf, ok := ret.Get(0).(func(dtos.InputVacancy, int, multipart.File) (*dtos.ResVacancy, []string, error)); ok {
		return rf(newVacancy, UserID, file)
	}
	if rf, ok := ret.Get(0).(func(dtos.InputVacancy, int, multipart.File) *dtos.ResVacancy); ok {
		r0 = rf(newVacancy, UserID, file)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.ResVacancy)
		}
	}

	if rf, ok := ret.Get(1).(func(dtos.InputVacancy, int, multipart.File) []string); ok {
		r1 = rf(newVacancy, UserID, file)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	if rf, ok := ret.Get(2).(func(dtos.InputVacancy, int, multipart.File) error); ok {
		r2 = rf(newVacancy, UserID, file)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindAllVacancies provides a mock function with given fields: page, size, searchAndFilter, ownerID, status
func (_m *Usecase) FindAllVacancies(page int, size int, searchAndFilter dtos.SearchAndFilter, ownerID int, status string) ([]dtos.ResVacancy, int64) {
	ret := _m.Called(page, size, searchAndFilter, ownerID, status)

	var r0 []dtos.ResVacancy
	var r1 int64
	if rf, ok := ret.Get(0).(func(int, int, dtos.SearchAndFilter, int, string) ([]dtos.ResVacancy, int64)); ok {
		return rf(page, size, searchAndFilter, ownerID, status)
	}
	if rf, ok := ret.Get(0).(func(int, int, dtos.SearchAndFilter, int, string) []dtos.ResVacancy); ok {
		r0 = rf(page, size, searchAndFilter, ownerID, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dtos.ResVacancy)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, dtos.SearchAndFilter, int, string) int64); ok {
		r1 = rf(page, size, searchAndFilter, ownerID, status)
	} else {
		r1 = ret.Get(1).(int64)
	}

	return r0, r1
}

// FindAllVolunteersByVacancyID provides a mock function with given fields: page, size, vacancyID, name
func (_m *Usecase) FindAllVolunteersByVacancyID(page int, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64) {
	ret := _m.Called(page, size, vacancyID, name)

	var r0 []dtos.ResRegistrantVacancy
	var r1 int64
	if rf, ok := ret.Get(0).(func(int, int, int, string) ([]dtos.ResRegistrantVacancy, int64)); ok {
		return rf(page, size, vacancyID, name)
	}
	if rf, ok := ret.Get(0).(func(int, int, int, string) []dtos.ResRegistrantVacancy); ok {
		r0 = rf(page, size, vacancyID, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dtos.ResRegistrantVacancy)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, int, string) int64); ok {
		r1 = rf(page, size, vacancyID, name)
	} else {
		r1 = ret.Get(1).(int64)
	}

	return r0, r1
}

// FindDetailVolunteers provides a mock function with given fields: vacancyID, volunteerID
func (_m *Usecase) FindDetailVolunteers(vacancyID int, volunteerID int) *dtos.ResRegistrantVacancy {
	ret := _m.Called(vacancyID, volunteerID)

	var r0 *dtos.ResRegistrantVacancy
	if rf, ok := ret.Get(0).(func(int, int) *dtos.ResRegistrantVacancy); ok {
		r0 = rf(vacancyID, volunteerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.ResRegistrantVacancy)
		}
	}

	return r0
}

// FindUserInVacancy provides a mock function with given fields: vacancyID, userID
func (_m *Usecase) FindUserInVacancy(vacancyID int, userID int) bool {
	ret := _m.Called(vacancyID, userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int, int) bool); ok {
		r0 = rf(vacancyID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindVacancyByID provides a mock function with given fields: vacancyID, ownerID
func (_m *Usecase) FindVacancyByID(vacancyID int, ownerID int) *dtos.ResVacancy {
	ret := _m.Called(vacancyID, ownerID)

	var r0 *dtos.ResVacancy
	if rf, ok := ret.Get(0).(func(int, int) *dtos.ResVacancy); ok {
		r0 = rf(vacancyID, ownerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtos.ResVacancy)
		}
	}

	return r0
}

// ModifyVacancy provides a mock function with given fields: vacancyData, file, oldData
func (_m *Usecase) ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) (bool, []string) {
	ret := _m.Called(vacancyData, file, oldData)

	var r0 bool
	var r1 []string
	if rf, ok := ret.Get(0).(func(dtos.InputVacancy, multipart.File, dtos.ResVacancy) (bool, []string)); ok {
		return rf(vacancyData, file, oldData)
	}
	if rf, ok := ret.Get(0).(func(dtos.InputVacancy, multipart.File, dtos.ResVacancy) bool); ok {
		r0 = rf(vacancyData, file, oldData)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(dtos.InputVacancy, multipart.File, dtos.ResVacancy) []string); ok {
		r1 = rf(vacancyData, file, oldData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

// ModifyVacancyStatus provides a mock function with given fields: input, oldData
func (_m *Usecase) ModifyVacancyStatus(input dtos.StatusVacancies, oldData dtos.ResVacancy) (bool, []string) {
	ret := _m.Called(input, oldData)

	var r0 bool
	var r1 []string
	if rf, ok := ret.Get(0).(func(dtos.StatusVacancies, dtos.ResVacancy) (bool, []string)); ok {
		return rf(input, oldData)
	}
	if rf, ok := ret.Get(0).(func(dtos.StatusVacancies, dtos.ResVacancy) bool); ok {
		r0 = rf(input, oldData)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(dtos.StatusVacancies, dtos.ResVacancy) []string); ok {
		r1 = rf(input, oldData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

// RegisterVacancy provides a mock function with given fields: newApply, userID
func (_m *Usecase) RegisterVacancy(newApply dtos.ApplyVacancy, userID int) (bool, []string) {
	ret := _m.Called(newApply, userID)

	var r0 bool
	var r1 []string
	if rf, ok := ret.Get(0).(func(dtos.ApplyVacancy, int) (bool, []string)); ok {
		return rf(newApply, userID)
	}
	if rf, ok := ret.Get(0).(func(dtos.ApplyVacancy, int) bool); ok {
		r0 = rf(newApply, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(dtos.ApplyVacancy, int) []string); ok {
		r1 = rf(newApply, userID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

// RemoveVacancy provides a mock function with given fields: vacancyID, oldData
func (_m *Usecase) RemoveVacancy(vacancyID int, oldData dtos.ResVacancy) error {
	ret := _m.Called(vacancyID, oldData)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, dtos.ResVacancy) error); ok {
		r0 = rf(vacancyID, oldData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStatusRegistrar provides a mock function with given fields: input, registrarID
func (_m *Usecase) UpdateStatusRegistrar(input dtos.StatusRegistrar, registrarID int) (bool, []string) {
	ret := _m.Called(input, registrarID)

	var r0 bool
	var r1 []string
	if rf, ok := ret.Get(0).(func(dtos.StatusRegistrar, int) (bool, []string)); ok {
		return rf(input, registrarID)
	}
	if rf, ok := ret.Get(0).(func(dtos.StatusRegistrar, int) bool); ok {
		r0 = rf(input, registrarID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(dtos.StatusRegistrar, int) []string); ok {
		r1 = rf(input, registrarID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
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

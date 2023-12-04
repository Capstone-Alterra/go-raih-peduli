// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	dtos "raihpeduli/features/news/dtos"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"

	news "raihpeduli/features/news"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// DeleteByID provides a mock function with given fields: newsID
func (_m *Repository) DeleteByID(newsID int) error {
	ret := _m.Called(newsID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(newsID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFile provides a mock function with given fields: filename
func (_m *Repository) DeleteFile(filename string) error {
	ret := _m.Called(filename)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTotalData provides a mock function with given fields:
func (_m *Repository) GetTotalData() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetTotalDataBySearchAndFilter provides a mock function with given fields: searchAndFilter
func (_m *Repository) GetTotalDataBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64 {
	ret := _m.Called(searchAndFilter)

	var r0 int64
	if rf, ok := ret.Get(0).(func(dtos.SearchAndFilter) int64); ok {
		r0 = rf(searchAndFilter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Insert provides a mock function with given fields: newNews
func (_m *Repository) Insert(newNews news.News) (*news.News, error) {
	ret := _m.Called(newNews)

	var r0 *news.News
	var r1 error
	if rf, ok := ret.Get(0).(func(news.News) (*news.News, error)); ok {
		return rf(newNews)
	}
	if rf, ok := ret.Get(0).(func(news.News) *news.News); ok {
		r0 = rf(newNews)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*news.News)
		}
	}

	if rf, ok := ret.Get(1).(func(news.News) error); ok {
		r1 = rf(newNews)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Paginate provides a mock function with given fields: pagination, searchAndFilter
func (_m *Repository) Paginate(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter) ([]news.News, error) {
	ret := _m.Called(pagination, searchAndFilter)

	var r0 []news.News
	var r1 error
	if rf, ok := ret.Get(0).(func(dtos.Pagination, dtos.SearchAndFilter) ([]news.News, error)); ok {
		return rf(pagination, searchAndFilter)
	}
	if rf, ok := ret.Get(0).(func(dtos.Pagination, dtos.SearchAndFilter) []news.News); ok {
		r0 = rf(pagination, searchAndFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]news.News)
		}
	}

	if rf, ok := ret.Get(1).(func(dtos.Pagination, dtos.SearchAndFilter) error); ok {
		r1 = rf(pagination, searchAndFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectBoockmarkByNewsAndOwnerID provides a mock function with given fields: newsID, ownerID
func (_m *Repository) SelectBoockmarkByNewsAndOwnerID(newsID int, ownerID int) (string, error) {
	ret := _m.Called(newsID, ownerID)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) (string, error)); ok {
		return rf(newsID, ownerID)
	}
	if rf, ok := ret.Get(0).(func(int, int) string); ok {
		r0 = rf(newsID, ownerID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(newsID, ownerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectBookmarkedNewsID provides a mock function with given fields: ownerID
func (_m *Repository) SelectBookmarkedNewsID(ownerID int) (map[int]string, error) {
	ret := _m.Called(ownerID)

	var r0 map[int]string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (map[int]string, error)); ok {
		return rf(ownerID)
	}
	if rf, ok := ret.Get(0).(func(int) map[int]string); ok {
		r0 = rf(ownerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[int]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(ownerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectByID provides a mock function with given fields: newsID
func (_m *Repository) SelectByID(newsID int) (*news.News, error) {
	ret := _m.Called(newsID)

	var r0 *news.News
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*news.News, error)); ok {
		return rf(newsID)
	}
	if rf, ok := ret.Get(0).(func(int) *news.News); ok {
		r0 = rf(newsID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*news.News)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(newsID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 news.News) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(news.News) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UploadFile provides a mock function with given fields: file
func (_m *Repository) UploadFile(file multipart.File) (string, error) {
	ret := _m.Called(file)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(multipart.File) (string, error)); ok {
		return rf(file)
	}
	if rf, ok := ret.Get(0).(func(multipart.File) string); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(multipart.File) error); ok {
		r1 = rf(file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

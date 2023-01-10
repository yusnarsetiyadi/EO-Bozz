// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	user "capstone-alta1/features/user"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the RepositoryInterface type
type UserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: input
func (_m *UserRepository) Create(input user.Core) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(user.Core) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *UserRepository) Delete(id uint) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindUser provides a mock function with given fields: email
func (_m *UserRepository) FindUser(email string) (user.Core, error) {
	ret := _m.Called(email)

	var r0 user.Core
	if rf, ok := ret.Get(0).(func(string) user.Core); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(user.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *UserRepository) GetAll() ([]user.Core, error) {
	ret := _m.Called()

	var r0 []user.Core
	if rf, ok := ret.Get(0).(func() []user.Core); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.Core)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *UserRepository) GetById(id uint) (user.Core, error) {
	ret := _m.Called(id)

	var r0 user.Core
	if rf, ok := ret.Get(0).(func(uint) user.Core); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(user.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: input, id
func (_m *UserRepository) Update(input user.Core, id uint) error {
	ret := _m.Called(input, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(user.Core, uint) error); ok {
		r0 = rf(input, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePassword provides a mock function with given fields: input, id
func (_m *UserRepository) UpdatePassword(input user.Core, id uint) error {
	ret := _m.Called(input, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(user.Core, uint) error); ok {
		r0 = rf(input, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRepository(t mockConstructorTestingTNewUserRepository) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

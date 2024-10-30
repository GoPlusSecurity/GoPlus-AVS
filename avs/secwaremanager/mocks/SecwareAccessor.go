// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	secwaremanager "goplus/avs/secwaremanager"

	mock "github.com/stretchr/testify/mock"

	types "goplus/shared/pkg/types"
)

// SecwareAccessor is an autogenerated mock type for the SecwareAccessor type
type SecwareAccessor struct {
	mock.Mock
}

// GetSecwareHealth provides a mock function with given fields: _a0
func (_m *SecwareAccessor) GetSecwareHealth(_a0 *secwaremanager.SecwareStatus) (secwaremanager.SecwareHealth, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetSecwareHealth")
	}

	var r0 secwaremanager.SecwareHealth
	var r1 error
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus) (secwaremanager.SecwareHealth, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus) secwaremanager.SecwareHealth); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(secwaremanager.SecwareHealth)
	}

	if rf, ok := ret.Get(1).(func(*secwaremanager.SecwareStatus) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecwareMeta provides a mock function with given fields: _a0
func (_m *SecwareAccessor) GetSecwareMeta(_a0 *secwaremanager.SecwareStatus) (secwaremanager.SecwareMeta, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetSecwareMeta")
	}

	var r0 secwaremanager.SecwareMeta
	var r1 error
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus) (secwaremanager.SecwareMeta, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus) secwaremanager.SecwareMeta); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(secwaremanager.SecwareMeta)
	}

	if rf, ok := ret.Get(1).(func(*secwaremanager.SecwareStatus) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTask provides a mock function with given fields: _a0, _a1
func (_m *SecwareAccessor) HandleTask(_a0 *secwaremanager.SecwareStatus, _a1 *types.SignedSecwareTask) (types.SignedSecwareResult, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for HandleTask")
	}

	var r0 types.SignedSecwareResult
	var r1 error
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus, *types.SignedSecwareTask) (types.SignedSecwareResult, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(*secwaremanager.SecwareStatus, *types.SignedSecwareTask) types.SignedSecwareResult); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(types.SignedSecwareResult)
	}

	if rf, ok := ret.Get(1).(func(*secwaremanager.SecwareStatus, *types.SignedSecwareTask) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSecwareAccessor creates a new instance of SecwareAccessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSecwareAccessor(t interface {
	mock.TestingT
	Cleanup(func())
}) *SecwareAccessor {
	mock := &SecwareAccessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

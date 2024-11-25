// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	exec "os/exec"

	mock "github.com/stretchr/testify/mock"
)

// CommandExecutor is an autogenerated mock type for the CommandExecutor type
type CommandExecutor struct {
	mock.Mock
}

// ExecCommand provides a mock function with given fields: name, arg
func (_m *CommandExecutor) ExecCommand(name string, arg ...string) *exec.Cmd {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecCommand")
	}

	var r0 *exec.Cmd
	if rf, ok := ret.Get(0).(func(string, ...string) *exec.Cmd); ok {
		r0 = rf(name, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exec.Cmd)
		}
	}

	return r0
}

// NewCommandExecutor creates a new instance of CommandExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommandExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommandExecutor {
	mock := &CommandExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

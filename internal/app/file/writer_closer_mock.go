// Code generated by mockery v2.38.0. DO NOT EDIT.

package file

import mock "github.com/stretchr/testify/mock"

// MockWriteCloser is an autogenerated mock type for the WriteCloser type
type MockWriteCloser struct {
	mock.Mock
}

type MockWriteCloser_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWriteCloser) EXPECT() *MockWriteCloser_Expecter {
	return &MockWriteCloser_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *MockWriteCloser) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWriteCloser_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockWriteCloser_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockWriteCloser_Expecter) Close() *MockWriteCloser_Close_Call {
	return &MockWriteCloser_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockWriteCloser_Close_Call) Run(run func()) *MockWriteCloser_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWriteCloser_Close_Call) Return(_a0 error) *MockWriteCloser_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWriteCloser_Close_Call) RunAndReturn(run func() error) *MockWriteCloser_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: p
func (_m *MockWriteCloser) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWriteCloser_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockWriteCloser_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - p []byte
func (_e *MockWriteCloser_Expecter) Write(p interface{}) *MockWriteCloser_Write_Call {
	return &MockWriteCloser_Write_Call{Call: _e.mock.On("Write", p)}
}

func (_c *MockWriteCloser_Write_Call) Run(run func(p []byte)) *MockWriteCloser_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockWriteCloser_Write_Call) Return(n int, err error) *MockWriteCloser_Write_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockWriteCloser_Write_Call) RunAndReturn(run func([]byte) (int, error)) *MockWriteCloser_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWriteCloser creates a new instance of MockWriteCloser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWriteCloser(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWriteCloser {
	mock := &MockWriteCloser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
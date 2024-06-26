// Code generated by mockery v2.38.0. DO NOT EDIT.

package file

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
)

// MockFileService_DownloadFileServer is an autogenerated mock type for the FileService_DownloadFileServer type
type MockFileService_DownloadFileServer struct {
	mock.Mock
}

type MockFileService_DownloadFileServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFileService_DownloadFileServer) EXPECT() *MockFileService_DownloadFileServer_Expecter {
	return &MockFileService_DownloadFileServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *MockFileService_DownloadFileServer) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// MockFileService_DownloadFileServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type MockFileService_DownloadFileServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *MockFileService_DownloadFileServer_Expecter) Context() *MockFileService_DownloadFileServer_Context_Call {
	return &MockFileService_DownloadFileServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *MockFileService_DownloadFileServer_Context_Call) Run(run func()) *MockFileService_DownloadFileServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_Context_Call) Return(_a0 context.Context) *MockFileService_DownloadFileServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_Context_Call) RunAndReturn(run func() context.Context) *MockFileService_DownloadFileServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *MockFileService_DownloadFileServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileService_DownloadFileServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type MockFileService_DownloadFileServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *MockFileService_DownloadFileServer_Expecter) RecvMsg(m interface{}) *MockFileService_DownloadFileServer_RecvMsg_Call {
	return &MockFileService_DownloadFileServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *MockFileService_DownloadFileServer_RecvMsg_Call) Run(run func(m interface{})) *MockFileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_RecvMsg_Call) Return(_a0 error) *MockFileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *MockFileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *MockFileService_DownloadFileServer) Send(_a0 *v1.DownloadFileResponse) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.DownloadFileResponse) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileService_DownloadFileServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockFileService_DownloadFileServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 *v1.DownloadFileResponse
func (_e *MockFileService_DownloadFileServer_Expecter) Send(_a0 interface{}) *MockFileService_DownloadFileServer_Send_Call {
	return &MockFileService_DownloadFileServer_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *MockFileService_DownloadFileServer_Send_Call) Run(run func(_a0 *v1.DownloadFileResponse)) *MockFileService_DownloadFileServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*v1.DownloadFileResponse))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_Send_Call) Return(_a0 error) *MockFileService_DownloadFileServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_Send_Call) RunAndReturn(run func(*v1.DownloadFileResponse) error) *MockFileService_DownloadFileServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *MockFileService_DownloadFileServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileService_DownloadFileServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type MockFileService_DownloadFileServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockFileService_DownloadFileServer_Expecter) SendHeader(_a0 interface{}) *MockFileService_DownloadFileServer_SendHeader_Call {
	return &MockFileService_DownloadFileServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *MockFileService_DownloadFileServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *MockFileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_SendHeader_Call) Return(_a0 error) *MockFileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *MockFileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *MockFileService_DownloadFileServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileService_DownloadFileServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type MockFileService_DownloadFileServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *MockFileService_DownloadFileServer_Expecter) SendMsg(m interface{}) *MockFileService_DownloadFileServer_SendMsg_Call {
	return &MockFileService_DownloadFileServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *MockFileService_DownloadFileServer_SendMsg_Call) Run(run func(m interface{})) *MockFileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_SendMsg_Call) Return(_a0 error) *MockFileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *MockFileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *MockFileService_DownloadFileServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileService_DownloadFileServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type MockFileService_DownloadFileServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockFileService_DownloadFileServer_Expecter) SetHeader(_a0 interface{}) *MockFileService_DownloadFileServer_SetHeader_Call {
	return &MockFileService_DownloadFileServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *MockFileService_DownloadFileServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *MockFileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_SetHeader_Call) Return(_a0 error) *MockFileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileService_DownloadFileServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *MockFileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *MockFileService_DownloadFileServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// MockFileService_DownloadFileServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type MockFileService_DownloadFileServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *MockFileService_DownloadFileServer_Expecter) SetTrailer(_a0 interface{}) *MockFileService_DownloadFileServer_SetTrailer_Call {
	return &MockFileService_DownloadFileServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *MockFileService_DownloadFileServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *MockFileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *MockFileService_DownloadFileServer_SetTrailer_Call) Return() *MockFileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockFileService_DownloadFileServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *MockFileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFileService_DownloadFileServer creates a new instance of MockFileService_DownloadFileServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileService_DownloadFileServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileService_DownloadFileServer {
	mock := &MockFileService_DownloadFileServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

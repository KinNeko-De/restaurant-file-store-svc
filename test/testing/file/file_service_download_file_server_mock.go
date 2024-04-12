// Code generated by mockery v2.38.0. DO NOT EDIT.

package file

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
)

// FileService_DownloadFileServer is an autogenerated mock type for the FileService_DownloadFileServer type
type FileService_DownloadFileServer struct {
	mock.Mock
}

type FileService_DownloadFileServer_Expecter struct {
	mock *mock.Mock
}

func (_m *FileService_DownloadFileServer) EXPECT() *FileService_DownloadFileServer_Expecter {
	return &FileService_DownloadFileServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *FileService_DownloadFileServer) Context() context.Context {
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

// FileService_DownloadFileServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type FileService_DownloadFileServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *FileService_DownloadFileServer_Expecter) Context() *FileService_DownloadFileServer_Context_Call {
	return &FileService_DownloadFileServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *FileService_DownloadFileServer_Context_Call) Run(run func()) *FileService_DownloadFileServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *FileService_DownloadFileServer_Context_Call) Return(_a0 context.Context) *FileService_DownloadFileServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_Context_Call) RunAndReturn(run func() context.Context) *FileService_DownloadFileServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *FileService_DownloadFileServer) RecvMsg(m interface{}) error {
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

// FileService_DownloadFileServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type FileService_DownloadFileServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *FileService_DownloadFileServer_Expecter) RecvMsg(m interface{}) *FileService_DownloadFileServer_RecvMsg_Call {
	return &FileService_DownloadFileServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *FileService_DownloadFileServer_RecvMsg_Call) Run(run func(m interface{})) *FileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_RecvMsg_Call) Return(_a0 error) *FileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *FileService_DownloadFileServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *FileService_DownloadFileServer) Send(_a0 *v1.DownloadFileResponse) error {
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

// FileService_DownloadFileServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type FileService_DownloadFileServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 *v1.DownloadFileResponse
func (_e *FileService_DownloadFileServer_Expecter) Send(_a0 interface{}) *FileService_DownloadFileServer_Send_Call {
	return &FileService_DownloadFileServer_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *FileService_DownloadFileServer_Send_Call) Run(run func(_a0 *v1.DownloadFileResponse)) *FileService_DownloadFileServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*v1.DownloadFileResponse))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_Send_Call) Return(_a0 error) *FileService_DownloadFileServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_Send_Call) RunAndReturn(run func(*v1.DownloadFileResponse) error) *FileService_DownloadFileServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *FileService_DownloadFileServer) SendHeader(_a0 metadata.MD) error {
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

// FileService_DownloadFileServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type FileService_DownloadFileServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_DownloadFileServer_Expecter) SendHeader(_a0 interface{}) *FileService_DownloadFileServer_SendHeader_Call {
	return &FileService_DownloadFileServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *FileService_DownloadFileServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *FileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_SendHeader_Call) Return(_a0 error) *FileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *FileService_DownloadFileServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *FileService_DownloadFileServer) SendMsg(m interface{}) error {
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

// FileService_DownloadFileServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type FileService_DownloadFileServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *FileService_DownloadFileServer_Expecter) SendMsg(m interface{}) *FileService_DownloadFileServer_SendMsg_Call {
	return &FileService_DownloadFileServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *FileService_DownloadFileServer_SendMsg_Call) Run(run func(m interface{})) *FileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_SendMsg_Call) Return(_a0 error) *FileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *FileService_DownloadFileServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *FileService_DownloadFileServer) SetHeader(_a0 metadata.MD) error {
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

// FileService_DownloadFileServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type FileService_DownloadFileServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_DownloadFileServer_Expecter) SetHeader(_a0 interface{}) *FileService_DownloadFileServer_SetHeader_Call {
	return &FileService_DownloadFileServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *FileService_DownloadFileServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *FileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_SetHeader_Call) Return(_a0 error) *FileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_DownloadFileServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *FileService_DownloadFileServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *FileService_DownloadFileServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// FileService_DownloadFileServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type FileService_DownloadFileServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_DownloadFileServer_Expecter) SetTrailer(_a0 interface{}) *FileService_DownloadFileServer_SetTrailer_Call {
	return &FileService_DownloadFileServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *FileService_DownloadFileServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *FileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_DownloadFileServer_SetTrailer_Call) Return() *FileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *FileService_DownloadFileServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *FileService_DownloadFileServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewFileService_DownloadFileServer creates a new instance of FileService_DownloadFileServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileService_DownloadFileServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileService_DownloadFileServer {
	mock := &FileService_DownloadFileServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
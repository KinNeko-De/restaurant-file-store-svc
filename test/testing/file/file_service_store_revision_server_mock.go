// Code generated by mockery v2.38.0. DO NOT EDIT.

package file

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
)

// FileService_StoreRevisionServer is an autogenerated mock type for the FileService_StoreRevisionServer type
type FileService_StoreRevisionServer struct {
	mock.Mock
}

type FileService_StoreRevisionServer_Expecter struct {
	mock *mock.Mock
}

func (_m *FileService_StoreRevisionServer) EXPECT() *FileService_StoreRevisionServer_Expecter {
	return &FileService_StoreRevisionServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *FileService_StoreRevisionServer) Context() context.Context {
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

// FileService_StoreRevisionServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type FileService_StoreRevisionServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *FileService_StoreRevisionServer_Expecter) Context() *FileService_StoreRevisionServer_Context_Call {
	return &FileService_StoreRevisionServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *FileService_StoreRevisionServer_Context_Call) Run(run func()) *FileService_StoreRevisionServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_Context_Call) Return(_a0 context.Context) *FileService_StoreRevisionServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_Context_Call) RunAndReturn(run func() context.Context) *FileService_StoreRevisionServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// Recv provides a mock function with given fields:
func (_m *FileService_StoreRevisionServer) Recv() (*v1.StoreRevisionRequest, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Recv")
	}

	var r0 *v1.StoreRevisionRequest
	var r1 error
	if rf, ok := ret.Get(0).(func() (*v1.StoreRevisionRequest, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *v1.StoreRevisionRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.StoreRevisionRequest)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FileService_StoreRevisionServer_Recv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Recv'
type FileService_StoreRevisionServer_Recv_Call struct {
	*mock.Call
}

// Recv is a helper method to define mock.On call
func (_e *FileService_StoreRevisionServer_Expecter) Recv() *FileService_StoreRevisionServer_Recv_Call {
	return &FileService_StoreRevisionServer_Recv_Call{Call: _e.mock.On("Recv")}
}

func (_c *FileService_StoreRevisionServer_Recv_Call) Run(run func()) *FileService_StoreRevisionServer_Recv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_Recv_Call) Return(_a0 *v1.StoreRevisionRequest, _a1 error) *FileService_StoreRevisionServer_Recv_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FileService_StoreRevisionServer_Recv_Call) RunAndReturn(run func() (*v1.StoreRevisionRequest, error)) *FileService_StoreRevisionServer_Recv_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *FileService_StoreRevisionServer) RecvMsg(m interface{}) error {
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

// FileService_StoreRevisionServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type FileService_StoreRevisionServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *FileService_StoreRevisionServer_Expecter) RecvMsg(m interface{}) *FileService_StoreRevisionServer_RecvMsg_Call {
	return &FileService_StoreRevisionServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *FileService_StoreRevisionServer_RecvMsg_Call) Run(run func(m interface{})) *FileService_StoreRevisionServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_RecvMsg_Call) Return(_a0 error) *FileService_StoreRevisionServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *FileService_StoreRevisionServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SendAndClose provides a mock function with given fields: _a0
func (_m *FileService_StoreRevisionServer) SendAndClose(_a0 *v1.StoreFileResponse) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendAndClose")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.StoreFileResponse) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FileService_StoreRevisionServer_SendAndClose_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendAndClose'
type FileService_StoreRevisionServer_SendAndClose_Call struct {
	*mock.Call
}

// SendAndClose is a helper method to define mock.On call
//   - _a0 *v1.StoreFileResponse
func (_e *FileService_StoreRevisionServer_Expecter) SendAndClose(_a0 interface{}) *FileService_StoreRevisionServer_SendAndClose_Call {
	return &FileService_StoreRevisionServer_SendAndClose_Call{Call: _e.mock.On("SendAndClose", _a0)}
}

func (_c *FileService_StoreRevisionServer_SendAndClose_Call) Run(run func(_a0 *v1.StoreFileResponse)) *FileService_StoreRevisionServer_SendAndClose_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*v1.StoreFileResponse))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_SendAndClose_Call) Return(_a0 error) *FileService_StoreRevisionServer_SendAndClose_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_SendAndClose_Call) RunAndReturn(run func(*v1.StoreFileResponse) error) *FileService_StoreRevisionServer_SendAndClose_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *FileService_StoreRevisionServer) SendHeader(_a0 metadata.MD) error {
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

// FileService_StoreRevisionServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type FileService_StoreRevisionServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_StoreRevisionServer_Expecter) SendHeader(_a0 interface{}) *FileService_StoreRevisionServer_SendHeader_Call {
	return &FileService_StoreRevisionServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *FileService_StoreRevisionServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *FileService_StoreRevisionServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_SendHeader_Call) Return(_a0 error) *FileService_StoreRevisionServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *FileService_StoreRevisionServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *FileService_StoreRevisionServer) SendMsg(m interface{}) error {
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

// FileService_StoreRevisionServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type FileService_StoreRevisionServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *FileService_StoreRevisionServer_Expecter) SendMsg(m interface{}) *FileService_StoreRevisionServer_SendMsg_Call {
	return &FileService_StoreRevisionServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *FileService_StoreRevisionServer_SendMsg_Call) Run(run func(m interface{})) *FileService_StoreRevisionServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_SendMsg_Call) Return(_a0 error) *FileService_StoreRevisionServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *FileService_StoreRevisionServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *FileService_StoreRevisionServer) SetHeader(_a0 metadata.MD) error {
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

// FileService_StoreRevisionServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type FileService_StoreRevisionServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_StoreRevisionServer_Expecter) SetHeader(_a0 interface{}) *FileService_StoreRevisionServer_SetHeader_Call {
	return &FileService_StoreRevisionServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *FileService_StoreRevisionServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *FileService_StoreRevisionServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_SetHeader_Call) Return(_a0 error) *FileService_StoreRevisionServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FileService_StoreRevisionServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *FileService_StoreRevisionServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *FileService_StoreRevisionServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// FileService_StoreRevisionServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type FileService_StoreRevisionServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *FileService_StoreRevisionServer_Expecter) SetTrailer(_a0 interface{}) *FileService_StoreRevisionServer_SetTrailer_Call {
	return &FileService_StoreRevisionServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *FileService_StoreRevisionServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *FileService_StoreRevisionServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *FileService_StoreRevisionServer_SetTrailer_Call) Return() *FileService_StoreRevisionServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *FileService_StoreRevisionServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *FileService_StoreRevisionServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewFileService_StoreRevisionServer creates a new instance of FileService_StoreRevisionServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileService_StoreRevisionServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileService_StoreRevisionServer {
	mock := &FileService_StoreRevisionServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

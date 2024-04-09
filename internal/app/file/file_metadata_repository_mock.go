// Code generated by mockery v2.38.0. DO NOT EDIT.

package file

import (
	context "context"

	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// MockFileMetadataRepository is an autogenerated mock type for the FileMetadataRepository type
type MockFileMetadataRepository struct {
	mock.Mock
}

type MockFileMetadataRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFileMetadataRepository) EXPECT() *MockFileMetadataRepository_Expecter {
	return &MockFileMetadataRepository_Expecter{mock: &_m.Mock}
}

// FetchFileMetadata provides a mock function with given fields: ctx, fileId
func (_m *MockFileMetadataRepository) FetchFileMetadata(ctx context.Context, fileId uuid.UUID) (FileMetadata, error) {
	ret := _m.Called(ctx, fileId)

	if len(ret) == 0 {
		panic("no return value specified for FetchFileMetadata")
	}

	var r0 FileMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (FileMetadata, error)); ok {
		return rf(ctx, fileId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) FileMetadata); ok {
		r0 = rf(ctx, fileId)
	} else {
		r0 = ret.Get(0).(FileMetadata)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, fileId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFileMetadataRepository_FetchFileMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchFileMetadata'
type MockFileMetadataRepository_FetchFileMetadata_Call struct {
	*mock.Call
}

// FetchFileMetadata is a helper method to define mock.On call
//   - ctx context.Context
//   - fileId uuid.UUID
func (_e *MockFileMetadataRepository_Expecter) FetchFileMetadata(ctx interface{}, fileId interface{}) *MockFileMetadataRepository_FetchFileMetadata_Call {
	return &MockFileMetadataRepository_FetchFileMetadata_Call{Call: _e.mock.On("FetchFileMetadata", ctx, fileId)}
}

func (_c *MockFileMetadataRepository_FetchFileMetadata_Call) Run(run func(ctx context.Context, fileId uuid.UUID)) *MockFileMetadataRepository_FetchFileMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockFileMetadataRepository_FetchFileMetadata_Call) Return(_a0 FileMetadata, _a1 error) *MockFileMetadataRepository_FetchFileMetadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileMetadataRepository_FetchFileMetadata_Call) RunAndReturn(run func(context.Context, uuid.UUID) (FileMetadata, error)) *MockFileMetadataRepository_FetchFileMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// StoreFileMetadata provides a mock function with given fields: ctx, fileMetadata
func (_m *MockFileMetadataRepository) StoreFileMetadata(ctx context.Context, fileMetadata FileMetadata) error {
	ret := _m.Called(ctx, fileMetadata)

	if len(ret) == 0 {
		panic("no return value specified for StoreFileMetadata")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, FileMetadata) error); ok {
		r0 = rf(ctx, fileMetadata)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileMetadataRepository_StoreFileMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreFileMetadata'
type MockFileMetadataRepository_StoreFileMetadata_Call struct {
	*mock.Call
}

// StoreFileMetadata is a helper method to define mock.On call
//   - ctx context.Context
//   - fileMetadata FileMetadata
func (_e *MockFileMetadataRepository_Expecter) StoreFileMetadata(ctx interface{}, fileMetadata interface{}) *MockFileMetadataRepository_StoreFileMetadata_Call {
	return &MockFileMetadataRepository_StoreFileMetadata_Call{Call: _e.mock.On("StoreFileMetadata", ctx, fileMetadata)}
}

func (_c *MockFileMetadataRepository_StoreFileMetadata_Call) Run(run func(ctx context.Context, fileMetadata FileMetadata)) *MockFileMetadataRepository_StoreFileMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(FileMetadata))
	})
	return _c
}

func (_c *MockFileMetadataRepository_StoreFileMetadata_Call) Return(_a0 error) *MockFileMetadataRepository_StoreFileMetadata_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFileMetadataRepository_StoreFileMetadata_Call) RunAndReturn(run func(context.Context, FileMetadata) error) *MockFileMetadataRepository_StoreFileMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFileMetadataRepository creates a new instance of MockFileMetadataRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileMetadataRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileMetadataRepository {
	mock := &MockFileMetadataRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

package file

import (
	"context"
	"testing"

	"github.com/google/uuid"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func CreateDownloadFileStream(t *testing.T) *MockFileService_DownloadFileServer {
	mockStream := NewMockFileService_DownloadFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateDownloadFileRequestFromUuid(t *testing.T, fileId uuid.UUID) *apiRestaurantFile.DownloadFileRequest {
	requestedFileId, err := apiProtobuf.ToProtobuf(fileId)
	require.Nil(t, err)
	return CreateDownloadFileRequest(t, requestedFileId)
}

func CreateDownloadFileRequest(t *testing.T, requestedFileId *apiProtobuf.Uuid) *apiRestaurantFile.DownloadFileRequest {
	request := &apiRestaurantFile.DownloadFileRequest{
		FileId: requestedFileId,
	}
	return request
}

func (mockStream *MockFileService_DownloadFileServer) SetupRecordDownloadedFile(t *testing.T) func() []byte {
	actualFile := make([]byte, 0)
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualFile = append(actualFile, response.GetChunk()...)
	}).Return(nil)
	return func() []byte {
		return actualFile
	}
}

func (mockStream *MockFileService_DownloadFileServer) SetupRecordStoredFileMetadata(t *testing.T) func() *apiRestaurantFile.StoredFile {
	actualStoredFileMetadata := &apiRestaurantFile.StoredFile{}
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualStoredFileMetadata = response.GetStoredFile()
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoredFile {
		return actualStoredFileMetadata
	}
}

func (mockStream *MockFileService_DownloadFileServer) SetupSendError(t *testing.T, err error) {
	mockStream.EXPECT().Send(mock.Anything).Return(err).Times(1)
}

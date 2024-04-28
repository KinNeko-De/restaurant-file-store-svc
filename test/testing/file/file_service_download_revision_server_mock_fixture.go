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

func CreateDownloadRevisionStream(t *testing.T) *MockFileService_DownloadRevisionServer {
	mockStream := NewMockFileService_DownloadRevisionServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateDownloadRevisionRequestFromUuid(t *testing.T, fileId uuid.UUID, revisionId uuid.UUID) *apiRestaurantFile.DownloadRevisionRequest {
	requestedFileId, err := apiProtobuf.ToProtobuf(fileId)
	require.Nil(t, err)
	requestedRevisionId, err := apiProtobuf.ToProtobuf(revisionId)
	require.Nil(t, err)
	return CreateDownloadRevisionRequest(t, requestedFileId, requestedRevisionId)
}

func CreateDownloadRevisionRequest(t *testing.T, requestedFileId *apiProtobuf.Uuid, requestedRevisionId *apiProtobuf.Uuid) *apiRestaurantFile.DownloadRevisionRequest {
	request := &apiRestaurantFile.DownloadRevisionRequest{
		FileId:     requestedFileId,
		RevisionId: requestedRevisionId,
	}
	return request
}

func (mockStream *MockFileService_DownloadRevisionServer) SetupRecordDownloadedFile(t *testing.T) func() []byte {
	actualFile := make([]byte, 0)
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualFile = append(actualFile, response.GetChunk()...)
	}).Return(nil)
	return func() []byte {
		return actualFile
	}
}

func (mockStream *MockFileService_DownloadRevisionServer) SetupRecordStoredFileMetadata(t *testing.T) func() *apiRestaurantFile.StoredFile {
	actualStoredFile := &apiRestaurantFile.StoredFile{}
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualStoredFile = response.GetStoredFile()
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoredFile {
		return actualStoredFile
	}
}

func (mockStream *MockFileService_DownloadRevisionServer) SetupSendError(t *testing.T, err error) {
	mockStream.EXPECT().Send(mock.Anything).Return(err).Times(1)
}

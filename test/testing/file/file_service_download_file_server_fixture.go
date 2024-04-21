package file

import (
	"context"
	"testing"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateDownloadFileStream(t *testing.T) *FileService_DownloadFileServer {
	mockStream := NewFileService_DownloadFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func SetupRecordDownloadedFile(t *testing.T, mockStream *FileService_DownloadFileServer) func() []byte {
	actualFile := make([]byte, 0)
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualFile = append(actualFile, response.GetChunk()...)
	}).Return(nil)
	return func() []byte {
		return actualFile
	}
}

func SetupRecordStoredFileMetadata(t *testing.T, mockStream *FileService_DownloadFileServer) func() *apiRestaurantFile.StoredFileMetadata {
	actualStoredFileMetadata := &apiRestaurantFile.StoredFileMetadata{}
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *apiRestaurantFile.DownloadFileResponse) {
		actualStoredFileMetadata = response.GetMetadata()
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoredFileMetadata {
		return actualStoredFileMetadata
	}
}

package file

import (
	context "context"
	"io"
	"testing"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateStoreFileStream(t *testing.T) *MockFileService_StoreFileServer {
	mockStream := NewMockFileService_StoreFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateMetadataStoreFileRequestFromFileName(t *testing.T, fileName string) *apiRestaurantFile.StoreFileRequest {
	sentFileName := &apiRestaurantFile.StoreFile{
		Name: fileName,
	}

	return CreateMetadataStoreFileRequest(t, sentFileName)
}

// TODO write test that uses this function
func CreateMetadataStoreFileRequest(t *testing.T, storeFile *apiRestaurantFile.StoreFile) *apiRestaurantFile.StoreFileRequest {
	metadata := &apiRestaurantFile.StoreFileRequest{
		Part: &apiRestaurantFile.StoreFileRequest_StoreFile{
			StoreFile: storeFile,
		},
	}
	return metadata
}

func CreateChunkStoreFileRequest(t *testing.T, chunk []byte) *apiRestaurantFile.StoreFileRequest {
	chunkRequest := &apiRestaurantFile.StoreFileRequest{
		Part: &apiRestaurantFile.StoreFileRequest_Chunk{
			Chunk: chunk,
		},
	}
	return chunkRequest
}

func (mockStream *MockFileService_StoreFileServer) SetupSuccessfulSend(t *testing.T, sentFileName string, fileChunks [][]byte) {
	metadata := CreateMetadataStoreFileRequestFromFileName(t, sentFileName)
	mockStream.SetupSendMetadata(t, metadata)
	mockStream.SetupSendFile(t, fileChunks)
	mockStream.SetupSendEndOfFile(t)
}

func (mockStream *MockFileService_StoreFileServer) SetupSendMetadata(t *testing.T, metadata *apiRestaurantFile.StoreFileRequest) {
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)
}

func (mockStream *MockFileService_StoreFileServer) SetupSendFile(t *testing.T, fileChunks [][]byte) {
	for _, chunk := range fileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}
}

func (mockStream *MockFileService_StoreFileServer) SetupSendError(t *testing.T, sendError error) {
	mockStream.EXPECT().Recv().Return(nil, sendError).Times(1)
}

func (mockStream *MockFileService_StoreFileServer) SetupSendEndOfFile(t *testing.T) {
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
}

func (mockStream *MockFileService_StoreFileServer) SetupSendAndClose(t *testing.T) func() *apiRestaurantFile.StoreFileResponse {
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *apiRestaurantFile.StoreFileResponse) {
		actualResponse = response
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoreFileResponse {
		return actualResponse
	}
}

func (mockStream *MockFileService_StoreFileServer) SetupSendAndCloseError(t *testing.T, closeError error) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Return(closeError).Times(1)
}

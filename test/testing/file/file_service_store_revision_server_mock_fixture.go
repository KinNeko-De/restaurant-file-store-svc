package file

import (
	context "context"
	"io"
	"testing"

	"github.com/google/uuid"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateStoreRevisionStream(t *testing.T) *MockFileService_StoreRevisionServer {
	mockStream := NewMockFileService_StoreRevisionServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateMetadataStoreRevisionRequestFromFileName(t *testing.T, fileId uuid.UUID, fileName string) *apiRestaurantFile.StoreRevisionRequest {
	storeRevision := &apiRestaurantFile.StoreRevision{
		FileId: &apiProtobuf.Uuid{
			Value: fileId.String(),
		},
		StoreFile: &apiRestaurantFile.StoreFile{
			Name: fileName,
		},
	}

	return CreateMetadataRevisionFileRequest(t, storeRevision)
}

// TODO write test that uses this function
func CreateMetadataRevisionFileRequest(t *testing.T, storeRevision *apiRestaurantFile.StoreRevision) *apiRestaurantFile.StoreRevisionRequest {
	request := &apiRestaurantFile.StoreRevisionRequest{
		Part: &apiRestaurantFile.StoreRevisionRequest_StoreRevision{
			StoreRevision: storeRevision,
		},
	}
	return request
}

func CreateMetadataStoreRevisionRequest(t *testing.T, fileId uuid.UUID, fileName string) *apiRestaurantFile.StoreRevisionRequest {
	metadata := &apiRestaurantFile.StoreRevisionRequest{
		Part: &apiRestaurantFile.StoreRevisionRequest_StoreRevision{
			StoreRevision: &apiRestaurantFile.StoreRevision{
				FileId: &apiProtobuf.Uuid{
					Value: fileId.String(),
				},
				StoreFile: &apiRestaurantFile.StoreFile{
					Name: fileName,
				},
			},
		},
	}
	return metadata
}

func CreateChunkStoreRevisionRequest(t *testing.T, chunk []byte) *apiRestaurantFile.StoreRevisionRequest {
	chunkRequest := &apiRestaurantFile.StoreRevisionRequest{
		Part: &apiRestaurantFile.StoreRevisionRequest_Chunk{
			Chunk: chunk,
		},
	}
	return chunkRequest
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSuccessfulSend(t *testing.T, existingFileId uuid.UUID, sentFileName string, fileChunks [][]byte) {
	metadata := CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName)
	mockStream.SetupSendMetadata(t, metadata)
	mockStream.SetupSendFile(t, fileChunks)
	mockStream.SetupSendEndOfFile(t)
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendMetadata(t *testing.T, metadata *apiRestaurantFile.StoreRevisionRequest) {
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendFile(t *testing.T, fileChunks [][]byte) {
	for _, chunk := range fileChunks {
		chunkRequest := CreateChunkStoreRevisionRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendError(t *testing.T, sendError error) {
	mockStream.EXPECT().Recv().Return(nil, sendError).Times(1)
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendEndOfFile(t *testing.T) {
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendAndClose(t *testing.T) func() *apiRestaurantFile.StoreFileResponse {
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *apiRestaurantFile.StoreFileResponse) {
		actualResponse = response
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoreFileResponse {
		return actualResponse
	}
}

func (mockStream *MockFileService_StoreRevisionServer) SetupSendAndCloseError(t *testing.T, sendError error) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Return(sendError).Times(1)
}

package file

import (
	"io"
	"testing"

	"github.com/google/uuid"
)

func CreateValidStoreFileStreamThatAbortsOnFileWrite(t *testing.T, fileName string, successfulWritenfileChunks [][]byte) *MockFileService_StoreFileServer {
	mockStream := CreateStoreFileStream(t)

	metadata := CreateMetadataStoreFileRequestFromFileName(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	return mockStream
}

func CreateValidStoreRevisionStreamThatAbortsOnFileWrite(t *testing.T, fileId uuid.UUID, fileName string, successfulWritenfileChunks [][]byte) *MockFileService_StoreRevisionServer {
	// TODO reorganzie/rename this as is has nothing to do with file write
	mockStream := CreateStoreRevisionStream(t)

	metadata := CreateMetadataStoreRevisionRequest(t, fileId, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreRevisionRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	return mockStream
}

func CreateValidStoreFileStreamThatAbortsOnFileClose(t *testing.T, fileName string, successfulWritenfileChunks [][]byte) *MockFileService_StoreFileServer {
	mockStream := CreateStoreFileStream(t)

	metadata := CreateMetadataStoreFileRequestFromFileName(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func CreateValidStoreRevisionStreamThatAbortsOnFileClose(t *testing.T, fileId uuid.UUID, fileName string, successfulWritenfileChunks [][]byte) *MockFileService_StoreRevisionServer {
	mockStream := CreateStoreRevisionStream(t)

	metadata := CreateMetadataStoreRevisionRequest(t, fileId, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreRevisionRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

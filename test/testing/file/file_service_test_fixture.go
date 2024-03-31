package file

import (
	"context"
	"io"
	"testing"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateFileStream(t *testing.T) *FileService_StoreFileServer {
	mockStream := NewFileService_StoreFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateMetadataRequest(t *testing.T, fileName string) *v1.StoreFileRequest {
	metadata := &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Name{
			Name: fileName,
		},
	}
	return metadata
}

func CreateChunkRequest(t *testing.T, chunk []byte) *v1.StoreFileRequest {
	chunkRequest := &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Chunk{
			Chunk: chunk,
		},
	}
	return chunkRequest
}

func CreateValidFileStream(t *testing.T, fileName string, fileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := CreateFileStream(t)

	metadata := CreateMetadataRequest(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range fileChunks {
		chunkRequest := CreateChunkRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func CreateValidFileStreamThatAbortsOnFileOperation(t *testing.T, fileName string, successfulWritenfileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := CreateFileStream(t)

	metadata := CreateMetadataRequest(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	return mockStream
}

func SetupAndRecordSuccessfulResponse(t *testing.T, mockStream *FileService_StoreFileServer, actualResponse **v1.StoreFileResponse) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *v1.StoreFileResponse) {
		*actualResponse = response
	}).Return(nil).Times(1)
}

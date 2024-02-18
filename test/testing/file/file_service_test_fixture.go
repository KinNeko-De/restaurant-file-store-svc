package file

import (
	"context"
	"io"
	"testing"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateValidFileStream(t *testing.T, fileName string, fileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := NewFileService_StoreFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	var metadata = &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Name{
			Name: fileName,
		},
	}
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)
	for _, chunk := range fileChunks {
		var chunkRequest = &v1.StoreFileRequest{
			File: &v1.StoreFileRequest_Chunk{
				Chunk: chunk,
			},
		}
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func SetupSuccessfulResponse(t *testing.T, mockStream *FileService_StoreFileServer, actualResponse **v1.StoreFileResponse) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *v1.StoreFileResponse) {
		*actualResponse = response
	}).Return(nil).Times(1)
}

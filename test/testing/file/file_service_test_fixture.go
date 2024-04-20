package file

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func CreateStoreFileStream(t *testing.T) *FileService_StoreFileServer {
	mockStream := NewFileService_StoreFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateStoreRevisionStream(t *testing.T) *FileService_StoreRevisionServer {
	mockStream := NewFileService_StoreRevisionServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

func CreateMetadataStoreFileRequest(t *testing.T, fileName string) *v1.StoreFileRequest {
	metadata := &v1.StoreFileRequest{
		Part: &v1.StoreFileRequest_StoreFile{
			StoreFile: &v1.StoreFile{
				Name: fileName,
			},
		},
	}
	return metadata
}

func CreateMetadataStoreRevisionRequest(t *testing.T, fileId uuid.UUID, fileName string) *v1.StoreRevisionRequest {
	metadata := &v1.StoreRevisionRequest{
		Part: &v1.StoreRevisionRequest_StoreRevision{
			StoreRevision: &v1.StoreRevision{
				FileId: &protobuf.Uuid{
					Value: fileId.String(),
				},
				StoreFile: &v1.StoreFile{
					Name: fileName,
				},
			},
		},
	}
	return metadata
}

func CreateChunkStoreFileRequest(t *testing.T, chunk []byte) *v1.StoreFileRequest {
	chunkRequest := &v1.StoreFileRequest{
		Part: &v1.StoreFileRequest_Chunk{
			Chunk: chunk,
		},
	}
	return chunkRequest
}

func CreateChunkStoreRevisionRequest(t *testing.T, chunk []byte) *v1.StoreRevisionRequest {
	chunkRequest := &v1.StoreRevisionRequest{
		Part: &v1.StoreRevisionRequest_Chunk{
			Chunk: chunk,
		},
	}
	return chunkRequest
}

func CreateDownloadFileRequest(t *testing.T, fileId *protobuf.Uuid) *v1.DownloadFileRequest {
	request := &v1.DownloadFileRequest{
		FileId: fileId,
	}
	return request
}

func CreateValidStoreFileStream(t *testing.T, fileName string, fileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := CreateStoreFileStream(t)

	metadata := CreateMetadataStoreFileRequest(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range fileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func CreateValidStoreRevisionStream(t *testing.T, existingFileId uuid.UUID, fileName string, fileChunks [][]byte) *FileService_StoreRevisionServer {
	mockStream := CreateStoreRevisionStream(t)

	metadata := CreateMetadataStoreRevisionRequest(t, existingFileId, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range fileChunks {
		chunkRequest := CreateChunkStoreRevisionRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func CreateValidStoreFileStreamThatAbortsOnFileWrite(t *testing.T, fileName string, successfulWritenfileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := CreateStoreFileStream(t)

	metadata := CreateMetadataStoreFileRequest(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	return mockStream
}

func CreateValidStoreRevisionStreamThatAbortsOnFileWrite(t *testing.T, fileId uuid.UUID, fileName string, successfulWritenfileChunks [][]byte) *FileService_StoreRevisionServer {
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

func CreateValidStoreFileStreamThatAbortsOnFileClose(t *testing.T, fileName string, successfulWritenfileChunks [][]byte) *FileService_StoreFileServer {
	mockStream := CreateStoreFileStream(t)

	metadata := CreateMetadataStoreFileRequest(t, fileName)
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)

	for _, chunk := range successfulWritenfileChunks {
		chunkRequest := CreateChunkStoreFileRequest(t, chunk)
		mockStream.EXPECT().Recv().Return(chunkRequest, nil).Times(1)
	}

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	return mockStream
}

func CreateValidStoreRevisionStreamThatAbortsOnFileClose(t *testing.T, fileId uuid.UUID, fileName string, successfulWritenfileChunks [][]byte) *FileService_StoreRevisionServer {
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

func SetupAndRecordSuccessfulStoreFileResponse(t *testing.T, mockStream *FileService_StoreFileServer, actualResponse **v1.StoreFileResponse) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *v1.StoreFileResponse) {
		*actualResponse = response
	}).Return(nil).Times(1)
}

func SetupAndRecordSuccessfulStoreRevisionResponse(t *testing.T, mockStream *FileService_StoreRevisionServer, actualResponse **v1.StoreFileResponse) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *v1.StoreFileResponse) {
		*actualResponse = response
	}).Return(nil).Times(1)
}

func CreateDownloadFileStream(t *testing.T) *FileService_DownloadFileServer {
	mockStream := NewFileService_DownloadFileServer(t)

	ctx := context.Background()
	mockStream.EXPECT().Context().Return(ctx).Maybe()

	return mockStream
}

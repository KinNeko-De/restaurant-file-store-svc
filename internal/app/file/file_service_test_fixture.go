package file

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func createSut(t *testing.T, mockFileRepository *MockFileRepository, mockFileMetadataRepository *MockFileMetadataRepository) FileServiceServer {
	sut := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository
	return sut
}

func createValidFileStream(t *testing.T, fileName string, fileChunks [][]byte) *FileService_StoreFileServer {
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

func setupSuccessfulResponse(t *testing.T, mockStream *FileService_StoreFileServer, actualResponse **v1.StoreFileResponse) {
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *v1.StoreFileResponse) {
		*actualResponse = response
	}).Return(nil).Times(1)
}

func createWriterCloserMock(t *testing.T, writtenChunks [][]byte) *MockWriteCloser {
	fileWriter := &MockWriteCloser{}
	for _, chunk := range writtenChunks {
		fileWriter.EXPECT().Write(chunk).Return(len(chunk), nil).Times(1)
	}
	fileWriter.EXPECT().Close().Return(nil).Times(1)
	return fileWriter
}

func createFileRepositoryMock(t *testing.T, fileWriter *MockWriteCloser, generatedFileId **uuid.UUID) *MockFileRepository {
	mockFileRepository := &MockFileRepository{}
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), 0).
		Run(func(ctx context.Context, fileId uuid.UUID, chunkSize int) { *generatedFileId = &fileId }).
		Return(fileWriter, nil).
		Times(1)

	return mockFileRepository
}

func createFileMetadataRepositoryMock(t *testing.T) *MockFileMetadataRepository {
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	return mockFileMetadataRepository
}

func splitIntoChunks(b []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(b); i += chunkSize {
		end := i + chunkSize

		if end > len(b) {
			end = len(b)
		}

		chunks = append(chunks, b[i:end])
	}

	return chunks
}

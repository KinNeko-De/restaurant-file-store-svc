//go:build unit

package file

import (
	context "context"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	ioFixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/io"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupAndRecordSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{sentFile})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)
	assert.NotEqual(t, uuid.Nil, generatedFileId)
	assert.Equal(t, uuid.Version(0x4), generatedFileId.Version())
	assert.Equal(t, uuid.RFC4122, generatedFileId.Variant())

	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expectedFileExtension, storedFileMetadata.Revisions[0].Extension)
	assert.NotNil(t, storedFileMetadata.Revisions[0].CreatedAt)

	assert.Equal(t, generatedFileId.String(), actualResponse.StoredFile.Id.Value)
	assert.Equal(t, generatedRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()[0:512]
	sentFileName := "test.pdf"
	expectedSize := uint64(512)
	expectedMediaType := "application/pdf"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupAndRecordSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{sentFile})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)

	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)

	assert.NotNil(t, storedFileMetadata)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)
}

func TestStoreFile_FileDataIsSentInMultipleChunks_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"
	expectedSize := uint64(51124)
	expectedMediaType := "application/pdf"
	expectedFileExtension := ".pdf"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, chunks)
	fixture.SetupAndRecordSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloser(t, chunks)
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)
	assert.Nil(t, actualError)

	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expectedFileExtension, storedFileMetadata.Revisions[0].Extension)
	assert.NotNil(t, storedFileMetadata.Revisions[0].CreatedAt)
}

func TestStoreFile_CommunicationError_MetadataRequest_RetryIsRequested(t *testing.T) {
	mockStream := fixture.CreateFileStream(t)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{})
	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_CommunicationError_ChunckRequest_RetryIsRequested(t *testing.T) {
	mockStream := fixture.CreateFileStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataRequest(t, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{})
	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_CommunicationError_SendAndClose_RetryIsRequested(t *testing.T) {
	file := fixture.TextFile()

	mockStream := fixture.CreateFileStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataRequest(t, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(fixture.CreateChunkRequest(t, file), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(mock.Anything).Return(errors.New("ups..someting went wrong")).Times(1)
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{file})
	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.Contains(t, actualStatus.Message(), "response")
	assert.NotNil(t, storedFileMetadata) // TODO: Decide how to clean up this, maybe add metrics to track this; maybe add a small saga?
}

func TestStoreFile_InvalidRequest_MetadataIsMissing_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateFileStream(t)
	firstRequest := fixture.CreateChunkRequest(t, fixture.TextFile())
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_InvalidRequest_MetadataIsSentTwice_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateFileStream(t)
	firstRequest := fixture.CreateMetadataRequest(t, "test.txt")
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)
	secondRequest := fixture.CreateMetadataRequest(t, "test2.txt")
	mockStream.EXPECT().Recv().Return(secondRequest, nil).Times(1)

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata

	fileWriter := ioFixture.CreateWriterCloser(t, [][]byte{})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_FileCreatingError_RetryRequested(t *testing.T) {
	err := errors.New("Error creating file")
	sentFileName := "test.txt"

	var storedFileMetadata *FileMetadata
	mockStream := fixture.CreateValidFileStreamThatAbortsOnFileWrite(t, sentFileName, [][]byte{})

	mockFileRepository := &MockFileRepository{}
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), mock.IsType(uuid.New())).Return(nil, err).Times(1)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "create")
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_FileWritingError_RetryRequested(t *testing.T) {
	err := errors.New("Error writing file")
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	mockStream := fixture.CreateValidFileStreamThatAbortsOnFileWrite(t, sentFileName, [][]byte{sentFile})
	fileWriter := ioFixture.CreateWriterCloserRanIntoWriteError(t, [][]byte{}, err)
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "write")
	assert.Nil(t, storedFileMetadata)
}

func TestStoreFile_FileClosingError_RetryRequested(t *testing.T) {
	err := errors.New("Error closing file")
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	mockStream := fixture.CreateValidFileStreamThatAbortsOnFileClose(t, sentFileName, [][]byte{sentFile})
	fileWriter := ioFixture.CreateWriterCloserRanIntoCloseError(t, [][]byte{sentFile}, err)
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId, &generatedRevisionId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t, &storedFileMetadata)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "close")
	assert.Nil(t, storedFileMetadata)
}

func createSut(t *testing.T, mockFileRepository *MockFileRepository, mockFileMetadataRepository *MockFileMetadataRepository) FileServiceServer {
	t.Helper()
	sut := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository
	return sut
}

func createFileRepositoryMock(t *testing.T, fileWriter *ioFixture.MockWriteCloser, generatedFileId **uuid.UUID, generatedRevisionId **uuid.UUID) *MockFileRepository {
	t.Helper()
	mockFileRepository := &MockFileRepository{}
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), mock.IsType(uuid.New())).
		Run(func(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) {
			*generatedFileId = &fileId
			*generatedRevisionId = &revisionId
		}).
		Return(fileWriter, nil).
		Times(1)

	return mockFileRepository
}

func createFileMetadataRepositoryMock(t *testing.T, storedFileMetadata **FileMetadata) *MockFileMetadataRepository {
	t.Helper()
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	mockFileMetadataRepository.EXPECT().StoreFileMetadata(mock.Anything, mock.IsType(FileMetadata{})).
		Run(func(ctx context.Context, fileMetadata FileMetadata) { *storedFileMetadata = &fileMetadata }).
		Return(nil).
		Times(1)

	return mockFileMetadataRepository
}

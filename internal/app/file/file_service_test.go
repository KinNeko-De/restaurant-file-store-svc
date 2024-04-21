//go:build unit

package file

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	ioFixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/io"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredFileMetadata := mockFileMetadataRepository.setupStoreFileMetadata(t)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)
	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()

	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.Equal(t, uuid.Version(0x4), actualStoredFileId.Version())
	assert.Equal(t, uuid.RFC4122, actualStoredFileId.Variant())
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, uuid.Version(0x4), actualStoredRevisionId.Version())
	assert.Equal(t, uuid.RFC4122, actualStoredRevisionId.Variant())

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expectedFileExtension, storedFileMetadata.Revisions[0].Extension)
	assert.NotNil(t, storedFileMetadata.Revisions[0].CreatedAt)

	assert.Equal(t, actualStoredFileId.String(), actualResponse.StoredFile.Id.Value)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreRevision_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	existingFileId := uuid.New()
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredRevision := mockFileMetadataRepository.setupStoreRevisionMetadata(t, existingFileId)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.Nil(t, actualError)
	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)

	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, uuid.Version(0x4), actualStoredRevisionId.Version())
	assert.Equal(t, uuid.RFC4122, actualStoredRevisionId.Variant())

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	actualStoredRevision := recordStoredRevision()
	assert.NotNil(t, actualStoredRevision)
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedSize, actualStoredRevision.Size)
	assert.Equal(t, expectedMediaType, actualStoredRevision.MediaType)
	assert.Equal(t, expectedFileExtension, actualStoredRevision.Extension)
	assert.NotNil(t, actualStoredRevision.CreatedAt)
}

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()[0:512]
	chunks := [][]byte{sentFile}
	sentFileName := "test.pdf"
	expectedSize := uint64(512)
	expectedMediaType := "application/pdf"
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredFileMetadata := mockFileMetadataRepository.setupStoreFileMetadata(t)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)

	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)

	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	existingFileId := uuid.New()
	sentFile := fixture.PdfFile()[0:512]
	chunks := [][]byte{sentFile}
	sentFileName := "test.pdf"
	expectedSize := uint64(512)
	expectedMediaType := "application/pdf"
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredRevision := mockFileMetadataRepository.setupStoreRevisionMetadata(t, existingFileId)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.Nil(t, actualError)

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)

	actualStoredRevision := recordStoredRevision()
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedSize, actualStoredRevision.Size)
	assert.Equal(t, expectedMediaType, actualStoredRevision.MediaType)

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreFile_FileDataIsSentInMultipleChunks_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"
	expectedSize := uint64(51124)
	expectedMediaType := "application/pdf"
	expectedFileExtension := ".pdf"
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredFileMetadata := mockFileMetadataRepository.setupStoreFileMetadata(t)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedSize, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expectedFileExtension, storedFileMetadata.Revisions[0].Extension)
	assert.NotNil(t, storedFileMetadata.Revisions[0].CreatedAt)

	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileDataIsSentInMultipleChunks_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	existingFileId := uuid.New()
	sentFile := fixture.PdfFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"
	expectedSize := uint64(51124)
	expectedMediaType := "application/pdf"
	expectedFileExtension := ".pdf"
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	recordActualResponse := mockStream.SetupSendAndClose(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoreRevision := mockFileMetadataRepository.setupStoreRevisionMetadata(t, existingFileId)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.Nil(t, actualError)

	actualResponse := recordActualResponse()
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	actualStoredRevision := recordStoreRevision()
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedSize, actualStoredRevision.Size)
	assert.Equal(t, expectedMediaType, actualStoredRevision.MediaType)
	assert.Equal(t, expectedFileExtension, actualStoredRevision.Extension)
	assert.NotNil(t, actualStoredRevision.CreatedAt)

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreFile_CommunicationError_MetadataRequest_RetryIsRequested(t *testing.T) {
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
}

func TestStoreRevision_CommunicationError_MetadataRequest_RetryIsRequested(t *testing.T) {
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
}

func TestStoreFile_CommunicationError_ChunckRequest_RetryIsRequested(t *testing.T) {
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, ioFixture.NewMockWriteCloser(t))

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_CommunicationError_ChunckRequest_RetryIsRequested(t *testing.T) {
	existingFileId := uuid.New()
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreRevisionRequest(t, existingFileId, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, errors.New("ups..someting went wrong")).Times(1)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, ioFixture.NewMockWriteCloser(t))

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreFile_CommunicationError_SendAndClose_RetryIsRequested(t *testing.T) {
	file := fixture.TextFile()
	chunks := [][]byte{file}

	mockStream := fixture.CreateStoreFileStream(t)
	// TODO use fixture here
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(fixture.CreateChunkStoreFileRequest(t, file), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(mock.Anything).Return(errors.New("ups..someting went wrong")).Times(1)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredFileMetadata := mockFileMetadataRepository.setupStoreFileMetadata(t)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.Contains(t, actualStatus.Message(), "response")
	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata) // TODO: Decide how to clean up this, maybe add metrics to track this; maybe add a small saga?
	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_CommunicationError_SendAndClose_RetryIsRequested(t *testing.T) {
	file := fixture.TextFile()
	chunks := [][]byte{file}
	existingFileId := uuid.New()
	mockStream := fixture.CreateStoreRevisionStream(t)
	// TODO use fixture here
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreRevisionRequest(t, existingFileId, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(fixture.CreateChunkStoreRevisionRequest(t, file), nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(mock.Anything).Return(errors.New("ups..someting went wrong")).Times(1)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	recordStoredRevision := mockFileMetadataRepository.setupStoreRevisionMetadata(t, existingFileId)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")

	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.Contains(t, actualStatus.Message(), "response")
	assert.NotNil(t, recordStoredRevision()) // TODO: Decide how to clean up this, maybe add metrics to track this; maybe add a small saga? and should the user retry this? the revision is already stored

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreFile_InvalidRequest_MetadataIsMissing_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateStoreFileStream(t)
	firstRequest := fixture.CreateChunkStoreFileRequest(t, fixture.TextFile())
	// TODO use fixture here
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
}

func TestStoreRevision_InvalidRequest_MetadataIsMissing_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateStoreRevisionStream(t)
	firstRequest := fixture.CreateChunkStoreRevisionRequest(t, fixture.TextFile())
	// TODO use fixture here
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
}

func TestStoreFile_InvalidRequest_MetadataIsSentTwice_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateStoreFileStream(t)
	firstRequest := fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt")
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)
	secondRequest := fixture.CreateMetadataStoreFileRequestFromFileName(t, "test2.txt")
	mockStream.EXPECT().Recv().Return(secondRequest, nil).Times(1)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredFileId := mockFileRepository.setupCreateFileNewFile(t, ioFixture.NewMockWriteCloser(t))

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_InvalidRequest_MetadataIsSentTwice_FileIsRejected(t *testing.T) {
	existingFileId := uuid.New()
	// TODO use fixture here
	mockStream := fixture.CreateStoreRevisionStream(t)
	firstRequest := fixture.CreateMetadataStoreRevisionRequest(t, existingFileId, "test.txt")
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)
	secondRequest := fixture.CreateMetadataStoreRevisionRequest(t, existingFileId, "test2.txt")
	mockStream.EXPECT().Recv().Return(secondRequest, nil).Times(1)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, ioFixture.NewMockWriteCloser(t))

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.NotEqual(t, uuid.Nil, recordStoredRevisionId())
}

func TestStoreFile_FileCreatingError_RetryRequested(t *testing.T) {
	err := errors.New("Error creating file")
	sentFileName := "test.txt"
	chunks := [][]byte{}
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), mock.IsType(uuid.New())).Return(nil, err).Times(1)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "create")
}

func TestStoreRevision_FileCreatingError_RetryRequested(t *testing.T) {
	err := errors.New("Error creating file")
	sentFileName := "test.txt"
	existingFileId := uuid.New()
	chunks := [][]byte{}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	// TODO use fixture here
	mockFileRepository.EXPECT().CreateFile(mock.Anything, existingFileId, mock.IsType(uuid.New())).Return(nil, err).Times(1)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "create")
}

func TestStoreFile_FileWritingError_RetryRequested(t *testing.T) {
	err := errors.New("Error writing file")
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	chunks := [][]byte{sentFile}
	writtenChunks := [][]byte{}
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupWrite(t, writtenChunks)
	fileWriter.SetupWriteError(t, err)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredIds := mockFileRepository.setupCreateFileNewFile(t, fileWriter)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "write")
	actualStoredFileId, actualStoredRevisionId := recordStoredIds()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileWritingError_RetryRequested(t *testing.T) {
	err := errors.New("Error writing file")
	sentFileName := "test.txt"
	sentFile := fixture.TextFile()
	existingFileId := uuid.New()
	chunks := [][]byte{sentFile}
	writtenChunks := [][]byte{}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupWrite(t, writtenChunks)
	fileWriter.SetupWriteError(t, err)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "write")
	assert.NotEqual(t, uuid.Nil, recordStoredRevisionId())
}

func TestStoreFile_FileClosingError_RetryRequested(t *testing.T) {
	err := errors.New("Error closing file")
	sentFileName := "test.txt"
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupWrite(t, chunks)
	fileWriter.SetupCloseError(t, err)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredIds := mockFileRepository.setupCreateFileNewFile(t, fileWriter)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "close")
	actualStoredFileId, actualStoredRevisionId := recordStoredIds()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileClosingError_RetryRequested(t *testing.T) {
	err := errors.New("Error closing file")
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"
	existingFileId := uuid.New()
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupWrite(t, chunks)
	fileWriter.SetupCloseError(t, err)
	mockFileRepository := NewMockFileRepository(t)
	recordStoreRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)

	sut := createSut(t, mockFileRepository, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "close")
	assert.NotEqual(t, uuid.Nil, recordStoreRevisionId())
}

func TestStoreFile_StoreFileMetadataThrowsError_RetryRequested(t *testing.T) {
	err := errors.New("Error contains possible sensitive information")
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	chunks := [][]byte{sentFile}
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredIds := mockFileRepository.setupCreateFileNewFile(t, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryMockStoreFileMetadataReturnsError(t, err)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	assert.False(t, errors.Is(actualError, err), "Error should not be passed to the client")
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	actualStoredFileId, actualStoredRevisionId := recordStoredIds()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileIdNotFound(t *testing.T) {
	err := errors.New("file id not matching")
	existingFileId := uuid.New()
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryMockStoreRevisionReturnsError(t, existingFileId, err)
	mockFileMetadataRepository.EXPECT().NoMatchError().Return(err).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	assert.False(t, errors.Is(actualError, err), "Error should not be passed to the client")
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.NotFound, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), existingFileId.String())
	assert.NotEqual(t, uuid.Nil, recordStoredRevisionId())
}

func TestStoreRevision_FileIdIsNil(t *testing.T) {
	request := &v1.StoreRevisionRequest{
		Part: &v1.StoreRevisionRequest_StoreRevision{
			StoreRevision: &v1.StoreRevision{
				FileId: nil,
			},
		},
	}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.EXPECT().Recv().Return(request, nil).Times(1)

	sut := createSut(t, nil, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "mandatory")
}

func TestStoreRevision_FileIdIsInvalid(t *testing.T) {
	invalidUuid := "433b4b7c-4b1e-4b1e4b1e4b1e"

	request := &v1.StoreRevisionRequest{
		Part: &v1.StoreRevisionRequest_StoreRevision{
			StoreRevision: &v1.StoreRevision{
				FileId: &protobuf.Uuid{
					Value: invalidUuid,
				},
			},
		},
	}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.EXPECT().Recv().Return(request, nil).Times(1)

	sut := createSut(t, nil, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), invalidUuid)
}

func TestStoreRevision_StoreFileMetadataThrowsError_RetryRequested(t *testing.T) {
	err := errors.New("Error contains possible sensitive information")
	existingFileId := uuid.New()
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	chunks := [][]byte{sentFile}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, sentFileName))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	fileWriter := ioFixture.NewMockWriteCloser(t)
	fileWriter.SetupSuccessfulWrite(t, chunks)
	mockFileRepository := NewMockFileRepository(t)
	recordStoredRevisionId := mockFileRepository.setupCreateFileNewRevision(t, existingFileId, fileWriter)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryMockStoreRevisionReturnsError(t, existingFileId, err)
	mockFileMetadataRepository.EXPECT().NoMatchError().Return(errors.New("not expected error")).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	assert.False(t, errors.Is(actualError, err), "Error should not be passed to the client")
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
	assert.NotEqual(t, uuid.Nil, recordStoredRevisionId())
}

func TestDownloadFile_FileIdIsNil(t *testing.T) {
	request := fixture.CreateDownloadFileRequest(t, nil)

	mockStream := fixture.CreateDownloadFileStream(t)
	sut := createSut(t, nil, NewMockFileMetadataRepository(t))
	err := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, err)
	actualStatus, ok := status.FromError(err)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "fileId")
	assert.Contains(t, actualStatus.Message(), "mandatory")
}

func TestDownloadFile_FileIdIsInvalid(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"Empty", ""},
		{"InvalidFormat", "433b4b7c-4b1e-4b1e4b1e4b1e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileId := &protobuf.Uuid{
				Value: test.uuid,
			}
			request := fixture.CreateDownloadFileRequest(t, fileId)
			mockStream := fixture.CreateDownloadFileStream(t)

			sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
			err := sut.DownloadFile(request, mockStream)

			assert.NotNil(t, err)
			actualStatus, ok := status.FromError(err)
			require.True(t, ok, "Expected a gRPC status error")
			require.NotNil(t, actualStatus)
			assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), "fileId")
			assert.Contains(t, actualStatus.Message(), "not a valid uuid")
			assert.Contains(t, actualStatus.Message(), test.uuid)
		})
	}
}

func TestDownloadFile_FileIdNotFound(t *testing.T) {
	tests := []struct {
		name          string
		notFoundError error
	}{
		{"NotWrappedEror", errors.New("file not found")},
		{"JoinedError", errors.Join(errors.New("file not found"), errors.New("wrappedError"))},
		{"WrappedError", fmt.Errorf("wrapper error %w", errors.New("file not found"))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileId := uuid.New()
			request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
			mockFileMetadataRepository := NewMockFileMetadataRepository(t)
			mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(FileMetadata{}, test.notFoundError).Times(1)
			mockFileMetadataRepository.EXPECT().NotFoundError().Return(test.notFoundError).Times(1)
			mockStream := fixture.CreateDownloadFileStream(t)

			sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
			actualError := sut.DownloadFile(request, mockStream)

			assert.NotNil(t, actualError)
			actualStatus, ok := status.FromError(actualError)
			require.True(t, ok, "Expected a gRPC status error")
			require.NotNil(t, actualStatus)
			assert.Equal(t, codes.NotFound, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), fileId.String())
		})
	}
}

func TestDownloadFile_ErrorFetchingMetadataThatIsNotSameAsNotFoundError(t *testing.T) {
	fileId := uuid.New()
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(FileMetadata{}, errors.New("ups..someting went wrong")).Times(1)
	mockFileMetadataRepository.EXPECT().NotFoundError().Return(errors.New("file not found")).Times(1)
	mockStream := fixture.CreateDownloadFileStream(t)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.Internal, actualStatus.Code())
}

func TestDownloadFile_SendMetadataFails(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
	sendErr := errors.New("send error due to network connection as example")

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        revisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(sendErr).Times(1)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.Internal, actualStatus.Code())
}

func TestDownloadFile_FindingTheFileBytesFails(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
	openErr := errors.New("open error because file disapperred most likey due someone fuckeled around manually")

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        revisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.EXPECT().OpenFile(mock.Anything, fileId, revisionId).Return(nil, openErr).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.Internal, actualStatus.Code())
}

func TestDownloadFile_ReadingTheFileBytesFails(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
	readErr := errors.New("read error due to network connection as example")

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        revisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := NewMockFileRepository(t)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupReadError(t, readErr)
	mockFileRepository.EXPECT().OpenFile(mock.Anything, fileId, revisionId).Return(mockReader, nil).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.Internal, actualStatus.Code())
}

func TestDownloadFile_ClosingTheFileBytesFails_ErrorIsNotReportedToClient(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)
	file := fixture.TextFile()
	closeErr := errors.New("close error due to network connection as example")

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        revisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := NewMockFileRepository(t)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockReader.SetupCloseError(t, closeErr)
	mockFileRepository.EXPECT().OpenFile(mock.Anything, fileId, revisionId).Return(mockReader, nil).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualFile := make([]byte, 0)
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *v1.DownloadFileResponse) {
		actualFile = append(actualFile, response.GetChunk()...)
	}).Return(nil)
	actualError := sut.DownloadFile(request, mockStream)

	assert.Nil(t, actualError)
	assert.Equal(t, file, actualFile)
}

func TestDownloadFile_SendFileBytesFails(t *testing.T) {
	fileId := uuid.New()
	revisionId := uuid.New()
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
	file := fixture.TextFile()
	sendErr := errors.New("read error due to network connection as example")

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        revisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1)     // metadata are sent
	mockStream.EXPECT().Send(mock.Anything).Return(sendErr).Times(1) // file bytes are sent
	mockFileRepository := NewMockFileRepository(t)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockFileRepository.EXPECT().OpenFile(mock.Anything, fileId, revisionId).Return(mockReader, nil).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.Internal, actualStatus.Code())
}

func TestDownloadFile_LatestRevisionIsDownloaded_FileIsSplittedIntoChunks(t *testing.T) {
	fileId := uuid.New()
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
	fileThatIsBiggerThanTheMaxChunkSizeForGrpc := fixture.PdfFile()

	firstRevision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain",
		Size:      1024,
		CreatedAt: time.Now().UTC().Add(-time.Hour),
	}

	latestedRevision := Revision{
		Id:        uuid.New(),
		Extension: ".pdf",
		MediaType: "application/pdf",
		Size:      2048,
		CreatedAt: time.Now().UTC(),
	}

	fileMetadata := FileMetadata{
		Id:        fileId,
		Revisions: []Revision{firstRevision, latestedRevision},
	}

	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFileMetadataRepositoryToFetchMetadata(t, fileId, fileMetadata)
	mockFileRepository := NewMockFileRepository(t)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupSuccessfulRead(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc)
	mockFileRepository.EXPECT().OpenFile(mock.Anything, fileId, latestedRevision.Id).Return(mockReader, nil).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFileMetadata := mockStream.SetupRecordStoredFileMetadata(t)
	recordDownloadedFile := mockStream.SetupRecordDownloadedFile(t)

	actualError := sut.DownloadFile(request, mockStream)

	assert.Nil(t, actualError)
	actualStoredFileMetadata := recordStoredFileMetadata()
	assert.Equal(t, latestedRevision.Extension, actualStoredFileMetadata.Extension)
	assert.Equal(t, latestedRevision.MediaType, actualStoredFileMetadata.MediaType)
	assert.Equal(t, latestedRevision.Size, actualStoredFileMetadata.Size)
	assert.Equal(t, latestedRevision.CreatedAt, actualStoredFileMetadata.CreatedAt.AsTime())
	assert.Equal(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc, recordDownloadedFile())
}

func createSut(t *testing.T, mockFileRepository *MockFileRepository, mockFileMetadataRepository *MockFileMetadataRepository) FileServiceServer {
	t.Helper()
	sut := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository
	return sut
}

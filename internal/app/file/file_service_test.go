//go:build unit

package file

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	ioFixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"

	expectedMetada := ExpectedMetadata{
		Size:      4,
		MediaType: "text/plain; charset=utf-8",
		Extension: ".txt",
	}

	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSuccessfulSend(t, sentFileName, chunks)
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
	assert.Equal(t, actualStoredFileId.String(), actualResponse.StoredFile.Id.Value)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
	assertMetadata(t, actualResponse.StoredFile.Metadata, expectedMetada)
	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedMetada.Size, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMetada.MediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expectedMetada.Extension, storedFileMetadata.Revisions[0].Extension)
	assert.NotNil(t, storedFileMetadata.Revisions[0].CreatedAt)

	assert.Equal(t, actualStoredFileId.String(), actualResponse.StoredFile.Id.Value)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreRevision_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	existingFileId := uuid.New()
	sentFile := fixture.TextFile()
	chunks := [][]byte{sentFile}
	sentFileName := "test.txt"

	expectedMetada := ExpectedMetadata{
		Size:      4,
		MediaType: "text/plain; charset=utf-8",
		Extension: ".txt",
	}

	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSuccessfulSend(t, existingFileId, sentFileName, chunks)
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
	assertMetadata(t, actualResponse.StoredFile.Metadata, expectedMetada)
	actualStoredRevision := recordStoredRevision()
	assert.NotNil(t, actualStoredRevision)
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedMetada.Size, actualStoredRevision.Size)
	assert.Equal(t, expectedMetada.MediaType, actualStoredRevision.MediaType)
	assert.Equal(t, expectedMetada.Extension, actualStoredRevision.Extension)
	assert.NotNil(t, actualStoredRevision.CreatedAt)
}

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()[0:512]
	chunks := [][]byte{sentFile}
	sentFileName := "test.pdf"

	expectedMetadata := ExpectedMetadata{
		Size:      512,
		MediaType: "application/pdf",
		Extension: ".pdf",
	}

	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSuccessfulSend(t, sentFileName, chunks)
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
	assertMetadata(t, actualResponse.StoredFile.Metadata, expectedMetadata)
	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expectedMetadata.Size, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expectedMetadata.MediaType, storedFileMetadata.Revisions[0].MediaType)

	actualStoredFileId, actualStoredRevisionId := recordStoredFileId()
	assert.NotEqual(t, uuid.Nil, actualStoredFileId)
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
}

func TestStoreRevision_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	existingFileId := uuid.New()
	sentFile := fixture.PdfFile()[0:512]
	chunks := [][]byte{sentFile}
	sentFileName := "test.pdf"

	expectedMetadata := ExpectedMetadata{
		Size:      512,
		MediaType: "application/pdf",
		Extension: ".pdf",
	}

	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSuccessfulSend(t, existingFileId, sentFileName, chunks)
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
	assertMetadata(t, actualResponse.StoredFile.Metadata, expectedMetadata)
	actualStoredRevision := recordStoredRevision()
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedMetadata.Size, actualStoredRevision.Size)
	assert.Equal(t, expectedMetadata.MediaType, actualStoredRevision.MediaType)

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreFile_FileDataIsSentInMultipleChunks_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"

	expecedMetada := ExpectedMetadata{
		Size:      51124,
		MediaType: "application/pdf",
		Extension: ".pdf",
	}

	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSuccessfulSend(t, sentFileName, chunks)
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
	assertMetadata(t, actualResponse.StoredFile.Metadata, expecedMetada)
	storedFileMetadata := recordStoredFileMetadata()
	assert.NotNil(t, storedFileMetadata)
	assert.NotNil(t, storedFileMetadata.Id)
	assert.NotNil(t, storedFileMetadata.Revisions)
	assert.Len(t, storedFileMetadata.Revisions, 1)
	assert.NotNil(t, storedFileMetadata.Revisions[0].Id)
	assert.Equal(t, expecedMetada.Size, storedFileMetadata.Revisions[0].Size)
	assert.Equal(t, expecedMetada.MediaType, storedFileMetadata.Revisions[0].MediaType)
	assert.Equal(t, expecedMetada.Extension, storedFileMetadata.Revisions[0].Extension)
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

	expectedMetada := ExpectedMetadata{
		Size:      51124,
		MediaType: "application/pdf",
		Extension: ".pdf",
	}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSuccessfulSend(t, existingFileId, sentFileName, chunks)
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
	assertMetadata(t, actualResponse.StoredFile.Metadata, expectedMetada)
	actualStoredRevision := recordStoreRevision()
	assert.NotNil(t, actualStoredRevision.Id)
	assert.Equal(t, expectedMetada.Size, actualStoredRevision.Size)
	assert.Equal(t, expectedMetada.MediaType, actualStoredRevision.MediaType)
	assert.Equal(t, expectedMetada.Extension, actualStoredRevision.Extension)
	assert.NotNil(t, actualStoredRevision.CreatedAt)

	actualStoredRevisionId := recordStoredRevisionId()
	assert.NotEqual(t, uuid.Nil, actualStoredRevisionId)
	assert.Equal(t, actualStoredRevisionId.String(), actualResponse.StoredFile.RevisionId.Value)
}

func TestStoreFile_FileNameIsInvalid_FileIsRejected(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{"Empty", ""},
		{"NoFileExtension", "ineedanextension"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStream := fixture.CreateStoreFileStream(t)
			mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, test.fileName))

			sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
			actualError := sut.StoreFile(mockStream)

			assert.NotNil(t, actualError)
			actualStatus, ok := status.FromError(actualError)
			assert.True(t, ok, "Expected a gRPC status error")
			assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), test.fileName)
		})
	}
}

func TestStoreRevision_FileNameIsInvalid_FileIsRejected(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{"Empty", ""},
		{"NoFileExtension", "ineedanextension"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			existingFileId := uuid.New()
			mockStream := fixture.CreateStoreRevisionStream(t)
			mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, test.fileName))

			sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
			actualError := sut.StoreRevision(mockStream)

			assert.NotNil(t, actualError)
			actualStatus, ok := status.FromError(actualError)
			assert.True(t, ok, "Expected a gRPC status error")
			assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), test.fileName)
		})
	}
}

func TestStoreFile_CommunicationError_MetadataRequest_RetryIsRequested(t *testing.T) {
	communicationError := errors.New("ups..someting went wrong")
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendError(t, communicationError)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
}

func TestStoreRevision_CommunicationError_MetadataRequest_RetryIsRequested(t *testing.T) {
	communicationError := errors.New("ups..someting went wrong")
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendError(t, communicationError)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.Internal, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "retry")
}

func TestStoreFile_CommunicationError_ChunckRequest_RetryIsRequested(t *testing.T) {
	communicationError := errors.New("ups..someting went wrong")
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt"))
	mockStream.SetupSendError(t, communicationError)
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
	communicationError := errors.New("ups..someting went wrong")
	existingFileId := uuid.New()
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, "test.txt"))
	mockStream.SetupSendError(t, communicationError)
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
	closeError := errors.New("ups..someting went wrong")
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt"))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	mockStream.SetupSendAndCloseError(t, closeError)
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
	communicationError := errors.New("ups..someting went wrong")
	file := fixture.TextFile()
	chunks := [][]byte{file}
	existingFileId := uuid.New()
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, "test.txt"))
	mockStream.SetupSendFile(t, chunks)
	mockStream.SetupSendEndOfFile(t)
	mockStream.SetupSendAndCloseError(t, communicationError)
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
	chunks := [][]byte{fixture.TextFile()}
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendFile(t, chunks)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreFile(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
}

func TestStoreRevision_InvalidRequest_MetadataIsMissing_FileIsRejected(t *testing.T) {
	chunks := [][]byte{fixture.TextFile()}
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendFile(t, chunks)

	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
}

func TestStoreFile_InvalidRequest_MetadataIsSentTwice_FileIsRejected(t *testing.T) {
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, "test.txt"))
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreFileRequestFromFileName(t, "test2.txt"))
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
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, "test.txt"))
	mockStream.SetupSendMetadata(t, fixture.CreateMetadataStoreRevisionRequestFromFileName(t, existingFileId, "test2.txt"))
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
	mockFileRepository.setupCreateFileError(t, err)

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
	mockFileRepository.setupCreateFileError(t, err)

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
	mockFileMetadataRepository.setupFileMetadataRepositoryMockStoreRevisionReturnsError(t, existingFileId, err, err)

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
	storeRevision := &apiRestaurantFile.StoreRevision{
		FileId: nil,
	}
	request := fixture.CreateMetadataStoreRevisionRequest(t, storeRevision)
	mockStream := fixture.CreateStoreRevisionStream(t)
	mockStream.SetupSendMetadata(t, request)

	sut := createSut(t, nil, NewMockFileMetadataRepository(t))
	actualError := sut.StoreRevision(mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	assert.True(t, ok, "Expected a gRPC status error")
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "mandatory")
}

func TestStoreRevision_FileIdIsInvalid(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"Empty", ""},
		{"InvalidFormat", "433b4b7c-4b1e-4b1e4b1e4b1e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storeRevision := &apiRestaurantFile.StoreRevision{
				FileId: &protobuf.Uuid{
					Value: test.uuid,
				},
			}

			request := fixture.CreateMetadataStoreRevisionRequest(t, storeRevision)
			mockStream := fixture.CreateStoreRevisionStream(t)
			mockStream.SetupSendMetadata(t, request)

			sut := createSut(t, nil, NewMockFileMetadataRepository(t))
			actualError := sut.StoreRevision(mockStream)

			assert.NotNil(t, actualError)
			actualStatus, ok := status.FromError(actualError)
			assert.True(t, ok, "Expected a gRPC status error")
			assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), test.uuid)
		})
	}
}

func TestStoreRevision_StoreFileMetadataThrowsError_RetryRequested(t *testing.T) {
	err := errors.New("Error contains possible sensitive information")
	noMatchErr := errors.New("not expected error")
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
	mockFileMetadataRepository.setupFileMetadataRepositoryMockStoreRevisionReturnsError(t, existingFileId, err, noMatchErr)

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
	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	err := sut.DownloadFile(request, mockStream)

	assert.NotNil(t, err)
	actualStatus, ok := status.FromError(err)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "fileId")
	assert.Contains(t, actualStatus.Message(), "mandatory")
}

func TestDownloadRevision_FileIdIsNil(t *testing.T) {
	otherId, _ := apiProtobuf.ToProtobuf(uuid.New())
	request := fixture.CreateDownloadRevisionRequest(t, nil, otherId)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	err := sut.DownloadRevision(request, mockStream)

	assert.NotNil(t, err)
	actualStatus, ok := status.FromError(err)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "fileId")
	assert.Contains(t, actualStatus.Message(), "mandatory")
}

func TestDownloadRevision_RevisionIdIsNil(t *testing.T) {
	otherId, _ := apiProtobuf.ToProtobuf(uuid.New())
	request := fixture.CreateDownloadRevisionRequest(t, otherId, nil)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
	err := sut.DownloadRevision(request, mockStream)

	assert.NotNil(t, err)
	actualStatus, ok := status.FromError(err)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), "revisionId")
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

func TestDownloadRevision_FileIdIsInvalid(t *testing.T) {
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
			otherId, _ := apiProtobuf.ToProtobuf(uuid.New())
			request := fixture.CreateDownloadRevisionRequest(t, fileId, otherId)
			mockStream := fixture.CreateDownloadRevisionStream(t)

			sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
			err := sut.DownloadRevision(request, mockStream)

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

func TestDownloadRevision_RevisionIdIsInvalid(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"Empty", ""},
		{"InvalidFormat", "433b4b7c-4b1e-4b1e4b1e4b1e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			revisionId := &protobuf.Uuid{
				Value: test.uuid,
			}
			otherId, _ := apiProtobuf.ToProtobuf(uuid.New())
			request := fixture.CreateDownloadRevisionRequest(t, otherId, revisionId)
			mockStream := fixture.CreateDownloadRevisionStream(t)

			sut := createSut(t, NewMockFileRepository(t), NewMockFileMetadataRepository(t))
			err := sut.DownloadRevision(request, mockStream)

			assert.NotNil(t, err)
			actualStatus, ok := status.FromError(err)
			require.True(t, ok, "Expected a gRPC status error")
			require.NotNil(t, actualStatus)
			assert.Equal(t, codes.InvalidArgument, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), "revisionId")
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
			mockFileMetadataRepository.setupFetchFileMetadataReturnsError(t, fileId, test.notFoundError, test.notFoundError)
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

func TestDownloadRevision_FileIdNotFound(t *testing.T) {
	tests := []struct {
		name        string
		notFoundErr error
	}{
		{"NotWrappedEror", errors.New("file not found")},
		{"JoinedError", errors.Join(errors.New("file not found"), errors.New("wrappedError"))},
		{"WrappedError", fmt.Errorf("wrapper error %w", errors.New("file not found"))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileId := uuid.New()
			revisionId := uuid.New()
			request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileId, revisionId)
			mockFileMetadataRepository := NewMockFileMetadataRepository(t)
			mockFileMetadataRepository.setupFetchFileMetadataReturnsError(t, fileId, test.notFoundErr, test.notFoundErr)
			mockStream := fixture.CreateDownloadRevisionStream(t)

			sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
			actualError := sut.DownloadRevision(request, mockStream)

			assert.NotNil(t, actualError)
			actualStatus, ok := status.FromError(actualError)
			require.True(t, ok, "Expected a gRPC status error")
			require.NotNil(t, actualStatus)
			assert.Equal(t, codes.NotFound, actualStatus.Code())
			assert.Contains(t, actualStatus.Message(), fileId.String())
		})
	}
}

func TestDownloadRevision_RevisionIdNotFound(t *testing.T) {
	fileId := uuid.New()
	requestedRevisionId := uuid.New()
	existingRevisionId := uuid.New()
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileId, requestedRevisionId)

	fileMetadata := FileMetadata{
		Id: fileId,
		Revisions: []Revision{
			{
				Id:        existingRevisionId,
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}
	mockStream := fixture.CreateDownloadRevisionStream(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, codes.NotFound, actualStatus.Code())
	assert.Contains(t, actualStatus.Message(), requestedRevisionId.String())
}

func TestDownloadFile_ErrorFetchingMetadataThatIsNotSameAsNotFoundError(t *testing.T) {
	fileId := uuid.New()
	fetchErr := errors.New("ups..someting went wrong")
	notFoundErr := errors.New("file not found")
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadataReturnsError(t, fileId, fetchErr, notFoundErr)
	mockStream := fixture.CreateDownloadFileStream(t)
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileId)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
}

func TestDownloadRevision_ErrorFetchingMetadataThatIsNotSameAsNotFoundError(t *testing.T) {
	fetchErr := errors.New("ups..someting went wrong")
	notFoundErr := errors.New("file not found")

	fileId := uuid.New()
	revisionId := uuid.New()
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileId, revisionId)

	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadataReturnsError(t, fileId, fetchErr, notFoundErr)
	mockStream := fixture.CreateDownloadRevisionStream(t)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
}

func TestDownloadFile_SendMetadataFails(t *testing.T) {
	sendErr := errors.New("send error due to network connection as example")

	fileMetadata := FileMetadata{
		Id: uuid.New(),
		Revisions: []Revision{
			{
				Id:        uuid.New(),
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileMetadata.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	mockStream.SetupSendError(t, sendErr)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
}

func TestDownloadRevision_SendMetadataFails(t *testing.T) {
	sendErr := errors.New("send error due to network connection as example")

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileMetadata.Id, revision.Id)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	mockStream.SetupSendError(t, sendErr)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)

	sut := createSut(t, NewMockFileRepository(t), mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
}

func TestDownloadFile_FindingTheFileBytesFails(t *testing.T) {
	openErr := errors.New("open error because file disapperred most likey due someone fuckeled around manually")

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileMetadata.Id, revision.Id)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileError(t, fileMetadata.Id, revision.Id, openErr)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadRevision_FindingTheFileBytesFails(t *testing.T) {
	openErr := errors.New("open error because the file has disappeared, most likely because someone manually fuckeled around")

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileMetadata.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileError(t, fileMetadata.Id, revision.Id, openErr)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadFile_ReadingTheFileBytesFails(t *testing.T) {
	readErr := errors.New("read error due to network connection as example")

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileMetadata.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupReadError(t, readErr)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadRevision_ReadingTheFileBytesFails(t *testing.T) {
	readErr := errors.New("read error due to network connection as example")

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileMetadata.Id, revision.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupReadError(t, readErr)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadFile_ClosingTheFileBytesFails_ErrorIsNotReportedToClient(t *testing.T) {
	closeErr := errors.New("close error due to network connection as example")
	file := fixture.TextFile()

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileMetadata.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	recordDownloadedFile := mockStream.SetupRecordDownloadedFile(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockReader.SetupCloseError(t, closeErr)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.Nil(t, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
	assert.Equal(t, file, recordDownloadedFile())
}

func TestDownloadRevision_ClosingTheFileBytesFails_ErrorIsNotReportedToClient(t *testing.T) {
	closeErr := errors.New("close error due to network connection as example")
	file := fixture.TextFile()

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileMetadata.Id, revision.Id)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	recordDownloadedFile := mockStream.SetupRecordDownloadedFile(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockReader.SetupCloseError(t, closeErr)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assert.Nil(t, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
	assert.Equal(t, file, recordDownloadedFile())
}

func TestDownloadFile_SendFileBytesFails(t *testing.T) {
	sendErr := errors.New("read error due to network connection as example")
	file := fixture.TextFile()

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadFileRequestFromUuid(t, fileMetadata.Id)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockStream.SetupSendError(t, sendErr)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadRevision_SendFileBytesFails(t *testing.T) {
	sendErr := errors.New("read error due to network connection as example")
	file := fixture.TextFile()

	revision := Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain; charset=utf-8",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}
	fileMetadata := FileMetadata{
		Id:        uuid.New(),
		Revisions: []Revision{revision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, revision)
	request := fixture.CreateDownloadRevisionRequestFromUuid(t, fileMetadata.Id, revision.Id)

	mockStream := fixture.CreateDownloadRevisionStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	mockStream.SetupSendError(t, sendErr)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupRead(t, file)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, revision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assertGrpcError(t, codes.Internal, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
}

func TestDownloadFile_LatestRevisionIsDownloaded_FileIsSplittedIntoChunks(t *testing.T) {
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
		Id:        uuid.New(),
		Revisions: []Revision{firstRevision, latestedRevision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, latestedRevision)
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileMetadata.Id)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	recordDownloadedFile := mockStream.SetupRecordDownloadedFile(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupSuccessfulRead(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, latestedRevision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadFile(request, mockStream)

	assert.Nil(t, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
	assert.Equal(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc, recordDownloadedFile())
}

func TestDownloadRevision_FirstRevisionIsDownloaded_FileIsSplittedIntoChunks(t *testing.T) {
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
		Id:        uuid.New(),
		Revisions: []Revision{firstRevision, latestedRevision},
	}
	expectedStoredFile := expectRevision(fileMetadata.Id, firstRevision)
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileMetadata.Id)
	requestedRevisionId, _ := apiProtobuf.ToProtobuf(firstRevision.Id)
	request := fixture.CreateDownloadRevisionRequest(t, requestedFileId, requestedRevisionId)

	mockStream := fixture.CreateDownloadFileStream(t)
	recordStoredFile := mockStream.SetupRecordStoredFileMetadata(t)
	recordDownloadedFile := mockStream.SetupRecordDownloadedFile(t)
	mockFileMetadataRepository := NewMockFileMetadataRepository(t)
	mockFileMetadataRepository.setupFetchFileMetadata(t, fileMetadata)
	mockReader := ioFixture.NewMockReadCloser(t)
	mockReader.SetupSuccessfulRead(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc)
	mockFileRepository := NewMockFileRepository(t)
	mockFileRepository.setupOpenFileExistingFile(t, fileMetadata.Id, firstRevision.Id, mockReader)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.DownloadRevision(request, mockStream)

	assert.Nil(t, actualError)
	assertExistingFile(t, recordStoredFile(), expectedStoredFile)
	assert.Equal(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc, recordDownloadedFile())
}

func createSut(t *testing.T, mockFileRepository *MockFileRepository, mockFileMetadataRepository *MockFileMetadataRepository) FileServiceServer {
	t.Helper()
	sut := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository
	return sut
}

type ExpectedMetadata struct {
	Size      uint64
	MediaType string
	Extension string
}

type ExpectedExistingFile struct {
	Id         uuid.UUID
	RevisionId uuid.UUID
	Metadata   ExpectedMetadata
	CreatedAt  time.Time
}

func expectRevision(fileId uuid.UUID, expectedRevision Revision) ExpectedExistingFile {
	return ExpectedExistingFile{
		Id:         fileId,
		RevisionId: expectedRevision.Id,
		Metadata: ExpectedMetadata{
			Size:      expectedRevision.Size,
			MediaType: expectedRevision.MediaType,
			Extension: expectedRevision.Extension,
		},
		CreatedAt: expectedRevision.CreatedAt,
	}
}

func assertExistingFile(t *testing.T, actualStoredFile *apiRestaurantFile.StoredFile, expectedFile ExpectedExistingFile) {
	t.Helper()
	assert.NotNil(t, actualStoredFile)
	assert.Equal(t, expectedFile.Id.String(), actualStoredFile.Id.GetValue())
	assert.Equal(t, expectedFile.RevisionId.String(), actualStoredFile.RevisionId.GetValue())
	assertMetadata(t, actualStoredFile.Metadata, expectedFile.Metadata)
	assert.Equal(t, expectedFile.CreatedAt, actualStoredFile.Metadata.CreatedAt.AsTime())
}

func assertMetadata(t *testing.T, actualMetadata *apiRestaurantFile.StoredFile_Metadata, expectedMetadata ExpectedMetadata) {
	t.Helper()
	assert.NotNil(t, actualMetadata)
	assert.Equal(t, expectedMetadata.Size, actualMetadata.Size)
	assert.Equal(t, expectedMetadata.MediaType, actualMetadata.MediaType)
	assert.Equal(t, expectedMetadata.Extension, actualMetadata.Extension)
	assert.NotNil(t, actualMetadata.CreatedAt)
}

func assertGrpcError(t *testing.T, expectedCode codes.Code, actualError error) *status.Status {
	t.Helper()
	assert.NotNil(t, actualError)
	actualStatus, ok := status.FromError(actualError)
	require.True(t, ok, "Expected a gRPC status error")
	require.NotNil(t, actualStatus)
	assert.Equal(t, expectedCode, actualStatus.Code())
	return actualStatus
}

//go:build unit

package file

import (
	context "context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
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
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"

	var generatedFileId *uuid.UUID
	var generatedRevisionId *uuid.UUID
	var storedFileMetadata *FileMetadata
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream := fixture.CreateValidStoreFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupAndRecordSuccessfulStoreFileResponse(t, mockStream, &actualResponse)
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
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream := fixture.CreateValidStoreFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupAndRecordSuccessfulStoreFileResponse(t, mockStream, &actualResponse)
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
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream := fixture.CreateValidStoreFileStream(t, sentFileName, chunks)
	fixture.SetupAndRecordSuccessfulStoreFileResponse(t, mockStream, &actualResponse)
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
	mockStream := fixture.CreateStoreFileStream(t)
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
	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreFileRequest(t, "test.txt"), nil).Times(1)
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

	mockStream := fixture.CreateStoreFileStream(t)
	mockStream.EXPECT().Recv().Return(fixture.CreateMetadataStoreFileRequest(t, "test.txt"), nil).Times(1)
	mockStream.EXPECT().Recv().Return(fixture.CreateChunkStoreFileRequest(t, file), nil).Times(1)
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
	mockStream := fixture.CreateStoreFileStream(t)
	firstRequest := fixture.CreateChunkStoreFileRequest(t, fixture.TextFile())
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
	mockStream := fixture.CreateStoreFileStream(t)
	firstRequest := fixture.CreateMetadataStoreFileRequest(t, "test.txt")
	mockStream.EXPECT().Recv().Return(firstRequest, nil).Times(1)
	secondRequest := fixture.CreateMetadataStoreFileRequest(t, "test2.txt")
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
	mockStream := fixture.CreateValidStoreFileStreamThatAbortsOnFileWrite(t, sentFileName, [][]byte{})

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
	mockStream := fixture.CreateValidStoreFileStreamThatAbortsOnFileWrite(t, sentFileName, [][]byte{sentFile})
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
	mockStream := fixture.CreateValidStoreFileStreamThatAbortsOnFileClose(t, sentFileName, [][]byte{sentFile})
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

func TestDownloadFile_FileIdIsNil(t *testing.T) {
	request := &v1.DownloadFileRequest{
		FileId: nil,
	}

	mockStream := fixture.CreateDownloadFileStream(t)
	sut := createSut(t, nil, nil)
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
			request := &apiRestaurantFile.DownloadFileRequest{
				FileId: &protobuf.Uuid{
					Value: test.uuid,
				},
			}

			mockStream := fixture.CreateDownloadFileStream(t)
			sut := createSut(t, nil, nil)
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
			requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
			request := fixture.CreateDownloadFileRequest(t, requestedFileId)

			mockFileMetadataRepository := &MockFileMetadataRepository{}
			mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(FileMetadata{}, test.notFoundError).Times(1)
			mockFileMetadataRepository.EXPECT().NotFoundError().Return(test.notFoundError).Times(1)

			sut := createSut(t, nil, mockFileMetadataRepository)
			mockStream := fixture.CreateDownloadFileStream(t)
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
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)

	mockFileMetadataRepository := &MockFileMetadataRepository{}
	mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(FileMetadata{}, errors.New("ups..someting went wrong")).Times(1)
	mockFileMetadataRepository.EXPECT().NotFoundError().Return(errors.New("file not found")).Times(1)

	sut := createSut(t, nil, mockFileMetadataRepository)
	mockStream := fixture.CreateDownloadFileStream(t)
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
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
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
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(sendErr).Times(1)
	mockFileRepository := &MockFileRepository{}

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
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
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
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
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := &MockFileRepository{}
	mockFileRepository.EXPECT().ReadFile(mock.Anything, fileId, revisionId).Return(nil, openErr).Times(1)

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
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
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
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := &MockFileRepository{}
	readCloser := ioFixture.CreateReadCloserRanIntoReadError(t, readErr)
	mockFileRepository.EXPECT().ReadFile(mock.Anything, fileId, revisionId).Return(readCloser, nil).Times(1)

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
	requestedFileId, _ := apiProtobuf.ToProtobuf(fileId)
	request := fixture.CreateDownloadFileRequest(t, requestedFileId)
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
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1) // metadata are sent
	mockFileRepository := &MockFileRepository{}
	readCloser := ioFixture.CreateReadCloserRanIntoCloseError(t, file, closeErr)
	mockFileRepository.EXPECT().ReadFile(mock.Anything, fileId, revisionId).Return(readCloser, nil).Times(1)

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
	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockStream.EXPECT().Send(mock.Anything).Return(nil).Times(1)     // metadata are sent
	mockStream.EXPECT().Send(mock.Anything).Return(sendErr).Times(1) // file bytes are sent
	mockFileRepository := &MockFileRepository{}
	mockReader := ioFixture.CreateReadCloser(t, file)
	mockFileRepository.EXPECT().ReadFile(mock.Anything, fileId, revisionId).Return(mockReader, nil).Times(1)

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

	mockFileMetadataRepository := &MockFileMetadataRepository{}
	setupFileMetadataRepositoryToFetchMetadata(t, mockFileMetadataRepository, fileId, fileMetadata)
	mockFileRepository := &MockFileRepository{}
	mockReader := ioFixture.CreateReadCloser(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc)
	mockFileRepository.EXPECT().ReadFile(mock.Anything, fileId, latestedRevision.Id).Return(mockReader, nil).Times(1)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	mockStream := fixture.CreateDownloadFileStream(t)
	mockStream.EXPECT().Send(mock.MatchedBy(func(response *apiRestaurantFile.DownloadFileResponse) bool {
		actualMetadata := response.GetMetadata()
		return actualMetadata.Extension == latestedRevision.Extension &&
			actualMetadata.MediaType == latestedRevision.MediaType &&
			actualMetadata.Size == latestedRevision.Size &&
			actualMetadata.CreatedAt.AsTime().Equal(latestedRevision.CreatedAt)
	})).Return(nil).Times(1)

	actualFile := make([]byte, 0)
	mockStream.EXPECT().Send(mock.Anything).Run(func(response *v1.DownloadFileResponse) {
		actualFile = append(actualFile, response.GetChunk()...)
	}).Return(nil)
	actualError := sut.DownloadFile(request, mockStream)

	assert.Nil(t, actualError)
	assert.Equal(t, fileThatIsBiggerThanTheMaxChunkSizeForGrpc, actualFile)
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

func setupFileMetadataRepositoryToFetchMetadata(t *testing.T, mockFileMetadataRepository *MockFileMetadataRepository, fileId uuid.UUID, fileMetadata FileMetadata) {
	t.Helper()
	mockFileMetadataRepository.EXPECT().FetchFileMetadata(mock.Anything, fileId).Return(fileMetadata, nil).Times(1)
	mockFileMetadataRepository.EXPECT().NotFoundError().Return(errors.New("not expected error")).Times(1)
}

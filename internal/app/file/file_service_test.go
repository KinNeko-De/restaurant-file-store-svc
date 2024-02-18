package file

import (
	context "context"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	ioFixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/io"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"

	var generatedFileId *uuid.UUID
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloserMock(t, [][]byte{sentFile})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t)

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

	assert.Equal(t, generatedFileId.String(), actualResponse.StoredFile.Id.Value)
}

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsExact512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()[0:512]
	sentFileName := "test.pdf"
	expectedSize := uint64(512)
	expectedMediaType := "application/pdf"

	var generatedFileId *uuid.UUID
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, [][]byte{sentFile})
	fixture.SetupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloserMock(t, [][]byte{sentFile})
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t)

	sut := createSut(t, mockFileRepository, mockFileMetadataRepository)
	actualError := sut.StoreFile(mockStream)

	assert.Nil(t, actualError)

	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
}

func TestStoreFile_FileDataIsSentInMultipleChunks_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.PdfFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"
	expectedSize := uint64(51124)
	expectedMediaType := "application/pdf"
	expectedFileExtension := ".pdf"

	var generatedFileId *uuid.UUID
	var actualResponse *v1.StoreFileResponse
	mockStream := fixture.CreateValidFileStream(t, sentFileName, chunks)
	fixture.SetupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := ioFixture.CreateWriterCloserMock(t, chunks)
	mockFileRepository := createFileRepositoryMock(t, fileWriter, &generatedFileId)
	mockFileMetadataRepository := createFileMetadataRepositoryMock(t)

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
}

func createSut(t *testing.T, mockFileRepository *MockFileRepository, mockFileMetadataRepository *MockFileMetadataRepository) FileServiceServer {
	sut := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository
	return sut
}

func createFileRepositoryMock(t *testing.T, fileWriter *ioFixture.MockWriteCloser, generatedFileId **uuid.UUID) *MockFileRepository {
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

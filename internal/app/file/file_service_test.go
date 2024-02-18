package file

import (
	"testing"

	"github.com/google/uuid"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	// target2 := [6]int{}
	target := make([]int, 6)
	first := []int{1, 2, 3}
	second := []int{4, 5}
	size := 0
	copy(target[size:], first)
	size += len(first)
	copy(target[size:], second)
	size += len(second)
	target = target[:size]
}

/*
func TestFile(t *testing.T) {
	file, _ := os.Open("protobuf.pdf")
	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	_, _ = bufio.NewReader(file).Read(bs)
	variable := ""
	for _, b := range bs {
		variable += strconv.FormatUint(uint64(b), 10)
		variable += ", "
	}

	output, _ := os.Create("test.txt")
	output.WriteString(variable)
	output.Close()
}
*/

func TestStoreFile_FileDataIsSentInOneChunk_FileSizeIsSmallerThan512SniffBytes(t *testing.T) {
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"

	var generatedFileId *uuid.UUID
	var actualResponse *v1.StoreFileResponse
	mockStream := createValidFileStream(t, sentFileName, [][]byte{sentFile})
	setupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := createWriterCloserMock(t, [][]byte{sentFile})
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
	mockStream := createValidFileStream(t, sentFileName, [][]byte{sentFile})
	setupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := createWriterCloserMock(t, [][]byte{sentFile})
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
	chunks := splitIntoChunks(sentFile, 256)
	sentFileName := "test.pdf"
	expectedSize := uint64(51124)
	expectedMediaType := "application/pdf"
	expectedFileExtension := ".pdf"

	var generatedFileId *uuid.UUID
	var actualResponse *v1.StoreFileResponse
	mockStream := createValidFileStream(t, sentFileName, chunks)
	setupSuccessfulResponse(t, mockStream, &actualResponse)
	fileWriter := createWriterCloserMock(t, chunks)
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

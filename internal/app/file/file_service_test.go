package file

import (
	"io"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/internal/testing/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArray(t *testing.T) {
	//target := [6]int{}
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

func TestStoreFile_FileDataIsSentInOneChunk(t *testing.T) {
	sentFile := fixture.TextFile()
	sentFileName := "test.txt"
	expectedSize := uint64(4)
	expectedMediaType := "text/plain; charset=utf-8"
	expectedFileExtension := ".txt"
	expectedRevision := int64(1)

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
	assert.Equal(t, expectedRevision, actualResponse.StoredFile.Revision)
	assert.NotNil(t, actualResponse.StoredFileMetadata)
	assert.Equal(t, expectedSize, actualResponse.StoredFileMetadata.Size)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFileMetadata.MediaType)
	assert.Equal(t, expectedFileExtension, actualResponse.StoredFileMetadata.Extension)
	assert.NotNil(t, actualResponse.StoredFileMetadata.CreatedAt)

	assert.Equal(t, generatedFileId.String(), actualResponse.StoredFile.Id.Value)

}

func TestStoreFile_PdfFile(t *testing.T) {
	mockStream := NewFileService_StoreFileServer(t)

	var metadata = &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Name{
			Name: "test.pdf",
		},
	}
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)
	var chunk = &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Chunk{
			Chunk: fixture.PdfFile(),
		},
	}
	mockStream.EXPECT().Recv().Return(chunk, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(mock.MatchedBy(func(response *v1.StoreFileResponse) bool {
		// TODO ID not null and timestamp not null
		return response.GetStoredFile().GetRevision() == 1 &&
			response.GetStoredFileMetadata().GetSize() == 51124 &&
			response.GetStoredFileMetadata().GetMediaType() == "application/pdf" &&
			response.GetStoredFileMetadata().GetExtension() == ".pdf"
	})).Return(nil).Times(1)

	server := FileServiceServer{}
	actualError := server.StoreFile(mockStream)
	assert.Nil(t, actualError)
}

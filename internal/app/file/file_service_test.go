package file

import (
	"context"
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

func TestStoreFile_TextFile(t *testing.T) {
	sentFile := fixture.TextFile()

	expected := CreateExpectedResponse(4, "text/plain; charset=utf-8", ".txt")
	mockStream := CreateMockStream(t, "test.txt", [][]byte{sentFile}, expected)
	fileWriter := &MockWriteCloser{}
	fileWriter.On("Write", mock.IsType(make([]byte, 0))).Return(nil)
	fileWriter.On("Close").Return(nil)
	mockFileRepository := &MockFileRepository{}
	mockFileRepository.EXPECT().CreateFile(mock.Anything, mock.IsType(uuid.New()), 0).Return(fileWriter, nil).Times(1)
	mockFileMetadataRepository := &MockFileMetadataRepository{}

	server := FileServiceServer{}
	FileRepositoryInstance = mockFileRepository
	FileMetadataRepositoryInstance = mockFileMetadataRepository

	actualError := server.StoreFile(mockStream)

	assert.Nil(t, actualError)
	// assert.Equal(t, sentFile, storedFile)
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

func CreateMockStream(t *testing.T, fileName string, fileChunks [][]byte, argumentMatcher interface{}) *FileService_StoreFileServer {
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
	mockStream.EXPECT().SendAndClose(argumentMatcher).Return(nil).Times(1)

	return mockStream
}

func CreateExpectedResponse(expectedSize uint64, expectedMediaType string, expectedFileExtension string) interface{} {
	return mock.MatchedBy(func(response *v1.StoreFileResponse) bool {
		// TODO ID not null and timestamp not null
		return response.GetStoredFile().GetRevision() == 1 &&
			response.GetStoredFileMetadata().GetSize() == expectedSize &&
			response.GetStoredFileMetadata().GetMediaType() == expectedMediaType &&
			response.GetStoredFileMetadata().GetExtension() == expectedFileExtension
	})
}

type MockWriteCloser struct {
	mock.Mock
}

func (m *MockWriteCloser) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return len(p), args.Error(0)
}

func (m *MockWriteCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

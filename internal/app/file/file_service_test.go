package file

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"testing"

	v1 "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
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

func TestFile(t *testing.T) {
	file, _ := os.Open("test.txt")
	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	_, _ = bufio.NewReader(file).Read(bs)
	variable := ""
	for _, b := range bs {
		variable += strconv.FormatUint(uint64(b), 10)
		variable += ", "
	}
}

func TestStoreFile(t *testing.T) {
	mockStream := NewFileService_StoreFileServer(t)

	var metadata = &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Name{
			Name: "test.txt",
		},
	}
	mockStream.EXPECT().Recv().Return(metadata, nil).Times(1)
	var chunk = &v1.StoreFileRequest{
		File: &v1.StoreFileRequest_Chunk{
			Chunk: []byte{116, 101, 115, 116}, // TODO test.txt with content "test". extract to fixture
		},
	}
	mockStream.EXPECT().Recv().Return(chunk, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	mockStream.EXPECT().SendAndClose(mock.MatchedBy(func(response *v1.StoreFileResponse) bool {
		// TODO ID not null and timestamp not null
		return response.GetStoredFile().GetRevision() == 1 &&
			response.GetStoredFileMetadata().GetSize() == 4 &&
			response.GetStoredFileMetadata().GetMediaType() == "text/plain; charset=utf-8" &&
			response.GetStoredFileMetadata().GetExtension() == ".txt"
	})).Return(nil).Times(1)
	/*
		mockStream.EXPECT().SendAndClose(&v1.StoreFileResponse{
			StoredFile: &v1.StoredFile{
				Id:       &apiProtobuf.Uuid{Value: "test"},
				Revision: 1,
			},
			StoredFileMetadata: &v1.StoredFileMetadata{
				CreatedAt: timestamppb.Now(),
				Size:      4,
				MediaType: "text/plain; charset=utf-8",
				Extension: ".txt",
			},
		})
	*/

	server := FileServiceServer{}
	server.StoreFile(mockStream)
}

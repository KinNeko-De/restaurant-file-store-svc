//go:build acceptance

package main

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStoreFile_StoreRevision_DownloadSpecificRevision_DownloadLatestRevision(t *testing.T) {
	conn, dialErr := grpc.Dial("localhost:42985", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, dialErr)
	defer conn.Close()
	client := apiRestaurantFile.NewFileServiceClient(conn)

	fileName := "test.txt"
	fileName2 := "hello.txt"
	expectedExtension := ".txt"
	expectedMediaType := "text/plain; charset=utf-8"
	sentFile := fixture.TextFile()
	sentFile2 := fixture.HelloWorldFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	chunks2 := fixture.SplitIntoChunks(sentFile2, 256)
	expectedSize := uint64(len(sentFile))
	expectedSize2 := uint64(len(sentFile2))

	storeFileResponse := CreateFile(t, client, fileName, chunks, expectedExtension, expectedMediaType, expectedSize)
	storeRevisionResponse := StoreRevision(t, client, storeFileResponse.StoredFile.Id, fileName2, chunks2, expectedExtension, expectedMediaType, expectedSize2)
	DownloadLatestRevision(t, client, storeRevisionResponse, sentFile2)
	DownloadPreviousRevision(t, client, storeFileResponse, sentFile)
}

func CreateFile(t *testing.T, client apiRestaurantFile.FileServiceClient, fileName string, fileChunks [][]byte, expectedExtension string, expectedMediaType string, expectedSize uint64) *apiRestaurantFile.StoreFileResponse {
	uploadStream, uploadErr := client.StoreFile(context.Background())
	require.Nil(t, uploadErr)

	var metadata = &apiRestaurantFile.StoreFileRequest{
		Part: &apiRestaurantFile.StoreFileRequest_StoreFile{
			StoreFile: &apiRestaurantFile.StoreFile{
				Name: fileName,
			},
		},
	}
	uploadStream.Send(metadata)

	for _, chunk := range fileChunks {
		var chunkRequest = &apiRestaurantFile.StoreFileRequest{
			Part: &apiRestaurantFile.StoreFileRequest_Chunk{
				Chunk: chunk,
			},
		}
		uploadStream.Send(chunkRequest)
	}

	actualResponse, uploadErr := uploadStream.CloseAndRecv()
	require.Nil(t, uploadErr)
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFile.Metadata)
	assert.NotNil(t, actualResponse.StoredFile.Metadata.CreatedAt)
	assert.Equal(t, expectedExtension, actualResponse.StoredFile.Metadata.Extension)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFile.Metadata.MediaType)
	assert.Equal(t, expectedSize, actualResponse.StoredFile.Metadata.Size)

	return actualResponse
}

func StoreRevision(t *testing.T, client apiRestaurantFile.FileServiceClient, fileId *apiProtobuf.Uuid, fileName string, fileChunks [][]byte, expectedExtension string, expectedMediaType string, expectedSize uint64) *apiRestaurantFile.StoreFileResponse {
	uploadStream, uploadErr := client.StoreRevision(context.Background())
	require.Nil(t, uploadErr, "failed to store revision: %w", uploadErr)

	var metadata = &apiRestaurantFile.StoreRevisionRequest{
		Part: &apiRestaurantFile.StoreRevisionRequest_StoreRevision{
			StoreRevision: &apiRestaurantFile.StoreRevision{
				FileId: fileId,
				StoreFile: &apiRestaurantFile.StoreFile{
					Name: fileName,
				},
			},
		},
	}
	uploadStream.Send(metadata)

	for _, chunk := range fileChunks {
		var chunkRequest = &apiRestaurantFile.StoreRevisionRequest{
			Part: &apiRestaurantFile.StoreRevisionRequest_Chunk{
				Chunk: chunk,
			},
		}
		uploadStream.Send(chunkRequest)
	}

	actualResponse, uploadErr := uploadStream.CloseAndRecv()
	require.Nil(t, uploadErr)
	assert.NotNil(t, actualResponse)
	assert.NotNil(t, actualResponse.StoredFile)
	assert.Equal(t, fileId.Value, actualResponse.StoredFile.Id.Value)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.RevisionId)
	assert.NotNil(t, actualResponse.StoredFile.Metadata)
	assert.NotNil(t, actualResponse.StoredFile.Metadata.CreatedAt)
	assert.Equal(t, expectedExtension, actualResponse.StoredFile.Metadata.Extension)
	assert.Equal(t, expectedMediaType, actualResponse.StoredFile.Metadata.MediaType)
	assert.Equal(t, expectedSize, actualResponse.StoredFile.Metadata.Size)

	return actualResponse
}

func DownloadLatestRevision(t *testing.T, client apiRestaurantFile.FileServiceClient, storeFileResponse *apiRestaurantFile.StoreFileResponse, exptectedFile []byte) {
	downloadStream, downloadErr := client.DownloadFile(context.Background(), &apiRestaurantFile.DownloadFileRequest{
		FileId: storeFileResponse.StoredFile.Id,
	})
	require.Nil(t, downloadErr)
	require.NotNil(t, downloadStream)
	downloadResponse, err := downloadStream.Recv()
	require.Nil(t, err)
	require.NotNil(t, downloadResponse)
	downloadedFile := downloadResponse.GetStoredFile()
	receivedFile := RecordDownloadedFile(t, downloadStream)

	require.NotNil(t, downloadedFile)
	assert.NotNil(t, downloadedFile.Metadata.CreatedAt)
	assert.WithinDuration(t, storeFileResponse.StoredFile.Metadata.CreatedAt.AsTime(), downloadedFile.Metadata.CreatedAt.AsTime(), time.Millisecond)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.Extension, downloadedFile.Metadata.Extension)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.MediaType, downloadedFile.Metadata.MediaType)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.Size, downloadedFile.Metadata.Size)
	assert.Equal(t, exptectedFile, receivedFile)
}

func DownloadPreviousRevision(t *testing.T, client apiRestaurantFile.FileServiceClient, storeFileResponse *apiRestaurantFile.StoreFileResponse, sentFile []byte) {
	downloadStream, downloadErr := client.DownloadRevision(context.Background(), &apiRestaurantFile.DownloadRevisionRequest{
		FileId:     storeFileResponse.StoredFile.Id,
		RevisionId: storeFileResponse.StoredFile.RevisionId,
	})
	require.Nil(t, downloadErr)
	require.NotNil(t, downloadStream)
	downloadResponse, err := downloadStream.Recv()
	require.Nil(t, err)
	require.NotNil(t, downloadResponse)
	downloadedFile := downloadResponse.GetStoredFile()
	receivedFile := RecordDownloadedFile(t, downloadStream)

	require.NotNil(t, downloadedFile)
	assert.NotNil(t, downloadedFile.Metadata.CreatedAt)
	assert.WithinDuration(t, storeFileResponse.StoredFile.Metadata.CreatedAt.AsTime(), downloadedFile.Metadata.CreatedAt.AsTime(), time.Millisecond)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.Extension, downloadedFile.Metadata.Extension)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.MediaType, downloadedFile.Metadata.MediaType)
	assert.Equal(t, storeFileResponse.StoredFile.Metadata.Size, downloadedFile.Metadata.Size)
	assert.Equal(t, sentFile, receivedFile)

}

func RecordDownloadedFile(t *testing.T, downloadStream apiRestaurantFile.FileService_DownloadFileClient) []byte {
	var receivedFile []byte
	for {
		downloadResponse, err := downloadStream.Recv()
		if err == io.EOF {
			break
		}
		require.Nil(t, err)
		chunk := downloadResponse.GetChunk()
		require.NotNil(t, chunk)
		receivedFile = append(receivedFile, chunk...)
	}
	return receivedFile
}

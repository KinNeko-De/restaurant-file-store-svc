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
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStoreFile(t *testing.T) {
	fileName := "test.txt"
	expectedExtension := ".txt"
	expectedMediaType := "text/plain; charset=utf-8"
	sentFile := fixture.TextFile()
	expectedSize := uint64(len(sentFile))
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	startTime := time.Now()

	conn, dialErr := grpc.Dial("localhost:42985", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, dialErr)
	defer conn.Close()

	client := apiRestaurantFile.NewFileServiceClient(conn)
	ctx := context.Background()
	uploadStream, uploadErr := client.StoreFile(ctx)
	require.Nil(t, uploadErr)

	var metadata = &apiRestaurantFile.StoreFileRequest{
		Part: &apiRestaurantFile.StoreFileRequest_StoreFile{
			StoreFile: &apiRestaurantFile.StoreFile{
				Name: fileName,
			},
		},
	}
	uploadStream.Send(metadata)

	for _, chunk := range chunks {
		var chunkRequest = &apiRestaurantFile.StoreFileRequest{
			Part: &apiRestaurantFile.StoreFileRequest_Chunk{
				Chunk: chunk,
			},
		}
		uploadStream.Send(chunkRequest)
	}

	actualResponse, uploadErr := uploadStream.CloseAndRecv()
	require.Nil(t, uploadErr)
	duration := time.Since(startTime)
	t.Logf("Call duration: %s", duration)

	assert.NotNil(t, actualResponse)

	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.RevisionId)

	downloadStream, downloadErr := client.DownloadFile(ctx, &apiRestaurantFile.DownloadFileRequest{
		FileId: actualResponse.StoredFile.Id,
	})

	require.Nil(t, downloadErr)
	require.NotNil(t, downloadStream)

	downloadResponse, err := downloadStream.Recv()
	require.Nil(t, err)
	require.NotNil(t, downloadResponse)

	downloadMetadata := downloadResponse.GetMetadata()

	require.NotNil(t, downloadMetadata)
	assert.NotNil(t, downloadMetadata.CreatedAt)
	assert.WithinDuration(t, actualResponse.StoredFileMetadata.CreatedAt.AsTime(), downloadMetadata.CreatedAt.AsTime(), time.Millisecond)
	assert.Equal(t, expectedExtension, downloadMetadata.Extension)
	assert.Equal(t, actualResponse.StoredFileMetadata.Extension, downloadMetadata.Extension)
	assert.Equal(t, expectedMediaType, downloadMetadata.MediaType)
	assert.Equal(t, actualResponse.StoredFileMetadata.MediaType, downloadMetadata.MediaType)
	assert.Equal(t, expectedSize, downloadMetadata.Size)
	assert.Equal(t, actualResponse.StoredFileMetadata.Size, downloadMetadata.Size)

	var receivedFile []byte
	for {
		downloadResponse, err := downloadStream.Recv()
		if err != io.EOF {
			break
		}
		chunk := downloadResponse.GetChunk()
		require.NotNil(t, chunk)
		receivedFile = append(receivedFile, chunk...)
	}

	assert.Equal(t, sentFile, receivedFile)
}

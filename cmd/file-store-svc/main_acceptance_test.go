//go:build acceptance

package main

import (
	"context"
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
	sentFile := fixture.TextFile()
	chunks := fixture.SplitIntoChunks(sentFile, 256)
	startTime := time.Now()

	conn, dialErr := grpc.Dial("localhost:42985", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, dialErr)
	defer conn.Close()

	client := apiRestaurantFile.NewFileServiceClient(conn)
	ctx := context.Background()
	stream, err := client.StoreFile(ctx)
	require.Nil(t, err)

	var metadata = &apiRestaurantFile.StoreFileRequest{
		Part: &apiRestaurantFile.StoreFileRequest_StoreFile{
			StoreFile: &apiRestaurantFile.StoreFile{
				Name: fileName,
			},
		},
	}
	stream.Send(metadata)

	for _, chunk := range chunks {
		var chunkRequest = &apiRestaurantFile.StoreFileRequest{
			Part: &apiRestaurantFile.StoreFileRequest_Chunk{
				Chunk: chunk,
			},
		}
		stream.Send(chunkRequest)
	}

	actualResponse, err := stream.CloseAndRecv()
	require.Nil(t, err)
	duration := time.Since(startTime)
	t.Logf("Call duration: %s", duration)

	assert.NotNil(t, actualResponse)

	assert.NotNil(t, actualResponse.StoredFile)
	assert.NotNil(t, actualResponse.StoredFile.Id)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.Id)
	assert.NotNil(t, actualResponse.StoredFile.RevisionId)
	assert.NotEqual(t, uuid.Nil, actualResponse.StoredFile.RevisionId)

	downloadStream, downloadError := client.DownloadFile(ctx, &apiRestaurantFile.DownloadFileRequest{
		FileId: actualResponse.StoredFile.Id,
	})

	require.Nil(t, downloadError)
	require.NotNil(t, downloadStream)

	metadata2, err2 := downloadStream.Recv()
	require.Nil(t, err2)
	require.NotNil(t, metadata2)

	// TODO assert more
}

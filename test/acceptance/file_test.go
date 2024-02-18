package acceptance

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"

	fixture "github.com/kinneko-de/restaurant-file-store-svc/test/testing/file"
)

func TestStoreFile(t *testing.T) {
	requestFileMetadata := &apiRestaurantFile.StoreFileRequest{
		File: &apiRestaurantFile.StoreFileRequest_Name{
			Name: "test.txt",
		},
	}
	requestFile := &apiRestaurantFile.StoreFileRequest{
		File: &apiRestaurantFile.StoreFileRequest_Chunk{
			Chunk: fixture.TextFile(),
		},
	}

	conn, err := grpc.Dial("localhost:3110", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	defer conn.Close()
	client := apiRestaurantFile.NewFileServiceClient(conn)
	stream, err := client.StoreFile(context.Background())
	require.Nil(t, err)
	stream.Send(requestFileMetadata)
	stream.Send(requestFile)
	acutalResponse, err := stream.CloseAndRecv()

	assert.Nil(t, err, "error should be nil but got %w", err)
	assert.NotNil(t, acutalResponse)
	assert.NotNil(t, acutalResponse.StoredFile)
	assert.NotNil(t, acutalResponse.StoredFile.Id)
	assert.NotEqual(t, uuid.Nil.String(), acutalResponse.StoredFile.Id.Value)
	assert.NotNil(t, acutalResponse.StoredFile.RevisionId)
	assert.NotEqual(t, uuid.Nil.String(), acutalResponse.StoredFile.RevisionId.Value)
}

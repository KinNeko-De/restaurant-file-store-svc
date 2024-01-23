package file

import (
	context "context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type GoogleCloudStorage struct {
	Client *storage.Client
}

// initialize implements FileStorage.
func (g *GoogleCloudStorage) Initialize() error {
	return g.CreateClient()
}

func (g *GoogleCloudStorage) CreateClient() error {
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx)
	if err != nil {
		return err
	}
	g.Client = client
	return nil
}

func (g *GoogleCloudStorage) CreateFile(ctx context.Context, revisionId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	bucket := g.Client.Bucket("kinneko-de")
	object := bucket.Object(revisionId.String())
	writer := object.NewWriter(ctx)
	writer.ChunkSize = chunkSize
	// TODO: checksum control
	return writer, nil
}

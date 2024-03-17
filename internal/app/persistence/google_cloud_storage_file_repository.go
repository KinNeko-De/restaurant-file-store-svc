package persistence

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type GoogleCloudStorageFileRepository struct {
	Client *storage.Client
}

func (g *GoogleCloudStorageFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	bucket := g.Client.Bucket("kinneko-de")
	objectname := fileId.String() + "/" + revisionId.String()
	object := bucket.Object(objectname)
	writer := object.NewWriter(ctx)
	writer.ChunkSize = chunkSize
	// TODO: checksum control
	return writer, nil
}

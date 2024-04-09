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

func (g *GoogleCloudStorageFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.WriteCloser, error) {
	bucket := g.Client.Bucket("kinneko-de")
	objectname := fileId.String() + "/" + revisionId.String()
	object := bucket.Object(objectname)
	writer := object.NewWriter(ctx)
	// TODO: chunksize must be set to slightly larger than the maximum size of the file, for that we need to know the size of the file before writing it
	// TODO: checksum control
	return writer, nil
}

func (g *GoogleCloudStorageFileRepository) ReadFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.ReadCloser, error) {
	bucket := g.Client.Bucket("kinneko-de")
	objectname := fileId.String() + "/" + revisionId.String()
	object := bucket.Object(objectname)
	reader, err := object.NewReader(ctx)
	return reader, err
}

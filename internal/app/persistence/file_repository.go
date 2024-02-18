package persistence

import (
	"context"
	"io"
	"os"

	"github.com/google/uuid"
)

type LocalStorageRepository struct {
}

func (*LocalStorageRepository) CreateFile(ctx context.Context, fileId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	return os.Create(fileId.String())
}

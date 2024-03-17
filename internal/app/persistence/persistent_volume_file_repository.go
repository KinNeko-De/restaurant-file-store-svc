package persistence

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type PersistentVolumeFileRepository struct {
}

func (g *PersistentVolumeFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	panic("not implemented")
}

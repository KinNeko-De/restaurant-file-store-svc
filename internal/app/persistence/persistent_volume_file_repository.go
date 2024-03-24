package persistence

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
)

type PersistentVolumeFileRepository struct {
}

func (g *PersistentVolumeFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	pathName := fileId.String()
	pathAndFile := path.Join(pathName, revisionId.String())
	err := os.MkdirAll(pathName, os.ModePerm)
	if err != nil {
		return nil, err
	}
	writer, err := os.Create(pathAndFile)
	return writer, err
}

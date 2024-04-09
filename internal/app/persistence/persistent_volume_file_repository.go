package persistence

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
)

type PersistentVolumeFileRepository struct {
	StoragePath string
}

func (g *PersistentVolumeFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.WriteCloser, error) {
	fileFolder := path.Join(g.StoragePath, fileId.String())
	fileLocation := path.Join(fileFolder, revisionId.String())
	err := os.MkdirAll(fileFolder, os.ModePerm)
	if err != nil {
		return nil, err
	}
	writer, err := os.Create(fileLocation)
	return writer, err
}

func (g *PersistentVolumeFileRepository) ReadFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.ReadCloser, error) {
	fileLocation := path.Join(g.StoragePath, fileId.String(), revisionId.String())
	reader, err := os.Open(fileLocation)
	return reader, err
}

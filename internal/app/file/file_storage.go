package file

import (
	context "context"
	"io"

	"github.com/google/uuid"
)

var (
	Storage FileStorage = &GoogleCloudStorage{}
)

type FileStorage interface {
	Initialize() error
	CreateFile(ctx context.Context, fileId uuid.UUID, chunkSize int) (io.WriteCloser, error)
}

package file

import (
	context "context"
	"io"

	"github.com/google/uuid"
)

var (
	FileRepositoryInstance FileRepository
)

type FileRepository interface {
	CreateFile(ctx context.Context, fileId uuid.UUID, chunkSize int) (io.WriteCloser, error)
}

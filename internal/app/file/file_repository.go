package file

import (
	context "context"
	"io"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

var (
	FileRepositoryInstance FileRepository = &persistence.GoogleCloudStorage{}
)

type FileRepository interface {
	Initialize() error
	CreateFile(ctx context.Context, fileId uuid.UUID, chunkSize int) (io.WriteCloser, error)
}

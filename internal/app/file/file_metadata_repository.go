package file

import (
	context "context"

	"github.com/google/uuid"
)

var (
	FileMetadataRepositoryInstance FileMetadataRepository
)

type FileMetadataRepository interface {
	StoreFileMetadata(ctx context.Context, fileMetadata FileMetadata) error
	FetchFileMetadata(ctx context.Context, fileId uuid.UUID) (FileMetadata, error)
}

package file

import context "context"

var (
	FileMetadataRepositoryInstance FileMetadataRepository
)

type FileMetadataRepository interface {
	StoreFileMetadata(ctx context.Context, fileMetadata FileMetadata) error
}

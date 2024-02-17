package file

import context "context"

var (
	FileMetadataRepositoryInstance FileMetadataRepository
)

type FileMetadataRepository interface {
	CreateFileMetadata(context.Context, FileMetadata) error
}

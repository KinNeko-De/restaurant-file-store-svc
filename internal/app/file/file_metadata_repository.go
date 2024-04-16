package file

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	FileMetadataRepositoryInstance FileMetadataRepository
)

type FileMetadataRepository interface {
	StoreFileMetadata(ctx context.Context, fileMetadata FileMetadata) error
	FetchFileMetadata(ctx context.Context, fileId uuid.UUID) (FileMetadata, error)
	NotFoundError() error
}

func fetchMetadata(ctx context.Context, requestedFileId uuid.UUID, scopedLogger zerolog.Logger) (FileMetadata, error) {
	fileMetadata, err := FileMetadataRepositoryInstance.FetchFileMetadata(ctx, requestedFileId)
	if errors.Is(err, FileMetadataRepositoryInstance.NotFoundError()) {
		scopedLogger.Err(err).Msg("file not found")
		return FileMetadata{}, status.Error(codes.NotFound, "file with id '"+requestedFileId.String()+"' not found.")
	}
	if err != nil {
		scopedLogger.Err(err).Msg("error fetching file metadata")
		return FileMetadata{}, status.Error(codes.Internal, "error fetching file metadata. please retry the request")
	}
	return fileMetadata, nil
}

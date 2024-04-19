package file

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	FileRepositoryInstance FileRepository = nil
)

type FileRepository interface {
	CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.WriteCloser, error)
	OpenFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (io.ReadCloser, error)
}

func writeFile(stream ChunckStream, ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID) (uint64, []byte, error) {
	fileWriter, err := FileRepositoryInstance.CreateFile(ctx, fileId, revisionId)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to create file")
		return 0, nil, status.Error(codes.Internal, "failed to create file. please retry the request")
	}

	totalFileSize, sniff, err := receiveChunks(stream, fileWriter)
	if err != nil {
		return 0, nil, err
	}

	closeErr := fileWriter.Close()
	if closeErr != nil {
		logger.Logger.Err(closeErr).Msg("failed to close file")
		return 0, nil, status.Error(codes.Internal, "failed to close file. please retry the request")
	}

	return totalFileSize, sniff, nil
}

package file

import (
	"io"

	"github.com/google/uuid"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func sendFile(stream apiRestaurantFile.FileService_DownloadFileServer, requestedFileId uuid.UUID, revisionId uuid.UUID, scopedLogger zerolog.Logger) error {
	fileReader, err := FileRepositoryInstance.OpenFile(stream.Context(), requestedFileId, revisionId)
	if err != nil { // if the file is not found, we have an internal error in consistence of our data. that information should not be exposed to the client
		scopedLogger.Err(err).Msg("error reading file")
		return status.Error(codes.Internal, "error reading file. please retry the request")
	}
	err = sendChunks(fileReader, stream)
	if err != nil {
		scopedLogger.Err(err).Msg("error sending file chunks")
		return status.Error(codes.Internal, "error sending file chunks. please retry the request")
	}
	err = fileReader.Close()
	if err != nil {
		scopedLogger.Err(err).Msg("closing file fail. error is ignored")
	}
	return nil
}

func sendChunks(fileReader io.ReadCloser, stream apiRestaurantFile.FileService_DownloadFileServer) error {
	maxSizeToRead := make([]byte, 16*1024)
	for {
		readBytes, err := fileReader.Read(maxSizeToRead)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&apiRestaurantFile.DownloadFileResponse{
			Part: &apiRestaurantFile.DownloadFileResponse_Chunk{
				Chunk: maxSizeToRead[:readBytes],
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func sendMetadata(stream apiRestaurantFile.FileService_DownloadFileServer, revision Revision) error {
	return stream.Send(&apiRestaurantFile.DownloadFileResponse{
		Part: &apiRestaurantFile.DownloadFileResponse_Metadata{
			Metadata: &apiRestaurantFile.StoredFileMetadata{
				CreatedAt: timestamppb.New(revision.CreatedAt),
				Size:      revision.Size,
				MediaType: revision.MediaType,
				Extension: revision.Extension,
			},
		},
	})
}

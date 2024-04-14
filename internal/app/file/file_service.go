package file

import (
	"context"
	"errors"
	"io"
	"reflect"

	"github.com/google/uuid"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ apiRestaurantFile.FileServiceServer = &FileServiceServer{}
)

type FileServiceServer struct {
	apiRestaurantFile.UnimplementedFileServiceServer
}

func (s *FileServiceServer) StoreFile(stream apiRestaurantFile.FileService_StoreFileServer) error {
	storeFile, err := receiveMetadata(stream)
	if err != nil {
		return err
	}

	createdFileMetadata, err := createFile(stream, storeFile.Name)
	if err != nil {
		return err
	}

	response, err := createStoreFileResponse(createdFileMetadata)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to to create response")
		return status.Error(codes.Internal, "failed to create response. please retry the request")
	}

	err = stream.SendAndClose(response)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to send response")
		return status.Error(codes.Internal, "failed to send response. please retry the request")
	}

	return nil
}

func createFile(stream apiRestaurantFile.FileService_StoreFileServer, fileName string) (*FileMetadata, error) {
	fileId := uuid.New()
	revisionId := uuid.New()

	totalFileSize, sniff, err := writeFile(stream, fileId, revisionId)
	if err != nil {
		return nil, err
	}

	createdRevision := newRevision(revisionId, fileName, totalFileSize, sniff)
	createdFileMetadata := newFileMetadata(fileId, createdRevision)

	err = FileMetadataRepositoryInstance.StoreFileMetadata(stream.Context(), createdFileMetadata)
	return &createdFileMetadata, err
}

func writeFile(stream apiRestaurantFile.FileService_StoreFileServer, fileId uuid.UUID, revisionId uuid.UUID) (uint64, []byte, error) {
	fileWriter, err := FileRepositoryInstance.CreateFile(stream.Context(), fileId, revisionId)
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

func receiveChunks(stream apiRestaurantFile.FileService_StoreFileServer, f io.WriteCloser) (uint64, []byte, error) {
	var totalFileSize uint64 = 0
	var sniffByteCount uint64 = 0
	sniff := make([]byte, sniffSize)
	for {
		finished, chunk, err := receiveChunk(stream)
		if err != nil {
			return 0, nil, err
		}
		if finished {
			break
		}
		totalFileSize += uint64(len(chunk))

		if sniffByteCount < sniffSize {
			missingBytes := min(sniffSize-sniffByteCount, uint64(len(chunk)))
			copy(sniff[sniffByteCount:], chunk[:missingBytes])
			sniffByteCount += missingBytes
		}

		_, err = f.Write(chunk)
		if err != nil {
			logger.Logger.Err(err).Msg("failed to write chunk to file")
			return 0, nil, status.Error(codes.Internal, "failed to write file. please retry the request")
		}
	}

	return totalFileSize, sniff[:sniffByteCount], nil
}

func receiveMetadata(stream apiRestaurantFile.FileService_StoreFileServer) (*apiRestaurantFile.StoreFile, error) {
	firstRequest, err := stream.Recv()
	if err != nil {
		logger.Logger.Err(err).Msg("receiving message failed")
		return nil, status.Errorf(codes.Internal, "receiving message failed. please retry the request")
	}

	msg := firstRequest.GetStoreFile()
	if msg == nil {
		return nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_Name' expected. Actual value is "+reflect.TypeOf(firstRequest.Part).String()+".")
	}
	return msg, nil
}

func receiveChunk(stream apiRestaurantFile.FileService_StoreFileServer) (bool, []byte, error) {
	request, err := stream.Recv()
	if err == io.EOF {
		return true, nil, nil
	}
	if err != nil {
		logger.Logger.Err(err).Msg("failed to receive chunk")
		return false, nil, status.Errorf(codes.Internal, "failed to receive chunk. please retry the request")
	}

	msg := request.GetChunk()
	if msg == nil {
		return false, nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_Chunk' expected. Actual value is "+reflect.TypeOf(request.Part).String()+".")
	}
	return false, msg, nil
}

func createStoreFileResponse(createdFileMetadata *FileMetadata) (*apiRestaurantFile.StoreFileResponse, error) {
	revison := createdFileMetadata.LatestRevision()

	fileUuid, err := apiProtobuf.ToProtobuf(createdFileMetadata.Id)
	if err != nil {
		return nil, err
	}

	revisionUuid, err := apiProtobuf.ToProtobuf(revison.Id)
	if err != nil {
		return nil, err
	}

	var response = &apiRestaurantFile.StoreFileResponse{
		StoredFile: &apiRestaurantFile.StoredFile{
			Id:         fileUuid,
			RevisionId: revisionUuid,
		},
		StoredFileMetadata: &apiRestaurantFile.StoredFileMetadata{
			CreatedAt: timestamppb.New(revison.CreatedAt),
			Size:      revison.Size,
			MediaType: revison.MediaType,
			Extension: revison.Extension,
		},
	}
	return response, nil
}

func (s *FileServiceServer) DownloadFile(request *apiRestaurantFile.DownloadFileRequest, stream apiRestaurantFile.FileService_DownloadFileServer) error {
	requestedFileId, err := getRequestedFileId(request)
	if err != nil {
		return err
	}
	scopedLogger := logger.Logger.With().Str("fileId", requestedFileId.String()).Logger()

	fileMetadata, err := fetchMetadata(stream.Context(), requestedFileId, scopedLogger)
	if err != nil {
		return err
	}
	revision := fileMetadata.LatestRevision()
	err = sendMetadata(stream, revision)
	if err != nil {
		scopedLogger.Err(err).Msg("error sending file metadata")
		return status.Error(codes.Internal, "error sending file metadata. please retry the request")
	}

	err = sendFile(stream, requestedFileId, revision.Id, scopedLogger)
	if err != nil {
		return err
	}

	return nil
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

func getRequestedFileId(request *apiRestaurantFile.DownloadFileRequest) (uuid.UUID, error) {
	requested := request.GetFileId()
	if requested == nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "fileId is mandatory. Please provide a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}
	fileId, err := apiProtobuf.ToUuid(requested)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "fileId '"+requested.String()+"' is not a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}
	return fileId, nil
}

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

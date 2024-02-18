package file

import (
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/google/uuid"
	apiProtobuf "github.com/kinneko-de/api-contract/golang/kinnekode/protobuf"
	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
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

const sniffSize = 512 // defined by the net/http package

func (s *FileServiceServer) StoreFile(stream apiRestaurantFile.FileService_StoreFileServer) error {
	metaData, err := receiveMetadata(stream)
	if err != nil {
		return err
	}

	totalFileSize, sniff, fileId, revisionId, err := writeFile(stream)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(sniff)
	extension := filepath.Ext(metaData.Name)
	createdAt := time.Now().UTC()

	// TODO Store file metadata

	response, err := createStoreFileResponse(fileId, revisionId, totalFileSize, contentType, extension, createdAt)
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

func writeFile(stream apiRestaurantFile.FileService_StoreFileServer) (uint64, []byte, uuid.UUID, uuid.UUID, error) {
	fileId := uuid.New()
	revisionId := uuid.New()
	f, err := FileRepositoryInstance.CreateFile(stream.Context(), fileId, 0)
	if err != nil {
		logger.Logger.Err(err).Msg("ferror while creating file")
		return 0, nil, uuid.Nil, uuid.Nil, status.Error(codes.Internal, "failed to write file. please retry the request")
	}
	defer f.Close()

	totalFileSize, sniff, err := receiveChunks(stream, f)
	if err != nil {
		return 0, nil, uuid.Nil, uuid.Nil, err
	}

	return totalFileSize, sniff, fileId, revisionId, nil
}

func receiveChunks(stream apiRestaurantFile.FileService_StoreFileServer, f io.WriteCloser) (uint64, []byte, error) {
	var totalFileSize uint64 = 0
	var sniffByteCount uint64 = 0
	sniff := make([]byte, sniffSize)
	for {
		finished, chunkMessage, err := receiveChunk(stream)
		if err != nil {
			return 0, nil, err
		}
		if finished {
			break
		}
		totalFileSize += uint64(len(chunkMessage.Chunk))

		if sniffByteCount < sniffSize {
			missingBytes := min(sniffSize-sniffByteCount, uint64(len(chunkMessage.Chunk)))
			copy(sniff[sniffByteCount:], chunkMessage.Chunk[:missingBytes])
			sniffByteCount += missingBytes
		}

		_, err = f.Write(chunkMessage.Chunk)
		if err != nil {
			logger.Logger.Err(err).Msg("failed to write chunk to file")
			return 0, nil, status.Error(codes.Internal, "failed to write chunk to file. please retry the request")
		}
	}

	return totalFileSize, sniff[:sniffByteCount], nil
}

func receiveMetadata(stream apiRestaurantFile.FileService_StoreFileServer) (*apiRestaurantFile.StoreFileRequest_Name, error) {
	firstRequest, err := stream.Recv()
	if err != nil {
		logger.Logger.Err(err).Msg("receiving message failed")
		return nil, status.Errorf(codes.Internal, "receiving message failed. please retry the request")
	}
	msg, ok := firstRequest.File.(*apiRestaurantFile.StoreFileRequest_Name)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_Name' expected. Actual value is "+reflect.TypeOf(firstRequest.File).String()+".")
	}
	return msg, nil
}

func receiveChunk(stream apiRestaurantFile.FileService_StoreFileServer) (bool, *apiRestaurantFile.StoreFileRequest_Chunk, error) {
	request, err := stream.Recv()
	if err == io.EOF {
		return true, nil, nil
	}
	if err != nil {
		logger.Logger.Err(err).Msg("failed to receive chunk")
		return false, nil, status.Errorf(codes.Internal, "failed to receive chunk. please retry the request")
	}
	msg, ok := request.File.(*apiRestaurantFile.StoreFileRequest_Chunk)
	if !ok {
		return false, nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_Chunk' expected. Actual value is "+reflect.TypeOf(request.File).String()+".")
	}
	return false, msg, nil
}

func createStoreFileResponse(fileId uuid.UUID, revisionId uuid.UUID, totalFileSize uint64, contentType string, extension string, createdAt time.Time) (*apiRestaurantFile.StoreFileResponse, error) {
	fileUuid, err := apiProtobuf.ToProtobuf(fileId)
	if err != nil {
		return nil, err
	}

	revisionUuid, err := apiProtobuf.ToProtobuf(revisionId)
	if err != nil {
		return nil, err
	}

	var response = &apiRestaurantFile.StoreFileResponse{
		StoredFile: &apiRestaurantFile.StoredFile{
			Id:         fileUuid,
			RevisionId: revisionUuid,
		},
		StoredFileMetadata: &apiRestaurantFile.StoredFileMetadata{
			CreatedAt: timestamppb.New(createdAt),
			Size:      totalFileSize,
			MediaType: contentType,
			Extension: extension,
		},
	}
	return response, nil
}

func (s *FileServiceServer) DownloadFile(request *apiRestaurantFile.DownloadFileRequest, stream apiRestaurantFile.FileService_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method not implemented")
}

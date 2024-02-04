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

// - receive filename
// - read extension
// - create file in storage
// - receive chunk
// - write chunk to file
// - write chunk to snifarray if snifarray < 512 byte
// - increase size
// - receive chunk gain until end
// - create file id
// - write new document to database
// - write response
func (s *FileServiceServer) StoreFile(stream apiRestaurantFile.FileService_StoreFileServer) error {
	metaData, err := receiveMetadata(stream)
	if err != nil {
		return err
	}

	const sniffSize = 512
	var totalFileSize uint64 = 0
	var sniffByteCount uint64 = 0
	sniff := make([]byte, sniffSize)
	ctx := stream.Context()

	fileId := uuid.New()
	f, err := FileRepositoryInstance.CreateFile(ctx, fileId, 0)
	if err != nil {
		return err
	}

	for {
		finished, chunkMessage, err := receiveChunk(stream)
		if finished {
			break
		}
		if err != nil {
			return err
		}
		totalFileSize += uint64(len(chunkMessage.Chunk))

		if sniffByteCount < sniffSize {
			missingBytes := min(sniffSize-sniffByteCount, uint64(len(chunkMessage.Chunk)))
			copy(sniff[sniffByteCount:], chunkMessage.Chunk[:missingBytes])
			sniffByteCount += missingBytes
		}

		_, err = f.Write(chunkMessage.Chunk)
		if err != nil {
			return status.Error(codes.Internal, "failed to write chunk to file")
		}
	}

	f.Close()

	contentType := http.DetectContentType(sniff[:sniffByteCount])
	extension := filepath.Ext(metaData.Name)

	response, err := createResponse(fileId, totalFileSize, contentType, extension)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to convert google uuid to protobuf uuid")
		return status.Error(codes.Internal, "failed to convert google uuid to protobuf uuid")
	}
	err = stream.SendAndClose(response)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to send response")
		return status.Error(codes.Internal, "failed to send response")
	}

	return nil
}

func createResponse(fileId uuid.UUID, totalFileSize uint64, contentType string, extension string) (*apiRestaurantFile.StoreFileResponse, error) {
	fileUuid, err := apiProtobuf.ToProtobuf(fileId)
	if err != nil {
		return nil, err
	}
	var response = &apiRestaurantFile.StoreFileResponse{
		StoredFile: &apiRestaurantFile.StoredFile{
			Id:       fileUuid,
			Revision: 1,
		},
		StoredFileMetadata: &apiRestaurantFile.StoredFileMetadata{
			CreatedAt: timestamppb.New(time.Now()),
			Size:      totalFileSize,
			MediaType: contentType,
			Extension: extension,
		},
	}
	return response, nil
}

func receiveMetadata(stream apiRestaurantFile.FileService_StoreFileServer) (*apiRestaurantFile.StoreFileRequest_Name, error) {
	firstRequest, err := stream.Recv()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "receiving message failed")
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
		return false, nil, err
	}
	msg, ok := request.File.(*apiRestaurantFile.StoreFileRequest_Chunk)
	if !ok {
		return false, nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_Chunk' expected. Actual value is "+reflect.TypeOf(request.File).String()+".")
	}
	return false, msg, nil
}

func (s *FileServiceServer) DownloadFile(request *apiRestaurantFile.DownloadFileRequest, stream apiRestaurantFile.FileService_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method not implemented")
}

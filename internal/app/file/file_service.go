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
	name := metaData.Name
	extension := filepath.Ext(name)

	fileId := uuid.New()
	var totalFileSize uint64 = 0
	sniff := make([]byte, 512)

	ctx := stream.Context()
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

		if totalFileSize < 512 {
			missingBytes := 512 - totalFileSize
			remaingBytesInChunk := uint64(len(chunkMessage.Chunk))
			if remaingBytesInChunk < missingBytes {
				missingBytes = remaingBytesInChunk
			}
			copy(sniff[totalFileSize:], chunkMessage.Chunk[:missingBytes])
		}

		totalFileSize += uint64(len(chunkMessage.Chunk))

		_, err = f.Write(chunkMessage.Chunk)
		if err != nil {
			return status.Error(codes.Internal, "failed to write chunk to file")
		}
	}

	sniffByteCount := totalFileSize
	if sniffByteCount > 512 {
		sniffByteCount = 512
	}
	contentType := http.DetectContentType(sniff[:sniffByteCount])
	f.Close()

	fileUuid, err := apiProtobuf.ToProtobuf(fileId)
	if err != nil {
		return status.Error(codes.Internal, "failed to convert google uuid to protobuf uuid")
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
	err = stream.SendAndClose(response)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to send response")
		return status.Error(codes.Internal, "failed to send response")
	}

	return nil
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

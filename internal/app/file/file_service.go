package file

import (
	"io"
	"os"
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

var _ apiRestaurantFile.FileServiceServer = &FileServiceServer{}

type FileServiceServer struct {
	apiRestaurantFile.UnimplementedFileServiceServer
}

func (s *FileServiceServer) StoreFile(stream apiRestaurantFile.FileService_StoreFileServer) error {
	metaData, err := receiveMetadata(stream)
	if err != nil {
		return err
	}
	name := metaData.Name

	// Create a new file
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		// The client has finished sending data
		finished, chunkMessage, err := receiveChunk(stream)
		if finished {
			break
		}
		if err != nil {
			return err
		}

		_, err = f.Write(chunkMessage.Chunk)
		if err != nil {
			return status.Error(codes.Internal, "failed to write chunk to file")
		}
	}

	// Send the response to the client
	fileId := uuid.New()
	// todo read metadata from file

	contractFileId, err := apiProtobuf.ToProtobuf(fileId)
	var response = &apiRestaurantFile.StoreFileResponse{
		StoredFile: &apiRestaurantFile.StoredFile{
			Id:       contractFileId,
			Revision: 1,
		},
		StoredFileMetadata: &apiRestaurantFile.StoredFileMetadata{
			CreatedAt: timestamppb.New(time.Now()),
			Size:      1212,
			MediaType: "image/jpeg",
			Extension: "jpg",
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
	// TODO: Implement this method
	return nil
}

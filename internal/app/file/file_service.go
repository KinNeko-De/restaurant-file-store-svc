package file

import (
	"io"
	"reflect"

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

func (s *FileServiceServer) StoreFile(stream apiRestaurantFile.FileService_StoreFileServer) error {
	fileName, err := receiveMetadata(stream)
	if err != nil {
		return err
	}

	createdFileMetadata, err := createFile(stream, fileName.Name)
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

	totalFileSize, sniff, err := writeFile(stream, fileId)
	if err != nil {
		return nil, err
	}

	createdRevision := newRevision(fileName, totalFileSize, sniff)
	createdFileMetadata := newFileMetadata(fileId, createdRevision)

	err = FileMetadataRepositoryInstance.StoreFileMetadata(stream.Context(), createdFileMetadata)
	return createdFileMetadata, err
}

func writeFile(stream apiRestaurantFile.FileService_StoreFileServer, fileId uuid.UUID) (uint64, []byte, error) {
	fileWriter, err := FileRepositoryInstance.CreateFile(stream.Context(), fileId, 0)
	if err != nil {
		logger.Logger.Err(err).Msg("ferror while creating file")
		return 0, nil, status.Error(codes.Internal, "failed to write file. please retry the request")
	}
	defer fileWriter.Close()
	totalFileSize, sniff, err := receiveChunks(stream, fileWriter)
	if err != nil {
		return 0, nil, err
	}

	return totalFileSize, sniff, nil
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

func createStoreFileResponse(createdFileMetadata *FileMetadata) (*apiRestaurantFile.StoreFileResponse, error) {
	fileUuid, err := apiProtobuf.ToProtobuf(createdFileMetadata.Id)
	if err != nil {
		return nil, err
	}

	revisionUuid, err := apiProtobuf.ToProtobuf(createdFileMetadata.Revisions[0].Id)
	if err != nil {
		return nil, err
	}

	var response = &apiRestaurantFile.StoreFileResponse{
		StoredFile: &apiRestaurantFile.StoredFile{
			Id:         fileUuid,
			RevisionId: revisionUuid,
		},
		StoredFileMetadata: &apiRestaurantFile.StoredFileMetadata{
			CreatedAt: timestamppb.New(createdFileMetadata.CreatedAt),
			Size:      createdFileMetadata.Revisions[0].Size,
			MediaType: createdFileMetadata.Revisions[0].MediaType,
			Extension: createdFileMetadata.Revisions[0].Extension,
		},
	}
	return response, nil
}

func (s *FileServiceServer) DownloadFile(request *apiRestaurantFile.DownloadFileRequest, stream apiRestaurantFile.FileService_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method not implemented")
}

package file

import (
	"errors"
	"path/filepath"

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
	storeFile, err := receiveMetadata(stream)
	if err != nil {
		return err
	}
	extension := filepath.Ext(storeFile.Name)
	if extension == "" {
		return status.Error(codes.InvalidArgument, "fileName '"+storeFile.Name+"'is invalid. fileName must have a dot and an extension. e.g. 'example.txt', '.txt'")
	}

	createdFileMetadata, err := createFile(stream, extension)
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

func (s *FileServiceServer) StoreRevision(stream apiRestaurantFile.FileService_StoreRevisionServer) error {
	storeRevision, err := receiveRevisionMetadata(stream)
	if err != nil {
		return err
	}

	fileId, err := getFileId(storeRevision)
	if err != nil {
		return err
	}

	extension := filepath.Ext(storeRevision.StoreFile.Name)
	if extension == "" {
		return status.Error(codes.InvalidArgument, "fileName '"+storeRevision.StoreFile.Name+"'is invalid. fileName must have a dot and an extension. e.g. 'example.txt', '.txt'")

	}

	createdFileMetadata, err := createRevision(stream, fileId, extension)
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

func getFileId(request *apiRestaurantFile.StoreRevision) (uuid.UUID, error) {
	requested := request.FileId
	if requested == nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "fileId is mandatory. Please provide a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}

	fileId, err := apiProtobuf.ToUuid(requested)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "fileId '"+requested.String()+"' is not a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}

	return fileId, nil
}

func createFile(stream apiRestaurantFile.FileService_StoreFileServer, extension string) (*FileMetadata, error) {
	fileId := uuid.New()
	revisionId := uuid.New()

	totalFileSize, sniff, err := writeFile(&StoreFile{stream}, stream.Context(), fileId, revisionId)
	if err != nil {
		return nil, err
	}

	createdRevision := newRevision(revisionId, extension, totalFileSize, sniff)
	createdFileMetadata := newFileMetadata(fileId, createdRevision)

	err = FileMetadataRepositoryInstance.StoreFileMetadata(stream.Context(), createdFileMetadata)
	if err != nil {
		logger.Logger.Err(err).Msg("failed to store file metadata")
		return nil, status.Error(codes.Internal, "failed to store file metadata. please retry the request")
	}

	return &createdFileMetadata, err
}

func createRevision(stream apiRestaurantFile.FileService_StoreRevisionServer, fileId uuid.UUID, extension string) (*FileMetadata, error) {
	revisionId := uuid.New()

	totalFileSize, sniff, err := writeFile(&StoreRevision{stream}, stream.Context(), fileId, revisionId)
	if err != nil {
		return nil, err
	}

	createdRevision := newRevision(revisionId, extension, totalFileSize, sniff)
	err = FileMetadataRepositoryInstance.StoreRevision(stream.Context(), fileId, createdRevision)
	if err != nil {
		if errors.Is(err, FileMetadataRepositoryInstance.NoMatchError()) {
			return nil, status.Error(codes.NotFound, "file with id '"+fileId.String()+"' not found.")
		}
		logger.Logger.Err(err).Msg("failed to store revision metadata")
		return nil, status.Error(codes.Internal, "failed to store file metadata. please retry the request")
	}

	createdFileMetadata := newFileMetadata(fileId, createdRevision)
	return &createdFileMetadata, err
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

func (s *FileServiceServer) DownloadRevision(request *apiRestaurantFile.DownloadRevisionRequest, stream apiRestaurantFile.FileService_DownloadRevisionServer) error {
	requestedFileId, requestedRevisionId, err := getRequestedFile(request)
	if err != nil {
		return err
	}

	scopedLogger := logger.Logger.With().Str("fileId", requestedFileId.String()).Str("revisionId", requestedRevisionId.String()).Logger()

	fileMetadata, err := fetchMetadata(stream.Context(), requestedFileId, scopedLogger)
	if err != nil {
		return err
	}
	revision, err := fileMetadata.GetRevision(requestedRevisionId)
	if err != nil {
		return status.Error(codes.NotFound, "revision with id '"+requestedRevisionId.String()+"' not found.")
	}

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

func getRequestedFile(request *apiRestaurantFile.DownloadRevisionRequest) (uuid.UUID, uuid.UUID, error) {
	requestedFileId := request.GetFileId()
	if requestedFileId == nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.InvalidArgument, "fileId is mandatory. Please provide a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}
	fileId, err := apiProtobuf.ToUuid(requestedFileId)
	if err != nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.InvalidArgument, "fileId '"+requestedFileId.String()+"' is not a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}

	requestedRevisionId := request.GetRevisionId()
	if requestedRevisionId == nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.InvalidArgument, "revisionId is mandatory. Please provide a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}
	revsionId, err := apiProtobuf.ToUuid(requestedRevisionId)
	if err != nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.InvalidArgument, "revisionId '"+requestedRevisionId.String()+"' is not a valid uuid. The uuid must be in the following format: 12345678-90ab-cdef-1234-567890abcef0")
	}

	return fileId, revsionId, nil
}

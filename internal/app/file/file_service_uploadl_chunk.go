package file

import (
	"reflect"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChunckStream interface {
	Recv() (ChunckRequest, error)
}

type StoreFile struct {
	stream apiRestaurantFile.FileService_StoreFileServer
}

func (s StoreFile) Recv() (ChunckRequest, error) {
	request, err := s.stream.Recv()
	return StoreFileRequest{request: request}, err
}

type StoreRevision struct {
	stream apiRestaurantFile.FileService_StoreRevisionServer
}

func (s StoreRevision) Recv() (ChunckRequest, error) {
	request, err := s.stream.Recv()
	return StoreRevisionRequest{request: request}, err
}

type ChunckRequest interface {
	GetChunk() ([]byte, error)
}

type StoreFileRequest struct {
	request *apiRestaurantFile.StoreFileRequest
}

func (w StoreFileRequest) GetChunk() ([]byte, error) {
	chunk := w.request.GetChunk()
	if chunk == nil {
		return nil, status.Errorf(codes.InvalidArgument, "PartCase of type 'fileServiceApi.StoreFileRequest_Chunk' expected. Actual value is "+reflect.TypeOf(w.request.Part).String()+".")
	}
	return chunk, nil
}

type StoreRevisionRequest struct {
	request *apiRestaurantFile.StoreRevisionRequest
}

func (w StoreRevisionRequest) GetChunk() ([]byte, error) {
	chunk := w.request.GetChunk()
	if chunk == nil {
		return nil, status.Errorf(codes.InvalidArgument, "PartCase of type 'fileServiceApi.StoreRevisionRequest_Chunk' expected. Actual value is "+reflect.TypeOf(w.request.Part).String()+".")
	}
	return chunk, nil
}

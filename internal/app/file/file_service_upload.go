package file

import (
	"io"
	"reflect"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func receiveChunks(stream ChunckStream, f io.WriteCloser) (uint64, []byte, error) {
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
		return nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreFileRequest_StoreFile' expected. Actual value is "+reflect.TypeOf(firstRequest.Part).String()+".")
	}
	return msg, nil
}

func receiveRevisionMetadata(stream apiRestaurantFile.FileService_StoreRevisionServer) (*apiRestaurantFile.StoreRevision, error) {
	firstRequest, err := stream.Recv()
	if err != nil {
		logger.Logger.Err(err).Msg("receiving message failed")
		return nil, status.Errorf(codes.Internal, "receiving message failed. please retry the request")
	}

	msg := firstRequest.GetStoreRevision()
	if msg == nil {
		return nil, status.Errorf(codes.InvalidArgument, "FileCase of type 'fileServiceApi.StoreRevisionRequest_StoreRevision' expected. Actual value is "+reflect.TypeOf(firstRequest.Part).String()+".")
	}
	return msg, nil
}

func receiveChunk(stream ChunckStream) (bool, []byte, error) {
	request, err := stream.Recv()
	if err == io.EOF {
		return true, nil, nil
	}
	if err != nil {
		logger.Logger.Err(err).Msg("failed to receive chunk")
		return false, nil, status.Errorf(codes.Internal, "failed to receive chunk. please retry the request")
	}

	msg, err := request.GetChunk()
	if err != nil {
		return false, nil, err
	}
	return false, msg, nil
}

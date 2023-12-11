package file

import (
	fileServiceApi "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
)

type FileServiceServer struct {
	fileServiceApi.UnimplementedFileServiceServer
}

func (s *FileServiceServer) StoreFile(stream fileServiceApi.FileService_StoreFileServer) error {
	return nil
}

func (s *FileServiceServer) DownloadFile(request *fileServiceApi.DownloadFileRequest, stream fileServiceApi.FileService_DownloadFileServer) error {
	return nil
}

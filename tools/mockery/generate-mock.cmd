pushd windows_x86_64

mockery --dir "%GOPATH%/pkg/mod/github.com/kinneko-de/api-contract/golang/kinnekode/restaurant@v0.0.3-store-files.7/file/v1" --name FileService_StoreFileServer --structname MockFileService_StoreFileServer --filename file_service_store_file_server_mock.go --with-expecter --output ../../../test/testing/file --outpkg file
mockery --dir "%GOPATH%/pkg/mod/github.com/kinneko-de/api-contract/golang/kinnekode/restaurant@v0.0.3-store-files.7/file/v1" --name FileService_StoreRevisionServer --structname MockFileService_StoreRevisionServer --filename file_service_store_revision_server_mock.go --with-expecter --output ../../../test/testing/file --outpkg file
mockery --dir "%GOPATH%/pkg/mod/github.com/kinneko-de/api-contract/golang/kinnekode/restaurant@v0.0.3-store-files.7/file/v1" --name FileService_DownloadFileServer --structname MockFileService_DownloadFileServer --filename file_service_download_file_server_mock.go --with-expecter --output ../../../test/testing/file --outpkg file
mockery --dir "%GOPATH%/pkg/mod/github.com/kinneko-de/api-contract/golang/kinnekode/restaurant@v0.0.3-store-files.7/file/v1" --name FileService_DownloadRevisionServer --structname MockFileService_DownloadRevisionServer --filename file_service_download_revision_server_mock.go --with-expecter --output ../../../test/testing/file --outpkg file
mockery --dir "../../../internal/app/file" --name FileMetadataRepository --filename file_metadata_repository_mock.go --with-expecter --inpackage
mockery --dir "../../../internal/app/file" --name FileRepository --filename file_repository_mock.go --with-expecter --inpackage
mockery --dir ".." --name WriteCloser --structname MockWriteCloser --filename write_closer_mock.go --with-expecter --output ../../../test/testing/io --outpkg io
mockery --dir ".." --name ReadCloser --structname MockReadCloser --filename read_closer_mock.go --with-expecter --output ../../../test/testing/io --outpkg io

popd

pause
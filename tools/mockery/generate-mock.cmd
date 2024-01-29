pushd windows_x86_64

mockery --dir "C:\Users\nils\go\pkg\mod\github.com\kinneko-de\api-contract\golang\kinnekode\restaurant@v0.0.2-store-files.14\file\v1" --name FileService_StoreFileServer --filename file_service_store_file_server_mock.go --with-expecter --output ../../../internal/app/file --outpkg file

popd
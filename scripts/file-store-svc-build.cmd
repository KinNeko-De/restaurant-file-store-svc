pushd ..

go mod download

echo 0.0.1-local > build/version.txt

set GOARCH=amd64
set GOOS=linux
go build -o bin/app cmd/file-store-svc/main.go

echo set_by_ci > build/version.txt

popd

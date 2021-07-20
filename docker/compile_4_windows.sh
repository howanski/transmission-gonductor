#!/bin/bash
docker run --rm -v "$PWD"/../:/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=amd64 -e CGO_ENABLED=1 -e GIN_MODE=release golang:1.16 go build -v -o gonductor-win64.exe
# docker run --rm -v "$PWD"/../:/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=386 -e CGO_ENABLED=1 -e GIN_MODE=release golang:1.16 go build -v -o gonductor-win32.exe
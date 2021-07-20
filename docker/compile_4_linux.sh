#!/bin/bash
docker run --rm -v "$PWD"/../:/usr/src/myapp -w /usr/src/myapp -e GOOS=linux -e GOARCH=amd64 -e GIN_MODE=release golang:1.16 go build -v -o gonductor-64
# docker run --rm -v "$PWD"/../:/usr/src/myapp -w /usr/src/myapp -e GOOS=linux -e GOARCH=386 -e GIN_MODE=release golang:1.16 go build -v -o gonductor-32
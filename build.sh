#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

dep ensure

protoc -I . protos/ingestor.proto --go_out=plugins=grpc:.

go build -o build/ingestor ingestor/main.go
go build -o build/keeper keeper/main.go
go build -o build/grpc grpc/main.go

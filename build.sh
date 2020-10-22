#!/usr/bin/env bash

go fmt ./...
go vet ./...
if [ $? -ne 0 ] ; then exit; fi

PROJECT=github.com/dmalix/financelime-rest-api/cmd

VERSION=$(cat VERSION)
#COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMPILER="$(go version)"

set -e
export GOFLAGS="-mod=vendor"
go build \
  -ldflags="-s -w -X '${PROJECT}/config.version=${VERSION}' -X '${PROJECT}/config.commit=${COMMIT}' -X '${PROJECT}/config.buildTime=${BUILD_TIME}' -X '${PROJECT}/config.compiler=${COMPILER}'" \
  -o bin/financelime-rest-api cmd/main.go
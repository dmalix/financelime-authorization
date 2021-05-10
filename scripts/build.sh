#!/usr/bin/env bash

PROJECT=github.com/dmalix/financelime-authorization

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMPILER="$(go version)"

set -e
export GOFLAGS="-mod=vendor"

go build \
  -ldflags="-s -w -X '${PROJECT}/app.versionNumber=${VERSION}' -X '${PROJECT}/app.versionCommit=${COMMIT}' -X '${PROJECT}/app.versionBuildTime=${BUILD_TIME}' -X '${PROJECT}/server.versionCompiler=${COMPILER}'" \
  -o bin/financelime-auth main.go
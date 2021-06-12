#!/usr/bin/env bash

PROJECT=github.com/dmalix/authorization-service

DEVELOPMENT_MODE=true

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMPILER="$(go version)"

set -e
export GOFLAGS="-mod=vendor"

go build \
  -ldflags="-s -w  -X '${PROJECT}/config.developmentMode=${DEVELOPMENT_MODE}' -X '${PROJECT}/config.versionNumber=${VERSION}' -X '${PROJECT}/config.versionCommit=${COMMIT}' -X '${PROJECT}/config.versionBuildTime=${BUILD_TIME}' -X '${PROJECT}/config.versionCompiler=${COMPILER}'" \
  -o bin/financelime-auth main.go
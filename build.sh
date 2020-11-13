#!/usr/bin/env bash
# Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
# Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
# License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html

go fmt ./...
go vet ./...
if [ $? -ne 0 ] ; then exit; fi

PROJECT=github.com/dmalix/financelime-authorization

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMPILER="$(go version)"

set -e
export GOFLAGS="-mod=vendor"
go build \
  -ldflags="-s -w -X '${PROJECT}/server.version=${VERSION}' -X '${PROJECT}/server.commit=${COMMIT}' -X '${PROJECT}/server.buildTime=${BUILD_TIME}' -X '${PROJECT}/server.compiler=${COMPILER}'" \
  -o bin/financelime-auth cmd/main.go
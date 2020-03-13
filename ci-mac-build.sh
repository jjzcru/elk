#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

GOOS=darwin
GOARCH=amd64

cd $MODULE_PATH

NAME=elk_v${VERSION}_${GOOS}_${GOARCH}
BIN_PATH=$BUILD_PATH/$NAME

go build -o $BIN_PATH
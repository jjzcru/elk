#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

GOOS=darwin

# Build for 386
GOARCH=386
cd $MODULE_PATH
NAME=elk

BIN_PATH=$BUILD_PATH/$NAME
go build -o $BIN_PATH

cd $BUILD_PATH
ZIP_PATH=${BIN_PATH}_v${VERSION}_${GOOS}_${GOARCH}.zip

zip $ZIP_PATH $NAME
rm $NAME

# Build for amd64
GOARCH=amd64
cd $MODULE_PATH
NAME=elk

BIN_PATH=$BUILD_PATH/$NAME
go build -o $BIN_PATH

cd $BUILD_PATH
ZIP_PATH=${BIN_PATH}_v${VERSION}_${GOOS}_${GOARCH}.zip

zip $ZIP_PATH $NAME
rm $NAME
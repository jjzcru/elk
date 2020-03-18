#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

declare -A platforms
platforms[linux,0]=amd64
platforms[linux,1]=386
platforms[linux,2]=arm
platforms[linux,3]=arm64
platforms[solaris,0]=amd64

for key in "${!platforms[@]}"; do
    GOOS=${key::-2}
    GOARCH=${platforms[$key]}
    cd $MODULE_PATH
    NAME=elk

    BIN_PATH=$BUILD_PATH/$NAME
    go build -o $BIN_PATH

    cd $BUILD_PATH
    ZIP_PATH=${BIN_PATH}_v${VERSION}_${GOOS}_${GOARCH}.zip

    zip $ZIP_PATH $NAME
    rm $NAME
done
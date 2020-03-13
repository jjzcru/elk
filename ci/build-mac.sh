#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

declare -A platforms
platforms[darwin,0]=386
platforms[darwin,1]=amd64

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
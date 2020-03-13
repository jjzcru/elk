#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

declare -A platforms
platforms[windows,0]=386
platforms[windows,1]=amd64

for key in "${!platforms[@]}"; do
    GOOS=${key::-2}
    GOARCH=${platforms[$key]}
    cd $MODULE_PATH
    NAME=elk

    BIN_PATH=$BUILD_PATH/${NAME}
    go build -o ${BIN_PATH}.exe

    cd $BUILD_PATH
    ZIP_PATH=${BIN_PATH}_v${VERSION}_${GOOS}_${GOARCH}.zip

    # ip $ZIP_PATH ${NAME}.exe
    # rm ${NAME}.exe
done
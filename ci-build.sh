#!/bin/bash

VERSION=$(<VERSION)
BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

declare -A platforms
platforms[linux,0]=amd64
platforms[linux,1]=386
platforms[darwin,0]=amd64
platforms[windows,0]=amd64
platforms[windows,1]=386

for key in "${!platforms[@]}"; do
    GOOS=${key::-2}
    GOARCH=${platforms[$key]}
    cd $MODULE_PATH
    NAME=elk_v${VERSION}_${GOOS}_${GOARCH}
    

    BIN_PATH=$BUILD_PATH/$NAME
    if [ $GOOS == "windows" ]
    then
        go build -o ${BIN_PATH}.exe
    else
        go build -o $BIN_PATH
    fi
done
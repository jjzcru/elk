#!/bin/bash

BASE_PATH=$(pwd)
BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)

COMMIT=$(git rev-parse --short HEAD)
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))

day=$(date +'%a')
month=$(date +'%b')
fill_date=$(date +'%d_%T_%Y')

DATE="${day^}_${month^}_${fill_date}"

declare -A platforms
platforms[linux,0]=amd64
platforms[linux,1]=386
platforms[linux,2]=arm
platforms[linux,3]=arm64
platforms[darwin,0]=amd64
platforms[windows,0]=amd64
platforms[windows,1]=386
platforms[solaris,0]=amd64

for key in "${!platforms[@]}"; do
    GOOS=${key::-2}
    GOARCH=${platforms[$key]}
    cd $MODULE_PATH
    NAME=elk_${VERSION}_${GOOS}_${GOARCH}

    BIN_PATH=$BUILD_PATH/$NAME
    if [ $GOOS == "windows" ]
    then
        go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE" -o ${BIN_PATH}.exe
    else
        go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE" -o $BIN_PATH
    fi
    
    cd $BUILD_PATH
    if [ $GOOS == "windows" ]
    then
        zip $BIN_PATH.zip ${NAME}.exe
        rm ${NAME}.exe
    else
        zip $BIN_PATH.zip $NAME
        rm $NAME
    fi
done
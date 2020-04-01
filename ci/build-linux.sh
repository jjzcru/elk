#!/bin/bash

BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

day=$(date +'%a')
month=$(date +'%b')
fill_date=$(date +'%d_%T_%Y')

DATE="${day^}_${month^}_${fill_date}"

declare -A platforms
platforms[linux,0]=amd64
platforms[linux,1]=386
platforms[linux,2]=arm
platforms[linux,3]=arm64
platforms[solaris,0]=amd64

echo "BUILT DETAILS"
echo "VERSION: $VERSION"
echo "COMMIT: $COMMIT"
echo "DATE: $DATE"
echo "GO VERSION: $GOVERSION"

for key in "${!platforms[@]}"; do
    GOOS=${key::-2}
    GOARCH=${platforms[$key]}
    cd $MODULE_PATH
    NAME=elk

    BIN_PATH=$BUILD_PATH/$NAME
    go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE -X main.goVersion=$GOVERSION" -o $BIN_PATH

    cd $BUILD_PATH
    ZIP_PATH=${BIN_PATH}_${VERSION}_${GOOS}_${GOARCH}.zip

    zip $ZIP_PATH $NAME
    rm $NAME
done
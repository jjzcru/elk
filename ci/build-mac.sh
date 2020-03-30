#!/bin/bash

BUILD_PATH=$(pwd)/bin
MODULE_PATH=$(pwd)/cmd/elk

GOOS=darwin

day=$(date +'%a')
month=$(date +'%b')
fill_date=$(date +'%d %T %Y')

DATE="${day} ${month} ${fill_date}"

echo "BUILT"
echo "VERSION: $VERSION"
echo "COMMIT: $COMMIT"
echo "DATE: $DATE"

# Build for 386
GOARCH=386
cd $MODULE_PATH
NAME=elk

BIN_PATH=$BUILD_PATH/$NAME
go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE -X main.goVersion=$GOVERSION" -o $BIN_PATH

cd $BUILD_PATH
ZIP_PATH=${BIN_PATH}_${VERSION}_${GOOS}_${GOARCH}.zip

zip $ZIP_PATH $NAME
rm $NAME

# Build for amd64
GOARCH=amd64
cd $MODULE_PATH
NAME=elk

BIN_PATH=$BUILD_PATH/$NAME
go build -o $BIN_PATH

cd $BUILD_PATH
ZIP_PATH=${BIN_PATH}_${VERSION}_${GOOS}_${GOARCH}.zip

zip $ZIP_PATH $NAME
rm $NAME
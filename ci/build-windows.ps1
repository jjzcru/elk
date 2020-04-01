$DATE = Get-Date -UFormat "%a_%b_%d_%T_%Y"
$COMMIT = $env:COMMIT
$VERSION = $env:VERSION

$BUILD_PATH = "$((Get-Item -Path ".\").FullName)\bin"
$MODULE_PATH = "$((Get-Item -Path ".\").FullName)\cmd\elk"

$env:GOOS = "windows"
$GOOS = "windows"
$NAME = "elk"

cls

echo "BUILT DETAILS"
echo "VERSION: $VERSION"
echo "COMMIT: $COMMIT"
echo "DATE: $DATE"
echo "GOVERSION: $GOVERSION"

# 386
$env:GOARCH = "386"
$GOARCH = "386"
cd $MODULE_PATH

$BIN_PATH = "$BUILD_PATH\$NAME"
echo "ARCH: $($GOARCH)"
echo "--------------------------"
echo "Building $($GOARCH) binary"
go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE -X main.goVersion=$GOVERSION" -o "$BIN_PATH.exe"
echo "Build successful"

cd "$BUILD_PATH"
$ZIP_PATH = "$($BIN_PATH)_$($VERSION)_$($GOOS)_$($GOARCH).zip"

echo "Compressing $($GOARCH) binary"
compress-archive "$BIN_PATH.exe" "$ZIP_PATH" -Force
rm "$NAME.exe"
echo "Compress successful"
echo "--------------------------"
echo ""

# amd64
$env:GOARCH = "amd64"
$GOARCH = "amd64"
cd $MODULE_PATH

$BIN_PATH = "$BUILD_PATH\$NAME"
echo "ARCH: $($GOARCH)"
echo "--------------------------"
echo "Building $($GOARCH) binary"
go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE -X main.goVersion=$GOVERSION" -o "$BIN_PATH.exe"
echo "Build successful"

cd "$BUILD_PATH"
$ZIP_PATH = "$($BIN_PATH)_$($VERSION)_$($GOOS)_$($GOARCH).zip"

echo "Compressing $($GOARCH) binary"
compress-archive "$BIN_PATH.exe" "$ZIP_PATH" -Force
rm "$NAME.exe"
echo "Compress successful"
echo "--------------------------"
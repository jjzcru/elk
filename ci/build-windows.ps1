$VERSION = Get-Content .\VERSION -Raw
$BASE_PATH = (Get-Item -Path ".\").FullName
$BUILD_PATH = "$((Get-Item -Path ".\").FullName)\bin"
$MODULE_PATH = "$((Get-Item -Path ".\").FullName)\cmd\elk"

$env:GOOS = "windows"
$GOOS = "windows"

$NAME = "elk"

# 386
$env:GOARCH = "386"
$GOARCH = "386"
cd $MODULE_PATH

$BIN_PATH = "$BUILD_PATH\$NAME"
echo $BIN_PATH
go build -o "$BIN_PATH.exe"

cd "$BUILD_PATH"
$ZIP_PATH = "$BIN_PATH_v$VERSION_$GOOS_$GOARCH.zip"
compress-archive "$NAME.exe" $ZIP_PATH
# rm ${NAME}.exe

# amd64
#$env:GOARCH = "amd64"
#$GOARCH = "amd64"
#cd $MODULE_PATH

#$BIN_PATH = "$($BUILD_PATH)/$($NAME))"
#go build -o $BIN_PATH.exe

#cd $BUILD_PATH
#$ZIP_PATH = "$BIN_PATH_v$VERSION_$GOOS_$GOARCH.zip"
#compress-archive "$NAME.exe" $ZIP_PATH
#rm ${NAME}.exe


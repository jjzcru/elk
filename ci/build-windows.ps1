

$BUILD_PATH = "$((Get-Item -Path ".\").FullName)\bin"
$MODULE_PATH = "$((Get-Item -Path ".\").FullName)\cmd\elk"

$env:GOOS = "windows"
$GOOS = "windows"
$NAME = "elk"

cls

# 386
$env:GOARCH = "386"
$GOARCH = "386"
cd $MODULE_PATH

$BIN_PATH = "$BUILD_PATH\$NAME"
echo "ARCH: $($GOARCH)"
echo "--------------------------"
echo "Building $($GOARCH) binary"
go build -o "$BIN_PATH.exe"
echo "Build successful"

cd "$BUILD_PATH"
$ZIP_PATH = "$($BIN_PATH)_$($VERSION)_$($GOOS)_$($GOARCH).zip"

echo "Compressing $($GOARCH) binary"
compress-archive "$BIN_PATH.exe" "$ZIP_PATH" -Force
rm "$NAME.exe"
echo "Compress successful"
echo "--------------------------"
echo ""
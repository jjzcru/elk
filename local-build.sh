COMMIT=$(git rev-parse --short HEAD)
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))

day=$(date +'%a')
month=$(date +'%b')
fill_date=$(date +'%d_%T_%Y')

DATE="${day}_${month}_${fill_date}"

go build -ldflags "-X main.v=$VERSION -X main.o=$GOOS -X main.arch=$GOARCH -X main.commit=$COMMIT -X main.date=$DATE" -o ./bin .
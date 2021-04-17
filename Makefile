linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-extldflags "-static" -s -w' -o out/linux/htmltoebook .
windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags '-extldflags "-static" -s -w' -o out/windows/htmltoebook.exe .
darwin:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags '-extldflags "-static" -s -w' -o out/darwin/htmltoebook .

all: linux windows darwin

# Makefile used to build releases
goget:
	go version
	go env
	go mod download
linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-extldflags "-static" -s -w' -o out/linux/htmltoebook .
windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags '-extldflags "-static" -s -w' -o out/windows/htmltoebook.exe .
darwin:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags '-extldflags "-static" -s -w' -o out/darwin/htmltoebook .

zips:
	cd out/linux && pwd && zip release-linux.zip htmltoebook
	cd out/windows && zip release-windows.zip htmltoebook.exe
	cd out/darwin && zip release-darwin.zip htmltoebook
	git log --pretty=format:"%s" > changelog.txt
	cat changelog.txt

all: goget linux windows darwin zips

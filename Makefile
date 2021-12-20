VER := $(shell git describe --tag)
SHA1 := $(shell git rev-parse HEAD)
NOW := $(shell date -u +'%Y-%m-%d_%TZ')

build:
	@echo Building the binary
	go build -ldflags "-X config.AppVersion=$(VER) -X config.Sha1Version=$(SHA1) -X config.BuildTime=$(NOW)" -o ./bin/note-server ./cmd

test:
	go test -v ./...

run:
	echo "executing the application"
	go run ./cmd/.

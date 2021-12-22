VER := $(shell git describe --tag)
SHA1 := $(shell git rev-parse HEAD)
NOW := $(shell date -u +'%Y-%m-%d_%TZ')

build:
	@echo Building the binary
	go build -ldflags "-X github.com/sksmith/note-server/config.AppVersion=$(VER)\
		-X github.com/sksmith/note-server/config.Sha1Version=$(SHA1)\
		-X github.com/sksmith/note-server/config.BuildTime=$(NOW)"\
		-o ./bin/note-server ./cmd

test:
	go test -v ./...

run:
	echo "executing the application"
	go run ./cmd/.

docker:
	@echo Building the docker image
	docker build \
		--build-arg VER=$(VER) \
		--build-arg SHA1=$(SHA1) \
		--build-arg NOW=$(NOW) .
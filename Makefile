.PHONY: all build tidy up down logs

export GOOS = linux
export GOARCH = amd64
export CGO_ENABLED = 0

export BUILD_VERSION := $(shell git describe --always --tags --abbrev=8)
export BUILD_TIME := $(shell date +%Y-%m-%dT%T%z)

all:
	@drone exec --trusted

# build all cmd/ programs

build: $(shell ls -d cmd/* | sed -e 's/cmd\//build./')
	@echo OK.

build.%: SERVICE=$*
build.%:
	go build -o build/$(SERVICE)-$(GOOS)-$(GOARCH) -ldflags "-X 'main.BuildVersion=$(BUILD_VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" ./cmd/$(SERVICE)/*.go

tidy:
	go mod tidy > /dev/null 2>&1
	go mod download > /dev/null 2>&1
	go fmt ./... > /dev/null 2>&1

# rpc generators

rpc: $(shell ls -d rpc/* | sed -e 's/\//./g')
	@echo OK.

rpc.%: SERVICE=$*
rpc.%:
	@echo '> protoc gen for $(SERVICE)'
	@protoc --proto_path=$(GOPATH)/src:. -Irpc/$(SERVICE) -I/opt/googleapis --go_out=paths=source_relative:. rpc/$(SERVICE)/*.proto
	@protoc --proto_path=$(GOPATH)/src:. -Irpc/$(SERVICE) -I/opt/googleapis --twirp_out=paths=source_relative:. rpc/$(SERVICE)/$(SERVICE).proto

up:
	@docker-compose up -d --remove-orphans

down:
	@docker-compose down --remove-orphans

logs:
	@docker-compose logs -f
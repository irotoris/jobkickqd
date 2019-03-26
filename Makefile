NAME     := jobkickqd
VERSION  := v0.0.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := -X 'main.version=$(VERSION)' -X 'main.Revision=$(REVISION)'
CURRENT  := $(shell pwd)

.PHONY: deps
## Install dependencies
deps:
	go get golang.org/x/lint
	go get ./...

.PHONY: build
## Build binaries
build: deps
	go build -ldflags "$(LDFLAGS)" -o build/$(NAME)

.PHONY: build-linux-amd64
## Build binaries for Linux(AMD64)
cross-build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/$(NAME)-linux-adm64 main.go

.PHONY: install
## compile and install
install:
	go install -ldflags "$(LDFLAGS)"

.PHONY: test
## Run tests
test: deps
	go test -v ./...

.PHONY: clean
## Clean
clean:
	go clean

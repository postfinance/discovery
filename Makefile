# This file is generated with create-go-app: do not edit.
.PHONY: build clean test snapshot all install help setup download-goreleaser download-golangci-lint download

# special target to export all variables
.EXPORT_ALL_VARIABLES:

PROJECT_RELEASE := 1
DOCKER_IMAGE_PREFIX := local
CI_PROJECT_URL := http://localhost
INSTALL_DIR := $(HOME)/bin
GO_VERSION := $(shell (go version | awk '{print $$3;}'))
BINARIES := $(shell find ./dist ! -name '*goreleaserdocker*' -path '*_linux_*' -type f -executable)
GORELEASER := 0.156.1
GOLANGCI := 1.36.0
PROTOC_BUF := 1.12.0
PROTOC := 21.12
PROTOC_GEN_GO := 1.28.1
PROTOC_GEN_GRPC_GO := 1.2.0
PROTOC_GEN_GRPC_GATEWAY := 2.15.0

## build: build the binaries only
build:
	goreleaser build --clean --snapshot

## snapshot: create a snapshot release
snapshot:
	goreleaser release --snapshot --clean --skip-sign

## clean: cleanup
clean:
	rm -rf ./dist
	rm -rf ./downloads
	rm -rf ./bin

## download:  installs goreleaser and golangci-lint in the correct version to ~/bin
download: download-goreleaser download-golangci-lint

all: build

## install: installs resulting binaries to $HOME/bin
install: build
	mkdir -p $(INSTALL_DIR) && install $(BINARIES) $(INSTALL_DIR)

## test: run linter and tests
test:
	golangci-lint run
	go test -v -count=1 ./...

## test-short: run test without linting and not in verbose mode
test-short:
	go test -count=1 ./...

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

download-goreleaser:
	mkdir -p downloads
	cd downloads && curl -s -L -o goreleaser_Linux_x86_64.tar.gz "https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER)/goreleaser_Linux_x86_64.tar.gz"
	cd downloads && tar -xvf goreleaser_Linux_x86_64.tar.gz
	mkdir -p $(INSTALL_DIR) && install downloads/goreleaser $(INSTALL_DIR)
	rm -rf downloads

download-golangci-lint:
	mkdir -p downloads
	cd downloads && curl -s -L -o golangci-lint-linux-amd64.tar.gz "https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI)/golangci-lint-$(GOLANGCI)-linux-amd64.tar.gz"
	cd downloads && tar -xvf golangci-lint-linux-amd64.tar.gz
	mkdir -p $(INSTALL_DIR) && install downloads/golangci-lint-$(GOLANGCI)-linux-amd64/golangci-lint $(INSTALL_DIR)
	rm -rf downloads

protoc-lint: bin/buf
	./bin/buf lint

## protoc-generate-go: create go protoc stubs
protoc-generate-go: bin/buf bin/protoc bin/protoc-gen-go-grpc bin/protoc-gen-grpc-gateway bin/protoc-gen-go
	./bin/buf generate --path=proto/postfinance

bin:
	mkdir bin

bin/protoc: bin
	curl -sLo protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC)/protoc-$(PROTOC)-linux-x86_64.zip
	unzip -o protoc.zip bin/protoc
	touch bin/protoc
	rm protoc.zip

bin/buf: bin
	curl -sLo - https://github.com/bufbuild/buf/releases/download/v$(PROTOC_BUF)/buf-Linux-x86_64.tar.gz | tar -C bin --strip-components=2 -xvzf - buf/bin/buf
	chmod +x bin/buf
	touch bin/buf

bin/protoc-gen-go: bin
	curl -sLo - https://github.com/protocolbuffers/protobuf-go/releases/download/v$(PROTOC_GEN_GO)/protoc-gen-go.v$(PROTOC_GEN_GO).linux.amd64.tar.gz | tar -xvz -C bin
	chmod +x bin/protoc-gen-go
	touch bin/protoc-gen-go

bin/protoc-gen-go-grpc: bin
	curl -sLo - https://github.com/grpc/grpc-go/releases/download/cmd%2Fprotoc-gen-go-grpc%2Fv$(PROTOC_GEN_GRPC_GO)/protoc-gen-go-grpc.v$(PROTOC_GEN_GRPC_GO).linux.amd64.tar.gz | tar -C bin -xvzf - ./protoc-gen-go-grpc
	chmod +x bin/protoc-gen-go-grpc
	touch bin/protoc-gen-go-grpc

bin/protoc-gen-grpc-gateway: bin
	curl -sLo bin/protoc-gen-grpc-gateway https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v$(PROTOC_GEN_GRPC_GATEWAY)/protoc-gen-grpc-gateway-v$(PROTOC_GEN_GRPC_GATEWAY)-linux-x86_64
	chmod +x bin/protoc-gen-grpc-gateway

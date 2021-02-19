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

## build: build the binaries only
build:
	goreleaser build --rm-dist --snapshot

## snapshot: create a snapshot release
snapshot:
	goreleaser release --snapshot --rm-dist --skip-sign

## clean: cleanup
clean:
	rm -rf ./dist
	rm -rf ./downloads

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

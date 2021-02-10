# This file is generated with create-go-app: do not edit.
.PHONY: build clean test snapshot all install help

# special target to export all variables
.EXPORT_ALL_VARIABLES:

PROJECT_RELEASE := 1
DOCKER_IMAGE_PREFIX := local
CI_PROJECT_URL := http://localhost
INSTALL_DIR := $(HOME)/bin
GO_VERSION := $(shell (go version | awk '{print $$3;}'))
BINARIES := $(shell find ./dist ! -name '*goreleaserdocker*' -path '*_linux_*' -type f -executable)

## build: build the binaries only
build:
	goreleaser build --rm-dist --snapshot

## snapshot: create a snapshot release
snapshot:
	goreleaser release --snapshot --rm-dist --skip-sign

## clean: cleanup
clean:
	rm -rf ./dist

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

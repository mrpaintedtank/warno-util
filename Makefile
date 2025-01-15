# Makefile for building the warno-util binary

# go vars
GO_VERSION := 1.23.4
GOPATH ?= $(shell go env GOPATH)

# build vars
VERSION ?= 0.0.1
BINARY_NAME := warno-util
MAIN_FILE := cmd/warno-util/main.go
OUTPUT_DIR := bin

# Build the binary for Windows
GOOS_WINDOWS := windows
GOARCH_WINDOWS := amd64

.PHONY: help
## help: show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.DEFAULT_GOAL := help

.PHONY: all
## all: clean, tidy, format, lint and build the application
all: clean tidy fmt lint build

.PHONY: build
# build: build warno-util for local use
build:
	@echo "Building $(BINARY_NAME) for Windows..."
	set GOOS=$(GOOS_WINDOWS)
	set GOARCH=$(GOARCH_WINDOWS)
	go build -ldflags "-X main.version=$(VERSION)" -o $(OUTPUT_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)
	@echo "Build completed: $(OUTPUT_DIR)/$(BINARY_NAME).exe"

.PHONY: clean
## clean: remove build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)\$(BINARY_NAME).exe
	rm -rf $(OUTPUT_DIR)\$(BINARY_NAME)
	@echo "Clean completed"

.PHONY: run
## run: run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(OUTPUT_DIR)/$(BINARY_NAME).exe
	@echo "Run completed"

.PHONY: tidy
## tidy: tidy Go modules
tidy:
	go mod tidy

.PHONY: fmt
## fmt: format Go source code
fmt:
	go fmt ./...

.PHONY: lint
## lint: run golangci-lint
lint: install_lint
	@echo $(pwd)
	@golangci-lint run cmd/... pkg/...

GOLANGCI_VERSION := v1.63.4
.PHONY: install_lint
## install_lint: install golangci-lint if not present
install_lint:
	@echo "Checking if golangci-lint is installed..."
	@which golangci-lint || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_VERSION))
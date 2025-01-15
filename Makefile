# Makefile for building the warno-util binary

# go vars
GO_VERSION := 1.23.4
GOPATH ?= $(shell go env GOPATH)

# build vars
VERSION ?= dev
BINARY_NAME := warno-util
MAIN_FILE := cmd/warno-util/main.go
OUTPUT_DIR := bin

## TODO goreleaser pipeline
# Build the binary for Windows
GOOS_WINDOWS := windows
GOARCH_WINDOWS := amd64

# Default target
.PHONY: all
all: clean tidy fmt lint build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for Windows..."
	set GOOS=$(GOOS_WINDOWS)
	set GOARCH=$(GOARCH_WINDOWS)
	go build -ldflags "-X main.version=$(VERSION)" -o $(OUTPUT_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)
	@echo "Build completed: $(OUTPUT_DIR)/$(BINARY_NAME).exe"

# Clean the output directory
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)\$(BINARY_NAME).exe
	rm -rf $(OUTPUT_DIR)\$(BINARY_NAME)
	@echo "Clean completed"

# Run the binary
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(OUTPUT_DIR)/$(BINARY_NAME).exe
	@echo "Run completed"

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: install_lint
	@echo $(pwd)
	@golangci-lint run cmd/... pkg/...


GOLANGCI_VERSION := v1.63.4
# checks if golangci-lint is installed, if not installs it
.PHONY: install_lint
install_lint:
	@echo "Checking if golangci-lint is installed..."
	@which golangci-lint || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_VERSION))


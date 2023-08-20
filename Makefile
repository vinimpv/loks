# Makefile for loks-cli

# Variables
BINARY_NAME=loks
GOBUILD=go build
GOTEST=go test
GOCLEAN=go clean
GOGET=go get
REMOVE=rm
BUILDPATH=./build/

# All target is used to build the binary
all: test build

# Build the binary
build:
	$(GOBUILD) -o $(BUILDPATH)$(BINARY_NAME) -v

# Default test target
test:
	$(GOTEST) -v ./...

# Clean the binary
clean:
	$(GOCLEAN)
	$(REMOVE) $(BUILDPATH)$(BINARY_NAME)

# Run the binary
run: build
	$(BUILDPATH)$(BINARY_NAME)

# Get all dependencies
deps:
	$(GOGET) github.com/spf13/cobra
	$(GOGET) go.etcd.io/bbolt


# PHONY is used to specify non-file targets
.PHONY: all build test clean run deps

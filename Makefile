# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=puff
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all build test clean run deps help reload

all: test build

build: ## Build the application
	$(GOBUILD) -o $(BINARY_NAME) -v

test: ## Run tests
	$(GOTEST) -v

clean: ## Clean up the project directory
	$(GOCLEAN)
	rm -f $(BINARY_NAME) \
	rm -f $(BINARY_UNIX) \
    $(GOCMD) mod tidy

run: ## Run demo app locally
	$(GOBUILD) -o $(BINARY_NAME) -v examples/demo.go ./$(BINARY_NAME)

reload: ## Run demo app locally with reload
	air --build.cmd "$(GOBUILD) -o $(BINARY_NAME) examples/demo.go" --build.bin "./$(BINARY_NAME)"

build-linux: ## Build the application for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

help: ## show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <command>\ncommands:\033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

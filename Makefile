# Personal Finance API Makefile

BINARY_NAME=personal-finance-api
MAIN_PATH=cmd/server/main.go
BUILD_DIR=build

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

.PHONY: all build clean test deps run dev docker-build docker-run docker-stop help

all: deps test build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	go test -v ./internal/controllers

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

run: deps
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

dev:
	@echo "Starting development server..."
	air

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(BINARY_NAME) $(BINARY_NAME):latest

docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

help:
	@echo "Available commands:"
	@echo "  all          - Download deps, run tests, and build"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build files"
	@echo "  test         - Run all existing tests"
	@echo "  deps         - Download dependencies"
	@echo "  run          - Run the application"
	@echo "  dev          - Run in development mode"
	@echo "  fmt          - Format code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-stop  - Stop Docker container"
	@echo "  help         - Show this help message"
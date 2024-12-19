# Makefile

# Variables
APP_NAME := main
DOCKER_COMPOSE_FILE := docker-compose.yaml
GO_BUILD_FLAGS := GOOS=linux GOARCH=amd64
BUILD_DIR := build
APP_BINARY := $(BUILD_DIR)/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Build the Go application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GO_BUILD_FLAGS) go build -o $(APP_BINARY) ./cmd/main.go
	chmod +x $(APP_BINARY)

# Run the application locally
.PHONY: run
run: build
	$(APP_BINARY)

# Run tests
.PHONY: test
test:
	go test ./... -v

# Lint the codebase
.PHONY: lint
lint:
	golangci-lint run

# Run Docker Compose
.PHONY: compose-up
compose-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop Docker Compose
.PHONY: compose-down
compose-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Restart Docker Compose
.PHONY: compose-restart
compose-restart: compose-down compose-up

# Tail application logs
.PHONY: logs
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Build Docker image for the application
.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):latest .

# Help menu
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  all              Build the application (default)"
	@echo "  clean            Remove build artifacts"
	@echo "  build            Compile the application"
	@echo "  run              Run the application locally"
	@echo "  test             Run unit tests"
	@echo "  lint             Lint the codebase"
	@echo "  compose-up       Start Docker Compose services"
	@echo "  compose-down     Stop Docker Compose services"
	@echo "  compose-restart  Restart Docker Compose services"
	@echo "  logs             Tail application logs"
	@echo "  docker-build     Build Docker image"
	@echo "  help             Display this help message"

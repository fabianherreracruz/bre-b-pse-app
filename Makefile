.PHONY: help build run test clean docker-up docker-down

help:
	@echo "BRE-B PSE Recaudos App - Available Commands"
	@echo "=========================================="
	@echo "make build         - Build the application"
	@echo "make run           - Run the application"
	@echo "make test          - Run tests"
	@echo "make clean         - Clean build files"
	@echo "make docker-up     - Start Docker services"
	@echo "make docker-down   - Stop Docker services"
	@echo "make docker-build  - Build Docker image"
	@echo "make dev           - Run in development mode"

build:
	@echo "Building application..."
	go build -o bre-b-pse-app ./cmd
	@echo "✅ Build complete"

run: build
	@echo "Running application..."
	./bre-b-pse-app

dev:
	@echo "Running in development mode..."
	go run ./cmd

test:
	@echo "Running tests..."
	go test -v -cover ./...

clean:
	@echo "Cleaning build files..."
	rm -f bre-b-pse-app
	rm -f *.xlsx
	@echo "✅ Clean complete"

docker-build:
	@echo "Building Docker image..."
	docker build -t bre-b-pse-app:latest .
	@echo "✅ Docker image built"

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "✅ Services started"
	@echo "API available at http://localhost:8080"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped"

docker-logs:
	docker-compose logs -f app

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

lint:
	@echo "Running linter..."
	golangci-lint run ./...

migrate:
	@echo "Running migrations..."
	go run ./cmd

env-setup:
	@echo "Setting up .env file..."
	cp .env.example .env
	@echo "✅ .env file created. Please configure it with your credentials."

.PHONY: help build run dev clean deps test

help:
	@echo "Sterling HMS Backend - Makefile Commands"
	@echo "========================================"
	@echo "make deps      - Download dependencies"
	@echo "make build     - Build the application"
	@echo "make run       - Run the built application"
	@echo "make dev       - Run in development mode with hot reload"
	@echo "make clean     - Remove build artifacts"
	@echo "make test      - Run tests"

deps:
	go mod download
	go mod tidy

build: deps
	@echo "Building application..."
	go build -o sterling-hms-backend cmd/main.go
	@echo "Build complete!"

run: build
	@echo "Starting server..."
	./sterling-hms-backend

dev: deps
	@echo "Installing air for hot reload..."
	go install github.com/cosmtrek/air@latest
	@echo "Starting in development mode with hot reload..."
	air

clean:
	@echo "Cleaning up..."
	rm -f sterling-hms-backend sterling-hms-backend.exe
	go clean
	@echo "Clean complete!"

test:
	go test ./...

db-init:
	psql -U postgres -d sterling_hms -f database/schema.sql
	@echo "Database initialized!"

db-create:
	psql -U postgres -c "CREATE DATABASE sterling_hms;" 2>/dev/null || echo "Database already exists"
	$(MAKE) db-init

# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o ./bin/pm4devs cmd/api/main.go

# Run the application
run: build
	./bin/pm4devs

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f ./bin/pm4devs

db-up:
	./database/pocketbase serve

# Live Reload
dev:
	air --build.cmd "go build -o bin/pm4devs cmd/api/main.go" --build.bin "./bin/pm4devs"

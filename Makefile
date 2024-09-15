# Variables for reuse
BINARY=bin/pm4devs-backend
CMD=cmd/main.go
DOCKER_COMPOSE=docker compose

# Targets
db-up:
	$(DOCKER_COMPOSE) up -d

db-down:
	$(DOCKER_COMPOSE) down

build:
	go build -o $(BINARY) $(CMD)

dev:
	air --build.cmd "go build -o $(BINARY) $(CMD)" --build.bin "./$(BINARY)"

run: build
	./$(BINARY)

setup:
	$(MAKE) db-up
	$(MAKE) build
	@go mod tidy
	$(MAKE) run

quick-setup:
	@echo "Setting Envs, database & build files"
	@source ./example.envrc && \
	goose up \
	go mod tidy && \
	$(MAKE) db-up && \
	$(MAKE) build && \
	./$(BINARY)

test:
	@go test ./...

test-v:
	@go test -v ./...


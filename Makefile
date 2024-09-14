db-up:
	docker compose up -d

db-down:
	docker compose down

build:
	go build -o bin/pm4devs-backend cmd/main.go

dev:
	air --build.cmd "go build -o bin/pm4devs-backend cmd/main.go" --build.bin "./bin/pm4devs-backend"

run: build
	./bin/pm4devs-backend

setup: db-up
	go mod tidy
	$(MAKE) build
	./bin/pm4devs-backend

quick-setup:
	@source ./example.envrc && \
	go mod tidy && \
	$(MAKE) db-up && \
	$(MAKE) build && \
	./bin/pm4devs-backend

test:
	@go test ./...

test-v:
	@go test -v ./...

db-up:
	docker compose up -d

db-down:
	docker compose down

build:
	go build -o bin/pm4devs-backend

dev:
	air

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

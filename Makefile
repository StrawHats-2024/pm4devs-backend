db-up:
	docker compose up -d

db-down:
	docker compose up -d

build:
	go build -o bin/pm4devs-backend

run: build
	./bin/pm4devs-backend

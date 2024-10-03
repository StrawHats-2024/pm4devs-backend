include .env

# Create API DSN
DSN="postgres://${DB_USER}:${DB_PASSWORD}@localhost/${DB_NAME}?sslmode=disable"
# DSN="postgres://postgres:passwrd@localhost/postgres?sslmode=disable"

# Create Tests DSN
TEST_DSN="postgres://${TEST_DB_USER}:${TEST_DB_PASSWORD}@localhost:5433/${TEST_DB_NAME}?sslmode=disable"

# Create Tests DSN for Docker
TEST_DSN_DOCKER="postgres://${TEST_DB_USER}:${TEST_DB_PASSWORD}@localhost/${TEST_DB_NAME}?sslmode=disable"

# ================================================================================ #
# HELPERS
# ================================================================================ #

# Helper command to print instructions when `make` or `make help` is run
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# ================================================================================ #
# DEVELOPMENT
# ================================================================================ #

## run: run API
.PHONY: run
run: db/start mig/up
	@echo "Starting server..."
	@go run ./cmd/api \
		-env=${ENV} \
		-port=${PORT} \
		-db-dsn=${DSN} \
		-smtp-host=${SMTP_HOST} \
		-smtp-port=${SMTP_PORT} \
		-smtp-username=${SMTP_USERNAME} \
		-smtp-password=${SMTP_PASSWORD} \
		-smtp-sender=${SMTP_SENDER}
	
## tests: run tests
#
#  note: -p=1 is used to force all packages to run serially to avoid database issues
#  note: flag parsing isn't working so an environment variable is used (see mocks/config.go)
.PHONY: tests
tests: db/start/tests
	@TEST_DSN=${TEST_DSN} go test -p=1 ./...

## tests/short: run tests skipping integration
#
#  The flags are required but the test database should never be accessed in short tests
.PHONY: tests/short
tests/short:
	@go test -short ./...

## tests/cover: run tests with code coverage
.PHONY: tests/cover
tests/cover: db/start/tests
	@TEST_DSN=${TEST_DSN} go test -cover -p=1 ./...


# ================================================================================ #
# DATABASE
# ================================================================================ #

## db/start: start the API database
.PHONY: db/start
db/start: db/stop/tests
	@echo "Starting API database..."
	@docker compose -p ${PROJECT_NAME} up -d postgres
	@while ! docker exec postgres pg_isready -h postgres -p 5432; do \
		echo "Waiting for PostgreSQL to start..."; \
		sleep 1; \
	done

## db/start/tests: start the Tests database
.PHONY: db/start/tests
db/start/tests: db/stop
	@echo "Starting Tests database..."
	@docker compose -p ${PROJECT_NAME} up -d postgres-tests
	@while ! docker exec postgres-tests pg_isready -h postgres-tests -p 5432; do \
		echo "Waiting for PostgreSQL to start..."; \
		sleep 1; \
	done
	@sleep 1

## db/stop: stop the API database
.PHONY: db/stop
db/stop:
	@echo "Stopping API database..."
	@docker compose -p ${PROJECT_NAME} stop postgres

## db/stop/tests: stop the Tests database
.PHONY: db/stop/tests
db/stop/tests:
	@echo "Stopping Tests database..."
	@docker compose -p ${PROJECT_NAME} stop postgres-tests


# ================================================================================ #
# PSQL
# ================================================================================ #

## sql: connect to the API database with psql
.PHONY: sql
sql: db/start
	@echo "Connecting to database..."
	@echo "Connected. Type 'exit' to exit."
	@docker exec -it postgres psql ${DSN}

## sql/tests: connect to the Tests database with psql
.PHONY: sql/tests
sql/tests: db/start/tests
	@echo "Connecting to database..."
	@echo "Connected. Type 'exit' to exit."
	@docker exec -it postgres-tests psql ${TEST_DSN_DOCKER}


# ================================================================================ #
# MIGRATIONS
# ================================================================================ #

## mig/new name=$1: create a new database migration
.PHONY: mig/new
mig/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
## mig/up: migrate to a specific version, or apply all migrations
.PHONY: mig/up
mig/up:
	@echo 'Running up migrations...'
	@migrate -path ./migrations -database ${DSN} up

## mig/test/up: migrate to a specific version, or apply all migrations
.PHONY: mig/test/up
mig/test/up:
	@echo 'Running up migrations...'
	@migrate -path ./migrations -database ${TEST_DSN} up

## mig/test/down: migrate to a specific version, or apply all migrations
.PHONY: mig/test/down
mig/test/down:
	@echo 'Running down migration...'
	@migrate -path ./migrations -database ${TEST_DSN} down

## mig/down: apply all down database migrations
.PHONY: mig/down
mig/down:
	@echo 'Running down migration...'
	@migrate -path ./migrations -database ${DSN} down

## mig/force version=$1: force the database to a migration version
#
#  If you have a bad migration, force to the highest version then run `mig/down`
.PHONY: mig/force
mig/force:
	@echo 'Forcing version to ${version}...'
	@migrate -path ./migrations -database ${DSN} force ${version}


# ================================================================================ #
# UTILS
# ================================================================================ #

## util/loc: lists the total lines of code
.PHONY: util/loc
util/loc:
	@go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}{{"\n"}}{{end}}' ./... | xargs wc -l | sort -n


# ================================================================================ #
# BUILD
# ================================================================================ #

## build: build the API
#
#  -ldflags='s' is used to strip symbol tables and DWARF debugging information
.PHONY: build
build:
	@echo "Building API..."
	@go build -ldflags='-s' -o=./bin/api ./cmd/api

## version: Output version of current binary
#
#  Requires a binary to exist at ./bin/api
.PHONY: version
version:
	@./bin/api -version

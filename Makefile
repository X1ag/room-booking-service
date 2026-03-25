.PHONY: env up down logs db-up db-down test-unit test-e2e test-with-db test-all test-cover

GO_TEST_ENV = set -a; \
	[ -f .env ] && . ./.env; \
	set +a; \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_SSLMODE=disable \
	MIGRATIONS_PATH=file://migrations

env:
	@if [ ! -f .env ]; then cp .env.example .env; fi

up: env
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

db-up: env
	docker compose up -d --wait postgres

db-down:
	docker compose stop postgres

test-unit:
	go test ./internal/... ./cmd/... -count=1

test-e2e: db-up
	@$(GO_TEST_ENV) go test ./tests/e2e -v -count=1

test-with-db: db-up
	@$(GO_TEST_ENV) go test ./... -count=1

test-all: test-with-db

test-cover: db-up
	@$(GO_TEST_ENV) go test ./... -coverprofile=coverage.out -count=1


.PHONY: env up down logs db-up db-down lint test-unit test-e2e test-with-db test-all test-cover

GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(shell command -v golangci-lint 2>/dev/null || echo "$(GOBIN)/golangci-lint")

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

lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "golangci-lint not found"; \
		echo "install it with: brew install golangci-lint"; \
		echo "or: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run ./...

test-unit:
	go test ./internal/... ./cmd/... -count=1

test-e2e: db-up
	@$(GO_TEST_ENV) go test ./tests/e2e -v -count=1

test-with-db: db-up
	@$(GO_TEST_ENV) go test ./... -count=1

test-all: test-with-db

test-cover: db-up
	@$(GO_TEST_ENV) go test ./... -coverprofile=coverage.out -count=1

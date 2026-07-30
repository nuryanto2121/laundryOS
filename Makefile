# =============================================================================
# LaundryOS — Makefile
# =============================================================================

.PHONY: up down restart logs ps build test clean \
        lint lint-api lint-gateway lint-dashboard \
        fmt fmt-api fmt-gateway fmt-dashboard \
        api gateway dashboard \
        migrate-up migrate-down

# -----------------------------------------------------------------------------
# Docker Compose
# -----------------------------------------------------------------------------

## First-time setup: create .env and install all local dependencies
bootstrap:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env from .env.example — edit it and set your secrets."; \
	fi
	cd apps/gateway && npm install
	cd apps/dashboard && npm install

## Start all services in detached mode
up:
	@if [ ! -f .env ]; then \
		echo ""; \
		echo "  ERROR: .env file not found."; \
		echo "  Run: cp .env.example .env"; \
		echo "  Then edit .env and set your secrets before retrying."; \
		echo ""; \
		exit 1; \
	fi
	docker compose up -d --build

## Stop all services
down:
	docker compose down

## Restart all services
restart:
	docker compose down && docker compose up -d --build

## Follow logs for all services
logs:
	docker compose logs -f

## Show running container status
ps:
	docker compose ps

## Build all Docker images without starting containers
build:
	docker compose build

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------

## Run all tests across all services
test: test-api

## Run Go API tests
test-api:
	cd apps/api && go test ./...

# DB_URL constructed from apps/api/.env variables for local migration runs.
DB_URL ?= postgres://postgres:secret@localhost:5432/laundryos?sslmode=disable

## Run database migrations (up)
migrate-up:
	migrate -path apps/api/migrations -database "$(DB_URL)" up

## Roll back database migrations (down)
migrate-down:
	migrate -path apps/api/migrations -database "$(DB_URL)" down

# Node 22 LTS path (keg-only on macOS — takes priority over system node)
NODE22_BIN := /usr/local/opt/node@22/bin
export PATH := $(NODE22_BIN):$(PATH)

# -----------------------------------------------------------------------------
# Lint
# -----------------------------------------------------------------------------

## Run all linters
lint: lint-api lint-gateway lint-dashboard

## Lint Go API (requires golangci-lint — install: brew install golangci-lint)
lint-api:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo ""; \
		echo "  ERROR: golangci-lint is not installed."; \
		echo "  Install it with: brew install golangci-lint"; \
		echo "  Or visit: https://golangci-lint.run/welcome/install/"; \
		echo ""; \
		exit 1; \
	fi
	cd apps/api && golangci-lint run ./...

## Lint Node Gateway
lint-gateway:
	cd apps/gateway && npx eslint src/

## Lint Vue Dashboard
lint-dashboard:
	cd apps/dashboard && npx eslint src/

# -----------------------------------------------------------------------------
# Format
# -----------------------------------------------------------------------------

## Format all services
fmt: fmt-api fmt-gateway fmt-dashboard

## Format Go API
fmt-api:
	cd apps/api && gofmt -w .

## Format Node Gateway
fmt-gateway:
	cd apps/gateway && npx prettier --write src/

## Format Vue Dashboard
fmt-dashboard:
	cd apps/dashboard && npx prettier --write src/

# -----------------------------------------------------------------------------
# Local dev (without Docker)
# -----------------------------------------------------------------------------

## Run Go API locally
api:
	cd apps/api && go run ./cmd/api

## Run Node Gateway locally
gateway:
	cd apps/gateway && node src/index.js

## Run Vue Dashboard locally
dashboard:
	cd apps/dashboard && npm run dev

# -----------------------------------------------------------------------------
# Clean
# -----------------------------------------------------------------------------

## Remove containers, volumes, and local build artifacts
clean:
	docker compose down -v --remove-orphans
	rm -f apps/api/bin/api
	rm -rf apps/dashboard/dist
	rm -rf apps/gateway/node_modules
	rm -rf apps/dashboard/node_modules

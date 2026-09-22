.DEFAULT_GOAL := all
.PHONY: *

help:
	@echo "help"

all: up_prod

init:
	cd backend && \
	go mod download && \
	migrate

generate: queries wire

beautify:
	cd backend && \
	go mod tidy && \
	go fmt ./...

ci: lint test build

up:
	docker compose -f compose.yaml up -d --build

down:
	docker compose -f compose.yaml down --remove-orphans

up_prod:
	docker compose -f compose.yaml -f compose.prod.yaml up -d --build

down_prod:
	docker compose -f compose.yaml -f compose.prod.yaml down --remove-orphans

up_test:
	docker compose -f compose.yaml -f compose.test.yaml up -d --build

down_test:
	docker compose -f compose.yaml -f compose.test.yaml down --remove-orphans

run:
	cd backend && go run ./cmd/api/main.go

run_dev:
	docker compose exec -itu root: app bash -c '/go/bin/dlv --listen=:40040 --headless=true --api-version=2 --accept-multiclient debug cmd/dev/main.go'

run_tests:
	docker compose exec -itu root: app bash -c 'go test ./...'

protos_add:
	git submodule add https://github.com/plezhaspace/protos.git

protos_update:
	git submodule update --remote

# make MIGRATION_NAME="migrationName" migration_create
migration_create:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable create -ext sql -dir /migrations -seq $(MIGRATION_NAME)

migration_up:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable up

migration_down:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable down -all

# make MIGRATION_VERSION="1" migration_force
migration_force:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable force $(MIGRATION_VERSION)

# queries:
# 	cd ./backend && sqlc generate
queries:
	cd ./backend && go tool sqlc generate

wire:
	cd ./backend && wire ./internal/platform/di
.PHONY: *
.DEFAULT_GOAL := run

api_port := 8080
api_url := localhost:$(api_port)
front_port := 3000
front_url := localhost:$(front_port)

test:
	cd backend && go test ./...

run:
	cd backend && APP_ENV=local go run ./cmd/api/main.go -config=../../backend/configs/config.yaml -env=../../backend/.env

init: generate_keys generate_code
	cd backend && \
	go mod download && \
	migrate

init_local: init
	cd backend && \
	[ ! -f ".env" ] && cp ".env.example" ".env"

generate_keys:
	cd backend && \
	openssl genpkey -algorithm Ed25519 -out ./keys/private.pem && \
	openssl pkey -in ./keys/backend/private.pem -pubout -out ./keys/public.pem

generate_code:
	cd backend && \
	go tool sqlc generate && \
	go tool wire ./cmd/consumer/ && \
	go tool wire ./internal/app/api/ && \
	go tool wire ./internal/platform/di/ && \
	go generate ./... && \
	go mod tidy

lint:
	cd backend && \
	go vet ./... && \
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

# name=foo make migration_create
migration_create:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable create -ext sql -dir /migrations -seq $(name)

migration_up:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable up

migration_down:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable down -all

# v=1 make migration_force
migration_force:
	docker run -v $(shell pwd)/backend/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://app:secret@localhost:5432/app?sslmode=disable force $(v)

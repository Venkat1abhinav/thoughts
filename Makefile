include .envrc

MIGRATIONS_PATH := ./cmd/migrate/migrations

DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: run build test fmt vet migrate-create migrate-up migrate-down migrate-force db-up db-down db-reset seed

run:
	go run ./cmd/api

build:
	go build -o ./bin/api ./cmd/api

test:
	go test ./...

fmt:
	gofumpt -w .

vet:
	go vet ./...

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" up

migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" down 1

migrate-force:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" force $(version)

db-up:
	podman compose up -d

db-down:
	podman compose down

db-reset:
	podman compose down -v
	podman compose up -d


seed:
	go run cmd/migrate/seed/main.go

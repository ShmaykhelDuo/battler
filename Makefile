POSTGRES_URL := postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable
REDIS_URL := redis://localhost:6379/0

.PHONY: start-dev
start-dev:
	docker compose -f compose.dev.yaml up -d

.PHONY: stop-dev
stop-dev:
	docker compose -f compose.dev.yaml down -v

.PHONY: new-migration
new-migration:
	goose create -dir migrations $(name) sql

.PHONY: migrate-up
migrate-up:
	goose postgres $(POSTGRES_URL) up -dir migrations $(name)

.PHONY: run
run:
	DB_CONN=$(POSTGRES_URL) CACHE_URL=$(REDIS_URL) go run ./cmd/app
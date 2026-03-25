DB_DSN ?= postgres://katerina:katerina@localhost:5432/katerina?sslmode=disable
MIGRATE = go run ./cmd/migrate -dsn "$(DB_DSN)" -dir migrations

.PHONY: up down build rebuild logs restart clean setup pull-model tidy \
        migrate migrate-down migrate-reset migrate-status migrate-create

## Start all services (detached)
up:
	docker compose up -d

## Stop all services
down:
	docker compose down

## Rebuild all images and restart
build:
	docker compose build

## Rebuild images and restart containers (use after Go code changes)
rebuild:
	docker compose build && docker compose up -d

## Follow logs (all services)
logs:
	docker compose logs -f

## Restart all services
restart:
	docker compose restart

## Remove containers, networks, and volumes (destructive)
clean:
	docker compose down -v --remove-orphans

## First-time setup: start stack + pull the LLM model
setup: up
	@echo "Waiting for Ollama to be ready..."
	@until docker compose exec ollama curl -sf http://localhost:11434/api/tags > /dev/null 2>&1; do sleep 2; done
	docker compose exec ollama ollama pull llama3
	@echo ""
	@echo "Setup complete. Open http://localhost:5173"
	@echo "Demo credentials: demo@katerina.local / password123"

## Pull / update the LLM model
pull-model:
	docker compose exec ollama ollama pull llama3

## Download Go dependencies (run once locally)
tidy:
	go mod tidy

## Apply all pending migrations
migrate:
	$(MIGRATE) up

## Roll back the last migration
migrate-down:
	$(MIGRATE) down

## Roll back all migrations
migrate-reset:
	$(MIGRATE) reset

## Show migration status
migrate-status:
	$(MIGRATE) status

## Create a new migration: make migrate-create name=add_something
migrate-create:
	$(MIGRATE) create $(name)

.PHONY: up down build rebuild logs restart clean setup pull-model tidy

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

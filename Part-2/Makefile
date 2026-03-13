.PHONY: help build up down logs clean dev test

# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml

help: ## Show help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build Docker images
	$(DOCKER_COMPOSE) build

up: ## Start application in production
	$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "Application started"
	@echo "   Frontend: http://localhost:3000"
	@echo "   Backend:  http://localhost:8080"

down: ## Stop application
	$(DOCKER_COMPOSE) down

logs: ## View logs from all services
	$(DOCKER_COMPOSE) logs -f

logs-backend: ## View backend logs
	$(DOCKER_COMPOSE) logs -f backend

logs-frontend: ## View frontend logs
	$(DOCKER_COMPOSE) logs -f frontend

clean: ## Clean containers, volumes and images
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

dev: ## Start in development mode
	$(DOCKER_COMPOSE_DEV) up --build

dev-down: ## Stop development mode
	$(DOCKER_COMPOSE_DEV) down

rebuild: ## Rebuild without cache
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

restart: ## Restart services
	$(DOCKER_COMPOSE) restart

ps: ## View container status
	$(DOCKER_COMPOSE) ps

health: ## Check backend health
	@curl -s http://localhost:8080/api/health | python -m json.tool || echo "Backend not available"

test-backend: ## Run backend tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test

install-backend: ## Install backend dependencies
	cd backend && go mod download

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

run-backend-local: ## Run backend locally
	cd backend && go run .

run-frontend-local: ## Run frontend locally
	cd frontend && npm run serve

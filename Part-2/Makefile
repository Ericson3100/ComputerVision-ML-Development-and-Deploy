.PHONY: help build up down logs clean dev test

# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Construir imágenes Docker
	$(DOCKER_COMPOSE) build

up: ## Iniciar aplicación en producción
	$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "✅ Aplicación iniciada"
	@echo "   Frontend: http://localhost:3000"
	@echo "   Backend:  http://localhost:8080"

down: ## Detener aplicación
	$(DOCKER_COMPOSE) down

logs: ## Ver logs de todos los servicios
	$(DOCKER_COMPOSE) logs -f

logs-backend: ## Ver logs del backend
	$(DOCKER_COMPOSE) logs -f backend

logs-frontend: ## Ver logs del frontend
	$(DOCKER_COMPOSE) logs -f frontend

clean: ## Limpiar contenedores, volúmenes e imágenes
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

dev: ## Iniciar en modo desarrollo
	$(DOCKER_COMPOSE_DEV) up --build

dev-down: ## Detener modo desarrollo
	$(DOCKER_COMPOSE_DEV) down

rebuild: ## Reconstruir sin caché
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

restart: ## Reiniciar servicios
	$(DOCKER_COMPOSE) restart

ps: ## Ver estado de contenedores
	$(DOCKER_COMPOSE) ps

health: ## Verificar health del backend
	@curl -s http://localhost:8080/api/health | python -m json.tool || echo "Backend no disponible"

test-backend: ## Ejecutar tests del backend
	cd backend && go test ./...

test-frontend: ## Ejecutar tests del frontend
	cd frontend && npm test

install-backend: ## Instalar dependencias del backend
	cd backend && go mod download

install-frontend: ## Instalar dependencias del frontend
	cd frontend && npm install

run-backend-local: ## Ejecutar backend localmente
	cd backend && go run .

run-frontend-local: ## Ejecutar frontend localmente
	cd frontend && npm run serve

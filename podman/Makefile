.PHONY: help podman-up podman-down podman-restart podman-logs podman-clean mongo-shell mongo-seed mongo-logs

# Get the directory where this Makefile is located
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Compose file location (relative to Makefile location)
COMPOSE_FILE := $(MAKEFILE_DIR)compose.yml

# Detect if podman-compose is available, otherwise use docker compose
COMPOSE_CMD := $(shell which podman-compose >/dev/null 2>&1 && echo "podman-compose" || echo "docker compose")

help: ## Show this help message
	@echo 'Usage: make -f podman/Makefile [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

podman-up: ## Start MongoDB container
	@echo "Starting Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) up -d
	@echo ""
	@echo "✅ MongoDB started successfully!"
	@echo ""
	@echo "📦 MongoDB is available at: localhost:27017"
	@echo "   Username: admin"
	@echo "   Password: admin123"
	@echo "   Database: pharmacy_modernization"
	@echo "   Connection: mongodb://admin:admin123@localhost:27017/pharmacy_modernization"

podman-down: ## Stop MongoDB container
	@echo "Stopping Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) down

podman-restart: ## Restart MongoDB container
	@echo "Restarting Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) restart

podman-logs: ## Show logs from MongoDB
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) logs -f

podman-clean: ## Stop container and remove volumes (WARNING: deletes all data)
	@echo "⚠️  This will remove all data from MongoDB volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(COMPOSE_CMD) -f $(COMPOSE_FILE) down -v; \
		echo "✅ Container and volumes removed"; \
	else \
		echo "❌ Cancelled"; \
	fi

# MongoDB service commands
mongo-logs: ## Show MongoDB logs
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) logs -f mongodb

mongo-shell: ## Connect to MongoDB shell
	podman exec -it mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

mongo-seed: ## Seed MongoDB with sample data
	@echo "🌱 Seeding MongoDB with sample patient data..."
	@cd $(MAKEFILE_DIR).. && go run ./cmd/seed

# Legacy Docker aliases (for backwards compatibility)
docker-up: podman-up ## Alias for podman-up (backwards compatibility)
docker-down: podman-down ## Alias for podman-down (backwards compatibility)
docker-restart: podman-restart ## Alias for podman-restart (backwards compatibility)
docker-logs: podman-logs ## Alias for podman-logs (backwards compatibility)
docker-clean: podman-clean ## Alias for podman-clean (backwards compatibility)

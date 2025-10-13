.PHONY: help podman-up podman-down podman-restart podman-logs podman-clean mongo-ui

# Compose file location
COMPOSE_FILE := compose.yml

# Detect if podman-compose is available, otherwise use docker compose
COMPOSE_CMD := $(shell which podman-compose >/dev/null 2>&1 && echo "podman-compose" || echo "docker compose")

help: ## Show this help message
	@echo 'Usage: make -f podman/Makefile [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

podman-up: ## Start all Podman containers
	@echo "Starting Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) up -d
	@echo ""
	@echo "✅ Containers started successfully!"
	@echo ""
	@echo "📦 MongoDB is available at: localhost:27017"
	@echo "   Username: admin"
	@echo "   Password: admin123"
	@echo ""
	@echo "🌐 Mongo Express UI is available at: http://localhost:8081"
	@echo "   Username: admin"
	@echo "   Password: admin123"
	@echo ""
	@echo "🔴 Redis is available at: localhost:6379"
	@echo "   Password: redis123"
	@echo ""
	@echo "🌐 Redis Commander UI is available at: http://localhost:8082"
	@echo "   Username: admin"
	@echo "   Password: admin123"
	@echo ""
	@echo "💾 Memcached is available at: localhost:11211"
	@echo "   No authentication required"

podman-down: ## Stop all Podman containers
	@echo "Stopping Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) down

podman-restart: ## Restart all Podman containers
	@echo "Restarting Podman containers..."
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) restart

podman-logs: ## Show logs from all containers
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) logs -f

podman-clean: ## Stop containers and remove volumes (WARNING: deletes all data)
	@echo "⚠️  This will remove all data from Podman volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(COMPOSE_CMD) -f $(COMPOSE_FILE) down -v; \
		echo "✅ Containers and volumes removed"; \
	else \
		echo "❌ Cancelled"; \
	fi

mongo-ui: ## Open MongoDB UI in browser
	@echo "Opening Mongo Express UI..."
	@open http://localhost:8081 || xdg-open http://localhost:8081 || echo "Please open http://localhost:8081 in your browser"

# Individual service commands
mongo-logs: ## Show MongoDB logs
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) logs -f mongodb

mongo-shell: ## Connect to MongoDB shell
	podman exec -it rxintake_mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

mongo-seed: ## Seed MongoDB with sample data
	@echo "🌱 Seeding MongoDB with sample patient data..."
	@cd .. && go run ./cmd/seed

redis-cli: ## Connect to Redis CLI
	podman exec -it rxintake_redis redis-cli -a redis123

redis-ui: ## Open Redis Commander UI in browser
	@echo "Opening Redis Commander UI..."
	@open http://localhost:8082 || xdg-open http://localhost:8082 || echo "Please open http://localhost:8082 in your browser"

memcached-stats: ## Show Memcached stats
	@echo "📊 Memcached Statistics:"
	@echo "stats" | nc localhost 11211 || echo "❌ Could not connect to Memcached. Make sure containers are running."

# Legacy Docker aliases (for backwards compatibility)
docker-up: podman-up ## Alias for podman-up (backwards compatibility)
docker-down: podman-down ## Alias for podman-down (backwards compatibility)
docker-restart: podman-restart ## Alias for podman-restart (backwards compatibility)
docker-logs: podman-logs ## Alias for podman-logs (backwards compatibility)
docker-clean: podman-clean ## Alias for podman-clean (backwards compatibility)

.PHONY: help docker-up docker-down docker-restart docker-logs docker-clean mongo-ui

# Docker compose file location
DOCKER_COMPOSE_FILE := docker-compose.yml

help: ## Show this help message
	@echo 'Usage: make -f docker/Makefile [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

docker-up: ## Start all Docker containers
	@echo "Starting Docker containers..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d
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

docker-down: ## Stop all Docker containers
	@echo "Stopping Docker containers..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

docker-restart: ## Restart all Docker containers
	@echo "Restarting Docker containers..."
	docker compose -f $(DOCKER_COMPOSE_FILE) restart

docker-logs: ## Show logs from all containers
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-clean: ## Stop containers and remove volumes (WARNING: deletes all data)
	@echo "⚠️  This will remove all data from Docker volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose -f $(DOCKER_COMPOSE_FILE) down -v; \
		echo "✅ Containers and volumes removed"; \
	else \
		echo "❌ Cancelled"; \
	fi

mongo-ui: ## Open MongoDB UI in browser
	@echo "Opening Mongo Express UI..."
	@open http://localhost:8081 || xdg-open http://localhost:8081 || echo "Please open http://localhost:8081 in your browser"

# Individual service commands
mongo-logs: ## Show MongoDB logs
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f mongodb

mongo-shell: ## Connect to MongoDB shell
	docker exec -it rxintake_mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

mongo-seed: ## Seed MongoDB with sample data
	@echo "🌱 Seeding MongoDB with sample patient data..."
	@cd .. && go run ./cmd/seed

redis-cli: ## Connect to Redis CLI
	docker exec -it rxintake_redis redis-cli -a redis123

redis-ui: ## Open Redis Commander UI in browser
	@echo "Opening Redis Commander UI..."
	@open http://localhost:8082 || xdg-open http://localhost:8082 || echo "Please open http://localhost:8082 in your browser"

memcached-stats: ## Show Memcached stats
	@echo "📊 Memcached Statistics:"
	@echo "stats" | nc localhost 11211 || echo "❌ Could not connect to Memcached. Make sure containers are running."


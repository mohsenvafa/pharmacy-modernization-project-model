# Docker Local Development Setup

This directory contains Docker configurations for local development.

## Prerequisites

- Docker Desktop installed and running
- Make utility installed

## Quick Start

### Start MongoDB with UI

```bash
make -f docker/Makefile docker-up
```

This will start:
- **MongoDB** on `localhost:27017`
- **Mongo Express UI** on `http://localhost:8081`
- **Redis** on `localhost:6379`
- **Redis Commander UI** on `http://localhost:8082`
- **Memcached** on `localhost:11211`

### Default Credentials

**MongoDB:**
- Username: `admin`
- Password: `admin123`
- Connection String: `mongodb://admin:admin123@localhost:27017/`
- Database: `rxintake`

**Mongo Express UI:**
- URL: http://localhost:8081
- Username: `admin`
- Password: `admin123`

**Redis:**
- Host: `localhost`
- Port: `6379`
- Password: `redis123`
- Connection String: `redis://:redis123@localhost:6379`

**Redis Commander UI:**
- URL: http://localhost:8082
- Username: `admin`
- Password: `admin123`

**Memcached:**
- Host: `localhost`
- Port: `11211`
- Connection String: `localhost:11211`
- No authentication required

## Available Commands

```bash
make -f docker/Makefile help           # Show all available commands
make -f docker/Makefile docker-up      # Start all containers
make -f docker/Makefile docker-down    # Stop all containers
make -f docker/Makefile docker-restart # Restart all containers
make -f docker/Makefile docker-logs    # View logs from all containers
make -f docker/Makefile docker-clean   # Stop containers and remove all data
make -f docker/Makefile mongo-ui       # Open MongoDB UI in browser
make -f docker/Makefile mongo-logs     # Show MongoDB logs
make -f docker/Makefile mongo-shell    # Connect to MongoDB shell
make -f docker/Makefile mongo-seed     # Seed MongoDB with sample patient data
make -f docker/Makefile redis-cli      # Connect to Redis CLI
make -f docker/Makefile redis-ui       # Open Redis Commander UI in browser
make -f docker/Makefile memcached-stats # Show Memcached statistics
```

## Services

### MongoDB
- **Port:** 27017
- **Version:** 7.0
- **Data:** Persisted in Docker volume `mongodb_data`

### Mongo Express
- **Port:** 8081
- **Purpose:** Web-based MongoDB admin interface
- **Features:** Browse databases, collections, execute queries

### Redis
- **Port:** 6379
- **Version:** 7 (Alpine)
- **Data:** Persisted in Docker volume `redis_data`
- **Authentication:** Password protected

### Redis Commander
- **Port:** 8082
- **Purpose:** Web-based Redis admin interface
- **Features:** Browse keys, execute commands, view data structures

### Memcached
- **Port:** 11211
- **Version:** 1.6 (Alpine)
- **Memory:** 64MB (configurable)
- **Authentication:** None (default Memcached behavior)

## Future Services

This setup can be extended to include:
- Kafka (with Kafka UI)
- PostgreSQL
- RabbitMQ
- Elasticsearch
- Other development dependencies

## Seeding Data

After starting MongoDB, seed it with sample patient data:

```bash
make -f docker/Makefile mongo-seed
```

This will populate the database with:
- **15 patients** (P001-P015)
- **13 addresses** (multiple patients have multiple addresses)
- **27 prescriptions** (various statuses: Active, Paused, Completed, Draft)

You can view the data in:
- Your application at http://localhost:8080
- Mongo Express UI at http://localhost:8081

## Notes

- Data is persisted in Docker volumes, so it survives container restarts
- Use `docker-clean` to completely reset all data (requires confirmation)
- All services run on a dedicated network `rxintake_network`
- The seed command will clear existing patients before inserting new ones


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

### Default Credentials

**MongoDB:**
- Username: `admin`
- Password: `admin123`
- Connection String: `mongodb://admin:admin123@localhost:27017/`

**Mongo Express UI:**
- URL: http://localhost:8081
- Username: `admin`
- Password: `admin123`

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

## Future Services

This setup can be extended to include:
- Kafka
- Redis
- Memcached
- Other development dependencies

## Notes

- Data is persisted in Docker volumes, so it survives container restarts
- Use `docker-clean` to completely reset all data (requires confirmation)
- All services run on a dedicated network `rxintake_network`


# Podman Local Development Setup

This directory contains Podman/Compose configurations for local development.

## Prerequisites

- **Podman** installed and running (or Docker as fallback)
- **podman-compose** installed (`pip install podman-compose`)
- Make utility installed

### Installing Podman

```bash
# macOS
brew install podman podman-compose
podman machine init
podman machine start

# Fedora/RHEL/CentOS
sudo dnf install podman podman-compose

# Ubuntu/Debian
sudo apt-get install podman
pip install podman-compose

# Windows
# Download from https://podman.io/getting-started/installation
```

## Quick Start

### Start MongoDB with UI

```bash
make -f podman/Makefile podman-up
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
make -f podman/Makefile help              # Show all available commands
make -f podman/Makefile podman-up         # Start all containers
make -f podman/Makefile podman-down       # Stop all containers
make -f podman/Makefile podman-restart    # Restart all containers
make -f podman/Makefile podman-logs       # View logs from all containers
make -f podman/Makefile podman-clean      # Stop containers and remove all data
make -f podman/Makefile mongo-ui          # Open MongoDB UI in browser
make -f podman/Makefile mongo-logs        # Show MongoDB logs
make -f podman/Makefile mongo-shell       # Connect to MongoDB shell
make -f podman/Makefile mongo-seed        # Seed MongoDB with sample patient data
make -f podman/Makefile redis-cli         # Connect to Redis CLI
make -f podman/Makefile redis-ui          # Open Redis Commander UI in browser
make -f podman/Makefile memcached-stats   # Show Memcached statistics
```

### Legacy Docker Commands

For backwards compatibility, Docker command aliases are available:

```bash
make -f podman/Makefile docker-up         # Alias for podman-up
make -f podman/Makefile docker-down       # Alias for podman-down
# ... (all docker-* commands map to podman-* commands)
```

## Services

### MongoDB
- **Port:** 27017
- **Version:** 7.0
- **Data:** Persisted in Podman volume `mongodb_data`

### Mongo Express
- **Port:** 8081
- **Purpose:** Web-based MongoDB admin interface
- **Features:** Browse databases, collections, execute queries

### Redis
- **Port:** 6379
- **Version:** 7 (Alpine)
- **Data:** Persisted in Podman volume `redis_data`
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
make -f podman/Makefile mongo-seed
```

This will populate the database with:
- **15 patients** (P001-P015)
- **13 addresses** (multiple patients have multiple addresses)
- **27 prescriptions** (various statuses: Active, Paused, Completed, Draft)

You can view the data in:
- Your application at http://localhost:8080
- Mongo Express UI at http://localhost:8081

## Podman vs Docker

This setup works with both Podman and Docker:

### Using Podman (Recommended)
```bash
# The Makefile automatically detects and uses podman-compose
make -f podman/Makefile podman-up
```

### Using Docker (Fallback)
```bash
# If podman-compose is not installed, it falls back to docker compose
make -f podman/Makefile podman-up  # Works with Docker too
# OR use the legacy aliases
make -f podman/Makefile docker-up
```

## Podman-specific Notes

### Rootless Containers
Podman runs containers without root by default, which is more secure:
```bash
# Check if running rootless
podman info | grep rootless
```

### SELinux Context
On systems with SELinux (Fedora, RHEL, CentOS), volumes may need :z or :Z suffix:
```yaml
volumes:
  - mongodb_data:/data/db:z  # For shared volumes
```

### Podman Machine (macOS/Windows)
On macOS and Windows, Podman uses a VM:
```bash
# Start the VM
podman machine start

# Check status
podman machine list

# Stop the VM
podman machine stop
```

## Notes

- Data is persisted in Podman volumes, so it survives container restarts
- Use `podman-clean` to completely reset all data (requires confirmation)
- All services run on a dedicated network `rxintake_network`
- The seed command will clear existing patients before inserting new ones
- Podman is more secure (rootless) and lightweight compared to Docker
- Compatible with Docker compose files (no changes needed)

## Troubleshooting

### Port Conflicts
If ports are already in use:
```bash
# Check what's using the port
sudo lsof -i :27017
# Or on Linux
ss -tulpn | grep 27017
```

### Volume Permissions
If you encounter permission issues:
```bash
# List volumes
podman volume ls

# Inspect volume
podman volume inspect mongodb_data

# Remove and recreate (WARNING: deletes data)
podman volume rm mongodb_data
```

### Container Not Starting
Check logs for specific container:
```bash
podman logs rxintake_mongodb
```

### Reset Everything
```bash
# Stop and remove everything
make -f podman/Makefile podman-clean

# Remove all volumes manually
podman volume prune
```

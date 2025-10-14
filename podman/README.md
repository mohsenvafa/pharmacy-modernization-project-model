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

### Start MongoDB

```bash
make -f podman/Makefile podman-up
```

This will start:
- **MongoDB** on `localhost:27017`

### Default Credentials

**MongoDB:**
- Username: `admin`
- Password: `admin123`
- Database: `pharmacy_modernization`
- Connection String: `mongodb://admin:admin123@localhost:27017/pharmacy_modernization`

## Available Commands

```bash
make -f podman/Makefile help           # Show all available commands
make -f podman/Makefile podman-up      # Start MongoDB container
make -f podman/Makefile podman-down    # Stop MongoDB container
make -f podman/Makefile podman-restart # Restart MongoDB container
make -f podman/Makefile podman-logs    # View MongoDB logs
make -f podman/Makefile podman-clean   # Stop container and remove all data
make -f podman/Makefile mongo-logs     # Show MongoDB logs
make -f podman/Makefile mongo-shell    # Connect to MongoDB shell
make -f podman/Makefile mongo-seed     # Seed MongoDB with sample patient data
```

### Legacy Docker Commands

For backwards compatibility, Docker command aliases are available:

```bash
make -f podman/Makefile docker-up      # Alias for podman-up
make -f podman/Makefile docker-down    # Alias for podman-down
# ... (all docker-* commands map to podman-* commands)
```

## Services

### MongoDB
- **Port:** 27017
- **Version:** 6.0
- **Container Name:** mongodb
- **Data:** Persisted in Podman volume `mongodb_data`
- **Healthcheck:** Automatic ping check every 10 seconds

## Seeding Data

After starting MongoDB, seed it with sample patient data:

```bash
make -f podman/Makefile mongo-seed
```

This will populate the database with:
- **15 patients** (P001-P015)
- **13 addresses** (multiple patients have multiple addresses)
- **27 prescriptions** (various statuses: Active, Paused, Completed, Draft)

You can view the data in your application at http://localhost:8080

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
- All services run on a dedicated network `pharmacy_modernization_network`
- The seed command will clear existing patients before inserting new ones
- Podman is more secure (rootless) and lightweight compared to Docker
- Compatible with Docker compose files (no changes needed)

## Troubleshooting

### Port Conflicts
If port 27017 is already in use:
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
Check logs for the container:
```bash
podman logs mongodb
# Or use the Makefile command
make -f podman/Makefile mongo-logs
```

### MongoDB Shell Connection Issues
If mongo-shell fails:
```bash
# Verify container is running
podman ps

# Check if mongosh is available in the container
podman exec -it mongodb which mongosh

# Try connecting manually
podman exec -it mongodb bash
```

### Reset Everything
```bash
# Stop and remove everything
make -f podman/Makefile podman-clean

# Remove all volumes manually
podman volume prune
```

## Configuration

The MongoDB instance is configured with:
- Authentication enabled (username/password required)
- Data persistence through volumes
- Default database: `pharmacy_modernization`
- Exposed on localhost:27017

To modify the configuration, edit `compose.yml` and update:
- Port mappings
- Environment variables
- Volume mounts
- Network settings

## Connection in Your Application

Your application should use this connection string:
```
mongodb://admin:admin123@localhost:27017/pharmacy_modernization
```

This is already configured in `internal/configs/app.yaml`.

## Future Services

This setup can be extended to include additional services if needed:
- Redis (caching)
- PostgreSQL (relational database)
- Elasticsearch (search)
- RabbitMQ (message queue)
- Kafka (event streaming)

# MongoDB Local Development - Quick Start

## 🚀 Getting Started in 3 Steps

### 1. Start MongoDB Container

```bash
make -f podman/Makefile podman-up
```

You'll see:
```
✅ MongoDB started successfully!

📦 MongoDB is available at: localhost:27017
   Username: admin
   Password: admin123
   Database: rxintake
   Connection: mongodb://admin:admin123@localhost:27017/rxintake
```

### 2. Seed the Database

```bash
make -f podman/Makefile mongo-seed
```

This will populate MongoDB with:
- 15 sample patients
- 13 addresses (multiple patients have multiple addresses)
- 27 prescriptions (various medications with different statuses)

### 3. Start Your Application

```bash
make dev
```

Your app is now connected to MongoDB and has real data to work with!

## 🔍 View Your Data

**Option 1: MongoDB Shell** (Command Line)
```bash
make -f podman/Makefile mongo-shell
```

Then query:
```javascript
use rxintake
db.patients.find().pretty()
db.addresses.find().pretty()
db.prescriptions.find().pretty()
```

**Option 2: Your Application** (Web UI)
- Open http://localhost:8080
- Navigate to the Patients page to see all 15 patients
- Click on any patient to see their addresses and prescriptions

**Option 3: MongoDB Compass** (GUI - Optional)
- Download [MongoDB Compass](https://www.mongodb.com/products/compass)
- Connect to: `mongodb://admin:admin123@localhost:27017/rxintake`
- Browse collections visually

## 🛠️ Common Tasks

**Stop MongoDB:**
```bash
make -f podman/Makefile podman-down
```

**View Logs:**
```bash
make -f podman/Makefile mongo-logs
```

**Reset and Re-seed:**
```bash
make -f podman/Makefile mongo-seed
```
(This clears existing data and inserts fresh sample data)

**Complete Reset (Remove all data):**
```bash
make -f podman/Makefile podman-clean
```

**Restart MongoDB:**
```bash
make -f podman/Makefile podman-restart
```

## ✅ Connection String

Your app is configured to connect to:
```
mongodb://admin:admin123@localhost:27017/rxintake
```

This is already set in `internal/configs/app.yaml`.

## 💡 Tips

- Data persists between container restarts (stored in Podman volumes)
- You can re-run `mongo-seed` anytime to reset to sample data
- Use MongoDB Compass for a visual interface (optional)
- The MongoDB shell (`mongo-shell`) is great for quick queries
- Check logs with `mongo-logs` if you encounter issues

## 🐳 Podman vs Docker

This setup works with both Podman and Docker. The Makefile automatically detects which one you have installed.

**If you have Docker instead:**
```bash
# These commands work the same way with Docker
make -f podman/Makefile podman-up    # Uses Docker if Podman not installed
# OR use the legacy aliases
make -f podman/Makefile docker-up
```

**Why Podman?**
- Rootless containers (more secure)
- No daemon required (lighter weight)
- Compatible with Docker commands and compose files
- Better integration with systemd

## 📋 Prerequisites

Make sure you have one of these installed:

**Podman (Recommended):**
```bash
# macOS
brew install podman podman-compose
podman machine init
podman machine start

# Fedora/RHEL
sudo dnf install podman podman-compose

# Ubuntu/Debian
sudo apt-get install podman
pip install podman-compose
```

**Docker (Alternative):**
```bash
# Download from https://docker.com
```

## 🔧 Troubleshooting

**Container won't start?**
```bash
# Check logs
make -f podman/Makefile podman-logs

# Check if port is in use
sudo lsof -i :27017
```

**Port already in use?**
```bash
# Check what's using port 27017
sudo lsof -i :27017
# Stop that service or change the port in compose.yml
```

**Can't connect to MongoDB shell?**
```bash
# Verify container is running
podman ps

# Check container logs
make -f podman/Makefile mongo-logs

# Try connecting manually
podman exec -it mongodb bash
```

**Permission errors?**
```bash
# On Podman, you might need to adjust SELinux context
# This is automatically handled in the compose file

# If using Podman on macOS/Windows, ensure machine is running
podman machine list
podman machine start
```

**Data not persisting?**
```bash
# Check volumes
podman volume ls

# Inspect MongoDB volume
podman volume inspect mongodb_data
```

**Need to reset everything?**
```bash
# This will delete all data
make -f podman/Makefile podman-clean

# Or manually remove volumes
podman volume rm mongodb_data
```

## 📊 Sample Data Details

After seeding, you'll have:

**Patients (15 total):**
- Patient IDs: P001 through P015
- Various ages, states, and phone numbers
- Realistic patient information

**Addresses (13 total):**
- Multiple address types (Home, Work, etc.)
- Some patients have multiple addresses
- Complete address information

**Prescriptions (27 total):**
- Statuses: Active, Paused, Completed, Draft
- Various medications and dosages
- Realistic prescription data

## 🎯 Next Steps

After setup:
1. Explore the patient data in your application
2. Test CRUD operations (Create, Read, Update, Delete)
3. Try the search functionality
4. Experiment with MongoDB queries in the shell
5. Build new features using the existing data

## 📚 Additional Resources

- [MongoDB Documentation](https://docs.mongodb.com/)
- [Podman Documentation](https://docs.podman.io/)
- [MongoDB Compass](https://www.mongodb.com/products/compass) - GUI tool
- [Studio 3T](https://studio3t.com/) - Another MongoDB GUI (commercial)

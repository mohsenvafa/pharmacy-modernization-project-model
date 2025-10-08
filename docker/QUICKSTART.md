# MongoDB Local Development - Quick Start

## 🚀 Getting Started in 3 Steps

### 1. Start MongoDB Container

```bash
make -f docker/Makefile docker-up
```

You'll see:
```
✅ Containers started successfully!

📦 MongoDB is available at: localhost:27017
   Username: admin
   Password: admin123

🌐 Mongo Express UI is available at: http://localhost:8081
   Username: admin
   Password: admin123

🔴 Redis is available at: localhost:6379
   Password: redis123

🌐 Redis Commander UI is available at: http://localhost:8082
   Username: admin
   Password: admin123
```

### 2. Seed the Database

```bash
make -f docker/Makefile mongo-seed
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

**Option 1: Mongo Express UI** (Easy)
- Open http://localhost:8081 in your browser
- Login with `admin`/`admin123`
- Navigate to `rxintake` database
- Browse collections: `patients`, `addresses`, `prescriptions`

**Option 2: MongoDB Shell** (Advanced)
```bash
make -f docker/Makefile mongo-shell
```

Then query:
```javascript
use rxintake
db.patients.find().pretty()
```

**Option 3: Your Application**
- Open http://localhost:8080
- Navigate to the Patients page to see all 15 patients
- Click on any patient to see their addresses and prescriptions

## 🛠️ Common Tasks

**Stop MongoDB:**
```bash
make -f docker/Makefile docker-down
```

**Reset and Re-seed:**
```bash
make -f docker/Makefile mongo-seed
```
(This clears existing data and inserts fresh sample data)

**View Logs:**
```bash
make -f docker/Makefile mongo-logs
```

**Complete Reset (Remove all data):**
```bash
make -f docker/Makefile docker-clean
```

## ✅ Connection String

Your app is now configured to connect to:
```
mongodb://admin:admin123@localhost:27017/rxintake
```

This is already set in `internal/configs/app.yaml`.

## 💡 Tips

- Data persists between container restarts (stored in Docker volumes)
- You can re-run `mongo-seed` anytime to reset to sample data
- Use Mongo Express UI to browse and edit data visually
- The sidebar in your app has a direct link to Mongo Express UI


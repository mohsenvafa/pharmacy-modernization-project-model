# VS Code / Cursor Debugging Setup

This directory contains debugging and task configurations for VS Code and Cursor.

## 🐛 Quick Start - Debug the Server

1. **Install Go extension** (if not already installed)
   - In VS Code/Cursor: Extensions → Search "Go" → Install

2. **Start debugging:**
   - Press `F5` or go to Run & Debug panel
   - Select "Debug Server" configuration
   - Click the green play button

3. **Set breakpoints:**
   - Click in the gutter (left of line numbers) to set breakpoints
   - The debugger will pause at breakpoints

## 📋 Available Debug Configurations

### 1. Debug Server
Basic server debugging without MongoDB auto-start.

**Use this when:**
- You've already started MongoDB manually
- You want to debug the server startup

**Steps:**
1. Start MongoDB: `make -f podman/Makefile podman-up`
2. Press `F5` and select "Debug Server"

### 2. Debug Server (with MongoDB)
Automatically starts MongoDB before debugging.

**Use this when:**
- You want everything to start automatically
- You're starting fresh

**Steps:**
1. Press `F5` and select "Debug Server (with MongoDB)"
2. MongoDB will start automatically

### 3. Debug Seed Script
Debug the database seeding process.

**Use this when:**
- Debugging seed data issues
- Testing patient/prescription creation

**Steps:**
1. Make sure MongoDB is running
2. Press `F5` and select "Debug Seed Script"

### 4. Debug Mock IRIS
Debug the mock external API server.

**Use this when:**
- Testing external service integrations
- Debugging mock responses

### 5. Attach to Running Process
Attach debugger to an already running process.

**Use this when:**
- You want to debug a running server
- You started the server outside the debugger

#### Understanding `${command:pickProcess}`

This configuration uses a special **VS Code variable substitution** that opens an interactive process picker.

**How it works:**
1. Start the debug config (F5 → select "Attach to Running Process")
2. VS Code shows a **searchable list** of all running processes
3. You **select the process** you want to debug
4. The debugger **attaches** to that process

**Visual example of the picker:**
```
┌─────────────────────────────────────────────────┐
│ Select the process to attach to:                │
├─────────────────────────────────────────────────┤
│ > server (PID: 12345) - ./cmd/server            │
│   chrome (PID: 11223)                           │
│   postgres (PID: 5432)                          │
│   go (PID: 98765) - go run ./cmd/server         │
│   node (PID: 45678)                             │
└─────────────────────────────────────────────────┘
```

Type to filter the list (e.g., "server" to find your server).

**Why use this?**
- Server is already running with important state
- Don't want to restart to debug
- Multiple instances running - need to pick one
- Long startup time
- Debugging production-like environment

### 6. Attach to Running Go Process
Same as above, but uses `${command:pickGoProcess}` which shows **only Go processes**.

**Comparison:**

| Configuration | Shows | Best For |
|--------------|-------|----------|
| Attach to Running Process | All system processes | When you're not sure |
| Attach to Running Go Process | Only Go processes | Finding your Go server quickly |

**Find your process PID:**
```bash
# Find Go server process
ps aux | grep "go run ./cmd/server"

# Or check what's on port 8080
lsof -i :8080

# Output shows:
# server    12345    mohsen
```

**Example workflow:**
```bash
# 1. Start server in terminal
go run ./cmd/server
# Server running on PID 12345

# 2. In VS Code: Press F5
# 3. Select "Attach to Running Go Process"
# 4. Type "server" to filter
# 5. Select your process
# 6. Debugger attached! Set breakpoints and debug live! 🎉
```

## 🎯 Setting Breakpoints

### In Code:
```go
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Click here in the gutter (left of line number) to set a breakpoint
    ctx := r.Context()
    
    // When execution hits here, debugger will pause
    patient, err := h.service.GetPatient(ctx, id)
    
    // You can inspect variables, step through code, etc.
    if err != nil {
        return err
    }
}
```

### Types of Breakpoints:
- **Normal Breakpoint**: Click in gutter (red dot)
- **Conditional Breakpoint**: Right-click in gutter → Add Conditional Breakpoint
  - Example: `id == "P001"` (only breaks when condition is true)
- **Logpoint**: Right-click in gutter → Add Logpoint
  - Example: `Patient: {id}` (logs without pausing)

## 🔍 Debug Controls

Once debugging:
- **F5** - Continue
- **F10** - Step Over (next line)
- **F11** - Step Into (enter function)
- **Shift+F11** - Step Out (exit function)
- **Ctrl+Shift+F5** - Restart
- **Shift+F5** - Stop

## 📊 Debug Panels

### Variables
View all variables in current scope:
- Local variables
- Function arguments
- Struct fields

### Watch
Add expressions to watch:
```
patient.Name
len(patients)
err != nil
```

### Call Stack
See the function call chain that led to current position.

### Debug Console
Execute Go expressions during debugging:
```go
p patient.Name
p len(addresses)
p fmt.Sprintf("ID: %s", patient.ID)
```

## ⚙️ Debug with Environment Variables

Edit `.vscode/launch.json` to add environment variables:

```json
{
  "env": {
    "PM_APP_ENV": "dev",
    "PM_DATABASE_MONGODB_URI": "mongodb://admin:admin123@localhost:27017",
    "PM_AUTH_JWT_SECRET": "my-debug-secret"
  }
}
```

## 🧪 Debug Tests

To debug a specific test:

1. **Navigate to the test file**
2. **Click "debug test" above the test function** (appears in code lens)
   
   Or add this configuration to launch.json:
   ```json
   {
     "name": "Debug Test",
     "type": "go",
     "request": "launch",
     "mode": "test",
     "program": "${fileDirname}",
     "args": [
       "-test.run",
       "TestName"
     ]
   }
   ```

## 🔧 Available Tasks

Tasks in `tasks.json` are **automated commands** you can run without typing them in terminal.

### What is tasks.json?

`tasks.json` defines a menu of common commands for your project:
- **Consistency**: Everyone runs the same commands
- **Discoverability**: New devs see available commands
- **Integration**: Works with debug configs and keybindings
- **Automation**: Can chain tasks together

### Running Tasks

**Method 1: Command Palette**
```
Cmd+Shift+P (Mac) / Ctrl+Shift+P (Win/Linux)
→ "Tasks: Run Task"
→ Select task from menu
```

**Method 2: Terminal Menu**
```
Terminal → Run Task → Select task
```

**Method 3: Build Task Shortcut**
```
Cmd+Shift+B (Mac) / Ctrl+Shift+B (Win/Linux)
→ Runs "generate-all" (default build task)
```

### Basic Tasks

| Task | What It Does |
|------|--------------|
| **start-mongodb** | Start MongoDB container (runs in background) |
| **stop-mongodb** | Stop MongoDB container |
| **seed-database** | Seed database with sample data (auto-starts MongoDB) |
| **generate-templ** | Generate Go code from templ files |
| **generate-graphql** | Generate GraphQL code from schemas |

### Composite Tasks (New!)

| Task | What It Does | How It Works |
|------|--------------|--------------|
| **generate-all** | Generate templ + GraphQL | Runs both in **parallel** |
| **dev-setup** | Start MongoDB + seed data | Runs in **sequence** |
| **clean-all** | Remove MongoDB volumes | Destructive cleanup |

### Task Features Explained

#### 1. `isBackground`
Runs task in background without blocking:
```json
{
  "label": "start-mongodb",
  "isBackground": true  // ← Runs in background
}
```
MongoDB starts and keeps running while you continue working.

#### 2. `dependsOn` - Task Dependencies
Run tasks in order or parallel:
```json
{
  "label": "seed-database",
  "dependsOn": ["start-mongodb"],  // ← Runs this first
  "dependsOrder": "sequence"        // ← Then runs seed
}
```

**Parallel execution:**
```json
{
  "label": "generate-all",
  "dependsOn": ["generate-templ", "generate-graphql"],
  "dependsOrder": "parallel"  // ← Both run simultaneously
}
```

#### 3. `group` - Build Tasks
Mark tasks as build tasks for quick access:
```json
{
  "label": "generate-all",
  "group": {
    "kind": "build",
    "isDefault": true  // ← Runs with Cmd+Shift+B
  }
}
```

Now press `Cmd+Shift+B` to run `generate-all` instantly!

#### 4. `presentation` - Output Control
Control how task output appears:
```json
{
  "presentation": {
    "reveal": "always",    // Show terminal output
    "panel": "shared"      // Reuse same terminal
  }
}
```

Options:
- `reveal`: `"always"`, `"silent"`, `"never"`
- `panel`: `"shared"`, `"dedicated"`, `"new"`

### Integration with Debug Configs

Tasks can run automatically before debugging:
```json
// In launch.json
{
  "name": "Debug Server (with MongoDB)",
  "preLaunchTask": "start-mongodb"  // ← Runs this task first!
}
```

Press F5 → MongoDB starts → Then debugger starts → All automatic! 🎉

### Real-World Examples

**Scenario 1: Fresh Start**
```
Run Task: "dev-setup"
→ Starts MongoDB
→ Seeds database with sample data
→ Ready to develop!
```

**Scenario 2: Code Generation**
```
Press: Cmd+Shift+B
→ Generates templ files
→ Generates GraphQL code
→ Both run in parallel (faster!)
```

**Scenario 3: Debug with Auto-Setup**
```
Press: F5 (Debug Server with MongoDB)
→ Task "start-mongodb" runs automatically
→ Debugger starts
→ Set breakpoints and debug!
```

### Task Output

Tasks run in VS Code's integrated terminal:
```
> Executing task: make -f podman/Makefile podman-up <

Starting Podman containers...
✅ MongoDB started successfully!

Terminal will be reused by tasks, press any key to close it.
```

## 🚀 Typical Debug Workflow

### Full Stack Development:

1. **Terminal 1** - Start MongoDB:
   ```bash
   make -f podman/Makefile podman-up
   ```

2. **Terminal 2** - Watch Tailwind CSS:
   ```bash
   make tailwind-watch
   ```

3. **Terminal 3** - Watch TypeScript:
   ```bash
   make watch-ts
   ```

4. **VS Code Debugger** - Debug Server (F5)

5. **Set breakpoints** in your handler/service code

6. **Make HTTP request** to your endpoint
   - Debugger pauses at breakpoint
   - Inspect variables
   - Step through code

### Quick Debug Cycle:

1. **Start everything:**
   ```bash
   make -f podman/Makefile podman-up  # MongoDB
   make dev                            # All watchers + server
   ```

2. **Stop server** (Ctrl+C)

3. **Start debugger** (F5)
   - Now you have hot-reload for frontend, debugger for backend

## 📝 Tips

### VS Code Variables Reference

You can use these special variables in `launch.json` configurations:

| Variable | What It Does | Example Use |
|----------|--------------|-------------|
| `${workspaceFolder}` | Path to workspace root | `"program": "${workspaceFolder}/cmd/server"` |
| `${file}` | Currently open file path | `"program": "${file}"` |
| `${fileDirname}` | Directory of current file | `"cwd": "${fileDirname}"` |
| `${fileBasename}` | Name of current file | Used in test configs |
| `${fileBasenameNoExtension}` | Filename without extension | Build output names |
| `${command:pickProcess}` | Opens process picker (all) | Attach to any process |
| `${command:pickGoProcess}` | Opens Go process picker | Attach to Go process |
| `${env:VAR_NAME}` | Environment variable value | `"port": "${env:PORT}"` |

**Example using variables:**
```json
{
  "name": "Debug Current File",
  "type": "go",
  "request": "launch",
  "mode": "auto",
  "program": "${file}",
  "cwd": "${workspaceFolder}",
  "env": {
    "APP_PORT": "${env:PORT}"
  }
}
```

**Alternative: Hardcoded Process ID**

Instead of `${command:pickProcess}`, you can hardcode a PID:
```json
{
  "name": "Attach to Process 12345",
  "type": "go",
  "request": "attach",
  "mode": "local",
  "processId": 12345  // ← Specific process ID
}
```

⚠️ **Note:** PIDs change every restart, so `pickProcess` is more flexible.

### Remote Debugging
If running on a remote server:
```bash
# On server
dlv debug --headless --listen=:2345 --api-version=2 ./cmd/server

# In launch.json
{
  "name": "Connect to Remote",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "/path/on/server",
  "port": 2345,
  "host": "your-server-ip"
}
```

### Performance Profiling
Add to launch.json:
```json
{
  "env": {
    "GODEBUG": "gctrace=1"
  },
  "buildFlags": "-race"  // Enable race detector
}
```

### Debug Logging
The app already has debug logging enabled in dev mode.
Check `internal/configs/app.yaml`:
```yaml
logging:
  level: debug  # Already set for development
  format: console
```

## 🆘 Troubleshooting

### "Cannot find package"
```bash
go mod download
go mod tidy
```

### "Delve not found"
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Breakpoints not working
1. Make sure you compiled with debug info (debugger does this automatically)
2. Check that you're debugging the right process
3. Ensure breakpoint is on an executable line (not comments/empty lines)

### Port already in use
```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process
kill <PID>
```

## 📚 Resources

- [Go Debugging in VS Code](https://github.com/golang/vscode-go/wiki/debugging)
- [Delve Debugger](https://github.com/go-delve/delve)
- [VS Code Debugging](https://code.visualstudio.com/docs/editor/debugging)


# PowerShell script to replace podman/Makefile commands for Windows compatibility
# 
# This script provides Windows users with equivalent functionality to the podman Makefile
# used on macOS/Linux systems for managing MongoDB and Memcached containers.
#
# Usage: .\podman\make.ps1 <command>
# Example: .\podman\make.ps1 podman-up
#          .\podman\make.ps1 podman-logs
#          .\podman\make.ps1 help
#
# If you encounter execution policy issues, run:
#   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"

# Get the directory where this script is located
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$PodmanDir = $ScriptDir
$ProjectRoot = Split-Path -Parent $PodmanDir

# Compose file location
$ComposeFile = Join-Path $PodmanDir "compose.yml"

# Load .env file from project root if it exists
$envFile = Join-Path $ProjectRoot ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            # Remove quotes if present
            $value = $value -replace '^["''](.*)["'']$', '$1'
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
}

# Always use podman compose
$ComposeCmd = "podman compose"

function Show-Help {
    Write-Host "Usage: .\podman\make.ps1 [command]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Available commands:" -ForegroundColor Cyan
    Write-Host "  podman-up        - Start MongoDB and Memcached containers"
    Write-Host "  podman-down      - Stop MongoDB and Memcached containers"
    Write-Host "  podman-restart   - Restart MongoDB and Memcached containers"
    Write-Host "  podman-logs      - Show logs from all containers"
    Write-Host "  podman-clean     - Stop containers and remove volumes (WARNING: deletes all data)"
    Write-Host "  mongo-logs       - Show MongoDB logs only"
    Write-Host "  mongo-shell      - Connect to MongoDB shell"
    Write-Host "  mongo-seed       - Seed MongoDB with sample patient data"
    Write-Host "  docker-up        - Alias for podman-up (backwards compatibility)"
    Write-Host "  docker-down      - Alias for podman-down (backwards compatibility)"
    Write-Host "  docker-restart   - Alias for podman-restart (backwards compatibility)"
    Write-Host "  docker-logs      - Alias for podman-logs (backwards compatibility)"
    Write-Host "  docker-clean     - Alias for podman-clean (backwards compatibility)"
}

function Start-PodmanContainers {
    Write-Host "Starting Podman containers..." -ForegroundColor Yellow
    
    Push-Location $PodmanDir
    try {
        $result = & $ComposeCmd -f $ComposeFile up -d 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ Failed to start containers" -ForegroundColor Red
            Write-Host $result
            exit 1
        }
        
        Write-Host ""
        Write-Host "✅ MongoDB started successfully!" -ForegroundColor Green
        Write-Host ""
        
        $mongoUser = $env:MONGO_ROOT_USERNAME
        $mongoPass = $env:MONGO_ROOT_PASSWORD
        $mongoDb = $env:MONGO_DATABASE
        
        Write-Host "📦 MongoDB is available at: localhost:27017" -ForegroundColor Cyan
        if ($mongoUser) {
            Write-Host "   Username: $mongoUser" -ForegroundColor Cyan
        }
        if ($mongoPass) {
            Write-Host "   Password: $mongoPass" -ForegroundColor Cyan
        }
        if ($mongoDb) {
            Write-Host "   Database: $mongoDb" -ForegroundColor Cyan
            $connString = "mongodb://"
            if ($mongoUser -and $mongoPass) {
                $connString += "$mongoUser`:$mongoPass@"
            }
            $connString += "localhost:27017/$mongoDb"
            Write-Host "   Connection: $connString" -ForegroundColor Cyan
        }
    } finally {
        Pop-Location
    }
}

function Stop-PodmanContainers {
    Write-Host "Stopping Podman containers..." -ForegroundColor Yellow
    
    Push-Location $PodmanDir
    try {
        & $ComposeCmd -f $ComposeFile down
    } finally {
        Pop-Location
    }
}

function Restart-PodmanContainers {
    Write-Host "Restarting Podman containers..." -ForegroundColor Yellow
    
    Push-Location $PodmanDir
    try {
        & $ComposeCmd -f $ComposeFile restart
    } finally {
        Pop-Location
    }
}

function Show-PodmanLogs {
    Push-Location $PodmanDir
    try {
        & $ComposeCmd -f $ComposeFile logs -f
    } finally {
        Pop-Location
    }
}

function Remove-PodmanContainersAndVolumes {
    Write-Host "⚠️  This will remove all data from MongoDB volumes!" -ForegroundColor Red
    $confirmation = Read-Host "Are you sure? [y/N]"
    
    if ($confirmation -eq "y" -or $confirmation -eq "Y") {
        Push-Location $PodmanDir
        try {
            & $ComposeCmd -f $ComposeFile down -v
            Write-Host "✅ Container and volumes removed" -ForegroundColor Green
        } finally {
            Pop-Location
        }
    } else {
        Write-Host "❌ Cancelled" -ForegroundColor Yellow
    }
}

function Show-MongoLogs {
    Push-Location $PodmanDir
    try {
        & $ComposeCmd -f $ComposeFile logs -f mongodb
    } finally {
        Pop-Location
    }
}

function Connect-MongoShell {
    $mongoUser = $env:MONGO_ROOT_USERNAME
    $mongoPass = $env:MONGO_ROOT_PASSWORD
    
    if (-not $mongoUser -or -not $mongoPass) {
        Write-Host "❌ MONGO_ROOT_USERNAME and MONGO_ROOT_PASSWORD must be set in .env file" -ForegroundColor Red
        exit 1
    }
    
    podman exec -it mongodb mongosh -u $mongoUser -p $mongoPass --authenticationDatabase admin
}

function Seed-MongoDB {
    Write-Host "🌱 Seeding MongoDB with sample patient data..." -ForegroundColor Yellow
    
    Push-Location $ProjectRoot
    try {
        go run ./cmd/seed
    } finally {
        Pop-Location
    }
}

# Main command dispatcher
switch ($Command.ToLower()) {
    "help" { Show-Help }
    "podman-up" { Start-PodmanContainers }
    "podman-down" { Stop-PodmanContainers }
    "podman-restart" { Restart-PodmanContainers }
    "podman-logs" { Show-PodmanLogs }
    "podman-clean" { Remove-PodmanContainersAndVolumes }
    "mongo-logs" { Show-MongoLogs }
    "mongo-shell" { Connect-MongoShell }
    "mongo-seed" { Seed-MongoDB }
    # Legacy Docker aliases
    "docker-up" { Start-PodmanContainers }
    "docker-down" { Stop-PodmanContainers }
    "docker-restart" { Restart-PodmanContainers }
    "docker-logs" { Show-PodmanLogs }
    "docker-clean" { Remove-PodmanContainersAndVolumes }
    default {
        Write-Host "❌ Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}


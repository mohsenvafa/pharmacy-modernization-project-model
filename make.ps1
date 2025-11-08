# PowerShell script to replace Makefile commands for Windows compatibility
# 
# This script provides Windows users with equivalent functionality to the Makefile
# used on macOS/Linux systems. All Makefile targets are available as PowerShell commands.
#
# Usage: .\make.ps1 <command>
# Example: .\make.ps1 setup
#          .\make.ps1 dev
#          .\make.ps1 help
#
# If you encounter execution policy issues, run:
#   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

param(
    [Parameter(Position=0)]
    [string]$Command = "help",
    
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$Args
)

$ErrorActionPreference = "Stop"

# Load .env file if it exists
$envFile = Join-Path $PSScriptRoot ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
}

# Determine Tailwind binary based on OS
if ($IsWindows -or $env:OS -eq "Windows_NT") {
    $TailwindBin = ".\bin\tailwindcss.exe"
} else {
    $TailwindBin = ".\bin\tailwindcss"
}

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Cyan
    Write-Host "  setup            - Initial setup (downloads Tailwind binary and installs npm dependencies)"
    Write-Host "  check-tools      - Check if required tools are installed"
    Write-Host "  tailwind-watch   - Watch Tailwind CSS files for changes"
    Write-Host "  dev-watch        - Run development server with templ in watch mode"
    Write-Host "  dev              - Run all watchers together (Tailwind, TypeScript, and dev server)"
    Write-Host "  mock-iris        - Run IRIS Mock Server (port 8881)"
    Write-Host "  build-iris-mock  - Build IRIS Mock Server binary"
    Write-Host "  build-ts         - Build TypeScript files"
    Write-Host "  watch-ts         - Watch TypeScript files for changes"
    Write-Host "  graphql-generate - Generate GraphQL code from schemas"
    Write-Host "  graphql-install  - Install gqlgen CLI tool"
    Write-Host "  podman-up        - Start MongoDB and Memcached containers"
    Write-Host "  podman-down      - Stop MongoDB and Memcached containers"
    Write-Host "  podman-logs      - Show container logs"
}

function Invoke-Setup {
    Write-Host "🔧 Running setup..." -ForegroundColor Yellow
    
    # Setup Tailwind
    Invoke-TailwindSetup
    
    # Install npm dependencies
    Write-Host "📦 Installing npm dependencies..." -ForegroundColor Yellow
    Push-Location "web"
    try {
        npm install
        if ($LASTEXITCODE -ne 0) {
            throw "npm install failed"
        }
    } finally {
        Pop-Location
    }
    
    Write-Host "✅ Setup complete!" -ForegroundColor Green
}

function Invoke-TailwindSetup {
    # Detect platform
    $PlatformPrefix = if ($IsWindows -or $env:OS -eq "Windows_NT") { "windows" }
                      elseif ($IsMacOS) { "macos" }
                      else { "linux" }
    
    # Detect architecture
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        # On Windows, check multiple possible architecture indicators
        $procArch = $env:PROCESSOR_ARCHITECTURE
        $procArchW6432 = $env:PROCESSOR_ARCHITEW6432
        
        if ($procArchW6432 -eq "ARM64" -or $procArch -eq "ARM64") {
            $Arch = "arm64"
        } elseif ($procArchW6432 -eq "AMD64" -or $procArch -eq "AMD64" -or $procArch -eq "x86_64") {
            $Arch = "x64"
        } else {
            # Fallback: check system info
            $sysInfo = (Get-CimInstance Win32_OperatingSystem).OSArchitecture
            if ($sysInfo -like "*64*") {
                $Arch = "x64"
            } else {
                $Arch = "x64" # Default to x64
            }
        }
    } elseif ($IsMacOS) {
        # On macOS
        $sysctl = sysctl -n machdep.cpu.brand_string 2>$null
        if ($LASTEXITCODE -eq 0 -and $sysctl -match "Apple") {
            $Arch = "arm64"
        } else {
            $Arch = "x64"
        }
    } else {
        # On Linux
        $unameM = uname -m 2>$null
        if ($unameM -eq "aarch64" -or $unameM -eq "arm64") {
            $Arch = "arm64"
        } else {
            $Arch = "x64"
        }
    }
    
    $Suffix = if ($PlatformPrefix -eq "windows") { ".exe" } else { "" }
    $TailwindVersion = "v3.4.17"
    $TailwindPlatform = "$PlatformPrefix-$Arch"
    $TailwindFilename = "tailwindcss-$TailwindPlatform$Suffix"
    $TailwindUrl = "https://github.com/tailwindlabs/tailwindcss/releases/download/$TailwindVersion/$TailwindFilename"
    
    $BinDir = Join-Path $PSScriptRoot "bin"
    $TailwindPath = Join-Path $BinDir "tailwindcss$Suffix"
    
    if (-not (Test-Path $TailwindPath)) {
        Write-Host "📥 Downloading Tailwind CLI $TailwindVersion ($TailwindPlatform)..." -ForegroundColor Yellow
        if (-not (Test-Path $BinDir)) {
            New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
        }
        
        try {
            Invoke-WebRequest -Uri $TailwindUrl -OutFile $TailwindPath -UseBasicParsing
            Write-Host "✅ Tailwind CLI downloaded successfully!" -ForegroundColor Green
        } catch {
            Write-Host "❌ Failed to download Tailwind CLI: $_" -ForegroundColor Red
            throw
        }
    } else {
        Write-Host "✅ Tailwind CLI already present at $TailwindPath" -ForegroundColor Green
    }
}

function Test-Tools {
    $BinDir = Join-Path $PSScriptRoot "bin"
    $Suffix = if ($IsWindows -or $env:OS -eq "Windows_NT") { ".exe" } else { "" }
    $TailwindPath = Join-Path $BinDir "tailwindcss$Suffix"
    $NodeModulesDir = Join-Path $PSScriptRoot "web" "node_modules"
    
    $Errors = @()
    
    if (-not (Test-Path $TailwindPath)) {
        $Errors += "Tailwind CLI missing at $TailwindPath. Run '.\make.ps1 setup' first."
    }
    
    if (-not (Test-Path $NodeModulesDir)) {
        $Errors += "Node modules missing. Run '.\make.ps1 setup' first."
    }
    
    if ($Errors.Count -gt 0) {
        foreach ($error in $Errors) {
            Write-Host "❌ $error" -ForegroundColor Red
        }
        exit 1
    }
}

function Start-TailwindWatch {
    Test-Tools
    Write-Host "🎨 Starting Tailwind CSS watcher..." -ForegroundColor Yellow
    Push-Location "web"
    try {
        npx tailwindcss -c tailwind.config.js -i styles/input.css -o public/app.css --watch
    } finally {
        Pop-Location
    }
}

function Start-DevWatch {
    Write-Host "🚀 Starting development server with templ watch..." -ForegroundColor Yellow
    templ generate -watch `
        -proxyport=7332 `
        -proxy="http://localhost:8080" `
        -cmd="go run -gcflags=all=-N -gcflags=all=-l ./cmd/server" `
        -open-browser=false
}

function Start-Dev {
    Write-Host "🚀 Starting development server with all watchers..." -ForegroundColor Cyan
    Write-Host "⚠️  Note: Press Ctrl+C to stop all watchers and the server" -ForegroundColor Yellow
    
    Test-Tools
    
    # Start Tailwind watcher in a separate process
    $tailwindProc = Start-Process -FilePath "npx" -ArgumentList "tailwindcss", "-c", "tailwind.config.js", "-i", "styles/input.css", "-o", "public/app.css", "--watch" -WorkingDirectory "web" -PassThru -WindowStyle Hidden
    
    # Start TypeScript watcher in a separate process
    Push-Location "web"
    $tsProc = Start-Process -FilePath "npm" -ArgumentList "run", "watch" -WorkingDirectory "web" -PassThru -WindowStyle Hidden
    Pop-Location
    
    # Store process IDs for cleanup
    $tailwindPid = $tailwindProc.Id
    $tsPid = $tsProc.Id
    
    try {
        # Start dev server in foreground
        Start-DevWatch
    } catch {
        # Cleanup on error
        Write-Host "`n🛑 Stopping watchers..." -ForegroundColor Yellow
        if ($tailwindProc -and -not $tailwindProc.HasExited) {
            Stop-Process -Id $tailwindPid -Force -ErrorAction SilentlyContinue
        }
        if ($tsProc -and -not $tsProc.HasExited) {
            Stop-Process -Id $tsPid -Force -ErrorAction SilentlyContinue
        }
        throw
    } finally {
        # Cleanup on exit
        Write-Host "`n🛑 Stopping watchers..." -ForegroundColor Yellow
        try {
            if ($tailwindProc -and -not $tailwindProc.HasExited) {
                Stop-Process -Id $tailwindPid -Force -ErrorAction SilentlyContinue
            }
        } catch {}
        try {
            if ($tsProc -and -not $tsProc.HasExited) {
                Stop-Process -Id $tsPid -Force -ErrorAction SilentlyContinue
            }
        } catch {}
    }
}

function Start-MockIris {
    Write-Host "🚀 Starting IRIS Mock Server on port 8881..." -ForegroundColor Cyan
    Write-Host "📍 Pharmacy API: http://localhost:8881/pharmacy/v1" -ForegroundColor Green
    Write-Host "📍 Billing API:  http://localhost:8881/billing/v1" -ForegroundColor Green
    Write-Host "📍 Stargate Auth: http://localhost:8881/oauth" -ForegroundColor Green
    go run ./cmd/iris_mock
}

function Build-IrisMock {
    Write-Host "🔨 Building IRIS Mock Server..." -ForegroundColor Yellow
    $exe = if ($IsWindows -or $env:OS -eq "Windows_NT") { ".exe" } else { "" }
    $output = "iris_mock$exe"
    go build -o $output ./cmd/iris_mock
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Built: ./$output" -ForegroundColor Green
    }
}

function Build-TypeScript {
    Write-Host "🔨 Building TypeScript..." -ForegroundColor Yellow
    Push-Location "web"
    try {
        npm run build
    } finally {
        Pop-Location
    }
}

function Watch-TypeScript {
    Write-Host "👀 Watching TypeScript files..." -ForegroundColor Yellow
    Push-Location "web"
    try {
        npm run watch
    } finally {
        Pop-Location
    }
}

function Invoke-GraphQLGenerate {
    Write-Host "🔄 Generating GraphQL code..." -ForegroundColor Yellow
    gqlgen generate
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ GraphQL code generated successfully!" -ForegroundColor Green
    }
}

function Install-GraphQLGen {
    Write-Host "📦 Installing gqlgen..." -ForegroundColor Yellow
    go install github.com/99designs/gqlgen@latest
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ gqlgen installed successfully!" -ForegroundColor Green
    }
}

function Start-PodmanContainers {
    $podmanScript = Join-Path $PSScriptRoot "podman" "make.ps1"
    if (Test-Path $podmanScript) {
        & $podmanScript podman-up
    } else {
        Write-Host "❌ podman/make.ps1 not found" -ForegroundColor Red
        exit 1
    }
}

function Stop-PodmanContainers {
    $podmanScript = Join-Path $PSScriptRoot "podman" "make.ps1"
    if (Test-Path $podmanScript) {
        & $podmanScript podman-down
    } else {
        Write-Host "❌ podman/make.ps1 not found" -ForegroundColor Red
        exit 1
    }
}

function Show-PodmanLogs {
    $podmanScript = Join-Path $PSScriptRoot "podman" "make.ps1"
    if (Test-Path $podmanScript) {
        & $podmanScript podman-logs
    } else {
        Write-Host "❌ podman/make.ps1 not found" -ForegroundColor Red
        exit 1
    }
}

# Main command dispatcher
switch ($Command.ToLower()) {
    "help" { Show-Help }
    "setup" { Invoke-Setup }
    "check-tools" { Test-Tools }
    "tailwind-watch" { Start-TailwindWatch }
    "dev-watch" { Start-DevWatch }
    "dev" { Start-Dev }
    "mock-iris" { Start-MockIris }
    "build-iris-mock" { Build-IrisMock }
    "build-ts" { Build-TypeScript }
    "watch-ts" { Watch-TypeScript }
    "graphql-generate" { Invoke-GraphQLGenerate }
    "graphql-install" { Install-GraphQLGen }
    "podman-up" { Start-PodmanContainers }
    "podman-down" { Stop-PodmanContainers }
    "podman-logs" { Show-PodmanLogs }
    default {
        Write-Host "❌ Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}


# build.ps1 - Build script for CodeAuthorityFirst
# This script builds the caf binary with version embedding using ldflags
# Supports Windows systems

param(
    [string]$OutputDir = ".",
    [switch]$Clean = $false
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = Join-Path $ScriptDir ".." ".." | Resolve-Path
$OutputPath = Join-Path $ProjectRoot $OutputDir

Write-Host "🔨 Building CodeAuthorityFirst..." -ForegroundColor Green
Push-Location $ProjectRoot

# Clean if requested
if ($Clean) {
    Write-Host "🧹 Cleaning previous builds..." -ForegroundColor Yellow
    Remove-Item "caf.exe" -ErrorAction SilentlyContinue
    Remove-Item "caf" -ErrorAction SilentlyContinue
}

# Get version from git or use fallback
try {
    if (Get-Command git -ErrorAction SilentlyContinue) {
        $GitVersion = (git describe --tags --dirty --always 2>$null) -replace "`r`n", ""
        $CommitHash = (git rev-parse HEAD 2>$null) -replace "`r`n", ""
        $IsDirty = (git diff-index --quiet HEAD -- 2>$null) -eq $false
        $BuildTime = Get-Date -UFormat "%Y-%m-%dT%H:%M:%SZ"
        
        $VERSION_FLAG = "-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GitVersion'"
        $VERSION_FLAG += " -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit=$CommitHash'"
        $VERSION_FLAG += " -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.datetime=$BuildTime'"
        $VERSION_FLAG += " -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.modified=$IsDirty'"
        
        Write-Host "📦 Version from git: $GitVersion" -ForegroundColor Cyan
        if ($IsDirty) {
            Write-Host "⚠️  Working directory has uncommitted changes" -ForegroundColor Yellow
        }
    } else {
        $VERSION_FLAG = "-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
        Write-Host "⚠️  Git not available, using fallback version" -ForegroundColor Yellow
    }
} catch {
    $VERSION_FLAG = "-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
    Write-Host "⚠️  Error getting git info, using fallback version" -ForegroundColor Yellow
}

Write-Host "🚀 Building with flags: $VERSION_FLAG" -ForegroundColor Blue

# Build the binary
try {
    go build -ldflags="$VERSION_FLAG" -o caf.exe ./cmd/caf
    
    if ($LASTEXITCODE -eq 0) {
        $BinaryPath = Join-Path $ProjectRoot "caf.exe"
        Write-Host "✅ Build successful!" -ForegroundColor Green
        Write-Host "📋 Binary created: $BinaryPath" -ForegroundColor Cyan
        
        # Test the binary
        Write-Host "🧪 Testing binary..." -ForegroundColor Yellow
        & $BinaryPath --version
        
        Write-Host "🎉 Build complete!" -ForegroundColor Green
    } else {
        Write-Host "❌ Build failed!" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ Build failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Pop-Location
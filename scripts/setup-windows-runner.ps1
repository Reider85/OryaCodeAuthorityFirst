#!/usr/bin/env pwsh

# Self-Hosted Runner Setup Script for Windows
# This script sets up a GitHub Actions self-hosted runner on Windows systems

param(
    [Parameter(Mandatory=$true)]
    [string]$GitHubToken,
    
    [string]$RunnerName = "hybrid-runner-windows",
    [string]$RunnerLabels = "windows,self-hosted,hybrid",
    [string]$RunnerGroup = "hybrid-runners"
)

Write-Host "🚀 Setting up self-hosted GitHub Actions runner" -ForegroundColor Green
Write-Host "Runner Name: $RunnerName" -ForegroundColor Yellow
Write-Host "Labels: $RunnerLabels" -ForegroundColor Yellow
Write-Host "Group: $RunnerGroup" -ForegroundColor Yellow

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "⚠️  This script should be run as an administrator" -ForegroundColor Red
    Write-Host "Please run PowerShell as an administrator and try again" -ForegroundColor Red
    exit 1
}

# Install Chocolatey if not present
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "📦 Installing Chocolatey..." -ForegroundColor Yellow
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
}

# Install dependencies
Write-Host "🔧 Installing dependencies..." -ForegroundColor Yellow
choco install -y docker-desktop git

# Start Docker Desktop
Write-Host "🐳 Starting Docker Desktop..." -ForegroundColor Yellow
Start-Process docker-desktop
Start-Sleep -Seconds 30

# Download and configure runner
Write-Host "🏃‍♂️ Downloading GitHub Actions runner..." -ForegroundColor Yellow
$runnerUrl = "https://github.com/actions/runner/releases/download/v2.313.0/actions-runner-win-x64-2.313.0.zip"
$runnerZip = "actions-runner-win-x64-2.313.0.zip"
$runnerDir = "actions-runner"

if (-not (Test-Path $runnerZip)) {
    Invoke-WebRequest -Uri $runnerUrl -OutFile $runnerZip
}

if (-not (Test-Path $runnerDir)) {
    Expand-Archive -Path $runnerZip -DestinationPath $runnerDir -Force
}

Set-Location $runnerDir

# Configure runner
Write-Host "⚙️  Configuring runner..." -ForegroundColor Yellow
.\config.cmd --url https://github.com/${env:GITHUB_REPOSITORY:-Reider85/OryaCodeAuthorityFirst} `
             --token $GitHubToken `
             --name $RunnerName `
             --labels $RunnerLabels `
             --runnergroup $RunnerGroup `
             --unattended

# Install as service
Write-Host "🔄 Installing runner service..." -ForegroundColor Yellow
.\svc.cmd install

# Start service
Write-Host "▶️  Starting runner service..." -ForegroundColor Yellow
.\svc.cmd start

Write-Host "✅ Self-hosted runner setup complete!" -ForegroundColor Green
Write-Host "Runner will appear in GitHub Actions settings within a few minutes" -ForegroundColor Green
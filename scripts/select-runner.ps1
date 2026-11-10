#!/usr/bin/env pwsh

# Runner Selection Utility for PowerShell
# This script determines whether to use GitHub-hosted or local runners
# based on GitHub Actions usage thresholds

param(
    [int]$Threshold = 80,  # Default threshold: 80%
    [string]$RunnerType = "hybrid"  # hybrid, github-only, local-only
)

Write-Host "🔍 Runner Selection Utility" -ForegroundColor Cyan
Write-Host "Threshold: ${Threshold}%" -ForegroundColor Yellow
Write-Host "Mode: ${RunnerType}" -ForegroundColor Yellow

# Check if we're in a GitHub Actions environment
if ($env:GITHUB_ACTIONS -eq "true") {
    Write-Host "📊 Running in GitHub Actions environment" -ForegroundColor Green
    
    # Build the Go runner manager
    Write-Host "🔨 Building runner-manager..." -ForegroundColor Yellow
    go build runner-manager.exe
    
    # Check GitHub usage
    Write-Host "📈 Checking GitHub Actions usage..." -ForegroundColor Yellow
    $result = & ./runner-manager.exe $Threshold
    $exitCode = $LASTEXITCODE
    
    if ($exitCode -eq 0) {
        Write-Host "✅ Using GitHub-hosted runners" -ForegroundColor Green
        Write-Output "::set-output name=runner-type::github"
        exit 0
    } else {
        Write-Host "🚨 GitHub usage threshold exceeded" -ForegroundColor Red
        if ($RunnerType -eq "hybrid" -or $RunnerType -eq "local-only") {
            Write-Host "🏃‍♂️ Switching to local runners" -ForegroundColor Green
            Write-Output "::set-output name=runner-type::local"
            exit 0
        } else {
            Write-Host "❌ Local runners not available in current mode" -ForegroundColor Red
            exit 1
        }
    }
} else {
    Write-Host "🏃‍♂️ Running locally - using local runners" -ForegroundColor Green
    Write-Output "::set-output name=runner-type::local"
    exit 0
}
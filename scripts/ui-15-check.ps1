# ui-15-check.ps1 - Automated 15-Item Critical UI Checklist Runner for Windows
# This script runs the automated 15-item critical UI checklist verification

param(
    [string]$ScreenPath = "",
    [switch]$AllScreens = $false,
    [string]$OutputDir = "",
    [switch]$Verbose = $false
)

# Set error preference
$ErrorActionPreference = "Stop"

# Import PowerShell module for colored output if available
if ($PSVersionTable.PSVersion.Major -ge 6) {
    $PSStyle = $PSStyle
}

# Function to print colored output
function Write-Status {
    param(
        [string]$Status,
        [string]$Message
    )
    
    if ($PSStyle) {
        switch ($Status) {
            "PASS" { Write-Host "$($PSStyle.Foreground.Green)✓ $Message$($PSStyle.Reset)" }
            "FAIL" { Write-Host "$($PSStyle.Foreground.Red)✗ $Message$($PSStyle.Reset)" }
            "WARN" { Write-Host "$($PSStyle.Foreground.Yellow)⚠ $Message$($PSStyle.Reset)" }
            "INFO" { Write-Host "$($PSStyle.Foreground.Blue)ℹ $Message$($PSStyle.Reset)" }
        }
    } else {
        switch ($Status) {
            "PASS" { Write-Host "✓ $Message" -ForegroundColor Green }
            "FAIL" { Write-Host "✗ $Message" -ForegroundColor Red }
            "WARN" { Write-Host "⚠ $Message" -ForegroundColor Yellow }
            "INFO" { Write-Host "ℹ $Message" -ForegroundColor Blue }
        }
    }
}

# Function to find UI screens
function Find-UIScreens {
    $screensDir = Join-Path $ProjectRoot "pkg\ui\screens"
    
    if (-not (Test-Path $screensDir -PathType Container)) {
        Write-Status "FAIL" "UI screens directory not found: $screensDir"
        exit 1
    }
    
    # Find all Go files in screens directory (excluding test files)
    $screens = Get-ChildItem -Path $screensDir -Filter "*.go" -File | 
               Where-Object { $_.Name -notlike "*_test.go" } |
               Sort-Object Name
    
    return $screens.FullName
}

# Function to run 15-item check for a single screen
function Check-Screen {
    param(
        [string]$ScreenFile
    )
    
    $screenName = [System.IO.Path]::GetFileNameWithoutExtension($ScreenFile)
    $screenDir = [System.IO.Path]::GetDirectoryName($ScreenFile)
    
    if ($OutputDir) {
        $outputFile = Join-Path $OutputDir "${screenName}_15_check_results.json"
    } else {
        $outputFile = Join-Path $screenDir "${screenName}_15_check_results.json"
    }
    
    Write-Status "INFO" "Checking $screenName..."
    
    # Run the Go program
    Set-Location $ProjectRoot
    $goRunArgs = @("$scriptRoot\ui_15_check.go", $ScreenFile, $outputFile)
    
    if ($Verbose) {
        Write-Status "INFO" "Running: go run $goRunArgs"
    }
    
    $result = & go run @goRunArgs 2>&1
    $exitCode = $LASTEXITCODE
    
    if ($exitCode -eq 0) {
        Write-Status "PASS" "$screenName check completed"
        return 0
    } else {
        if ($exitCode -eq 1) {
            Write-Status "FAIL" "$screenName failed critical checks"
        } elseif ($exitCode -eq 2) {
            Write-Status "WARN" "$screenName passed with warnings"
        } else {
            Write-Status "FAIL" "$screenName check failed with error code $exitCode"
        }
        return $exitCode
    }
}

# Function to generate summary report
function Generate-Summary {
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $summaryFile = Join-Path $ProjectRoot "ui_15_check_summary_$timestamp.json"
    $summaryReport = Join-Path $ProjectRoot "ui_15_check_summary_$timestamp.md"
    
    Write-Status "INFO" "Generating summary report..."
    
    # Create JSON summary
    $jsonSummary = @{
        timestamp = Get-Date -Format "o"
        total_screens = 0
        passed_screens = 0
        failed_screens = 0
        warned_screens = 0
        screens = @()
    } | ConvertTo-Json -Depth 10
    
    Set-Content -Path $summaryFile -Value $jsonSummary
    
    # Create Markdown summary
    $markdownSummary = @"
# 15-Item Critical UI Checklist Summary Report

**Generated:** $(Get-Date -Format "F")
**Total Screens Checked:** 0

## Results Overview

| Status | Count | Percentage |
|--------|-------|------------|
| ✅ PASS | 0 | 0% |
| ❌ FAIL | 0 | 0% |
| ⚠ WARN | 0 | 0% |

## Detailed Results

"@
    
    Set-Content -Path $summaryReport -Value $markdownSummary
    
    Write-Status "INFO" "Summary reports generated:"
    Write-Status "INFO" "  JSON: $summaryFile"
    Write-Status "INFO" "  Markdown: $summaryReport"
}

# Function to check if we're in the right directory
function Check-ProjectRoot {
    $goModPath = Join-Path $ProjectRoot "go.mod"
    $screensDir = Join-Path $ProjectRoot "pkg\ui\screens"
    
    if (-not (Test-Path $goModPath -PathType Leaf)) {
        Write-Status "FAIL" "Project root not found (go.mod missing)"
        exit 1
    }
    
    if (-not (Test-Path $screensDir -PathType Container)) {
        Write-Status "FAIL" "UI screens directory not found"
        exit 1
    }
}

# Main function
function Main {
    Write-Status "INFO" "Starting 15-Item Critical UI Checklist Verification"
    Write-Status "INFO" "Project Root: $ProjectRoot"
    Write-Host "================================================================"
    
    # Check if we're in the right project
    Check-ProjectRoot
    
    # Determine which screens to check
    if ($AllScreens) {
        $screens = Find-UIScreens
        $totalScreens = $screens.Count
    } elseif ($ScreenPath) {
        if (-not (Test-Path $ScreenPath -PathType Leaf)) {
            Write-Status "FAIL" "Screen file not found: $ScreenPath"
            exit 1
        }
        $screens = @($ScreenPath)
        $totalScreens = 1
    } else {
        Write-Status "FAIL" "Please specify either -ScreenPath or -AllScreens"
        exit 1
    }
    
    Write-Status "INFO" "Found $totalScreens screens to check"
    Write-Host ""
    
    # Initialize counters
    $passedScreens = 0
    $failedScreens = 0
    $warnedScreens = 0
    
    # Check each screen
    foreach ($screenFile in $screens) {
        $result = Check-Screen $screenFile
        if ($result -eq 0) {
            $passedScreens++
        } elseif ($result -eq 1) {
            $failedScreens++
        } elseif ($result -eq 2) {
            $warnedScreens++
        }
        Write-Host ""
    }
    
    # Generate summary
    Generate-Summary
    
    # Print final summary
    Write-Host "================================================================"
    Write-Status "INFO" "15-Item Critical UI Checklist Complete"
    Write-Status "INFO" "Screens Checked: $totalScreens"
    Write-Status "INFO" "Passed: $passedScreens"
    Write-Status "INFO" "Failed: $failedScreens"
    Write-Status "INFO" "Warnings: $warnedScreens"
    
    if ($failedScreens -gt 0) {
        Write-Status "FAIL" "$failedScreens screens failed critical checks"
        exit 1
    } elseif ($warnedScreens -gt 0) {
        Write-Status "WARN" "$warnedScreens screens have warnings"
        exit 2
    } else {
        Write-Status "PASS" "All screens passed the 15-item critical checklist"
        exit 0
    }
}

# Initialize variables
$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = Split-Path -Parent $scriptRoot

# Run main function
Main
<#
.SYNOPSIS
    Executes LSP testing for Criterion 3: LSP fallback builds symbol graph for Go repo up to 1000 files
.DESCRIPTION
    This script tests the LSP fallback functionality with large Go repositories
    to ensure symbol graph generation works correctly with up to 1000 files.
.NOTES
    File Name      : criterion3-lsp-test.ps1
    Prerequisites  : caf binary, git, gopls
    Author         : CodeAuthorityFirst Testing Team
    Version        : 1.0.0
#>

param(
    [Parameter(Mandatory=$false)]
    [string]$RepositoryPath = ".\repositories\large-go-project",
    
    [Parameter(Mandatory=$false)]
    [int]$TargetFileCount = 1000,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableSymbolValidation = $true,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableMemoryMonitoring = $true,
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDirectory = ".\reports\criterion3"
)

# Set up error handling
$ErrorActionPreference = "Stop"

# Create output directory
if (-not (Test-Path $OutputDirectory)) {
    New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
}

# Logging function
function Write-Log {
    param(
        [string]$Message,
        [string]$Level = "INFO"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logEntry = "[$timestamp] [$Level] $Message"
    
    Write-Host $logEntry
    Add-Content -Path "$OutputDirectory\test-log.txt" -Value $logEntry
}

# Initialize test environment
function Initialize-TestEnvironment {
    Write-Log "Initializing test environment..."
    
    # Check prerequisites
    $requiredTools = @("git", "caf", "gopls")
    foreach ($tool in $requiredTools) {
        if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
            Write-Log "Required tool '$tool' not found" -Level "ERROR"
            exit 1
        }
    }
    
    # Navigate to repository
    if (-not (Test-Path $RepositoryPath)) {
        Write-Log "Repository path '$RepositoryPath' not found" -Level "ERROR"
        exit 1
    }
    
    Set-Location $RepositoryPath
    
    # Initialize caf if not already initialized
    if (-not (Test-Path ".\.caf")) {
        Write-Log "Initializing caf in repository..."
        caf init --auto-detect | Out-File -Append "$OutputDirectory\initialization.log"
    }
    
    # Configure LSP adapter
    Write-Log "Configuring LSP adapter..."
    caf init --adapters universal-lsp | Out-File -Append "$OutputDirectory\lsp-config.log"
    
    # Verify LSP configuration
    $lspStatus = caf adapters list
    Add-Content -Path "$OutputDirectory\lsp-status.log" -Value $lspStatus
    
    Write-Log "Test environment initialized successfully"
}

# Create large test repository
function Create-LargeRepository {
    param(
        [int]$FileCount
    )
    
    Write-Log "Creating large repository with $FileCount files..."
    
    # Create basic structure
    $baseStructure = @(
        "cmd/app/main.go",
        "internal/config/config.go",
        "internal/handler/handler.go",
        "internal/service/service.go",
        "internal/model/model.go",
        "pkg/utils/utils.go",
        "pkg/database/database.go"
    )
    
    # Create base files
    foreach ($file in $baseStructure) {
        $dir = Split-Path $file -Parent
        if ($dir -ne "") {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
        
        $content = @"
package $(Split-Path (Split-Path $file -Leaf) -Leaf -Resolve)

import "fmt"

// Basic structure for $(Split-Path $file -Leaf -Resolve)
func BasicFunction() {
    fmt.Println("Basic function from $(Split-Path $file -Leaf -Resolve)")
}
"@
        
        Set-Content -Path $file -Value $content
    }
    
    # Create additional files to reach target count
    $createdFiles = $baseStructure.Count
    $fileNumber = 1
    
    while ($createdFiles -lt $FileCount) {
        $subDir = "internal/module$([math]::Ceiling($fileNumber / 100))"
        New-Item -ItemType Directory -Path $subDir -Force | Out-Null
        
        $fileName = "$subDir/file$fileNumber.go"
        $content = @"
package main

import "fmt"

// Function $fileNumber in module $([math]::Ceiling($fileNumber / 100))
func Function$fileNumber() {
    fmt.Println("Function $fileNumber")
}

// Helper function for $fileNumber
func Helper$fileNumber() {
    fmt.Println("Helper function $fileNumber")
}
"@
        
        Set-Content -Path $fileName -Value $content
        $createdFiles++
        $fileNumber++
        
        if ($fileNumber % 100 -eq 0) {
            Write-Log "Created $createdFiles files..."
        }
    }
    
    Write-Log "Repository created with $createdFiles files"
    
    # Initialize git repository
    if (-not (Test-Path ".git")) {
        git init
        git add .
        git commit -m "Initial commit with $FileCount files"
    }
    
    return $createdFiles
}

# Test symbol generation
function Test-SymbolGeneration {
    param(
        [int]$FileCount
    )
    
    Write-Log "Testing symbol generation with $FileCount files..."
    
    $testResults = @{
        FileCount = $FileCount
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        GenerationTime = 0
        SymbolCount = 0
        Success = $false
        Error = $null
        MemoryUsage = 0
        CPUUsage = 0
    }
    
    try {
        # Start memory monitoring
        $memoryMonitor = Start-Job -Name "MemoryMonitor" -ScriptBlock {
            param($OutputDir)
            while ($true) {
                $process = Get-Process -Name "gopls" -ErrorAction SilentlyContinue
                if ($process) {
                    $memory = $process.WorkingSet / 1MB
                    $cpu = $process.CPU
                    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                    "$timestamp,$memory,$cpu" | Out-File -Append "$OutputDir\memory-monitor.csv"
                }
                Start-Sleep -Seconds 1
            }
        } -ArgumentList $OutputDirectory
        
        # Measure symbol generation time
        $generationStart = Get-Date
        
        # Generate symbols using caf
        $symbolResult = caf lsp generate-symbols
        
        $generationEnd = Get-Date
        $testResults.GenerationTime = ($generationEnd - $generationStart).TotalMilliseconds
        
        # Stop memory monitoring
        Stop-Job -Name "MemoryMonitor"
        Remove-Job -Name "MemoryMonitor"
        
        # Get final memory usage
        $memoryData = Get-Content "$OutputDirectory\memory-monitor.csv" | Select-Object -Last 1
        if ($memoryData) {
            $memoryParts = $memoryData.Split(',')
            $testResults.MemoryUsage = [double]$memoryParts[1]
            $testResults.CPUUsage = [double]$memoryParts[2]
        }
        
        # Parse symbol result
        if ($symbolResult -match "symbols") {
            $testResults.Success = $true
            $testResults.SymbolCount = [regex]::Matches($symbolResult, "function").Count + [regex]::Matches($symbolResult, "type").Count
        } else {
            $testResults.Error = "Symbol generation failed: $symbolResult"
        }
        
        Write-Log "Symbol generation completed: Success=$($testResults.Success), Time=$($testResults.GenerationTime)ms, Symbols=$($testResults.SymbolCount)"
        
    }
    catch {
        $testResults.Error = $_.Exception.Message
        $testResults.Success = $false
        Write-Log "Symbol generation failed: $($_.Exception.Message)" -Level "ERROR"
    }
    
    return $testResults
}

# Validate symbol accuracy
function Validate-SymbolAccuracy {
    param(
        [int]$ExpectedSymbols,
        [int]$ActualSymbols
    )
    
    Write-Log "Validating symbol accuracy..."
    
    $validationResults = @{
        ExpectedSymbols = $ExpectedSymbols
        ActualSymbols = $ActualSymbols
        Accuracy = 0
        Success = $false
        Issues = @()
    }
    
    # Calculate accuracy
    if ($ExpectedSymbols -gt 0) {
        $validationResults.Accuracy = ($ActualSymbols / $ExpectedSymbols) * 100
    } else {
        $validationResults.Accuracy = 100
    }
    
    # Check accuracy thresholds
    $thresholds = @{
        MinAccuracy = 95  # 95% minimum accuracy
    }
    
    if ($validationResults.Accuracy -lt $thresholds.MinAccuracy) {
        $validationResults.Issues += "Symbol accuracy below threshold: $($validationResults.Accuracy)% < $($thresholds.MinAccuracy)%"
        $validationResults.Success = $false
    } else {
        $validationResults.Success = $true
    }
    
    # Generate validation report
    $validationResults | ConvertTo-Json | Out-File "$OutputDirectory\symbol-validation.json"
    
    Write-Log "Symbol accuracy validation completed"
    Write-Log "Accuracy: $($validationResults.Accuracy)%"
    Write-Log "Success: $($validationResults.Success)"
    
    return $validationResults
}

# Test performance requirements
function Test-PerformanceRequirements {
    param(
        [object]$TestResults
    )
    
    Write-Log "Testing performance requirements..."
    
    $performanceReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        GenerationTimeMs = $TestResults.GenerationTime
        MemoryUsageMB = $TestResults.MemoryUsage
        CPUUsage = $TestResults.CPUUsage
        Success = $false
        Violations = @()
    }
    
    # Check performance thresholds
    $thresholds = @{
        MaxGenerationTimeMs = 300000   # 5 minutes
        MaxMemoryUsageMB = 100         # 100MB
        MaxCPUUsage = 80               # 80%
    }
    
    # Check generation time
    if ($TestResults.GenerationTime -gt $thresholds.MaxGenerationTimeMs) {
        $performanceReport.Violations += "Generation time exceeded: $($TestResults.GenerationTime)ms > $($thresholds.MaxGenerationTimeMs)ms"
    }
    
    # Check memory usage
    if ($TestResults.MemoryUsage -gt $thresholds.MaxMemoryUsageMB) {
        $performanceReport.Violations += "Memory usage exceeded: $($TestResults.MemoryUsage)MB > $($thresholds.MaxMemoryUsageMB)MB"
    }
    
    # Check CPU usage
    if ($TestResults.CPUUsage -gt $thresholds.MaxCPUUsage) {
        $performanceReport.Violations += "CPU usage exceeded: $($TestResults.CPUUsage)% > $($thresholds.MaxCPUUsage)%"
    }
    
    # Determine success
    $performanceReport.Success = ($performanceReport.Violations.Count -eq 0)
    
    $performanceReport | ConvertTo-Json | Out-File "$OutputDirectory\performance-report.json"
    
    Write-Log "Performance testing completed"
    Write-Log "Success: $($performanceReport.Success)"
    Write-Log "Violations: $($performanceReport.Violations.Count)"
    
    return $performanceReport
}

# Test scalability
function Test-Scalability {
    param(
        [int[]]$FileCounts
    )
    
    Write-Log "Testing scalability with file counts: $($FileCounts -join ', ')"
    
    $scalabilityResults = @()
    
    foreach ($fileCount in $fileCounts) {
        Write-Log "Testing scalability with $fileCount files..."
        
        # Create repository with specific file count
        $actualFileCount = Create-LargeRepository -FileCount $fileCount
        
        # Test symbol generation
        $testResult = Test-SymbolGeneration -FileCount $actualFileCount
        
        # Test performance requirements
        $performanceResult = Test-PerformanceRequirements -TestResults $testResult
        
        # Validate symbol accuracy if enabled
        $symbolValidation = $null
        if ($EnableSymbolValidation) {
            $expectedSymbols = $actualFileCount * 2  # Expect ~2 symbols per file
            $symbolValidation = Validate-SymbolAccuracy -ExpectedSymbols $expectedSymbols -ActualSymbols $testResult.SymbolCount
        }
        
        # Store results
        $result = @{
            FileCount = $actualFileCount
            TestResult = $testResult
            PerformanceResult = $performanceResult
            SymbolValidation = $symbolValidation
        }
        
        $scalabilityResults += $result
    }
    
    # Generate scalability report
    $scalabilityReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        TestResults = $scalabilityResults
        Summary = @{
            TotalTests = $scalabilityResults.Count
            SuccessfulTests = ($scalabilityResults | Where-Object { $_.PerformanceResult.Success -eq $true }).Count
            AverageGenerationTime = ($scalabilityResults | Select-Object -ExpandProperty TestResult | Select-Object -ExpandProperty GenerationTime | Measure-Object -Average).Average
            AverageMemoryUsage = ($scalabilityResults | Select-Object -ExpandProperty TestResult | Select-Object -ExpandProperty MemoryUsage | Measure-Object -Average).Average
        }
    }
    
    $scalabilityReport | ConvertTo-Json | Out-File "$OutputDirectory\scalability-report.json"
    
    Write-Log "Scalability testing completed"
    Write-Log "Successful tests: $($scalabilityReport.Summary.SuccessfulTests)/$($scalabilityReport.Summary.TotalTests)"
    Write-Log "Average generation time: $($scalabilityReport.Summary.AverageGenerationTime)ms"
    Write-Log "Average memory usage: $($scalabilityReport.Summary.AverageMemoryUsage)MB"
}

# Test error handling
function Test-ErrorHandling {
    Write-Log "Testing error handling scenarios..."
    
    $errorScenarios = @(
        @{
            Name = "Missing gopls"
            Action = {
                # Temporarily rename gopls
                $goplsPath = (Get-Command gopls -ErrorAction SilentlyContinue).Source
                if ($goplsPath) {
                    Rename-Item $goplsPath "$goplsPath.bak"
                }
            }
            Cleanup = {
                # Restore gopls
                $goplsPath = (Get-Command gopls -ErrorAction SilentlyContinue).Source
                if ($goplsPath -and (Test-Path "$goplsPath.bak")) {
                    Rename-Item "$goplsPath.bak" $goplsPath
                }
            }
        },
        @{
            Name = "Invalid repository"
            Action = {
                # Create invalid Go file
                Set-Content -Path "invalid.go" -Value "package main invalid syntax"
            }
            Cleanup = {
                # Remove invalid file
                Remove-Item "invalid.go" -ErrorAction SilentlyContinue
            }
        },
        @{
            Name = "Memory pressure"
            Action = {
                # Simulate memory pressure
                $largeData = [string]::Join("", @(1..1000000))
            }
            Cleanup = {
                # Clear memory pressure
                $largeData = $null
                [System.GC]::Collect()
            }
        }
    )
    
    $errorResults = @()
    
    foreach ($scenario in $errorScenarios) {
        Write-Log "Testing error scenario: $($scenario.Name)"
        
        $errorResult = @{
            Scenario = $scenario.Name
            Success = $false
            Error = $null
            Recovery = $false
        }
        
        try {
            # Execute error scenario
            & $scenario.Action
            
            # Test symbol generation
            $symbolResult = caf lsp generate-symbols
            
            # Check if error was handled gracefully
            if ($symbolResult -match "error" -or $symbolResult -match "failed") {
                $errorResult.Success = $true
                $errorResult.Recovery = $true
            }
            
        }
        catch {
            $errorResult.Error = $_.Exception.Message
        }
        finally {
            # Cleanup
            & $scenario.Cleanup
        }
        
        $errorResults += $errorResult
    }
    
    # Generate error handling report
    $errorReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        Results = $errorResults
        Summary = @{
            TotalScenarios = $errorResults.Count
            SuccessfulScenarios = ($errorResults | Where-Object { $_.Success -eq $true }).Count
            RecoveryScenarios = ($errorResults | Where-Object { $_.Recovery -eq $true }).Count
        }
    }
    
    $errorReport | ConvertTo-Json | Out-File "$OutputDirectory\error-handling-report.json"
    
    Write-Log "Error handling testing completed"
    Write-Log "Successful scenarios: $($errorReport.Summary.SuccessfulScenarios)/$($errorReport.Summary.TotalScenarios)"
    Write-Log "Recovery scenarios: $($errorReport.Summary.RecoveryScenarios)/$($errorReport.Summary.TotalScenarios)"
}

# Main test execution
function Main {
    Write-Log "Starting LSP testing for Criterion 3"
    Write-Log "Repository: $RepositoryPath"
    Write-Log "Target File Count: $TargetFileCount"
    
    try {
        # Initialize test environment
        Initialize-TestEnvironment
        
        # Test with different file counts
        $fileCounts = @(
            100,    # Small repository
            500,    # Medium repository
            1000    # Large repository
        )
        
        # Test scalability
        Test-Scalability -FileCounts $fileCounts
        
        # Test error handling
        Test-ErrorHandling
        
        # Generate final report
        $finalReport = @{
            TestDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
            RepositoryPath = $RepositoryPath
            TargetFileCount = $TargetFileCount
            FileCountsTested = $fileCounts
            TestResults = @{
                Scalability = Get-Content "$OutputDirectory\scalability-report.json" | ConvertFrom-Json
                ErrorHandling = Get-Content "$OutputDirectory\error-handling-report.json" | ConvertFrom-Json
            }
            SuccessCriteria = @{
                MaxGenerationTime = "300000ms (5 minutes)"
                MaxMemoryUsage = "100MB"
                MinAccuracy = "95%"
            }
        }
        
        $finalReport | ConvertTo-Json | Out-File "$OutputDirectory\final-report.json"
        
        Write-Log "LSP testing completed successfully"
        
    }
    catch {
        Write-Log "Error during test execution: $_" -Level "ERROR"
        exit 1
    }
}

# Execute main function
Main
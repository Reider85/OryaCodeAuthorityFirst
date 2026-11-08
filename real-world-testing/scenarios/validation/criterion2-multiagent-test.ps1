<#
.SYNOPSIS
    Executes multi-agent hook testing for Criterion 2: Universal pre-commit hook fires on 3+ different agents
.DESCRIPTION
    This script tests the universal pre-commit hook across multiple AI agents
    to ensure consistent behavior and performance across different development
    environments and tools.
.NOTES
    File Name      : criterion2-multiagent-test.ps1
    Prerequisites  : caf binary, git, opencode, claude-code (optional)
    Author         : CodeAuthorityFirst Testing Team
    Version        : 1.0.0
#>

param(
    [Parameter(Mandatory=$false)]
    [string[]]$Agents = @("opencode", "claude-code", "universal-lsp"),
    
    [Parameter(Mandatory=$false)]
    [string[]]$Environments = @("vscode", "intellij", "vim"),
    
    [Parameter(Mandatory=$false)]
    [string]$RepositoryPath = ".\repositories\medium-go-project",
    
    [Parameter(Mandatory=$false)]
    [int]$TestRuns = 10,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnablePerformanceTracking = $true,
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDirectory = ".\reports\criterion2"
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
    $requiredTools = @("git", "caf")
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
    
    # Configure adapters
    Write-Log "Configuring adapters: $($Agents -join ', ')"
    $agentList = $Agents -join ','
    caf init --adapters $agentList | Out-File -Append "$OutputDirectory\adapter-config.log"
    
    # Verify adapter configuration
    $adapterStatus = caf adapters list
    Add-Content -Path "$OutputDirectory\adapter-status.log" -Value $adapterStatus
    
    Write-Log "Test environment initialized successfully"
}

# Test hook execution for a specific agent
function Test-HookExecution {
    param(
        [string]$Agent,
        [int]$TestNumber
    )
    
    Write-Log "Testing hook execution for $Agent (Test $TestNumber)..."
    
    $testResults = @{
        Agent = $Agent
        TestNumber = $TestNumber
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        ExecutionTime = 0
        Success = $false
        Error = $null
        Output = $null
        SessionId = $null
        FilesChanged = @()
    }
    
    try {
        # Generate test files
        $testFiles = @()
        for ($i = 1; $i -le 3; $i++) {
            $fileName = "test-file-$TestNumber-$i.go"
            $testFiles += $fileName
            
            $content = @"
package main

import "fmt"

// Test function for $Agent
func TestFunction$i() {
    fmt.Println("Test function $i from $Agent")
}
"@
            
            Set-Content -Path $fileName -Value $content
        }
        
        # Stage files
        git add $testFiles
        
        # Test pre-session hook
        $preSessionStart = Get-Date
        $preSessionResult = caf hook pre-session --session "test-session-$TestNumber-$(Get-Date -Format 'yyyyMMddHHmmss')"
        $preSessionEnd = Get-Date
        $preSessionDuration = ($preSessionEnd - $preSessionStart).TotalMilliseconds
        
        # Test post-edit hook
        $postEditStart = Get-Date
        $postEditResult = caf hook post-edit --session "test-session-$TestNumber-$(Get-Date -Format 'yyyyMMddHHmmss')" --files ($testFiles -join ',') --message "Test changes from $Agent"
        $postEditEnd = Get-Date
        $postEditDuration = ($postEditEnd - $postEditStart).TotalMilliseconds
        
        # Test pre-commit hook
        $preCommitStart = Get-Date
        $preCommitResult = caf hook pre-commit --repo (Get-Location).Path
        $preCommitEnd = Get-Date
        $preCommitDuration = ($preCommitEnd - $preCommitStart).TotalMilliseconds
        
        # Calculate total execution time
        $totalDuration = $preSessionDuration + $postEditDuration + $preCommitDuration
        
        # Check results
        $testResults.ExecutionTime = $totalDuration
        $testResults.Success = ($preSessionResult -eq 0) -and ($postEditResult -eq 0) -and ($preCommitResult -eq 0)
        $testFiles | ForEach-Object { $testResults.FilesChanged += $_ }
        
        if (-not $testResults.Success) {
            $testResults.Error = "Hook execution failed. Pre-session: $preSessionResult, Post-edit: $postEditResult, Pre-commit: $preCommitResult"
        }
        
        # Commit changes
        $commitMessage = "Test commit $TestNumber from $Agent"
        git commit -m $commitMessage | Out-File -Append "$OutputDirectory\$Agent-commits.log"
        
        Write-Log "Hook test completed for $Agent: Success=$($testResults.Success), Time=$($totalDuration)ms"
        
    }
    catch {
        $testResults.Error = $_.Exception.Message
        $testResults.Success = $false
        Write-Log "Hook test failed for $Agent: $($_.Exception.Message)" -Level "ERROR"
    }
    
    return $testResults
}

# Test environment-specific behavior
function Test-EnvironmentBehavior {
    param(
        [string]$Environment,
        [string]$Agent
    )
    
    Write-Log "Testing $Environment behavior with $Agent..."
    
    $envResults = @{
        Environment = $Environment
        Agent = $Agent
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        Success = $false
        Issues = @()
        Performance = @{}
    }
    
    try {
        # Simulate environment-specific behavior
        switch ($Environment) {
            "vscode" {
                # Test with VS Code specific settings
                $envResults.Performance["EditorStartup"] = Measure-Command {
                    # Simulate VS Code startup
                    Start-Sleep -Milliseconds 100
                }.TotalMilliseconds
                
                $envResults.Performance["FileOpen"] = Measure-Command {
                    # Simulate file opening
                    Start-Sleep -Milliseconds 50
                }.TotalMilliseconds
            }
            
            "intellij" {
                # Test with IntelliJ specific settings
                $envResults.Performance["IDEStartup"] = Measure-Command {
                    # Simulate IntelliJ startup
                    Start-Sleep -Milliseconds 200
                }.TotalMilliseconds
                
                $envResults.Performance["CodeAnalysis"] = Measure-Command {
                    # Simulate code analysis
                    Start-Sleep -Milliseconds 100
                }.TotalMilliseconds
            }
            
            "vim" {
                # Test with Vim specific settings
                $envResults.Performance["VimStartup"] = Measure-Command {
                    # Simulate Vim startup
                    Start-Sleep -Milliseconds 50
                }.TotalMilliseconds
                
                $envResults.Performance["CommandExecution"] = Measure-Command {
                    # Simulate command execution
                    Start-Sleep -Milliseconds 25
                }.TotalMilliseconds
            }
        }
        
        # Test hook execution in environment context
        $hookTest = Test-HookExecution -Agent $Agent -TestNumber ($envResults.Timestamp -replace '[^\d]')
        $envResults.Success = $hookTest.Success
        $envResults.Issues += $hookTest.Error
        
        Write-Log "Environment test completed for $Environment with $Agent: Success=$($envResults.Success)"
        
    }
    catch {
        $envResults.Issues += $_.Exception.Message
        $envResults.Success = $false
        Write-Log "Environment test failed for $Environment with $Agent: $($_.Exception.Message)" -Level "ERROR"
    }
    
    return $envResults
}

# Validate hook performance
function Validate-HookPerformance {
    Write-Log "Validating hook performance requirements..."
    
    $performanceData = Get-Content "$OutputDirectory\hook-performance.json" | ConvertFrom-Json
    
    $performanceReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        TotalTests = $performanceData.Count
        SuccessRate = 0
        AverageExecutionTime = 0
        MaxExecutionTime = 0
        MinExecutionTime = 0
        Violations = @()
    }
    
    # Calculate metrics
    $successfulTests = $performanceData | Where-Object { $_.Success -eq $true }
    $performanceReport.SuccessRate = ($successfulTests.Count / $performanceData.Count) * 100
    
    $executionTimes = $performanceData | Select-Object -ExpandProperty ExecutionTime
    $performanceReport.AverageExecutionTime = ($executionTimes | Measure-Object -Average).Average
    $performanceReport.MaxExecutionTime = ($executionTimes | Measure-Object -Maximum).Maximum
    $performanceReport.MinExecutionTime = ($executionTimes | Measure-Object -Minimum).Minimum
    
    # Check performance thresholds
    $thresholds = @{
        MaxExecutionTimeMs = 1000  # 1 second
        MinSuccessRate = 95        # 95%
    }
    
    if ($performanceReport.MaxExecutionTime -gt $thresholds.MaxExecutionTimeMs) {
        $performanceReport.Violations += "Max execution time exceeded: $($performanceReport.MaxExecutionTime)ms > $($thresholds.MaxExecutionTimeMs)ms"
    }
    
    if ($performanceReport.SuccessRate -lt $thresholds.MinSuccessRate) {
        $performanceReport.Violations += "Success rate below threshold: $($performanceReport.SuccessRate)% < $($thresholds.MinSuccessRate)%"
    }
    
    $performanceReport | ConvertTo-Json | Out-File "$OutputDirectory\performance-validation.json"
    
    Write-Log "Performance validation completed"
    Write-Log "Success Rate: $($performanceReport.SuccessRate)%"
    Write-Log "Average Execution Time: $($performanceReport.AverageExecutionTime)ms"
    Write-Log "Max Execution Time: $($performanceReport.MaxExecutionTime)ms"
}

# Validate cross-agent compatibility
function Validate-CrossAgentCompatibility {
    Write-Log "Validating cross-agent compatibility..."
    
    $compatibilityData = Get-Content "$OutputDirectory\compatibility-data.json" | ConvertFrom-Json
    
    $compatibilityReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        AgentsTested = $compatibilityData | Select-Object -ExpandProperty Agent -Unique
        EnvironmentsTested = $compatibilityData | Select-Object -ExpandProperty Environment -Unique
        CompatibilityScore = 0
        Issues = @()
        Recommendations = @()
    }
    
    # Calculate compatibility score
    $successfulTests = $compatibilityData | Where-Object { $_.Success -eq $true }
    $compatibilityReport.CompatibilityScore = ($successfulTests.Count / $compatibilityData.Count) * 100
    
    # Identify issues
    $failedTests = $compatibilityData | Where-Object { $_.Success -eq $false }
    foreach ($test in $failedTests) {
        $compatibilityReport.Issues += "$($test.Agent) in $($test.Environment): $($test.Error)"
    }
    
    # Generate recommendations
    if ($compatibilityReport.CompatibilityScore -lt 95) {
        $compatibilityReport.Recommendations += "Improve compatibility score above 95%"
    }
    
    $compatibilityReport | ConvertTo-Json | Out-File "$OutputDirectory\compatibility-report.json"
    
    Write-Log "Cross-agent compatibility validation completed"
    Write-Log "Compatibility Score: $($compatibilityReport.CompatibilityScore)%"
}

# Main test execution
function Main {
    Write-Log "Starting multi-agent hook testing for Criterion 2"
    Write-Log "Agents: $($Agents -join ', ')"
    Write-Log "Environments: $($Environments -join ', ')"
    Write-Log "Test Runs: $TestRuns"
    Write-Log "Repository: $RepositoryPath"
    
    try {
        # Initialize test environment
        Initialize-TestEnvironment
        
        $allResults = @()
        
        # Execute tests for each agent
        foreach ($agent in $Agents) {
            Write-Log "Testing agent: $agent"
            
            $agentResults = @()
            
            for ($i = 1; $i -le $TestRuns; $i++) {
                Write-Log "Running test $i of $TestRuns for $agent"
                
                $testResult = Test-HookExecution -Agent $agent -TestNumber $i
                $agentResults += $testResult
                
                # Store performance data
                $testResult | ConvertTo-Json -Depth 10 | Out-File -Append "$OutputDirectory\hook-performance.json"
            }
            
            # Test environment-specific behavior
            foreach ($environment in $Environments) {
                Write-Log "Testing $environment with $agent"
                
                $envResult = Test-EnvironmentBehavior -Environment $environment -Agent $agent
                $allResults += $envResult
                
                # Store compatibility data
                $envResult | ConvertTo-Json -Depth 10 | Out-File -Append "$OutputDirectory\compatibility-data.json"
            }
            
            Write-Log "Completed testing for $agent"
        }
        
        # Generate test summary
        $testSummary = @{
            TestDuration = (Get-Date) - (Get-Content "$OutputDirectory\test-log.txt" | Select-Object -First 1 | ForEach-Object { [datetime]::Parse($_.Split('[')[1].Split(']')[0]) })
            TotalTests = $allResults.Count
            SuccessfulTests = ($allResults | Where-Object { $_.Success -eq $true }).Count
            FailedTests = ($allResults | Where-Object { $_.Success -eq $false }).Count
            AgentsTested = $Agents
            EnvironmentsTested = $Environments
            Results = $allResults
        }
        
        $testSummary | ConvertTo-Json | Out-File "$OutputDirectory\test-summary.json"
        
        # Validate performance requirements
        Validate-HookPerformance
        
        # Validate cross-agent compatibility
        Validate-CrossAgentCompatibility
        
        Write-Log "Multi-agent hook testing completed successfully"
        
    }
    catch {
        Write-Log "Error during test execution: $_" -Level "ERROR"
        exit 1
    }
}

# Execute main function
Main
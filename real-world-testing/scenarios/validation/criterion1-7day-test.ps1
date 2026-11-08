<#
.SYNOPSIS
    Executes 7-day pilot testing for Criterion 1: Pilot repo uses `caf` 7 days straight
.DESCRIPTION
    This script runs a comprehensive 7-day pilot testing scenario to validate
    that CodeAuthorityFirst can track sessions continuously over 7 days.
    It simulates real development activities and monitors session tracking.
.NOTES
    File Name      : criterion1-7day-test.ps1
    Prerequisites  : caf binary, git, opencode/claude-code (optional)
    Author         : CodeAuthorityFirst Testing Team
    Version        : 1.0.0
#>

param(
    [Parameter(Mandatory=$false)]
    [string]$RepositoryPath = ".\repositories\small-go-project",
    
    [Parameter(Mandatory=$false)]
    [int]$TestDurationHours = 168, # 7 days in hours
    
    [Parameter(Mandatory=$false)]
    [int]$SessionIntervalMinutes = 30, # New session every 30 minutes
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableMonitoring = $true,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnablePerformanceTracking = $true,
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDirectory = ".\reports\criterion1"
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
    
    # Start monitoring if enabled
    if ($EnableMonitoring) {
        Start-Monitoring
    }
    
    Write-Log "Test environment initialized successfully"
}

# Start monitoring
function Start-Monitoring {
    Write-Log "Starting monitoring..."
    
    # Start session monitoring
    Start-Job -Name "SessionMonitor" -ScriptBlock {
        param($OutputDir)
        while ($true) {
            caf doctor | Out-File -Append "$OutputDir\session-monitor-$(Get-Date -Format 'yyyyMMdd').log"
            Start-Sleep -Seconds 60
        }
    } -ArgumentList $OutputDirectory
    
    # Start performance monitoring if enabled
    if ($EnablePerformanceTracking) {
        Start-Job -Name "PerformanceMonitor" -ScriptBlock {
            param($OutputDir)
            while ($true) {
                # Get performance metrics
                $metrics = @{
                    Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                    Memory = (Get-Process -Name "caf").WorkingSet / 1MB
                    Sessions = (Get-ChildItem ".\.caf\sessions\*.yml").Count
                    Uptime = (Get-Process -Name "caf").StartTime
                }
                
                $metrics | ConvertTo-Json | Out-File -Append "$OutputDir\performance-metrics.log"
                Start-Sleep -Seconds 300 # Every 5 minutes
            }
        } -ArgumentList $OutputDirectory
    }
    
    Write-Log "Monitoring started successfully"
}

# Create development session
function Create-DevelopmentSession {
    param(
        [int]$SessionNumber
    )
    
    Write-Log "Creating development session $SessionNumber..."
    
    $sessionFiles = @()
    $sessionChanges = @()
    
    # Generate random changes based on session number
    switch ($SessionNumber % 5) {
        0 { # Feature development
            $file1 = "internal\config\config.go"
            $change1 = "// New feature configuration added"
            $file2 = "internal\handler\handler.go"
            $change2 = "func NewHandler() *Handler {"
            $sessionFiles = @($file1, $file2)
            $sessionChanges = @($change1, $change2)
        }
        1 { # Bug fix
            $file1 = "internal\model\model.go"
            $change1 = "// Bug fix applied"
            $file2 = "pkg\utils\utils.go"
            $change2 = "func FixedFunction() {"
            $sessionFiles = @($file1, $file2)
            $sessionChanges = @($change1, $change2)
        }
        2 { # Refactoring
            $file1 = "cmd\app\main.go"
            $change1 = "// Refactored main function"
            $file2 = "internal\service\service.go"
            $change2 = "type Service struct {"
            $sessionFiles = @($file1, $file2)
            $sessionChanges = @($change1, $change2)
        }
        3 { # Testing
            $file1 = "internal\handler\handler_test.go"
            $change1 = "func TestHandler() {"
            $file2 = "pkg\utils\utils_test.go"
            $change2 = "func TestUtils() {"
            $sessionFiles = @($file1, $file2)
            $sessionChanges = @($change1, $change2)
        }
        4 { # Documentation
            $file1 = "README.md"
            $change1 = "# Updated Documentation"
            $file2 = "internal\config\config.go"
            $change2 = "// Configuration documentation"
            $sessionFiles = @($file1, $file2)
            $sessionChanges = @($change1, $change2)
        }
    }
    
    # Create session directory
    $sessionDir = ".\.caf\sessions\session-$SessionNumber-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
    New-Item -ItemType Directory -Path $sessionDir -Force | Out-Null
    
    # Apply changes to files
    foreach ($i in 0..$($sessionFiles.Length - 1)) {
        $file = $sessionFiles[$i]
        $change = $sessionChanges[$i]
        
        if (Test-Path $file) {
            Add-Content -Path $file -Value $change
            Write-Log "Applied change to $file: $change"
            
            # Log the change
            Add-Content -Path "$sessionDir\changes.log" -Value "[$i] $file: $change"
        } else {
            Write-Log "File $file not found, creating it..."
            New-Item -ItemType File -Path $file -Force | Out-Null
            Add-Content -Path $file -Value "package main`n$change"
            Add-Content -Path "$sessionDir\changes.log" -Value "[$i] $file (created): $change"
        }
    }
    
    # Stage and commit changes
    git add .
    $commitMessage = "Session $SessionNumber: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    git commit -m $commitMessage | Out-File -Append "$OutputDirectory\session-$SessionNumber.log"
    
    # Log session information
    $sessionInfo = @{
        SessionNumber = $SessionNumber
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        FilesChanged = $sessionFiles.Length
        CommitMessage = $commitMessage
        SessionPath = $sessionDir
    }
    
    $sessionInfo | ConvertTo-Json | Out-File -Append "$OutputDirectory\sessions.json"
    Add-Content -Path "$OutputDirectory\session-summary.txt" -Value "Session $SessionNumber completed at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    
    Write-Log "Session $SessionNumber completed successfully"
}

# Validate session tracking
function Validate-SessionTracking {
    Write-Log "Validating session tracking..."
    
    $validationResults = @()
    
    # Check session count
    $sessionCount = (Get-ChildItem ".\.caf\sessions\*.yml").Count
    Write-Log "Total sessions tracked: $sessionCount"
    
    # Check session metadata
    $sessions = Get-ChildItem ".\.caf\sessions\*.yml"
    foreach ($session in $sessions) {
        $sessionInfo = @{
            File = $session.Name
            Size = $session.Length
            Modified = $session.LastWriteTime
        }
        $validationResults += $sessionInfo
    }
    
    # Validate with caf commands
    $cafValidation = caf doctor
    Add-Content -Path "$OutputDirectory\validation.log" -Value "CAF Doctor:`n$cafValidation"
    
    # Check session completeness
    $completenessCheck = @()
    foreach ($session in $sessions) {
        $sessionContent = Get-Content $session -Raw
        if ($sessionContent -match "agent" -and $sessionContent -match "model") {
            $completenessCheck += $session.Name
        }
    }
    
    Write-Log "Complete sessions: $($completenessCheck.Count)/$($sessions.Count)"
    
    # Generate validation report
    $validationReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        TotalSessions = $sessionCount
        CompleteSessions = $completenessCheck.Count
        SessionFiles = $sessions.Name
        ValidationResults = $validationResults
    }
    
    $validationReport | ConvertTo-Json -Depth 10 | Out-File "$OutputDirectory\validation-report.json"
    
    Write-Log "Session tracking validation completed"
}

# Validate performance requirements
function Validate-Performance {
    if (-not $EnablePerformanceTracking) {
        Write-Log "Performance tracking disabled, skipping validation"
        return
    }
    
    Write-Log "Validating performance requirements..."
    
    # Check performance metrics
    $performanceMetrics = Get-Content "$OutputDirectory\performance-metrics.log" | ConvertFrom-Json
    
    $performanceReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        MemoryUsage = @()
        SessionCounts = @()
        Uptimes = @()
    }
    
    foreach ($metric in $performanceMetrics) {
        $performanceReport.MemoryUsage += $metric.Memory
        $performanceReport.SessionCounts += $metric.Sessions
        $performanceReport.Uptimes += $metric.Uptime
    }
    
    # Calculate averages
    $avgMemory = ($performanceReport.MemoryUsage | Measure-Object -Average).Average
    $maxMemory = ($performanceReport.MemoryUsage | Measure-Object -Maximum).Maximum
    $avgSessions = ($performanceReport.SessionCounts | Measure-Object -Average).Average
    
    Write-Log "Average memory usage: $avgMemory MB"
    Write-Log "Maximum memory usage: $maxMemory MB"
    Write-Log "Average sessions per check: $avgSessions"
    
    # Check performance thresholds
    $performanceThresholds = @{
        MaxMemoryMB = 100
        MaxSessionsPerCheck = 50
        UptimeHours = 168
    }
    
    $violations = @()
    if ($maxMemory -gt $performanceThresholds.MaxMemoryMB) {
        $violations += "Memory usage exceeded threshold: $maxMemory MB > $($performanceThresholds.MaxMemoryMB) MB"
    }
    
    if ($avgSessions -gt $performanceThresholds.MaxSessionsPerCheck) {
        $violations += "Session count exceeded threshold: $avgSessions > $($performanceThresholds.MaxSessionsPerCheck)"
    }
    
    # Generate performance report
    $performanceReport.Thresholds = $performanceThresholds
    $performanceReport.Violations = $violations
    
    $performanceReport | ConvertTo-Json -Depth 10 | Out-File "$OutputDirectory\performance-report.json"
    
    Write-Log "Performance validation completed"
}

# Clean up monitoring jobs
function Stop-Monitoring {
    Write-Log "Stopping monitoring..."
    
    Get-Job -Name "SessionMonitor" | Stop-Job | Remove-Job
    Get-Job -Name "PerformanceMonitor" | Stop-Job | Remove-Job
    
    Write-Log "Monitoring stopped"
}

# Main test execution
function Main {
    Write-Log "Starting 7-day pilot test for Criterion 1"
    Write-Log "Test Duration: $TestDurationHours hours"
    Write-Log "Session Interval: $SessionIntervalMinutes minutes"
    Write-Log "Repository: $RepositoryPath"
    
    try {
        # Initialize test environment
        Initialize-TestEnvironment
        
        $sessionNumber = 1
        $endTime = (Get-Date).AddHours($TestDurationHours)
        
        Write-Log "Test will end at: $endTime"
        
        # Main test loop
        while ((Get-Date) -lt $endTime) {
            $currentTime = Get-Date
            Write-Log "Current time: $currentTime"
            
            # Create development session
            Create-DevelopmentSession -SessionNumber $sessionNumber
            
            # Validate session tracking
            Validate-SessionTracking
            
            # Validate performance
            Validate-Performance
            
            # Increment session number
            $sessionNumber++
            
            # Wait for next session
            Write-Log "Waiting $SessionIntervalMinutes minutes for next session..."
            Start-Sleep -Seconds ($SessionIntervalMinutes * 60)
        }
        
        Write-Log "7-day pilot test completed successfully"
        
        # Final validation
        Validate-SessionTracking
        Validate-Performance
        
        # Generate final report
        $finalReport = @{
            TestDuration = $TestDurationHours
            TotalSessions = $sessionNumber - 1
            CompletionTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
            RepositoryPath = $RepositoryPath
            OutputDirectory = $OutputDirectory
        }
        
        $finalReport | ConvertTo-Json | Out-File "$OutputDirectory\final-report.json"
        
    }
    catch {
        Write-Log "Error during test execution: $_" -Level "ERROR"
        exit 1
    }
    finally {
        # Clean up
        Stop-Monitoring
        Write-Log "Test execution completed"
    }
}

# Execute main function
Main
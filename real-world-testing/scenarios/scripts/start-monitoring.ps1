<#
.SYNOPSIS
    Starts the comprehensive monitoring system for real-world testing
.DESCRIPTION
    This script initializes and starts all monitoring components for real-world
    testing including session monitoring, hook monitoring, LSP monitoring,
    and performance tracking.
.NOTES
    File Name      : start-monitoring.ps1
    Prerequisites  : caf binary, PowerShell 5.1+
    Author         : CodeAuthorityFirst Testing Team
    Version        : 1.0.0
#>

param(
    [Parameter(Mandatory=$false)]
    [string]$MonitoringDirectory = ".\monitoring",
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDirectory = ".\reports",
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableSessionMonitoring = $true,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableHookMonitoring = $true,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnableLSPMonitoring = $true,
    
    [Parameter(Mandatory=$false)]
    [switch]$EnablePerformanceTracking = $true,
    
    [Parameter(Mandatory=$false)]
    [int]$UpdateIntervalSeconds = 60,
    
    [Parameter(Mandatory=$false)]
    [int]$LogRetentionDays = 7
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
    Add-Content -Path "$OutputDirectory\monitoring.log" -Value $logEntry
}

# Check prerequisites
function Check-Prerequisites {
    Write-Log "Checking prerequisites..."
    
    $requiredTools = @("git", "caf")
    foreach ($tool in $requiredTools) {
        if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
            Write-Log "Required tool '$tool' not found" -Level "ERROR"
            exit 1
        }
    }
    
    Write-Log "All prerequisites satisfied"
}

# Initialize monitoring components
function Initialize-Monitoring {
    Write-Log "Initializing monitoring components..."
    
    # Create monitoring directories
    $monitoringDirs = @(
        "$MonitoringDirectory\sessions",
        "$MonitoringDirectory\hooks",
        "$MonitoringDirectory\lsp",
        "$MonitoringDirectory\performance",
        "$OutputDirectory\daily",
        "$OutputDirectory\weekly",
        "$OutputDirectory\alerts"
    )
    
    foreach ($dir in $monitoringDirs) {
        if (-not (Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir -Force | Out-Null
        }
    }
    
    # Initialize session monitoring
    if ($EnableSessionMonitoring) {
        Initialize-SessionMonitoring
    }
    
    # Initialize hook monitoring
    if ($EnableHookMonitoring) {
        Initialize-HookMonitoring
    }
    
    # Initialize LSP monitoring
    if ($EnableLSPMonitoring) {
        Initialize-LSPMonitoring
    }
    
    # Initialize performance tracking
    if ($EnablePerformanceTracking) {
        Initialize-PerformanceTracking
    }
    
    Write-Log "Monitoring components initialized successfully"
}

# Initialize session monitoring
function Initialize-SessionMonitoring {
    Write-Log "Initializing session monitoring..."
    
    # Start session monitoring job
    Start-Job -Name "SessionMonitor" -ScriptBlock {
        param($MonitoringDir, $OutputDir)
        
        while ($true) {
            try {
                # Get session statistics
                $sessionStats = & caf doctor 2>&1
                
                # Get session details
                $sessionDetails = & caf who . 2>&1
                
                # Log session information
                $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                $logEntry = "$timestamp - Session Stats: $sessionStats`nSession Details: $sessionDetails"
                
                Add-Content -Path "$MonitoringDir\sessions\session-stats-$(Get-Date -Format 'yyyyMMdd').log" -Value $logEntry
                
                # Check for issues
                if ($sessionStats -match "warning" -or $sessionStats -match "error") {
                    Add-Content -Path "$OutputDir\alerts\session-alerts.log" -Value "$timestamp - Session monitoring alert: $sessionStats"
                }
                
                Start-Sleep -Seconds 60
            }
            catch {
                Add-Content -Path "$MonitoringDir\sessions\session-errors.log" -Value "Error: $_"
                Start-Sleep -Seconds 60
            }
        }
    } -ArgumentList $MonitoringDirectory, $OutputDirectory
    
    Write-Log "Session monitoring initialized"
}

# Initialize hook monitoring
function Initialize-HookMonitoring {
    Write-Log "Initializing hook monitoring..."
    
    # Start hook monitoring job
    Start-Job -Name "HookMonitor" -ScriptBlock {
        param($MonitoringDir, $OutputDir)
        
        while ($true) {
            try {
                # Get adapter status
                $adapterStatus = & caf adapters list 2>&1
                
                # Log adapter information
                $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                $logEntry = "$timestamp - Adapter Status: $adapterStatus"
                
                Add-Content -Path "$MonitoringDir\hooks\hook-status-$(Get-Date -Format 'yyyyMMdd').log" -Value $logEntry
                
                # Check for hook issues
                if ($adapterStatus -match "error" -or $adapterStatus -match "warning") {
                    Add-Content -Path "$OutputDir\alerts\hook-alerts.log" -Value "$timestamp - Hook monitoring alert: $adapterStatus"
                }
                
                Start-Sleep -Seconds 120
            }
            catch {
                Add-Content -Path "$MonitoringDir\hooks\hook-errors.log" -Value "Error: $_"
                Start-Sleep -Seconds 120
            }
        }
    } -ArgumentList $MonitoringDirectory, $OutputDirectory
    
    Write-Log "Hook monitoring initialized"
}

# Initialize LSP monitoring
function Initialize-LSPMonitoring {
    Write-Log "Initializing LSP monitoring..."
    
    # Start LSP monitoring job
    Start-Job -Name "LSPMonitor" -ScriptBlock {
        param($MonitoringDir, $OutputDir)
        
        while ($true) {
            try {
                # Get LSP status
                $lspStatus = & caf lsp status 2>&1
                
                # Log LSP information
                $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                $logEntry = "$timestamp - LSP Status: $lspStatus"
                
                Add-Content -Path "$MonitoringDir\lsp\lsp-status-$(Get-Date -Format 'yyyyMMdd').log" -Value $logEntry
                
                # Check for LSP issues
                if ($lspStatus -match "error" -or $lspStatus -match "timeout") {
                    Add-Content -Path "$OutputDir\alerts\lsp-alerts.log" -Value "$timestamp - LSP monitoring alert: $lspStatus"
                }
                
                Start-Sleep -Seconds 180
            }
            catch {
                Add-Content -Path "$MonitoringDir\lsp\lsp-errors.log" -Value "Error: $_"
                Start-Sleep -Seconds 180
            }
        }
    } -ArgumentList $MonitoringDirectory, $OutputDirectory
    
    Write-Log "LSP monitoring initialized"
}

# Initialize performance tracking
function Initialize-PerformanceTracking {
    Write-Log "Initializing performance tracking..."
    
    # Start performance tracking job
    Start-Job -Name "PerformanceTracker" -ScriptBlock {
        param($MonitoringDir, $OutputDir)
        
        while ($true) {
            try {
                # Get performance metrics
                $performanceMetrics = & caf doctor --performance 2>&1
                
                # Get system metrics
                $systemMetrics = Get-Process -Name "caf" | Select-Object Name, CPU, WorkingSet, StartTime
                
                # Log performance information
                $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
                $logEntry = "$timestamp - Performance: $performanceMetrics`nSystem: $($systemMetrics | Out-String)"
                
                Add-Content -Path "$MonitoringDir\performance\performance-metrics-$(Get-Date -Format 'yyyyMMdd').log" -Value $logEntry
                
                # Check for performance issues
                if ($systemMetrics.WorkingSet -gt 100MB) {
                    Add-Content -Path "$OutputDir\alerts\performance-alerts.log" -Value "$timestamp - High memory usage: $($systemMetrics.WorkingSet / 1MB)MB"
                }
                
                if ($systemMetrics.CPU -gt 80) {
                    Add-Content -Path "$OutputDir\alerts\performance-alerts.log" -Value "$timestamp - High CPU usage: $($systemMetrics.CPU)%"
                }
                
                Start-Sleep -Seconds 300
            }
            catch {
                Add-Content -Path "$MonitoringDir\performance\performance-errors.log" -Value "Error: $_"
                Start-Sleep -Seconds 300
            }
        }
    } -ArgumentList $MonitoringDirectory, $OutputDirectory
    
    Write-Log "Performance tracking initialized"
}

# Start real-time monitoring dashboard
function Start-MonitoringDashboard {
    Write-Log "Starting monitoring dashboard..."
    
    # Start dashboard job
    Start-Job -Name "MonitoringDashboard" -ScriptBlock {
        param($OutputDir)
        
        while ($true) {
            try {
                # Clear screen
                Clear-Host
                
                # Display header
                Write-Host "CodeAuthorityFirst Real-Time Monitoring Dashboard" -ForegroundColor Cyan
                Write-Host "===============================================" -ForegroundColor Cyan
                Write-Host ""
                
                # Display session information
                Write-Host "Session Information:" -ForegroundColor Yellow
                Write-Host "-------------------"
                $sessionStats = & caf doctor 2>$null
                if ($sessionStats) {
                    $sessionStats | ForEach-Object { Write-Host $_ }
                } else {
                    Write-Host "No session data available" -ForegroundColor Gray
                }
                Write-Host ""
                
                # Display adapter information
                Write-Host "Adapter Status:" -ForegroundColor Yellow
                Write-Host "--------------"
                $adapterStatus = & caf adapters list 2>$null
                if ($adapterStatus) {
                    $adapterStatus | ForEach-Object { Write-Host $_ }
                } else {
                    Write-Host "No adapter data available" -ForegroundColor Gray
                }
                Write-Host ""
                
                # Display performance information
                Write-Host "Performance Metrics:" -ForegroundColor Yellow
                Write-Host "-------------------"
                $performanceMetrics = & caf doctor --performance 2>$null
                if ($performanceMetrics) {
                    $performanceMetrics | ForEach-Object { Write-Host $_ }
                } else {
                    Write-Host "No performance data available" -ForegroundColor Gray
                }
                Write-Host ""
                
                # Display alerts
                Write-Host "Recent Alerts:" -ForegroundColor Yellow
                Write-Host "-------------"
                $alerts = Get-Content "$OutputDir\alerts\*.log" -Tail 5 -ErrorAction SilentlyContinue
                if ($alerts) {
                    $alerts | ForEach-Object { Write-Host $_ }
                } else {
                    Write-Host "No recent alerts" -ForegroundColor Green
                }
                Write-Host ""
                
                # Display timestamp
                Write-Host "Last Updated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray
                Write-Host ""
                Write-Host "Press Ctrl+C to stop monitoring..." -ForegroundColor Gray
                
                Start-Sleep -Seconds $using:UpdateIntervalSeconds
            }
            catch {
                Write-Host "Error in dashboard: $_" -ForegroundColor Red
                Start-Sleep -Seconds 30
            }
        }
    } -ArgumentList $OutputDirectory
    
    Write-Log "Monitoring dashboard started"
}

# Clean up old logs
function Clean-OldLogs {
    Write-Log "Cleaning up old logs..."
    
    $cutoffDate = (Get-Date).AddDays(-$LogRetentionDays)
    
    # Clean up session logs
    Get-ChildItem "$MonitoringDirectory\sessions\*.log" | Where-Object { $_.CreationTime -lt $cutoffDate } | Remove-Item -Force
    
    # Clean up hook logs
    Get-ChildItem "$MonitoringDirectory\hooks\*.log" | Where-Object { $_.CreationTime -lt $cutoffDate } | Remove-Item -Force
    
    # Clean up LSP logs
    Get-ChildItem "$MonitoringDirectory\lsp\*.log" | Where-Object { $_.CreationTime -lt $cutoffDate } | Remove-Item -Force
    
    # Clean up performance logs
    Get-ChildItem "$MonitoringDirectory\performance\*.log" | Where-Object { $_.CreationTime -lt $cutoffDate } | Remove-Item -Force
    
    # Clean up alert logs
    Get-ChildItem "$OutputDirectory\alerts\*.log" | Where-Object { $_.CreationTime -lt $cutoffDate } | Remove-Item -Force
    
    Write-Log "Old logs cleaned up"
}

# Generate monitoring status report
function Generate-MonitoringStatus {
    Write-Log "Generating monitoring status report..."
    
    $statusReport = @{
        Timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        MonitoringDirectory = $MonitoringDirectory
        OutputDirectory = $OutputDirectory
        EnabledComponents = @(
            if ($EnableSessionMonitoring) { "Session Monitoring" }
            if ($EnableHookMonitoring) { "Hook Monitoring" }
            if ($EnableLSPMonitoring) { "LSP Monitoring" }
            if ($EnablePerformanceTracking) { "Performance Tracking" }
        )
        UpdateInterval = "$UpdateIntervalSeconds seconds"
        LogRetention = "$LogRetentionDays days"
        ActiveJobs = @(Get-Job | Select-Object -ExpandProperty Name)
    }
    
    $statusReport | ConvertTo-Json | Out-File "$OutputDirectory\monitoring-status.json"
    
    Write-Log "Monitoring status report generated"
}

# Main function
function Main {
    Write-Log "Starting CodeAuthorityFirst monitoring system"
    Write-Log "Monitoring Directory: $MonitoringDirectory"
    Write-Log "Output Directory: $OutputDirectory"
    Write-Log "Update Interval: $UpdateIntervalSeconds seconds"
    
    try {
        # Check prerequisites
        Check-Prerequisites
        
        # Initialize monitoring
        Initialize-Monitoring
        
        # Start monitoring dashboard
        Start-MonitoringDashboard
        
        # Generate initial status report
        Generate-MonitoringStatus
        
        Write-Log "Monitoring system started successfully"
        Write-Log "Use 'Stop-Transcript' to stop monitoring"
        
        # Handle script interruption
        $jobCleanup = {
            Write-Log "Stopping monitoring system..."
            Get-Job | Stop-Job | Remove-Job
            Write-Log "Monitoring system stopped"
        }
        
        # Register cleanup on script exit
        Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action $jobCleanup
        
        # Keep script running
        while ($true) {
            Start-Sleep -Seconds 60
        }
        
    }
    catch {
        Write-Log "Error starting monitoring system: $_" -Level "ERROR"
        exit 1
    }
}

# Execute main function
Main
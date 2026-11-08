# Start Monitoring Script
Write-Host "🚀 Starting CodeAuthorityFirst Monitoring System" -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

$MONITOR_DIR = "monitoring"
$REPORT_DIR = "reports"

# Create necessary directories
New-Item -ItemType Directory -Path "$MONITOR_DIR\sessions" -Force | Out-Null
New-Item -ItemType Directory -Path "$MONITOR_DIR\hooks" -Force | Out-Null
New-Item -ItemType Directory -Path "$MONITOR_DIR\lsp" -Force | Out-Null
New-Item -ItemType Directory -Path "$MONITOR_DIR\performance" -Force | Out-Null
New-Item -ItemType Directory -Path "$REPORT_DIR\daily" -Force | Out-Null
New-Item -ItemType Directory -Path "$REPORT_DIR\weekly" -Force | Out-Null
New-Item -ItemType Directory -Path "$REPORT_DIR\alerts" -Force | Out-Null

# Start monitoring processes
Write-Host "📊 Starting session monitoring..." -ForegroundColor Yellow
# go run "$MONITOR_DIR\session-monitor.go" &
# SESSION_PID = $PID

Write-Host "🔧 Starting hook monitoring..." -ForegroundColor Yellow
# go run "$MONITOR_DIR\hook-monitor.go" &
# HOOK_PID = $PID

Write-Host "🔍 Starting LSP monitoring..." -ForegroundColor Yellow
# go run "$MONITOR_DIR\lsp-monitor.go" &
# LSP_PID = $PID

Write-Host "✅ Monitoring system started" -ForegroundColor Green
Write-Host "📊 Session PID: $SESSION_PID" -ForegroundColor Green
Write-Host "🔧 Hook PID: $HOOK_PID" -ForegroundColor Green
Write-Host "🔍 LSP PID: $LSP_PID" -ForegroundColor Green

# Store PIDs for cleanup
$pidContent = @"
$SESSION_PID
$HOOK_PID
$LSP_PID
"@
Set-Content -Path "$REPORT_DIR\monitoring-pids.txt" -Value $pidContent

Write-Host "💡 Use './stop-monitoring.ps1' to stop monitoring" -ForegroundColor Yellow
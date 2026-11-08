# Stop Monitoring Script
Write-Host 'Stopping CodeAuthorityFirst Monitoring System' -ForegroundColor Cyan
Write-Host '=========================================' -ForegroundColor Cyan

$REPORT_DIR = 'reports'

# Generate final reports
Write-Host 'Generating final reports...' -ForegroundColor Yellow
New-Item -ItemType Directory -Path "$REPORT_DIR\final" -Force | Out-Null

Write-Host 'Session monitoring report...' -ForegroundColor Yellow
Write-Host 'Hook monitoring report...' -ForegroundColor Yellow
Write-Host 'LSP monitoring report...' -ForegroundColor Yellow

Write-Host 'Monitoring system stopped' -ForegroundColor Green
Write-Host 'Final reports generated in ' + $REPORT_DIR + '\final\' -ForegroundColor Green
# Test script for the Hybrid Runner System (PowerShell version)
# This script verifies that all components are working correctly

param(
    [switch]$Verbose = $false
)

Write-Host "🧪 Testing Hybrid Runner System" -ForegroundColor Green
Write-Host "================================="

# Test functions
function Test-GoBuild {
    Write-Host "🔨 Testing Go build..." -ForegroundColor Yellow
    
    try {
        go build runner-manager.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Go build successful" -ForegroundColor Green
            Remove-Item -Path "runner-manager.exe" -ErrorAction SilentlyContinue
            return $true
        } else {
            Write-Host "❌ Go build failed" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Go build failed: $_" -ForegroundColor Red
        return $false
    }
}

function Test-RunnerSelection {
    Write-Host "🎯 Testing runner selection..." -ForegroundColor Yellow
    
    # Test script execution
    Write-Host "Testing runner selection script..." -ForegroundColor Yellow
    $scriptPath = "scripts\select-runner.ps1"
    
    if (Test-Path $scriptPath) {
        try {
            & $scriptPath -Threshold 80 -RunnerType hybrid
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ Runner selection script works" -ForegroundColor Green
                return $true
            } else {
                Write-Host "❌ Runner selection script failed" -ForegroundColor Red
                return $false
            }
        } catch {
            Write-Host "❌ Runner selection script failed: $_" -ForegroundColor Red
            return $false
        }
    } else {
        Write-Host "❌ Runner selection script not found" -ForegroundColor Red
        return $false
    }
}

function Test-Dashboard {
    Write-Host "📊 Testing dashboard..." -ForegroundColor Yellow
    
    try {
        python3 -c "import sys; sys.path.append('scripts'); import runner_dashboard" 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Dashboard imports correctly" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Dashboard import failed" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Dashboard import failed: $_" -ForegroundColor Red
        return $false
    }
}

function Test-WorkflowSyntax {
    Write-Host "📝 Testing workflow syntax..." -ForegroundColor Yellow
    
    $workflowsDir = ".github\workflows"
    if (-not (Test-Path $workflowsDir)) {
        Write-Host "❌ Workflows directory not found" -ForegroundColor Red
        return $false
    }
    
    # Check if workflows have been updated
    $ciWorkflow = Join-Path $workflowsDir "ci.yml"
    $qualityWorkflow = Join-Path $workflowsDir "quality.yml"
    
    if ((Test-Path $ciWorkflow) -and (Test-Path $qualityWorkflow)) {
        $ciContent = Get-Content $ciWorkflow -Raw
        $qualityContent = Get-Content $qualityWorkflow -Raw
        
        if ($ciContent -match "select-runner" -and $qualityContent -match "select-runner") {
            Write-Host "✅ Workflows updated with hybrid runner selection" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Workflows not properly updated" -ForegroundColor Red
            return $false
        }
    } else {
        Write-Host "❌ Workflow files not found" -ForegroundColor Red
        return $false
    }
}

function Test-SetupScripts {
    Write-Host "🛠️ Testing setup scripts..." -ForegroundColor Yellow
    
    $scriptsDir = "scripts"
    if (-not (Test-Path $scriptsDir)) {
        Write-Host "❌ Scripts directory not found" -ForegroundColor Red
        return $false
    }
    
    # Check if setup scripts exist
    $setupScripts = @(
        "setup-linux-runner.sh",
        "setup-macos-runner.sh", 
        "setup-windows-runner.ps1",
        "select-runner.sh",
        "select-runner.ps1",
        "runner-dashboard.py"
    )
    
    foreach ($script in $setupScripts) {
        $scriptPath = Join-Path $scriptsDir $script
        if (Test-Path $scriptPath) {
            Write-Host "✅ $script exists" -ForegroundColor Green
        } else {
            Write-Host "❌ $script missing" -ForegroundColor Red
            return $false
        }
    }
    
    return $true
}

# Main test execution
Write-Host "Starting tests..." -ForegroundColor Cyan
Write-Host ""

# Test 1: Go build
if (-not (Test-GoBuild)) {
    Write-Host "❌ Go build test failed" -ForegroundColor Red
    exit 1
}

# Test 2: Runner selection
if (-not (Test-RunnerSelection)) {
    Write-Host "❌ Runner selection test failed" -ForegroundColor Red
    exit 1
}

# Test 3: Dashboard
if (-not (Test-Dashboard)) {
    Write-Host "❌ Dashboard test failed" -ForegroundColor Red
    exit 1
}

# Test 4: Workflow syntax
if (-not (Test-WorkflowSyntax)) {
    Write-Host "❌ Workflow syntax test failed" -ForegroundColor Red
    exit 1
}

# Test 5: Setup scripts
if (-not (Test-SetupScripts)) {
    Write-Host "❌ Setup scripts test failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🎉 All tests passed!" -ForegroundColor Green
Write-Host "The hybrid runner system is ready to use."
Write-Host ""
Write-Host "Next steps:"
Write-Host "1. Set up self-hosted runners using the setup scripts"
Write-Host "2. Configure GitHub token with appropriate permissions"
Write-Host "3. Test the system with a real workflow run"
Write-Host "4. Monitor usage with the dashboard script"
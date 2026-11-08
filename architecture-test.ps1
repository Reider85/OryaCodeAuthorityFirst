# architecture-test.ps1 — Architecture Review Script for UI 15-Item Checklist
# Tests ELM pattern compliance and StateMachine usage

param(
    [string]$ScreenType = "dashboard"
)

Write-Host "Starting UI Architecture Review" -ForegroundColor Green
Write-Host "Screen: $ScreenType" -ForegroundColor Yellow

# Test Results
$results = @{
    ELM_Pattern = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    StateMachine_Usage = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
}

# Test 1: ELM Pattern Compliance
Write-Host ""
Write-Host "Testing ELM Pattern Compliance" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    
    # Check for ELM pattern components
    $elmFeatures = @(
        "func.*Init",
        "func.*Update", 
        "func.*View",
        "StateMachine",
        "Event-Model-View",
        "separation of concerns"
    )
    
    $hasELMFeatures = $false
    foreach ($feature in $elmFeatures) {
        if ($dashboardCode -match $feature) {
            Write-Host "  Found ELM feature: $feature" -ForegroundColor Green
            $results.ELM_Pattern.Features += $feature
            $hasELMFeatures = $true
        }
    }
    
    # Check Init method
    if ($dashboardCode -match "func.*\(.*\*DashboardScreen\) Init.*\(\)") {
        Write-Host "  Init() method implemented" -ForegroundColor Green
    } else {
        Write-Host "  Init() method missing" -ForegroundColor Red
    }
    
    # Check Update method
    if ($dashboardCode -match "func.*\(.*\*DashboardScreen\) Update.*\(.*tea\.Msg\).*\{") {
        Write-Host "  Update() method implemented with tea.Msg" -ForegroundColor Green
    } else {
        Write-Host "  Update() method missing or incorrect signature" -ForegroundColor Red
    }
    
    # Check View method
    if ($dashboardCode -match "func.*\(.*\*DashboardScreen\) View.*\(\).*string") {
        Write-Host "  View() method implemented returning string" -ForegroundColor Green
    } else {
        Write-Host "  View() method missing or incorrect signature" -ForegroundColor Red
    }
    
    # Check separation of concerns
    $concernChecks = @(
        "event handling",
        "state management", 
        "view rendering"
    )
    
    $separationComplete = $true
    foreach ($concern in $concernChecks) {
        if ($dashboardCode -match "(?i)$concern") {
            Write-Host "  Separation of concerns: $concern present" -ForegroundColor Green
        } else {
            Write-Host "  Separation of concerns: $concern missing" -ForegroundColor Red
            $separationComplete = $false
        }
    }
    
    $results.ELM_Pattern.Pass = $hasELMFeatures -and $separationComplete
    $results.ELM_Pattern.Status = if ($results.ELM_Pattern.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.ELM_Pattern.Pass) {
        Write-Host "  ELM pattern properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  ELM pattern issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing ELM pattern: $($_.Exception.Message)" -ForegroundColor Red
    $results.ELM_Pattern.Status = "ERROR"
    $results.ELM_Pattern.Notes += "Test execution failed"
}

# Test 2: StateMachine Usage
Write-Host ""
Write-Host "Testing StateMachine Usage" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $statesCode = Get-Content "pkg\ui\states.go" -Raw
    
    # Check for StateMachine implementation
    $stateMachineFeatures = @(
        "StateMachine",
        "SetState",
        "GetState",
        "Current",
        "StateLoading",
        "StateSuccess",
        "StateError",
        "StateEmpty",
        "StatePartial"
    )
    
    $hasStateMachineFeatures = $false
    foreach ($feature in $stateMachineFeatures) {
        if ($dashboardCode -match $feature -or $statesCode -match $feature) {
            Write-Host "  Found StateMachine feature: $feature" -ForegroundColor Green
            $results.StateMachine_Usage.Features += $feature
            $hasStateMachineFeatures = $true
        }
    }
    
    # Check StateMachine embedding
    if ($dashboardCode -match "StateMachine.*struct") {
        Write-Host "  StateMachine properly embedded" -ForegroundColor Green
    } else {
        Write-Host "  StateMachine embedding issue" -ForegroundColor Red
    }
    
    # Check state transitions
    $stateTransitions = @(
        "SetState.*StateLoading",
        "SetState.*StateSuccess",
        "SetState.*StateError",
        "SetState.*StateEmpty",
        "SetState.*StatePartial"
    )
    
    $transitionsComplete = $true
    foreach ($transition in $stateTransitions) {
        if ($dashboardCode -match $transition) {
            Write-Host "  State transition found" -ForegroundColor Green
        } else {
            Write-Host "  State transition missing: $transition" -ForegroundColor Red
            $transitionsComplete = $false
        }
    }
    
    # Check state management in methods
    $stateManagementChecks = @(
        "Current.*State",
        "state.*machine",
        "state.*transition"
    )
    
    $stateManagementComplete = $true
    foreach ($check in $stateManagementChecks) {
        if ($dashboardCode -match $check) {
            Write-Host "  State management: $check present" -ForegroundColor Green
        } else {
            Write-Host "  State management: $check missing" -ForegroundColor Red
            $stateManagementComplete = $false
        }
    }
    
    # Check state isolation
    if ($dashboardCode -match "StateMachine.*{" -and $dashboardCode -match "Current.*int") {
        Write-Host "  State isolation properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  State isolation issues found" -ForegroundColor Red
    }
    
    $results.StateMachine_Usage.Pass = $hasStateMachineFeatures -and $transitionsComplete -and $stateManagementComplete
    $results.StateMachine_Usage.Status = if ($results.StateMachine_Usage.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.StateMachine_Usage.Pass) {
        Write-Host "  StateMachine properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  StateMachine issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing StateMachine usage: $($_.Exception.Message)" -ForegroundColor Red
    $results.StateMachine_Usage.Status = "ERROR"
    $results.StateMachine_Usage.Notes += "Test execution failed"
}

# Architecture Summary
Write-Host ""
Write-Host "Architecture Review Summary" -ForegroundColor Yellow
Write-Host "============================" -ForegroundColor Yellow

$passedTests = 0
$totalTests = 2

foreach ($test in $results.Keys) {
    $result = $results[$test]
    $statusEmoji = if ($result.Status -eq "PASS") { "PASS" } elseif ($result.Status -eq "FAIL") { "FAIL" } else { "ERROR" }
    
    Write-Host "$statusEmoji $test - $($result.Status)" -ForegroundColor $(if ($result.Status -eq "PASS") { "Green" } elseif ($result.Status -eq "FAIL") { "Red" } else { "Yellow" })
    
    if ($result.Features.Count -gt 0) {
        Write-Host "   Features: $($result.Features -join ', ')" -ForegroundColor Gray
    }
    
    if ($result.Notes.Count -gt 0) {
        foreach ($note in $result.Notes) {
            Write-Host "   Note: $note" -ForegroundColor Yellow
        }
    }
    
    if ($result.Status -eq "PASS") {
        $passedTests++
    }
}

# Overall Result
Write-Host ""
Write-Host "Overall Architecture Result: $passedTests/$totalTests tests passed" -ForegroundColor $(if ($passedTests -eq $totalTests) { "Green" } else { "Red" })

if ($passedTests -eq $totalTests) {
    Write-Host "All architecture requirements met!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some architecture requirements need attention" -ForegroundColor Red
    exit 1
}
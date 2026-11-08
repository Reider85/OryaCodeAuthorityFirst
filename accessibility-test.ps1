# accessibility-test.ps1 — Manual Accessibility Testing Script for UI 15-Item Checklist
# Tests INV-UI-009, INV-UI-010, Focus visibility requirements

param(
    [string]$ScreenType = "dashboard",
    [switch]$Verbose
)

# Set up test environment
$ErrorActionPreference = "Stop"
$VerbosePreference = if ($Verbose) { "Continue" } else { "SilentlyContinue" }

Write-Host "Starting UI Accessibility Testing" -ForegroundColor Green
Write-Host "Screen: $ScreenType" -ForegroundColor Yellow

# Test Results
$results = @{
    INV_UI_009_KeyboardOnly = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    INV_UI_010_HelpOverlay = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Focus_Visibility = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Symbol_Text_Pairing = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
}

# Test 1: INV-UI-009 100% Keyboard-Only Operation
Write-Host ""
Write-Host "Testing INV-UI-009: 100% Keyboard-Only Operation" -ForegroundColor Cyan

try {
    # Analyze keyboard navigation implementation
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $keysCode = Get-Content "pkg\ui\screens\key.go" -Raw
    $a11yCode = Get-Content "pkg\ui\a11y.go" -Raw
    
    # Check for keyboard navigation features
    $keyboardFeatures = @(
        "tea\.KeyTab",
        "tea\.KeyShiftTab",
        "tea\.KeyUp",
        "tea\.KeyDown", 
        "tea\.KeyEnter",
        "tea\.KeyEsc",
        "tea\.KeyCtrl+C",
        "CurrentWidget",
        "WidgetIDs",
        "navigation"
    )
    
    $hasKeyboardFeatures = $false
    foreach ($feature in $keyboardFeatures) {
        if ($dashboardCode -match $feature -or $keysCode -match $feature) {
            Write-Host "  Found keyboard feature: $feature" -ForegroundColor Green
            $results.INV_UI_009_KeyboardOnly.Features += $feature
            $hasKeyboardFeatures = $true
        }
    }
    
    # Check for mouse-only operations
    $mouseOperations = @(
        "mouse",
        "click",
        "double click",
        "right click",
        "drag",
        "drop"
    )
    
    $hasMouseOperations = $false
    foreach ($op in $mouseOperations) {
        if ($dashboardCode -match $op) {
            Write-Host "  WARNING: Potential mouse operation found: $op" -ForegroundColor Yellow
            $results.INV_UI_009_KeyboardOnly.Notes += "Potential mouse dependency: $op"
            $hasMouseOperations = $true
        }
    }
    
    # Check for keyboard shortcuts documentation
    if ($dashboardCode -match "helpSections") {
        Write-Host "  Help sections for keyboard shortcuts implemented" -ForegroundColor Green
    }
    
    # Check keyboard navigation completeness
    $navMethods = @(
        "Tab navigation",
        "Shift+Tab navigation", 
        "Arrow key navigation",
        "Enter activation",
        "Escape cancel"
    )
    
    $completeNavigation = $true
    foreach ($method in $navMethods) {
        if ($dashboardCode -match "(?i)$method") {
            Write-Host "  $method: Implemented" -ForegroundColor Green
        } else {
            Write-Host "  $method: Missing" -ForegroundColor Red
            $completeNavigation = $false
        }
    }
    
    $results.INV_UI_009_KeyboardOnly.Pass = $hasKeyboardFeatures -and !$hasMouseOperations -and $completeNavigation
    $results.INV_UI_009_KeyboardOnly.Status = if ($results.INV_UI_009_KeyboardOnly.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.INV_UI_009_KeyboardOnly.Pass) {
        Write-Host "  Full keyboard-only operation implemented" -ForegroundColor Green
    } else {
        Write-Host "  Keyboard-only operation incomplete" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing INV-UI-009: $($_.Exception.Message)" -ForegroundColor Red
    $results.INV_UI_009_KeyboardOnly.Status = "ERROR"
    $results.INV_UI_009_KeyboardOnly.Notes += "Test execution failed"
}

# Test 2: INV-UI-010 Help Overlay Functionality
Write-Host ""
Write-Host "Testing INV-UI-010: Help Overlay Functionality" -ForegroundColor Cyan

try {
    # Analyze help overlay implementation
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $helpCode = Get-Content "pkg\ui\components\helpoverlay.go" -Raw
    
    # Check help overlay features
    $helpFeatures = @(
        "ShowHelp",
        "helpSections",
        "NewHelpOverlay",
        "Keyboard Shortcuts",
        "toggle help",
        "question mark"
    )
    
    $hasHelpFeatures = $false
    foreach ($feature in $helpFeatures) {
        if ($dashboardCode -match $feature -or $helpCode -match $feature) {
            Write-Host "  Found help feature: $feature" -ForegroundColor Green
            $results.INV_UI_010_HelpOverlay.Features += $feature
            $hasHelpFeatures = $true
        }
    }
    
    # Check help activation
    if ($dashboardCode -match "tea\.KeyQuestionMark") {
        Write-Host "  Help activation with ? key implemented" -ForegroundColor Green
    }
    
    # Check help overlay accessibility
    if ($helpCode -match "keyboard.*accessible" -or $helpCode -match "a11y") {
        Write-Host "  Help overlay is keyboard accessible" -ForegroundColor Green
    }
    
    # Check help content completeness
    $helpContent = @(
        "navigation keys",
        "action keys", 
        "global shortcuts",
        "screen-specific help"
    )
    
    $completeHelp = $true
    foreach ($content in $helpContent) {
        if ($helpCode -match "(?i)$content") {
            Write-Host "  Help content: $content present" -ForegroundColor Green
        } else {
            Write-Host "  Help content: $content missing" -ForegroundColor Red
            $completeHelp = $false
        }
    }
    
    $results.INV_UI_010_HelpOverlay.Pass = $hasHelpFeatures -and $completeHelp
    $results.INV_UI_010_HelpOverlay.Status = if ($results.INV_UI_010_HelpOverlay.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.INV_UI_010_HelpOverlay.Pass) {
        Write-Host "  Comprehensive help overlay implemented" -ForegroundColor Green
    } else {
        Write-Host "  Help overlay incomplete" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing INV-UI-010: $($_.Exception.Message)" -ForegroundColor Red
    $results.INV_UI_010_HelpOverlay.Status = "ERROR"
    $results.INV_UI_010_HelpOverlay.Notes += "Test execution failed"
}

# Test 3: Focus Visibility
Write-Host ""
Write-Host "Testing Focus Visibility" -ForegroundColor Cyan

try {
    # Analyze focus visibility implementation
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $focusCode = Get-Content "pkg\ui\components\focusring.go" -Raw
    $a11yCode = Get-Content "pkg\ui\a11y.go" -Raw
    
    # Check focus indicators
    $focusFeatures = @(
        "CurrentWidget",
        "focusPrefix",
        "focus marker",
        "focus ring",
        "focus visibility"
    )
    
    $hasFocusFeatures = $false
    foreach ($feature in $focusFeatures) {
        if ($dashboardCode -match $feature -or $focusCode -match $feature) {
            Write-Host "  Found focus feature: $feature" -ForegroundColor Green
            $results.Focus_Visibility.Features += $feature
            $hasFocusFeatures = $true
        }
    }
    
    # Check focus marker implementation
    if ($dashboardCode -match "focusPrefix.*▸") {
        Write-Host "  Focus marker (▸) implemented" -ForegroundColor Green
    }
    
    # Check focus management
    if ($dashboardCode -match "CurrentWidget.*=") {
        Write-Host "  Focus management implemented" -ForegroundColor Green
    }
    
    # Check focus visibility across states
    $focusStates = @(
        "StateLoading",
        "StateSuccess", 
        "StateError",
        "StateEmpty",
        "StatePartial"
    )
    
    $focusInAllStates = $true
    foreach ($state in $focusStates) {
        if ($dashboardCode -match $state) {
            Write-Host "  Focus handling in $state: Present" -ForegroundColor Green
        } else {
            Write-Host "  Focus handling in $state: Missing" -ForegroundColor Red
            $focusInAllStates = $false
        }
    }
    
    $results.Focus_Visibility.Pass = $hasFocusFeatures -and $focusInAllStates
    $results.Focus_Visibility.Status = if ($results.Focus_Visibility.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Focus_Visibility.Pass) {
        Write-Host "  Focus visibility properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  Focus visibility issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing focus visibility: $($_.Exception.Message)" -ForegroundColor Red
    $results.Focus_Visibility.Status = "ERROR"
    $results.Focus_Visibility.Notes += "Test execution failed"
}

# Test 4: Symbol-Text Pairing
Write-Host ""
Write-Host "Testing Symbol-Text Pairing" -ForegroundColor Cyan

try {
    # Analyze symbol-text pairing implementation
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $a11yCode = Get-Content "pkg\ui\a11y.go" -Raw
    
    # Check for symbols with text
    $symbols = @(
        "✓",
        "⚠", 
        "✗",
        "🟢",
        "🔴",
        "🟡"
    )
    
    $hasSymbols = $false
    foreach ($symbol in $symbols) {
        if ($dashboardCode -match [regex]::Escape($symbol)) {
            Write-Host "  Found symbol: $symbol" -ForegroundColor Green
            $hasSymbols = $true
        }
    }
    
    # Check symbol-text pairing
    $symbolTextPairs = @(
        "✓.*success",
        "⚠.*warning", 
        "✗.*error",
        "🟢.*active",
        "🔴.*error",
        "🟡.*partial"
    )
    
    $hasSymbolText = $false
    foreach ($pair in $symbolTextPairs) {
        if ($dashboardCode -match $pair) {
            Write-Host "  Symbol-text pairing found" -ForegroundColor Green
            $hasSymbolText = $true
        }
    }
    
    # Check color-blind safety
    if ($a11yCode -match "color.*blind" -or $a11yCode -match "contrast") {
        Write-Host "  Color-blind safety features implemented" -ForegroundColor Green
    }
    
    # Check for symbols without text
    $lines = $dashboardCode -split "`n"
    $symbolsWithoutText = @()
    
    foreach ($line in $lines) {
        if ($line -match "[✓⚠✗🟢🔴🟡]") {
            if ($line -notmatch "(success|warning|error|active|partial|loading)" -and $line -notmatch "[a-zA-Z]") {
                $symbolsWithoutText += $line.Trim()
            }
        }
    }
    
    if ($symbolsWithoutText.Count -gt 0) {
        Write-Host "  WARNING: Symbols without text found:" -ForegroundColor Yellow
        foreach ($symbolLine in $symbolsWithoutText) {
            Write-Host "    $symbolLine" -ForegroundColor Yellow
            $results.Symbol_Text_Pairing.Notes += "Symbol without text: $symbolLine"
        }
    } else {
        Write-Host "  All symbols have text pairing" -ForegroundColor Green
    }
    
    $results.Symbol_Text_Pairing.Pass = $hasSymbols -and $hasSymbolText -and ($symbolsWithoutText.Count -eq 0)
    $results.Symbol_Text_Pairing.Status = if ($results.Symbol_Text_Pairing.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Symbol_Text_Pairing.Pass) {
        Write-Host "  Symbol-text pairing properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  Symbol-text pairing issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing symbol-text pairing: $($_.Exception.Message)" -ForegroundColor Red
    $results.Symbol_Text_Pairing.Status = "ERROR"
    $results.Symbol_Text_Pairing.Notes += "Test execution failed"
}

# Accessibility Summary
Write-Host ""
Write-Host "Accessibility Testing Summary" -ForegroundColor Yellow
Write-Host "=============================" -ForegroundColor Yellow

$passedTests = 0
$totalTests = 4

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
Write-Host "Overall Accessibility Result: $passedTests/$totalTests tests passed" -ForegroundColor $(if ($passedTests -eq $totalTests) { "Green" } else { "Red" })

if ($passedTests -eq $totalTests) {
    Write-Host "All accessibility requirements met!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some accessibility requirements need attention" -ForegroundColor Red
    exit 1
}
# error-handling-test.ps1 — Error Handling & UX Testing Script for UI 15-Item Checklist
# Tests consistent error states, universal keyboard shortcuts, and help overlay

param(
    [string]$ScreenType = "dashboard"
)

Write-Host "Starting UI Error Handling & UX Testing" -ForegroundColor Green
Write-Host "Screen: $ScreenType" -ForegroundColor Yellow

# Test Results
$results = @{
    Consistent_Error_States = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Universal_Keyboard_Shortcuts = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Help_Overlay_Documentation = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
}

# Test 1: Consistent Error States
Write-Host ""
Write-Host "Testing Consistent Error States" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $errorBoxCode = Get-Content "pkg\ui\components\errorbox.go" -Raw
    
    # Check for error state handling
    $errorFeatures = @(
        "StateError",
        "renderError",
        "LoadError",
        "LoadError",
        "error handling",
        "error recovery"
    )
    
    $hasErrorFeatures = $false
    foreach ($feature in $errorFeatures) {
        if ($dashboardCode -match $feature -or $errorBoxCode -match $feature) {
            Write-Host "  Found error handling feature: $feature" -ForegroundColor Green
            $results.Consistent_Error_States.Features += $feature
            $hasErrorFeatures = $true
        }
    }
    
    # Check error state transitions
    $errorTransitions = @(
        "SetState.*StateError",
        "StateMachine.*StateError",
        "handleLoadError",
        "handleDataLoaded"
    )
    
    $transitionsComplete = $true
    foreach ($transition in $errorTransitions) {
        if ($dashboardCode -match $transition) {
            Write-Host "  Error transition found: $transition" -ForegroundColor Green
        } else {
            Write-Host "  Error transition missing: $transition" -ForegroundColor Red
            $transitionsComplete = $false
        }
    }
    
    # Check error message consistency
    $errorMessages = @(
        "Error:",
        "Failed to",
        "Unable to",
        "Could not"
    )
    
    $messagesConsistent = $true
    foreach ($message in $errorMessages) {
        if ($dashboardCode -match $message -or $errorBoxCode -match $message) {
            Write-Host "  Error message pattern: $message found" -ForegroundColor Green
        } else {
            Write-Host "  Error message pattern: $message missing" -ForegroundColor Yellow
        }
    }
    
    # Check error recovery mechanisms
    $recoveryMechanisms = @(
        "retry",
        "refresh",
        "reload",
        "fallback"
    )
    
    $recoveryComplete = $true
    foreach ($mechanism in $recoveryMechanisms) {
        if ($dashboardCode -match $mechanism) {
            Write-Host "  Recovery mechanism: $mechanism found" -ForegroundColor Green
        } else {
            Write-Host "  Recovery mechanism: $mechanism missing" -ForegroundColor Yellow
            $recoveryComplete = $false
        }
    }
    
    # Check error state rendering
    if ($dashboardCode -match "renderError") {
        Write-Host "  Error state rendering implemented" -ForegroundColor Green
    } else {
        Write-Host "  Error state rendering missing" -ForegroundColor Red
        $transitionsComplete = $false
    }
    
    $results.Consistent_Error_States.Pass = $hasErrorFeatures -and $transitionsComplete -and $recoveryComplete
    $results.Consistent_Error_States.Status = if ($results.Consistent_Error_States.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Consistent_Error_States.Pass) {
        Write-Host "  Consistent error states implemented" -ForegroundColor Green
    } else {
        Write-Host "  Error state issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing consistent error states: $($_.Exception.Message)" -ForegroundColor Red
    $results.Consistent_Error_States.Status = "ERROR"
    $results.Consistent_Error_States.Notes += "Test execution failed"
}

# Test 2: Universal Keyboard Shortcuts
Write-Host ""
Write-Host "Testing Universal Keyboard Shortcuts" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $keysCode = Get-Content "pkg\ui\screens\key.go" -Raw
    
    # Check for universal keyboard shortcuts
    $shortcutFeatures = @(
        "global shortcuts",
        "universal keys",
        "common keys",
        "standard navigation"
    )
    
    $hasShortcutFeatures = $false
    foreach ($feature in $shortcutFeatures) {
        if ($dashboardCode -match $feature -or $keysCode -match $feature) {
            Write-Host "  Found shortcut feature: $feature" -ForegroundColor Green
            $results.Universal_Keyboard_Shortcuts.Features += $feature
            $hasShortcutFeatures = $true
        }
    }
    
    # Check for specific universal shortcuts
    $universalShortcuts = @(
        "tea\.KeyQuestionMark",   # Help
        "tea\.KeyCtrl+C",        # Quit
        "tea\.KeyEsc",           # Back/Cancel
        "tea\.KeyTab",           # Navigate forward
        "tea\.KeyShiftTab",      # Navigate backward
        "tea\.KeyEnter",         # Activate/Select
        "tea\.KeyR",             # Refresh
        "tea\.KeyUp",            # Move up
        "tea\.KeyDown"           # Move down
    )
    
    $shortcutsComplete = $true
    foreach ($shortcut in $universalShortcuts) {
        if ($dashboardCode -match $shortcut -or $keysCode -match $shortcut) {
            Write-Host "  Universal shortcut found: $shortcut" -ForegroundColor Green
        } else {
            Write-Host "  Universal shortcut missing: $shortcut" -ForegroundColor Red
            $shortcutsComplete = $false
        }
    }
    
    # Check for shortcut conflicts
    $shortcutConflicts = @(
        "duplicate shortcut",
        "conflicting key",
        "multiple handlers"
    )
    
    $conflictsFound = $false
    foreach ($conflict in $shortcutConflicts) {
        if ($dashboardCode -match $conflict -or $keysCode -match $conflict) {
            Write-Host "  WARNING: Potential shortcut conflict: $conflict" -ForegroundColor Yellow
            $results.Universal_Keyboard_Shortcuts.Notes += "Potential conflict: $conflict"
            $conflictsFound = $true
        }
    }
    
    # Check shortcut discoverability
    if ($dashboardCode -match "helpSections" -or $dashboardCode -match "Keyboard Shortcuts") {
        Write-Host "  Shortcut discoverability implemented" -ForegroundColor Green
    } else {
        Write-Host "  Shortcut discoverability missing" -ForegroundColor Yellow
    }
    
    # Check shortcut consistency across screens
    $screens = @("dashboard.go", "adapters.go", "onboarding.go")
    $consistencyComplete = $true
    
    foreach ($screen in $screens) {
        $screenPath = "pkg\ui\screens\$screen"
        if (Test-Path $screenPath) {
            $screenContent = Get-Content $screenPath -Raw
            foreach ($shortcut in $universalShortcuts) {
                if ($screenContent -match $shortcut) {
                    Write-Host "  Shortcut $shortcut consistent in $screen" -ForegroundColor Green
                } else {
                    Write-Host "  Shortcut $shortcut missing in $screen" -ForegroundColor Yellow
                    $consistencyComplete = $false
                }
            }
        }
    }
    
    $results.Universal_Keyboard_Shortcuts.Pass = $hasShortcutFeatures -and $shortcutsComplete -and $consistencyComplete
    $results.Universal_Keyboard_Shortcuts.Status = if ($results.Universal_Keyboard_Shortcuts.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Universal_Keyboard_Shortcuts.Pass) {
        Write-Host "  Universal keyboard shortcuts properly implemented" -ForegroundColor Green
    } else {
        Write-Host "  Keyboard shortcut issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing universal keyboard shortcuts: $($_.Exception.Message)" -ForegroundColor Red
    $results.Universal_Keyboard_Shortcuts.Status = "ERROR"
    $results.Universal_Keyboard_Shortcuts.Notes += "Test execution failed"
}

# Test 3: Help Overlay and Documentation
Write-Host ""
Write-Host "Testing Help Overlay and Documentation" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $helpCode = Get-Content "pkg\ui\components\helpoverlay.go" -Raw
    
    # Check help overlay features
    $helpFeatures = @(
        "ShowHelp",
        "helpSections",
        "NewHelpOverlay",
        "toggle help",
        "help content",
        "help navigation"
    )
    
    $hasHelpFeatures = $false
    foreach ($feature in $helpFeatures) {
        if ($dashboardCode -match $feature -or $helpCode -match $feature) {
            Write-Host "  Found help feature: $feature" -ForegroundColor Green
            $results.Help_Overlay_Documentation.Features += $feature
            $hasHelpFeatures = $true
        }
    }
    
    # Check help overlay functionality
    $helpFunctionality = @(
        "help activation",
        "help display",
        "help navigation",
        "help closing",
        "help content"
    )
    
    $functionalityComplete = $true
    foreach ($function in $helpFunctionality) {
        if ($helpCode -match $function) {
            Write-Host "  Help functionality: $function present" -ForegroundColor Green
        } else {
            Write-Host "  Help functionality: $function missing" -ForegroundColor Red
            $functionalityComplete = $false
        }
    }
    
    # Check help content completeness
    $helpContent = @(
        "keyboard shortcuts",
        "navigation keys",
        "action keys",
        "global commands",
        "screen-specific help"
    )
    
    $contentComplete = $true
    foreach ($content in $helpContent) {
        if ($helpCode -match "(?i)$content") {
            Write-Host "  Help content: $content present" -ForegroundColor Green
        } else {
            Write-Host "  Help content: $content missing" -ForegroundColor Red
            $contentComplete = $false
        }
    }
    
    # Check help documentation
    $docFiles = Get-ChildItem "docs\*.md" -File
    $docsComplete = $true
    
    foreach ($doc in $docFiles) {
        $content = Get-Content $doc -Raw
        if ($content -match "(?i)help" -and $content -match "(?i)keyboard") {
            Write-Host "  Documentation: $($doc.Name) contains help and keyboard info" -ForegroundColor Green
        } else {
            Write-Host "  Documentation: $($doc.Name) missing help or keyboard info" -ForegroundColor Yellow
            $docsComplete = $false
        }
    }
    
    # Check help accessibility
    if ($helpCode -match "accessible" -or $helpCode -match "a11y") {
        Write-Host "  Help accessibility features implemented" -ForegroundColor Green
    } else {
        Write-Host "  Help accessibility features missing" -ForegroundColor Yellow
    }
    
    $results.Help_Overlay_Documentation.Pass = $hasHelpFeatures -and $functionalityComplete -and $contentComplete -and $docsComplete
    $results.Help_Overlay_Documentation.Status = if ($results.Help_Overlay_Documentation.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Help_Overlay_Documentation.Pass) {
        Write-Host "  Comprehensive help overlay and documentation implemented" -ForegroundColor Green
    } else {
        Write-Host "  Help overlay issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing help overlay and documentation: $($_.Exception.Message)" -ForegroundColor Red
    $results.Help_Overlay_Documentation.Status = "ERROR"
    $results.Help_Overlay_Documentation.Notes += "Test execution failed"
}

# Error Handling & UX Summary
Write-Host ""
Write-Host "Error Handling & UX Testing Summary" -ForegroundColor Yellow
Write-Host "====================================" -ForegroundColor Yellow

$passedTests = 0
$totalTests = 3

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
Write-Host "Overall Error Handling & UX Result: $passedTests/$totalTests tests passed" -ForegroundColor $(if ($passedTests -eq $totalTests) { "Green" } else { "Red" })

if ($passedTests -eq $totalTests) {
    Write-Host "All error handling & UX requirements met!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some error handling & UX requirements need attention" -ForegroundColor Red
    exit 1
}
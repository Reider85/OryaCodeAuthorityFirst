# ui-self-review.ps1 - Quick UI self-review automation script for Windows
# This script helps developers quickly verify their UI implementations
# against the CodeAuthorityFirst UI self-review checklist

param(
    [string]$Path = ".",
    [switch]$Verbose = $false
)

# Set error preference
$ErrorActionPreference = "Stop"

# Function to print colored output
function Write-Status {
    param(
        [string]$Status,
        [string]$Message
    )
    
    switch ($Status) {
        "PASS" { Write-Host "$($PSStyle.Foreground.Green)✓ $Message$($PSStyle.Reset)" }
        "FAIL" { Write-Host "$($PSStyle.Foreground.Red)✗ $Message$($PSStyle.Reset)" }
        "WARN" { Write-Host "$($PSStyle.Foreground.Yellow)⚠ $Message$($PSStyle.Reset)" }
        "INFO" { Write-Host "$($PSStyle.Foreground.Blue)ℹ $Message$($PSStyle.Reset)" }
    }
}

# Function to check if file exists
function Test-FileExists {
    param(
        [string]$File,
        [string]$Description
    )
    
    if (Test-Path $File -PathType Leaf) {
        Write-Status "PASS" "$Description: $File"
        return $true
    } else {
        Write-Status "FAIL" "$Description: $File"
        return $false
    }
}

# Function to check if directory exists
function Test-DirectoryExists {
    param(
        [string]$Directory,
        [string]$Description
    )
    
    if (Test-Path $Directory -PathType Container) {
        Write-Status "PASS" "$Description: $Directory"
        return $true
    } else {
        Write-Status "FAIL" "$Description: $Directory"
        return $false
    }
}

# Function to check if command exists
function Test-CommandExists {
    param(
        [string]$Command,
        [string]$Description
    )
    
    if (Get-Command $Command -ErrorAction SilentlyContinue) {
        Write-Status "PASS" "$Description: $Command"
        return $true
    } else {
        Write-Status "FAIL" "$Description: $Command"
        return $false
    }
}

# Function to check Go build
function Test-GoBuild {
    param(
        [string]$Path,
        [string]$Name
    )
    
    try {
        $tempFile = [System.IO.Path]::GetTempFileName()
        $outputPath = $tempFile -replace '\.tmp$', '.exe'
        
        go build -o $outputPath $Path | Out-Null
        
        if ($LASTEXITCODE -eq 0) {
            Write-Status "PASS" "Go build successful: $Name"
            Remove-Item $outputPath -ErrorAction SilentlyContinue
            return $true
        } else {
            Write-Status "FAIL" "Go build failed: $Name"
            return $false
        }
    } catch {
        Write-Status "FAIL" "Go build failed: $Name - $_"
        return $false
    }
}

# Function to check performance
function Test-Performance {
    param(
        [string]$Path,
        [string]$Name,
        [int]$MaxTimeMs
    )
    
    try {
        $tempFile = [System.IO.Path]::GetTempFileName()
        $outputPath = $tempFile -replace '\.tmp$', '.exe'
        
        go build -o $outputPath $Path | Out-Null
        
        if ($LASTEXITCODE -eq 0) {
            $startTime = Get-Date
            & $outputPath | Out-Null
            $endTime = Get-Date
            
            $durationMs = [math]::Round(($endTime - $startTime).TotalMilliseconds * 1000) / 1000
            
            if ($durationMs -lt $MaxTimeMs) {
                Write-Status "PASS" "Performance $Name $($durationMs)ms < $($MaxTimeMs)ms"
                Remove-Item $outputPath -ErrorAction SilentlyContinue
                return $true
            } else {
                Write-Status "FAIL" "Performance $Name $($durationMs)ms >= $($MaxTimeMs)ms"
                Remove-Item $outputPath -ErrorAction SilentlyContinue
                return $false
            }
        } else {
            Write-Status "FAIL" "Performance test failed: build error"
            return $false
        }
    } catch {
        Write-Status "FAIL" "Performance test failed: $_"
        return $false
    }
}

# Function to check test coverage
function Test-TestCoverage {
    param(
        [string]$Path,
        [string]$Name
    )
    
    try {
        $testResult = go test -v $Path 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Status "PASS" "Tests passing: $Name"
            return $true
        } else {
            Write-Status "FAIL" "Tests failing: $Name"
            if ($Verbose) {
                Write-Host $testResult
            }
            return $false
        }
    } catch {
        Write-Status "FAIL" "Test check failed: $_"
        return $false
    }
}

# Function to check UI files
function Test-UIFiles {
    Write-Host "📁 Checking UI files..."
    
    $uiFiles = @(
        "docs/UI_SELF_REVIEW_CHECKLIST.md"
        "docs/UI_IMPLEMENTATION_EXAMPLES.md"
        "docs/UI_QUICK_REFERENCE.md"
        "pkg/ui/screens/ui_self_review_test.go"
    )
    
    $missingFiles = 0
    foreach ($file in $uiFiles) {
        $fullPath = Join-Path $Path $file
        if (-not (Test-FileExists $fullPath "UI documentation")) {
            $missingFiles++
        }
    }
    
    return ($missingFiles -eq 0)
}

# Function to check screen implementations
function Test-ScreenImplementations {
    Write-Host "🎬 Checking screen implementations..."
    
    $screens = @(
        "dashboard.go"
        "onboarding.go"
        "adapters.go"
    )
    
    $missingScreens = 0
    foreach ($screen in $screens) {
        $screenFile = Join-Path $Path "pkg\ui\screens\$screen"
        if (Test-Path $screenFile -PathType Leaf) {
            $content = Get-Content $screenFile -Raw
            
            if ($content -match "func.*Init\(\)" -and $content -match "func.*Update\(" -and $content -match "func.*View\(\)") {
                Write-Status "PASS" "Screen implementation: $screen"
            } else {
                Write-Status "FAIL" "Screen missing required methods: $screen"
                $missingScreens++
            }
        } else {
            Write-Status "WARN" "Screen file not found: $screen"
        }
    }
    
    return ($missingScreens -eq 0)
}

# Function to check performance requirements
function Test-PerformanceRequirements {
    Write-Host "⚡ Checking performance requirements..."
    
    $passed = 0
    $total = 0
    
    # Dashboard first paint <100ms
    $total++
    if (Test-Performance (Join-Path $Path "cmd\caf") "dashboard" 100) {
        $passed++
    }
    
    # Widget render <50ms
    $total++
    if (Test-Performance (Join-Path $Path "pkg\ui\screens\dashboard.go") "widget" 50) {
        $passed++
    }
    
    Write-Host "Performance: $passed/$total checks passed"
    return ($passed -eq $total)
}

# Function to check accessibility
function Test-Accessibility {
    Write-Host "♿ Checking accessibility requirements..."
    
    $uiFiles = @(
        "pkg\ui\screens\dashboard.go"
        "pkg\ui\screens\onboarding.go"
        "pkg\ui\screens\adapters.go"
    )
    
    $passed = 0
    $total = 0
    
    foreach ($file in $uiFiles) {
        $fullPath = Join-Path $Path $file
        if (Test-Path $fullPath -PathType Leaf) {
            $content = Get-Content $fullPath -Raw
            
            # Check for focus visibility
            if ($content -match "▸") {
                $passed++
                Write-Status "PASS" "Focus visibility: $file"
            } else {
                Write-Status "WARN" "Focus visibility missing: $file"
            }
            
            # Check for symbol-text pairing
            if ($content -match "✓.*[A-Za-z]" -or $content -match "⚠.*[A-Za-z]") {
                $passed++
                Write-Status "PASS" "Symbol-text pairing: $file"
            } else {
                Write-Status "WARN" "Symbol-text pairing missing: $file"
            }
        }
        $total += 2
    }
    
    Write-Host "Accessibility: $passed/$total checks passed"
    return ($passed -gt 0)
}

# Function to check keyboard shortcuts
function Test-KeyboardShortcuts {
    Write-Host "⌨️ Checking keyboard shortcuts..."
    
    $uiFiles = @(
        "pkg\ui\screens\dashboard.go"
        "pkg\ui\screens\onboarding.go"
        "pkg\ui\screens\adapters.go"
    )
    
    $passed = 0
    $total = 0
    
    foreach ($file in $uiFiles) {
        $fullPath = Join-Path $Path $file
        if (Test-Path $fullPath -PathType Leaf) {
            $content = Get-Content $fullPath -Raw
            
            # Check for help overlay
            if ($content -match "ShowHelp") {
                $passed++
                Write-Status "PASS" "Help overlay: $file"
            } else {
                Write-Status "WARN" "Help overlay missing: $file"
            }
            
            # Check for keyboard navigation
            if ($content -match "CurrentWidget" -or $content -match "CurrentItem") {
                $passed++
                Write-Status "PASS" "Keyboard navigation: $file"
            } else {
                Write-Status "WARN" "Keyboard navigation missing: $file"
            }
        }
        $total += 2
    }
    
    Write-Host "Keyboard shortcuts: $passed/$total checks passed"
    return ($passed -gt 0)
}

# Function to check documentation
function Test-Documentation {
    Write-Host "📚 Checking documentation..."
    
    $docFiles = @(
        "docs/UI_SELF_REVIEW_CHECKLIST.md"
        "docs/UI_IMPLEMENTATION_EXAMPLES.md"
        "docs/UI_QUICK_REFERENCE.md"
    )
    
    $missingSections = 0
    
    foreach ($doc in $docFiles) {
        $fullPath = Join-Path $Path $doc
        if (Test-Path $fullPath -PathType Leaf) {
            $content = Get-Content $fullPath -Raw
            
            if ($content -notmatch "Performance Requirements") {
                Write-Status "WARN" "Missing performance section: $doc"
                $missingSections++
            }
            
            if ($content -notmatch "Accessibility Compliance") {
                Write-Status "WARN" "Missing accessibility section: $doc"
                $missingSections++
            }
        }
    }
    
    return ($missingSections -eq 0)
}

# Function to check test coverage
function Test-TestCoverage {
    Write-Host "🧪 Checking test coverage..."
    
    $testFiles = @(
        "pkg\ui\screens\dashboard_test.go"
        "pkg\ui\screens\onboarding_test.go"
        "pkg\ui\screens\adapters_test.go"
        "pkg\ui\screens\ui_self_review_test.go"
    )
    
    $missingTests = 0
    
    foreach ($testFile in $testFiles) {
        $fullPath = Join-Path $Path $testFile
        if (Test-Path $fullPath -PathType Leaf) {
            $content = Get-Content $fullPath -Raw
            
            # Check for performance tests
            if ($content -notmatch "TestPerformance") {
                Write-Status "WARN" "Missing performance test: $testFile"
                $missingTests++
            }
            
            # Check for accessibility tests
            if ($content -notmatch "TestAccessibility") {
                Write-Status "WARN" "Missing accessibility test: $testFile"
                $missingTests++
            }
        } else {
            Write-Status "WARN" "Test file missing: $testFile"
            $missingTests++
        }
    }
    
    return ($missingTests -eq 0)
}

# Main function
function Main {
    Write-Host "🔍 CodeAuthorityFirst UI Self-Review"
    Write-Host "===================================="
    Write-Host
    
    # Set path to project root if not specified
    if ($Path -eq ".") {
        $Path = $PSScriptRoot
    }
    
    # Initialize counters
    $passedChecks = 0
    $totalChecks = 0
    
    # Run all checks
    Write-Host "📋 Running UI self-review checks..."
    Write-Host
    
    # 1. Check UI files
    $totalChecks++
    if (Test-UIFiles) {
        $passedChecks++
    }
    
    # 2. Check screen implementations
    $totalChecks++
    if (Test-ScreenImplementations) {
        $passedChecks++
    }
    
    # 3. Check performance requirements
    $totalChecks++
    if (Test-PerformanceRequirements) {
        $passedChecks++
    }
    
    # 4. Check accessibility
    $totalChecks++
    if (Test-Accessibility) {
        $passedChecks++
    }
    
    # 5. Check keyboard shortcuts
    $totalChecks++
    if (Test-KeyboardShortcuts) {
        $passedChecks++
    }
    
    # 6. Check documentation
    $totalChecks++
    if (Test-Documentation) {
        $passedChecks++
    }
    
    # 7. Check test coverage
    $totalChecks++
    if (Test-TestCoverage) {
        $passedChecks++
    }
    
    # Summary
    Write-Host
    Write-Host "📊 UI Self-Review Summary"
    Write-Host "========================"
    Write-Host "Passed: $passedChecks/$totalChecks checks"
    
    if ($passedChecks -eq $totalChecks) {
        Write-Host "$($PSStyle.Foreground.Green)🎉 All checks passed! UI implementation is ready for review.$($PSStyle.Reset)"
        exit 0
    } else {
        Write-Host "$($PSStyle.Foreground.Yellow)⚠ Some checks failed. Please address the issues above.$($PSStyle.Reset)"
        Write-Host "$($PSStyle.Foreground.Blue)💡 Use the detailed checklist for manual verification.$($PSStyle.Reset)"
        exit 1
    }
}

# Run main function
Main
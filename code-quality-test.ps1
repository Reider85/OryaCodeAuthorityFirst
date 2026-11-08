# code-quality-test.ps1 — Code Quality Assessment Script for UI 15-Item Checklist
# Tests Lipgloss optimization, file structure, testing coverage, and header comments

param(
    [string]$ScreenType = "dashboard"
)

Write-Host "Starting UI Code Quality Assessment" -ForegroundColor Green
Write-Host "Screen: $ScreenType" -ForegroundColor Yellow

# Test Results
$results = @{
    Lipgloss_Optimization = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    File_Structure = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Testing_Coverage = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
    Header_Comments = @{
        Status = "PENDING"
        Pass = $false
        Notes = @()
        Features = @()
    }
}

# Test 1: Lipgloss Pre-computed Styles
Write-Host ""
Write-Host "Testing Lipgloss Pre-computed Styles" -ForegroundColor Cyan

try {
    $dashboardCode = Get-Content "pkg\ui\screens\dashboard.go" -Raw
    $themeCode = Get-Content "pkg\ui\theme.go" -Raw
    
    # Check for Lipgloss optimization features
    $lipglossFeatures = @(
        "style.*Render",
        "pre-computed",
        "cached styles",
        "style pooling",
        "static styles"
    )
    
    $hasLipglossFeatures = $false
    foreach ($feature in $lipglossFeatures) {
        if ($dashboardCode -match $feature -or $themeCode -match $feature) {
            Write-Host "  Found Lipgloss optimization: $feature" -ForegroundColor Green
            $results.Lipgloss_Optimization.Features += $feature
            $hasLipglossFeatures = $true
        }
    }
    
    # Check for style creation patterns
    $stylePatterns = @(
        "theme\..*\.Render",
        "lipgloss\..*\.Style",
        "NewStyle",
        "CopyStyle"
    )
    
    $stylesOptimized = $true
    foreach ($pattern in $stylePatterns) {
        if ($themeCode -match $pattern) {
            Write-Host "  Style pattern found: $pattern" -ForegroundColor Green
        } else {
            Write-Host "  Style pattern missing: $pattern" -ForegroundColor Red
            $stylesOptimized = $false
        }
    }
    
    # Check for runtime style computation
    if ($dashboardCode -match "time\.Now.*Render" -or $dashboardCode -match "Render.*time") {
        Write-Host "  WARNING: Runtime style computation detected" -ForegroundColor Yellow
        $results.Lipgloss_Optimization.Notes += "Potential runtime style computation"
        $stylesOptimized = $false
    }
    
    # Check for style caching
    if ($themeCode -match "cache" -or $themeCode -match "pool") {
        Write-Host "  Style caching implemented" -ForegroundColor Green
    } else {
        Write-Host "  Style caching missing" -ForegroundColor Yellow
    }
    
    $results.Lipgloss_Optimization.Pass = $hasLipglossFeatures -and $stylesOptimized
    $results.Lipgloss_Optimization.Status = if ($results.Lipgloss_Optimization.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Lipgloss_Optimization.Pass) {
        Write-Host "  Lipgloss styles properly optimized" -ForegroundColor Green
    } else {
        Write-Host "  Lipgloss optimization issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing Lipgloss optimization: $($_.Exception.Message)" -ForegroundColor Red
    $results.Lipgloss_Optimization.Status = "ERROR"
    $results.Lipgloss_Optimization.Notes += "Test execution failed"
}

# Test 2: Required File Structure
Write-Host ""
Write-Host "Testing Required File Structure" -ForegroundColor Cyan

try {
    # Check for required UI files
    $requiredFiles = @(
        "pkg\ui\tui.go",
        "pkg\ui\states.go",
        "pkg\ui\screens\dashboard.go",
        "pkg\ui\screens\adapters.go",
        "pkg\ui\screens\onboarding.go",
        "pkg\ui\components\header.go",
        "pkg\ui\components\helpoverlay.go",
        "pkg\ui\components\focusring.go",
        "pkg\ui\theme.go",
        "pkg\ui\a11y.go"
    )
    
    $filesPresent = $true
    foreach ($file in $requiredFiles) {
        if (Test-Path $file) {
            Write-Host "  File present: $file" -ForegroundColor Green
        } else {
            Write-Host "  File missing: $file" -ForegroundColor Red
            $filesPresent = $false
            $results.File_Structure.Notes += "Missing file: $file"
        }
    }
    
    # Check for proper package declarations
    $screenFiles = Get-ChildItem "pkg\ui\screens\*.go" -File
    $packageCorrect = $true
    
    foreach ($file in $screenFiles) {
        $content = Get-Content $file -Raw
        if ($content -match "package screens") {
            Write-Host "  Package declaration correct: $($file.Name)" -ForegroundColor Green
        } else {
            Write-Host "  Package declaration incorrect: $($file.Name)" -ForegroundColor Red
            $packageCorrect = $false
        }
    }
    
    # Check for logical organization
    $directories = @(
        "pkg\ui\screens",
        "pkg\ui\components",
        "pkg\ui"
    )
    
    $organizationCorrect = $true
    foreach ($dir in $directories) {
        if (Test-Path $dir) {
            Write-Host "  Directory present: $dir" -ForegroundColor Green
        } else {
            Write-Host "  Directory missing: $dir" -ForegroundColor Red
            $organizationCorrect = $false
        }
    }
    
    $results.File_Structure.Pass = $filesPresent -and $packageCorrect -and $organizationCorrect
    $results.File_Structure.Status = if ($results.File_Structure.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.File_Structure.Pass) {
        Write-Host "  File structure properly organized" -ForegroundColor Green
    } else {
        Write-Host "  File structure issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing file structure: $($_.Exception.Message)" -ForegroundColor Red
    $results.File_Structure.Status = "ERROR"
    $results.File_Structure.Notes += "Test execution failed"
}

# Test 3: Unit Test Coverage
Write-Host ""
Write-Host "Testing Unit Test Coverage" -ForegroundColor Cyan

try {
    # Check for test files
    $testFiles = Get-ChildItem "pkg\ui\screens\*_test.go" -File
    $testCount = $testFiles.Count
    
    Write-Host "  Found $testCount test files" -ForegroundColor Green
    
    # Check test coverage areas
    $coverageAreas = @(
        "TestView",
        "TestUpdate", 
        "TestInit",
        "TestPerformance",
        "TestAccessibility",
        "TestStateMachine"
    )
    
    $coverageComplete = $true
    foreach ($area in $coverageAreas) {
        $foundCoverage = $false
        foreach ($testFile in $testFiles) {
            $content = Get-Content $testFile -Raw
            if ($content -match $area) {
                Write-Host "  Test coverage: $area found in $($testFile.Name)" -ForegroundColor Green
                $foundCoverage = $true
                break
            }
        }
        
        if (!$foundCoverage) {
            Write-Host "  Test coverage: $area missing" -ForegroundColor Yellow
            $coverageComplete = $false
        }
    }
    
    # Check for edge case testing
    $edgeCases = @(
        "error handling",
        "empty state",
        "loading state",
        "partial state"
    )
    
    $edgeCaseCoverage = $true
    foreach ($edgeCase in $edgeCases) {
        $foundEdgeCase = $false
        foreach ($testFile in $testFiles) {
            $content = Get-Content $testFile -Raw
            if ($content -match "(?i)$edgeCase") {
                Write-Host "  Edge case testing: $edgeCase found in $($testFile.Name)" -ForegroundColor Green
                $foundEdgeCase = $true
                break
            }
        }
        
        if (!$foundEdgeCase) {
            Write-Host "  Edge case testing: $edgeCase missing" -ForegroundColor Yellow
            $edgeCaseCoverage = $false
        }
    }
    
    $results.Testing_Coverage.Pass = $coverageComplete -and $edgeCaseCoverage
    $results.Testing_Coverage.Status = if ($results.Testing_Coverage.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Testing_Coverage.Pass) {
        Write-Host "  Comprehensive test coverage implemented" -ForegroundColor Green
    } else {
        Write-Host "  Test coverage issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing test coverage: $($_.Exception.Message)" -ForegroundColor Red
    $results.Testing_Coverage.Status = "ERROR"
    $results.Testing_Coverage.Notes += "Test execution failed"
}

# Test 4: Header Comment Completeness
Write-Host ""
Write-Host "Testing Header Comment Completeness" -ForegroundColor Cyan

try {
    # Check for header comments in UI files
    $uiFiles = Get-ChildItem "pkg\ui\*.go" -File
    $screenFiles = Get-ChildItem "pkg\ui\screens\*.go" -File
    $componentFiles = Get-ChildItem "pkg\ui\components\*.go" -File
    
    $allFiles = $uiFiles + $screenFiles + $componentFiles
    $headerComplete = $true
    
    foreach ($file in $allFiles) {
        $content = Get-Content $file -Raw
        
        # Check for required header sections
        $requiredHeaders = @(
            "// Purpose:",
            "// Layout:",
            "// Performance budget:",
            "// Dependencies:",
            "// Notes:"
        )
        
        $fileHeadersComplete = $true
        foreach ($header in $requiredHeaders) {
            if ($content -match $header) {
                Write-Host "  Header present: $header in $($file.Name)" -ForegroundColor Green
            } else {
                Write-Host "  Header missing: $header in $($file.Name)" -ForegroundColor Yellow
                $fileHeadersComplete = $false
            }
        }
        
        if (!$fileHeadersComplete) {
            Write-Host "  Incomplete header in $($file.Name)" -ForegroundColor Yellow
            $headerComplete = $false
            $results.Header_Comments.Notes += "Incomplete header in $($file.Name)"
        }
    }
    
    # Check for documentation completeness
    $docFiles = Get-ChildItem "docs\*.md" -File
    $docsComplete = $true
    
    foreach ($doc in $docFiles) {
        $content = Get-Content $doc -Raw
        if ($content.Length -lt 100) {
            Write-Host "  Documentation incomplete: $($doc.Name)" -ForegroundColor Yellow
            $docsComplete = $false
            $results.Header_Comments.Notes += "Incomplete documentation: $($doc.Name)"
        } else {
            Write-Host "  Documentation present: $($doc.Name)" -ForegroundColor Green
        }
    }
    
    $results.Header_Comments.Pass = $headerComplete -and $docsComplete
    $results.Header_Comments.Status = if ($results.Header_Comments.Pass) { "PASS" } else { "FAIL" }
    
    if ($results.Header_Comments.Pass) {
        Write-Host "  Header comments complete" -ForegroundColor Green
    } else {
        Write-Host "  Header comment issues found" -ForegroundColor Red
    }
    
} catch {
    Write-Host "  Error testing header comments: $($_.Exception.Message)" -ForegroundColor Red
    $results.Header_Comments.Status = "ERROR"
    $results.Header_Comments.Notes += "Test execution failed"
}

# Code Quality Summary
Write-Host ""
Write-Host "Code Quality Assessment Summary" -ForegroundColor Yellow
Write-Host "==============================" -ForegroundColor Yellow

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
Write-Host "Overall Code Quality Result: $passedTests/$totalTests tests passed" -ForegroundColor $(if ($passedTests -eq $totalTests) { "Green" } else { "Red" })

if ($passedTests -eq $totalTests) {
    Write-Host "All code quality requirements met!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some code quality requirements need attention" -ForegroundColor Red
    exit 1
}
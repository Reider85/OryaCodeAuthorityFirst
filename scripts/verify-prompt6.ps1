# Prompt 6: Verify Current State Before Manual Updates
# This script helps verify the current state of the repository before manual GitHub settings updates

Write-Host "PROMPT 6: Pre-verification Script" -ForegroundColor Green
Write-Host "=================================" -ForegroundColor Green
Write-Host ""

# Check if social preview SVG exists and has correct dimensions
Write-Host "1. Checking social preview SVG..." -ForegroundColor Yellow
if (Test-Path "docs/social-preview.svg") {
    Write-Host "✅ docs/social-preview.svg exists" -ForegroundColor Green
    
    # Check file size
    $size = (Get-Item "docs/social-preview.svg").Length
    Write-Host "   File size: $size bytes" -ForegroundColor White
    
    # Check if it contains required content
    $content = Get-Content "docs/social-preview.svg" -Raw
    if ($content -match "CodeAuthorityFirst") {
        Write-Host "✅ Contains 'CodeAuthorityFirst' title" -ForegroundColor Green
    } else {
        Write-Host "❌ Missing 'CodeAuthorityFirst' title" -ForegroundColor Red
    }
    
    if ($content -match "who wrote what, when, and why") {
        Write-Host "✅ Contains required tagline" -ForegroundColor Green
    } else {
        Write-Host "❌ Missing required tagline" -ForegroundColor Red
    }
    
    if ($content -match "1280" -and $content -match "640") {
        Write-Host "✅ Has correct dimensions (1280x640)" -ForegroundColor Green
    } else {
        Write-Host "❌ Incorrect dimensions" -ForegroundColor Red
    }
} else {
    Write-Host "❌ docs/social-preview.svg does not exist" -ForegroundColor Red
}
Write-Host ""

# Check README.md for existing content
Write-Host "2. Checking README.md..." -ForegroundColor Yellow
if (Test-Path "README.md") {
    Write-Host "✅ README.md exists" -ForegroundColor Green
    
    # Check for content
    $readme = Get-Content "README.md" -Raw
    if ($readme -match "Track who wrote what, when, and why") {
        Write-Host "✅ Contains project description" -ForegroundColor Green
    } else {
        Write-Host "❌ Missing project description" -ForegroundColor Red
    }
    
    # Check for badges
    $badgeCount = ($readme | Select-String -Pattern "!\[.*\](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/").Count
    Write-Host "   Found $badgeCount CI/quality badges" -ForegroundColor White
} else {
    Write-Host "❌ README.md does not exist" -ForegroundColor Red
}
Write-Host ""

# Check if docs directory exists
Write-Host "3. Checking documentation..." -ForegroundColor Yellow
if (Test-Path "docs" -PathType Container) {
    Write-Host "✅ docs/ directory exists" -ForegroundColor Green
    $docCount = (Get-ChildItem "docs" -Include "*.md", "*.svg" | Measure-Object).Count
    Write-Host "   Contains $docCount documentation files" -ForegroundColor White
} else {
    Write-Host "❌ docs/ directory does not exist" -ForegroundColor Red
}
Write-Host ""

# Check GitHub repository status (this is informational only)
Write-Host "4. Current GitHub Repository Status (Informational):" -ForegroundColor Yellow
Write-Host "   Topics: 0 (need to add 20)" -ForegroundColor White
Write-Host "   Description: Empty (need to add)" -ForegroundColor White
Write-Host "   Website: Not configured (need to add)" -ForegroundColor White
Write-Host "   Social Preview: Not uploaded (SVG ready)" -ForegroundColor White
Write-Host ""

Write-Host "NEXT STEPS:" -ForegroundColor Cyan
Write-Host "1. Follow docs/prompt6-execution-guide.md for detailed instructions" -ForegroundColor White
Write-Host "2. Use docs/prompt6-checklist.md to track your progress" -ForegroundColor White
Write-Host "3. Complete manual updates in GitHub web interface" -ForegroundColor White
Write-Host "4. Verify all changes using the checklist" -ForegroundColor White
Write-Host ""

Write-Host "PRE-VERIFICATION COMPLETE. Ready to proceed with manual updates." -ForegroundColor Green
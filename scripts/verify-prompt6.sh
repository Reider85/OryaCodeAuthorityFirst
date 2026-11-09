#!/bin/bash

# Prompt 6: Verify Current State Before Manual Updates
# This script helps verify the current state of the repository before manual GitHub settings updates

echo "🔍 Prompt 6: Pre-verification Script"
echo "===================================="
echo

# Check if social preview SVG exists and has correct dimensions
echo "1. Checking social preview SVG..."
if [ -f "docs/social-preview.svg" ]; then
    echo "✅ docs/social-preview.svg exists"
    
    # Check file size
    SIZE=$(stat -f%z "docs/social-preview.svg" 2>/dev/null || stat -c%s "docs/social-preview.svg" 2>/dev/null)
    echo "   File size: $SIZE bytes"
    
    # Check if it contains required content
    if grep -q "CodeAuthorityFirst" "docs/social-preview.svg"; then
        echo "✅ Contains 'CodeAuthorityFirst' title"
    else
        echo "❌ Missing 'CodeAuthorityFirst' title"
    fi
    
    if grep -q "who wrote what, when, and why" "docs/social-preview.svg"; then
        echo "✅ Contains required tagline"
    else
        echo "❌ Missing required tagline"
    fi
    
    if grep -q "1280" "docs/social-preview.svg" && grep -q "640" "docs/social-preview.svg"; then
        echo "✅ Has correct dimensions (1280x640)"
    else
        echo "❌ Incorrect dimensions"
    fi
else
    echo "❌ docs/social-preview.svg does not exist"
fi
echo

# Check README.md for existing content
echo "2. Checking README.md..."
if [ -f "README.md" ]; then
    echo "✅ README.md exists"
    
    # Check character count of description
    if grep -q "Track who wrote what, when, and why" README.md; then
        echo "✅ Contains project description"
    else
        echo "❌ Missing project description"
    fi
    
    # Check for badges
    BADGE_COUNT=$(grep -o "!\[.*\](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/" README.md | wc -l)
    echo "   Found $BADGE_COUNT CI/quality badges"
else
    echo "❌ README.md does not exist"
fi
echo

# Check if docs directory exists
echo "3. Checking documentation..."
if [ -d "docs" ]; then
    echo "✅ docs/ directory exists"
    DOC_COUNT=$(find docs -name "*.md" -o -name "*.svg" | wc -l)
    echo "   Contains $DOC_COUNT documentation files"
else
    echo "❌ docs/ directory does not exist"
fi
echo

# Check GitHub repository status (this is informational only)
echo "4. Current GitHub Repository Status (Informational):"
echo "   Topics: 0 (need to add 20)"
echo "   Description: Empty (need to add)"
echo "   Website: Not configured (need to add)"
echo "   Social Preview: Not uploaded (SVG ready)"
echo

echo "📋 Next Steps:"
echo "1. Follow docs/prompt6-execution-guide.md for detailed instructions"
echo "2. Use docs/prompt6-checklist.md to track your progress"
echo "3. Complete manual updates in GitHub web interface"
echo "4. Verify all changes using the checklist"
echo

echo "✅ Pre-verification complete. Ready to proceed with manual updates."
#!/bin/bash

# build.sh - Build script for CodeAuthorityFirst
# This script builds the caf binary with version embedding using ldflags
# Supports Unix/Linux/macOS systems

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "🔨 Building CodeAuthorityFirst..."

# Get version from git or use fallback
if command -v git &> /dev/null && [ -d .git ]; then
    GIT_VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo "unknown")
    VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GIT_VERSION'"
    echo "📦 Version from git: $GIT_VERSION"
else
    VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
    echo "⚠️  Git not available, using fallback version"
fi

# Build flags
LDFLAGS="$VERSION_FLAG -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit=$(git rev-parse HEAD 2>/dev/null || echo unknown)'"
LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.datetime=$(date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo unknown)'"
LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.modified=$(git diff-index --quiet HEAD -- 2>/dev/null || echo true)'"

echo "🚀 Building with flags: $LDFLAGS"

# Build the binary
go build -ldflags="$LDFLAGS" -o caf ./cmd/caf

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "📋 Binary created: $(pwd)/caf"
    
    # Test the binary
    echo "🧪 Testing binary..."
    ./caf --version
    
    echo "🎉 Build complete!"
else
    echo "❌ Build failed!"
    exit 1
fi
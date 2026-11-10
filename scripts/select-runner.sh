#!/bin/bash

# Runner Selection Utility
# This script determines whether to use GitHub-hosted or local runners
# based on GitHub Actions usage thresholds

set -e

# Configuration
THRESHOLD=${1:-80}  # Default threshold: 80%
RUNNER_TYPE=${2:-hybrid}  # hybrid, github-only, local-only

echo "🔍 Runner Selection Utility"
echo "Threshold: ${THRESHOLD}%"
echo "Mode: ${RUNNER_TYPE}"

# Check if we're in a GitHub Actions environment
if [[ -n "${GITHUB_ACTIONS}" ]]; then
    echo "📊 Running in GitHub Actions environment"
    
    # Build the Go runner manager
    echo "🔨 Building runner-manager..."
    go build runner-manager.go
    
    # Check GitHub usage
    echo "📈 Checking GitHub Actions usage..."
    if ./runner-manager "${THRESHOLD}"; then
        echo "✅ Using GitHub-hosted runners"
        echo "::set-output name=runner-type::github"
        exit 0
    else
        echo "🚨 GitHub usage threshold exceeded"
        if [[ "${RUNNER_TYPE}" == "hybrid" || "${RUNNER_TYPE}" == "local-only" ]]; then
            echo "🏃‍♂️ Switching to local runners"
            echo "::set-output name=runner-type::local"
            exit 0
        else
            echo "❌ Local runners not available in current mode"
            exit 1
        fi
    fi
else
    echo "🏃‍♂️ Running locally - using local runners"
    echo "::set-output name=runner-type::local"
    exit 0
fi
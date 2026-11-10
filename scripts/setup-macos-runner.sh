#!/bin/bash

# Self-Hosted Runner Setup Script for macOS
# This script sets up a GitHub Actions self-hosted runner on macOS systems

set -e

# Configuration
RUNNER_NAME="hybrid-runner-macos"
RUNNER_LABELS="macos,self-hosted,hybrid"
RUNNER_GROUP="hybrid-runners"
GITHUB_TOKEN="${1:-}"

if [[ -z "${GITHUB_TOKEN}" ]]; then
    echo "❌ Error: GITHUB_TOKEN environment variable not set"
    echo "Usage: $0 <github_token>"
    exit 1
fi

echo "🚀 Setting up self-hosted GitHub Actions runner"
echo "Runner Name: ${RUNNER_NAME}"
echo "Labels: ${RUNNER_LABELS}"
echo "Group: ${RUNNER_GROUP}"

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "📦 Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install dependencies
echo "🔧 Installing dependencies..."
brew install docker git

# Start Docker Desktop
echo "🐳 Starting Docker Desktop..."
open -a "Docker Desktop"
sleep 30

# Download and configure runner
echo "🏃‍♂️ Downloading GitHub Actions runner..."
cd /tmp
curl -o actions-runner-osx-x64-2.313.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.313.0/actions-runner-osx-x64-2.313.0.tar.gz
tar xzf ./actions-runner-osx-x64-2.313.0.tar.gz
cd actions-runner-*

# Configure runner
echo "⚙️  Configuring runner..."
./config.sh --url https://github.com/${GITHUB_REPOSITORY:-Reider85/OryaCodeAuthorityFirst} \
            --token ${GITHUB_TOKEN} \
            --name ${RUNNER_NAME} \
            --labels ${RUNNER_LABELS} \
            --runnergroup ${RUNNER_GROUP} \
            --unattended

# Install as service
echo "🔄 Installing runner service..."
sudo ./svc.sh install

# Start service
echo "▶️  Starting runner service..."
sudo ./svc.sh start

echo "✅ Self-hosted runner setup complete!"
echo "Runner will appear in GitHub Actions settings within a few minutes"
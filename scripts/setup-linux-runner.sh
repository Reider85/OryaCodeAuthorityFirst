#!/bin/bash

# Self-Hosted Runner Setup Script for Linux
# This script sets up a GitHub Actions self-hosted runner on Linux systems

set -e

# Configuration
RUNNER_NAME="hybrid-runner-linux"
RUNNER_LABELS="linux,self-hosted,hybrid"
GITHUB_TOKEN="${1:-}"
RUNNER_GROUP="hybrid-runners"

if [[ -z "${GITHUB_TOKEN}" ]]; then
    echo "❌ Error: GITHUB_TOKEN environment variable not set"
    echo "Usage: $0 <github_token>"
    exit 1
fi

echo "🚀 Setting up self-hosted GitHub Actions runner"
echo "Runner Name: ${RUNNER_NAME}"
echo "Labels: ${RUNNER_LABELS}"
echo "Group: ${RUNNER_GROUP}"

# Check if running as root
if [[ $EUID -eq 0 ]]; then
    echo "⚠️  Running as root - this is not recommended"
    echo "Consider running with a non-root user that has sudo privileges"
fi

# Update system
echo "📦 Updating system packages..."
sudo apt-get update
sudo apt-get install -y curl wget git

# Install dependencies
echo "🔧 Installing dependencies..."
sudo apt-get install -y docker.io docker-compose

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER

# Download and configure runner
echo "🏃‍♂️ Downloading GitHub Actions runner..."
cd /tmp
curl -o actions-runner-linux-x64-2.313.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.313.0/actions-runner-linux-x64-2.313.0.tar.gz
tar xzf ./actions-runner-linux-x64-2.313.0.tar.gz
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
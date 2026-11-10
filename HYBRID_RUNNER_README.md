# Hybrid GitHub Actions Runner System

This system implements a hybrid approach for GitHub Actions runners that automatically switches between GitHub-hosted and self-hosted runners based on usage thresholds to optimize costs and performance.

## 🚀 Overview

The hybrid runner system provides:

- **Cost Optimization**: Automatically switches to local runners when GitHub's free minutes are exhausted
- **Performance Improvement**: Reduces queue times for performance-critical workflows
- **High Availability**: Fallback mechanism when GitHub-hosted runners are unavailable
- **Flexible Configuration**: Configurable thresholds and runner types

## 📋 Features

- **Smart Runner Selection**: Uses GitHub API to monitor usage and make intelligent decisions
- **Multi-Platform Support**: Linux, macOS, and Windows self-hosted runners
- **Real-time Monitoring**: Dashboard to track usage and runner status
- **Automated Setup**: Scripts to easily deploy self-hosted runners
- **Seamless Integration**: Works with existing GitHub Actions workflows

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ GitHub API      │    │ Runner Manager  │    │ Runner Selection│
│ Monitoring      │◄──►│ (Go)            │◄──►│ Logic           │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ GitHub-Hosted   │    │ Self-Hosted     │    │ Workflows       │
│ Runners         │    │ Runners        │    │ (.github)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Installation

### 1. Set Up GitHub Token

```bash
export GITHUB_TOKEN=your_github_personal_access_token
```

### 2. Deploy Self-Hosted Runners

#### Linux Runner
```bash
chmod +x scripts/setup-linux-runner.sh
./scripts/setup-linux-runner.sh $GITHUB_TOKEN
```

#### macOS Runner
```bash
chmod +x scripts/setup-macos-runner.sh
./scripts/setup-macos-runner.sh $GITHUB_TOKEN
```

#### Windows Runner
```powershell
.\scripts\setup-windows-runner.ps1 -GitHubToken $GITHUB_TOKEN
```

### 3. Update Workflows

Modify your GitHub Actions workflows to use the hybrid runner selection:

```yaml
jobs:
  select-runner:
    name: Select Runner Type
    runs-on: ubuntu-latest
    outputs:
      runner-type: ${{ steps.select.outputs.runner-type }}
      runner-label: ${{ steps.select.outputs.runner-label }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Select Runner Type
        id: select
        run: |
          if [[ "${{ env.RUNNER_TYPE }}" == "hybrid" ]]; then
            echo "🔍 Using hybrid runner selection..."
            echo "::set-output name=runner-type::$(bash ./scripts/select-runner.sh 80 hybrid)"
            echo "::set-output name=runner-label::$(bash ./scripts/select-runner.sh 80 hybrid | cut -d'=' -f2)"
          else
            echo "::set-output name=runner-type::${{ env.RUNNER_TYPE }}"
            echo "::set-output name=runner-label::${{ env.RUNNER_TYPE }}"
          fi

  your-job:
    name: Your Job
    runs-on: ${{ needs.select-runner.outputs.runner-label }}
    needs: select-runner
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      # ... other steps
```

## 🔧 Configuration

### Environment Variables

- `RUNNER_THRESHOLD`: Usage threshold percentage (default: 80)
- `RUNNER_TYPE`: Runner selection mode (hybrid, github-only, local-only)

### Runner Labels

Self-hosted runners are configured with these labels:
- `linux`, `macos`, `windows`: Operating system
- `self-hosted`: Indicates self-hosted runner
- `hybrid`: Indicates hybrid mode capability

## 📊 Monitoring

### Dashboard

Run the monitoring dashboard to track usage and runner status:

```bash
python scripts/runner-dashboard.py
```

### GitHub API Integration

The system uses GitHub's API to monitor:
- Total minutes used
- Remaining free minutes
- Runner availability
- Usage breakdown by runner type

## 🚀 Usage Scenarios

### 1. Normal Operation (Below Threshold)
- Uses GitHub-hosted runners
- Fast setup and maintenance
- Access to latest tools and environments

### 2. High Usage (Above Threshold)
- Automatically switches to self-hosted runners
- No impact on CI/CD pipeline
- Cost optimization

### 3. GitHub Hosted Unavailable
- Fallback to self-hosted runners
- Continuous operation
- No service interruption

## 📁 File Structure

```
├── runner-manager.go          # Go utility for GitHub API monitoring
├── scripts/
│   ├── select-runner.sh       # Shell script for runner selection
│   ├── select-runner.ps1      # PowerShell script for runner selection
│   ├── setup-linux-runner.sh  # Linux runner setup script
│   ├── setup-macos-runner.sh # macOS runner setup script
│   ├── setup-windows-runner.ps1 # Windows runner setup script
│   └── runner-dashboard.py     # Python monitoring dashboard
├── .github/workflows/
│   ├── ci.yml                 # Updated CI workflow
│   ├── shared-runner-config.yml # Shared runner configuration
│   └── ...                    # Other workflows
└── README.md                  # This file
```

## 🔐 Security Considerations

1. **Token Management**: Store GitHub tokens as secrets
2. **Runner Security**: Use dedicated runner accounts
3. **Network Security**: Configure firewall rules for runner communication
4. **Access Control**: Limit runner access to specific repositories

## 🚨 Troubleshooting

### Common Issues

1. **Runner Not Registering**
   - Check GitHub token permissions
   - Verify network connectivity
   - Check runner logs

2. **Usage Monitoring Not Working**
   - Verify GitHub token has access to repository
   - Check API rate limits
   - Verify repository visibility

3. **Workflow Failures**
   - Check runner labels match workflow requirements
   - Verify runner availability
   - Check environment configuration

### Debug Commands

```bash
# Test runner manager
go build runner-manager.go
./runner-manager 80

# Check runner status
python scripts/runner-dashboard.py

# Test runner selection
bash scripts/select-runner.sh 80 hybrid
```

## 🔄 Maintenance

### Updating Runners

1. Download latest runner release
2. Update runner configuration
3. Restart runner service
4. Monitor for issues

### Monitoring

- Regularly check usage statistics
- Monitor runner health
- Update threshold based on usage patterns
- Scale runner resources as needed

## 📈 Performance Optimization

### Runner Scaling

- Add more runners during peak usage
- Use different runner types for different workloads
- Implement auto-scaling based on queue length

### Caching

- Implement artifact caching
- Use dependency caching
- Optimize build cache strategies

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the troubleshooting section
2. Review the GitHub Actions documentation
3. Open an issue in the repository
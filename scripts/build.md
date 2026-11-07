# Build Guide for CodeAuthorityFirst

This document explains how to build CodeAuthorityFirst with proper version embedding.

## Build Methods

### 1. Using Build Scripts (Recommended)

#### Unix/Linux/macOS
```bash
# Build with version embedding
./scripts/build.sh

# Clean and rebuild
./scripts/build.sh --clean
```

#### Windows
```powershell
# Build with version embedding
.\scripts\build.ps1

# Clean and rebuild
.\scripts\build.ps1 -Clean
```

### 2. Manual Build

#### With Version Embedding
```bash
# Get current git version
GIT_VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo "unknown")

# Build with ldflags
go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GIT_VERSION'" -o caf ./cmd/caf
```

#### Windows (Command Prompt)
```cmd
:: Get current git version
for /f "delims=" %%v in ('git describe --tags --dirty --always 2^>nul') do set GIT_VERSION=%%v
if "%GIT_VERSION%"=="" set GIT_VERSION=unknown

:: Build with ldflags
go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=%GIT_VERSION%" -o caf.exe ./cmd/caf
```

## Version Information

The build process embeds the following information:
- **Version**: Git tag or fallback version
- **Commit**: Git commit hash
- **Build Time**: UTC timestamp
- **Modified**: Whether working directory has uncommitted changes

### Version Display

After building, you can check the embedded version:
```bash
./caf --version
```

Example output:
```
CodeAuthorityFirst v0.4.1-mvp.1 (abc1234)
```

Or if the working directory is dirty:
```
CodeAuthorityFirst v0.4.1-mvp.1 (dirty:abc1234)
```

## CI/CD Integration

The CI workflow (`.github/workflows/ci.yml`) automatically builds with version embedding:
```yaml
- name: Build with version embedding
  run: |
    GIT_VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo "unknown")
    go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GIT_VERSION'" -o caf ./cmd/caf
```

## Development Builds

For development, you can build without version embedding:
```bash
go build -o caf ./cmd/caf
```

However, this will show "unknown" as the version.

## Troubleshooting

### Common Issues

1. **Git not available**: The build will use fallback version "unknown"
2. **Dirty working directory**: The version will include "(dirty:)" prefix
3. **No tags**: Git describe will use commit hash

### Version Extraction

The version information is extracted from:
1. **ldflags** (preferred): Embedded at build time
2. **runtime/debug**: Fallback if ldflags not available
3. **Default version**: "0.4.1-mvp.1" if neither above is available

## File Structure

```
scripts/
├── build.sh          # Unix/Linux/macOS build script
├── build.ps1         # Windows PowerShell build script
├── test_build.bat    # Windows batch test script
└── build.md          # This documentation

internal/
└── version/
    └── version.go     # Version management package

cmd/caf/
├── main.go          # Updated to use version package
└── dashboard_tui.go  # Updated to remove hardcoded versions
```

## Best Practices

1. **Always use build scripts** for consistent builds
2. **Commit changes** before building to get clean version
3. **Use CI/CD** for automated builds with version embedding
4. **Test the binary** after build with `--version` flag
5. **Keep version package** as the single source of truth
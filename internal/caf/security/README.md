# Security Package

The security package provides comprehensive security features for the CodeAuthorityFirst application, including sandboxed YAML parsing, secret scanning, audit logging, and security linting.

## Features

### 1. Sandbox YAML Parser
- Safe YAML parsing with protection against code execution
- Configurable validation modes (basic, strict, paranoid)
- Protection against YAML injection attacks
- Custom validation rules and blocked patterns

### 2. Secret Scanner
- Detects various types of secrets in YAML files:
  - API keys
  - JWT tokens
  - Passwords
  - Private keys
  - Certificates
  - Database URLs
  - AWS keys
- Configurable detection rules and severity thresholds
- Heuristic-based detection with entropy analysis
- Context-aware scanning

### 3. Audit Logger
- Comprehensive security event logging
- Multiple output formats (console, file, JSON)
- Configurable log levels and retention policies
- Structured logging with correlation IDs
- Automatic log rotation

### 4. Security Linter
- Static analysis for security vulnerabilities
- Detection of common security issues:
  - Hardcoded credentials
  - Insecure cryptographic usage
  - SQL injection risks
  - Path traversal vulnerabilities
  - Information exposure
- Customizable rules and severity levels

### 5. Security Manager
- Centralized coordination of all security components
- Configuration management
- Lifecycle management
- Comprehensive security analysis

## Installation

```bash
go get github.com/your-org/codeauthorityfirst/internal/caf/security
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/your-org/codeauthorityfirst/internal/caf/security"
)

func main() {
    // Create security manager
    manager := security.NewSecurityManager()
    
    // Initialize with default configuration
    err := manager.Initialize("")
    if err != nil {
        panic(err)
    }
    defer manager.Shutdown()
    
    // Parse YAML with security checks
    yamlContent := []byte(`
name: my_app
version: 1.0
config:
  database_url: postgresql://localhost:5432/mydb
  api_key: AKIA1234567890123456
`)
    
    var config map[string]interface{}
    err = manager.ParseYAML(yamlContent, &config, "app_config.yml")
    if err != nil {
        fmt.Printf("YAML validation failed: %v\n", err)
    } else {
        fmt.Println("YAML parsed successfully with security checks")
    }
    
    // Scan for secrets
    result, err := manager.ScanForSecrets(yamlContent, "app_config.yml")
    if err != nil {
        fmt.Printf("Secret scanning failed: %v\n", err)
    } else {
        fmt.Printf("Found %d potential secrets\n", len(result.Secrets))
    }
}
```

### Configuration

```go
// Create custom configuration
config := security.DefaultSecurityConfiguration()

// Configure YAML security
config.YAML.SandboxParser = true
config.YAML.ValidationMode = "strict"
config.YAML.MaxNestingDepth = 10

// Configure secret scanning
config.Secrets.Enabled = true
config.Secrets.SeverityThreshold = 0.8
config.Secrets.EnableHeuristics = true

// Configure audit logging
config.Audit.Enabled = true
config.Audit.LogLevel = "info"
config.Audit.OutputDir = "./audit_logs"

// Save configuration
err := config.SaveSecurityConfiguration("security.yml")
if err != nil {
    panic(err)
}

// Load configuration
loadedConfig, err := security.LoadSecurityConfiguration("security.yml")
if err != nil {
    panic(err)
}
```

### Advanced Usage

#### Custom Secret Detection Rules

```go
scanner := security.NewSecretScanner()

// Add custom rule for database URLs
scanner.AddRule(security.SecretRule{
    ID:          "custom_db_url",
    Name:        "Custom Database URL",
    Description: "Detects custom database URLs",
    Pattern:     regexp.MustCompile(`db_url:\s*["']?(postgresql|mysql|mongodb)://[^"'\s]+["']?`),
    Type:        security.SecretTypeDatabaseURL,
    Severity:    "high",
    Confidence:  0.8,
})
```

#### Custom Security Linting Rules

```go
linter := security.NewSecurityLinter()

// Add custom linting rule
linter.AddRule(security.SecurityLintRule{
    ID:          "custom_rule",
    Name:        "Custom Security Rule",
    Description: "Detects potentially dangerous patterns",
    Severity:    "high",
    Category:    "custom",
    CheckFunc: func(node *ast.File, fset *token.FileSet, issues []security.SecurityIssue) []security.SecurityIssue {
        // Custom logic to detect security issues
        return issues
    },
    Enabled: true,
})
```

#### Comprehensive Security Analysis

```go
// Perform comprehensive security analysis
result, err := manager.LintSecurity("/path/to/code", true)
if err != nil {
    panic(err)
}

// Display results
fmt.Println(result.FormatResult())

// Check for critical issues
if result.IsCritical() {
    fmt.Println("Critical issues found!")
}

// Get issues by severity
criticalIssues := result.GetCriticalIssues()
highIssues := result.GetHighSeverityIssues()
```

## API Reference

### SecurityManager

#### Methods

- `Initialize(configPath string) error` - Initialize the security manager
- `ParseYAML(content []byte, target interface{}, filePath string) error` - Parse YAML with security checks
- `ScanForSecrets(content []byte, filePath string) (*ScanResult, error)` - Scan content for secrets
- `LintFile(filePath string) ([]SecurityIssue, error)` - Lint a single file
- `LintDirectory(dirPath string) ([]SecurityIssue, error)` - Lint a directory
- `ValidateYAMLFile(filePath string) error` - Validate a YAML file
- `LoadConfiguration(configPath string) error` - Load configuration from file
- `SaveConfiguration(configPath string) error` - Save configuration to file
- `GetConfiguration() *SecurityConfiguration` - Get current configuration
- `UpdateConfiguration(updater func(*SecurityConfiguration)) error` - Update configuration
- `GetStats() map[string]interface{}` - Get statistics
- `Shutdown() error` - Shutdown the security manager

### SecurityConfiguration

#### Fields

- `YAML` - YAML security settings
- `Secrets` - Secret scanning settings
- `Audit` - Audit logging settings
- `Authentication` - Authentication settings

#### Methods

- `Validate() error` - Validate configuration
- `ToJSON() ([]byte, error)` - Convert to JSON
- `FromJSON(jsonData []byte) error` - Create from JSON

### SecretScanner

#### Methods

- `ScanYAML(content []byte, filePath string) *ScanResult` - Scan YAML for secrets
- `ScanFile(filePath string, content []byte) *ScanResult` - Scan file for secrets
- `AddRule(rule SecretRule)` - Add custom detection rule
- `RemoveRule(id string)` - Remove rule by ID
- `FilterBySeverity(secrets []Secret, minSeverity float64) []Secret` - Filter secrets by severity
- `FormatSecrets(secrets []Secret) string` - Format secrets for display

### AuditLogger

#### Methods

- `LogEvent(event SecurityEvent) error` - Log a security event
- `LogYAMLParse(filePath string, success bool, duration time.Duration, details map[string]interface{}) error` - Log YAML parsing
- `LogSecretScan(filePath string, secretsFound int, duration time.Duration, details map[string]interface{}) error` - Log secret scanning
- `LogConfigLoad(configPath string, success bool, duration time.Duration, details map[string]interface{}) error` - Log configuration loading
- `LogAuthCheck(userID, resource, action string, success bool, details map[string]interface{}) error` - Log authentication checks
- `LogDataAccess(userID, resource, action string, details map[string]interface{}) error` - Log data access
- `LogError(message string, details map[string]interface{}) error` - Log error
- `LogWarning(message string, details map[string]interface{}) error` - Log warning
- `LogInfo(message string, details map[string]interface{}) error` - Log info
- `GetStats() map[string]interface{}` - Get statistics
- `Shutdown() error` - Shutdown the logger

## Security Features

### YAML Security

The sandboxed YAML parser provides protection against:

- YAML injection attacks
- Code execution through custom tags
- Excessive nesting depth
- Large document sizes
- Blocked patterns and dangerous constructs
- Custom validation rules

### Secret Detection

The secret scanner detects:

- API keys (AWS, Google, Microsoft, etc.)
- JWT tokens
- Passwords
- Private keys (RSA, EC)
- Certificates
- Database connection strings
- OAuth secrets
- Custom patterns

### Security Linting

The security linter detects:

- Hardcoded credentials
- Insecure cryptographic algorithms
- SQL injection vulnerabilities
- Path traversal risks
- Information exposure in error messages
- Insecure random number generation
- Unchecked error returns

### Audit Logging

The audit logger provides:

- Comprehensive event logging
- Multiple output formats
- Configurable retention policies
- Structured logging with metadata
- Automatic log rotation
- Performance monitoring

## Configuration

### YAML Security Configuration

```yaml
security:
  yaml:
    sandbox_parser: true
    max_nesting_depth: 10
    max_document_size: 1048576
    allow_custom_tags: false
    validation_mode: strict
    required_fields: []
    blocked_patterns:
      - "<!\\[CDATA\\[.*?\\]\\]>"
      - "<\\?xml.*?\\?>"
      - "<!.*?>"
```

### Secret Scanning Configuration

```yaml
security:
  secrets:
    enabled: true
    max_line_length: 1000
    min_secret_length: 8
    max_secret_length: 512
    severity_threshold: 0.7
    enable_heuristics: true
    enable_entropy_check: true
    enable_context_check: true
    ignore_files:
      - ".git"
      - ".gitignore"
      - "node_modules"
```

### Audit Logging Configuration

```yaml
security:
  audit:
    enabled: true
    log_to_console: true
    log_to_file: true
    log_level: info
    max_file_size: 10485760
    max_files: 5
    retain_days: 30
    output_dir: audit_logs
    enable_json: true
    enable_structured: true
```

## Performance Considerations

- The security parser is optimized for performance with minimal overhead
- Secret scanning uses efficient pattern matching algorithms
- Audit logging supports asynchronous operation for better performance
- Configuration is cached for better performance
- Components can be disabled individually to reduce overhead

## Security Best Practices

1. **Enable all security features** in production environments
2. **Set appropriate severity thresholds** based on your risk tolerance
3. **Regularly audit logs** for security events
4. **Use strict validation mode** for maximum security
5. **Customize detection rules** for your specific environment
6. **Monitor performance** to ensure security features don't impact system performance
7. **Keep security components updated** with the latest security patches

## Troubleshooting

### Common Issues

1. **YAML parsing fails** - Check for blocked patterns or invalid YAML structure
2. **Secret scanner produces false positives** - Adjust severity threshold or add custom rules
3. **Audit logs not being written** - Check file permissions and disk space
4. **Security linting too slow** - Use file filtering and reduce scope

### Debug Mode

Enable debug logging for troubleshooting:

```go
logger := security.NewAuditLogger()
logger.LogLevel = "debug"
logger.Enabled = true
```

## Examples

See the `example.go` file for comprehensive usage examples demonstrating all security features.

## Testing

Run the test suite:

```bash
go test -v ./internal/caf/security/
```

Run benchmarks:

```bash
go test -bench=. ./internal/caf/security/
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.
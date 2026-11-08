package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/security"
)

func main() {
	fmt.Println("=== CodeAuthorityFirst Security Demo ===")
	
	// Create security manager
	manager := security.NewSecurityManager()
	
	// Create configuration directory
	configDir := ".caf"
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	
	// Create security configuration
	configPath := filepath.Join(configDir, "security.yml")
	config := security.DefaultSecurityConfiguration()
	
	// Enable all security features for demo
	config.YAML.SandboxParser = true
	config.YAML.ValidationMode = "strict"
	config.Secrets.Enabled = true
	config.Secrets.SeverityThreshold = 0.7
	config.Audit.Enabled = true
	config.Audit.LogToConsole = true
	config.Audit.OutputDir = "audit_logs"
	config.Audit.EnableJSON = true
	
	// Save configuration
	err = config.SaveSecurityConfiguration(configPath)
	if err != nil {
		log.Fatalf("Failed to save security configuration: %v", err)
	}
	
	// Initialize security manager
	err = manager.Initialize(configPath)
	if err != nil {
		log.Fatalf("Failed to initialize security manager: %v", err)
	}
	defer manager.Shutdown()
	
	fmt.Printf("Security manager initialized with configuration: %s\n", configPath)
	
	// Demo 1: YAML Security
	fmt.Println("\n=== Demo 1: YAML Security ===")
	demoYAMLSecurity(manager)
	
	// Demo 2: Secret Scanning
	fmt.Println("\n=== Demo 2: Secret Scanning ===")
	demoSecretScanning(manager)
	
	// Demo 3: Security Linting
	fmt.Println("\n=== Demo 3: Security Linting ===")
	demoSecurityLinting(manager)
	
	// Demo 4: Comprehensive Security Analysis
	fmt.Println("\n=== Demo 4: Comprehensive Security Analysis ===")
	demoSecurityAnalysis(manager)
	
	// Demo 5: Configuration Management
	fmt.Println("\n=== Demo 5: Configuration Management ===")
	demoConfigurationManagement(manager)
	
	// Demo 6: Real-world Usage
	fmt.Println("\n=== Demo 6: Real-world Usage ===")
	demoRealWorldUsage(manager)
	
	// Show final statistics
	fmt.Println("\n=== Final Statistics ===")
	showFinalStatistics(manager)
	
	fmt.Println("\n=== Security Demo Complete ===")
	fmt.Println("All security features have been demonstrated successfully!")
	fmt.Println("Check the 'audit_logs' directory for audit logs.")
}

func demoYAMLSecurity(manager *security.SecurityManager) {
	fmt.Println("Testing YAML parsing with security checks...")
	
	// Valid YAML
	validYAML := `
name: production_app
version: 2.0
config:
  database_host: localhost
  database_port: 5432
  max_connections: 100
  features:
    - user_management
    - reporting
    - audit_logging
  security:
    ssl_enabled: true
    encryption: aes256
`
	
	var configData map[string]interface{}
	err := manager.ParseYAML([]byte(validYAML), &configData, "production_config.yml")
	if err != nil {
		fmt.Printf("✗ Valid YAML parsing failed: %v\n", err)
	} else {
		fmt.Println("✓ Valid YAML parsed successfully")
		fmt.Printf("  Application: %v\n", configData["name"])
		fmt.Printf("  Version: %v\n", configData["version"])
	}
	
	// Invalid YAML with potential security issues
	invalidYAML := `
name: production_app
version: 2.0
config:
  database_url: postgresql://admin:!script alert('xss')@localhost:5432/prod_db
  features:
    - user_management
    - <!DOCTYPE html><script>alert('xss')</script>
`
	
	err = manager.ParseYAML([]byte(invalidYAML), &configData, "invalid_config.yml")
	if err != nil {
		fmt.Printf("✓ Invalid YAML correctly blocked: %v\n", err)
	} else {
		fmt.Println("✗ Invalid YAML should have been blocked")
	}
}

func demoSecretScanning(manager *security.SecurityManager) {
	fmt.Println("Testing secret scanning...")
	
	// Sample configuration with potential secrets
	configWithSecrets := `
# Production Configuration
app_name: my_production_app
version: 1.0

database:
  host: prod-db.example.com
  port: 5432
  name: production_db
  username: admin
  # This should be detected as a secret
  password: P@ssw0rd123!
  ssl_enabled: true
  
api:
  # This should be detected as an API key
  key: AKIAIOSFODNN7EXAMPLE
  timeout: 30
  rate_limit: 1000
  
security:
  # This should be detected as a JWT token
  jwt_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
  session_timeout: 3600
  
features:
  - authentication
  - authorization
  - logging
  - monitoring
`
	
	result, err := manager.ScanForSecrets([]byte(configWithSecrets), "production_config.yml")
	if err != nil {
		fmt.Printf("✗ Secret scanning failed: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Secret scanning completed. Found %d potential secrets:\n", len(result.Secrets))
	
	for i, secret := range result.Secrets {
		fmt.Printf("  %d. %s\n", i+1, secret.Type)
		fmt.Printf("     Line: %d, Column: %d\n", secret.Line, secret.Column)
		fmt.Printf("     Confidence: %.2f\n", secret.Confidence)
		fmt.Printf("     Context: %s\n", secret.Context)
		fmt.Printf("     Description: %s\n", secret.Description)
		fmt.Println()
	}
	
	if result.HasSecrets {
		fmt.Println("⚠️  Security alert: Potential secrets detected in configuration!")
	}
}

func demoSecurityLinting(manager *security.SecurityManager) {
	fmt.Println("Testing security linting...")
	
	// Create a temporary Go file with security issues
	goCode := `
package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	// Hardcoded password (security issue)
	password := "supersecretpassword123"
	
	// Insecure random number generation (security issue)
	sessionToken := rand.Intn(1000000)
	
	// Potential path traversal (security issue)
	userInput := "../../etc/passwd"
	file, err := os.Open(userInput)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	
	// SQL injection risk (security issue)
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userId)
	
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Session Token: %d\n", sessionToken)
	fmt.Printf("File: %s\n", file.Name())
	fmt.Printf("Query: %s\n", query)
}
`
	
	// Write test file
	testFile := "test_security.go"
	err := os.WriteFile(testFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Printf("✗ Failed to write test file: %v\n", err)
		return
	}
	defer os.Remove(testFile)
	
	// Lint the file
	issues, err := manager.LintFile(testFile)
	if err != nil {
		fmt.Printf("✗ Security linting failed: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Security linting completed. Found %d issues:\n", len(issues))
	
	for i, issue := range issues {
		fmt.Printf("  %d. %s\n", i+1, issue.Message)
		fmt.Printf("     Severity: %s\n", issue.Severity)
		fmt.Printf("     Category: %s\n", issue.Category)
		fmt.Printf("     File: %s:%d:%d\n", issue.File, issue.Line, issue.Column)
		fmt.Printf("     Suggestion: %s\n", issue.Suggestion)
		fmt.Println()
	}
	
	if len(issues) > 0 {
		fmt.Println("⚠️  Security issues detected in code!")
	}
}

func demoSecurityAnalysis(manager *security.SecurityManager) {
	fmt.Println("Testing comprehensive security analysis...")
	
	// Create test directory with multiple files
	testDir := "security_test_dir"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		fmt.Printf("✗ Failed to create test directory: %v\n", err)
		return
	}
	defer os.RemoveAll(testDir)
	
	// Create test files
	testFiles := map[string]string{
		"app_config.yml": `
name: test_app
version: 1.0
config:
  database_url: postgresql://admin:password@localhost:5432/testdb
  api_key: AKIA1234567890123456
`,
		"secrets.txt": `
api_key=AKIA1234567890123456
password=supersecretpassword
jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
`,
		"code.go": `
package main

import (
	"fmt"
	"os"
)

func main() {
	// Hardcoded secret
	secret := "AKIA1234567890123456"
	
	// Potential path traversal
	file := os.Open("/tmp/" + userInput)
	
	fmt.Println("Testing security analysis")
}
`,
	}
	
	// Write test files
	for filename, content := range testFiles {
		err := os.WriteFile(filepath.Join(testDir, filename), []byte(content), 0644)
		if err != nil {
			fmt.Printf("✗ Failed to write test file %s: %v\n", filename, err)
			return
		}
	}
	
	// Perform security analysis
	result, err := manager.LintSecurity(testDir, true)
	if err != nil {
		fmt.Printf("✗ Security analysis failed: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Security analysis completed. Analyzed %d files, found %d total issues:\n", len(testFiles), len(result.Issues))
	
	// Show summary
	fmt.Println("\nSummary:")
	for key, count := range result.Summary {
		fmt.Printf("  %s: %d\n", key, count)
	}
	
	// Show critical issues
	criticalIssues := result.GetCriticalIssues()
	if len(criticalIssues) > 0 {
		fmt.Printf("\nCritical Issues (%d):\n", len(criticalIssues))
		for i, issue := range criticalIssues {
			fmt.Printf("  %d. %s\n", i+1, issue.Message)
		}
	}
	
	// Show high severity issues
	highIssues := result.GetHighSeverityIssues()
	if len(highIssues) > 0 {
		fmt.Printf("\nHigh Severity Issues (%d):\n", len(highIssues))
		for i, issue := range highIssues {
			fmt.Printf("  %d. %s\n", i+1, issue.Message)
		}
	}
	
	fmt.Println("\nFormatted Analysis Report:")
	fmt.Println(result.FormatResult())
}

func demoConfigurationManagement(manager *security.SecurityManager) {
	fmt.Println("Testing configuration management...")
	
	// Show current configuration
	config := manager.GetConfiguration()
	fmt.Printf("Current configuration:\n")
	fmt.Printf("  YAML Sandbox Parser: %v\n", config.YAML.SandboxParser)
	fmt.Printf("  Secret Scanning Enabled: %v\n", config.Secrets.Enabled)
	fmt.Printf("  Audit Logging Enabled: %v\n", config.Audit.Enabled)
	fmt.Printf("  Severity Threshold: %.2f\n", config.Secrets.SeverityThreshold)
	fmt.Printf("  Log Level: %s\n", config.Audit.LogLevel)
	
	// Update configuration
	err := manager.UpdateConfiguration(func(c *security.SecurityConfiguration) {
		c.Secrets.SeverityThreshold = 0.9
		c.Audit.LogLevel = "warning"
		c.YAML.ValidationMode = "paranoid"
	})
	if err != nil {
		fmt.Printf("✗ Failed to update configuration: %v\n", err)
		return
	}
	
	fmt.Println("✓ Configuration updated successfully")
	
	// Show updated configuration
	updatedConfig := manager.GetConfiguration()
	fmt.Printf("Updated configuration:\n")
	fmt.Printf("  Severity Threshold: %.2f\n", updatedConfig.Secrets.SeverityThreshold)
	fmt.Printf("  Log Level: %s\n", updatedConfig.Audit.LogLevel)
	fmt.Printf("  Validation Mode: %s\n", updatedConfig.YAML.ValidationMode)
	
	// Save updated configuration
	updatedConfigPath := ".caf/updated_security.yml"
	err = manager.SaveConfiguration(updatedConfigPath)
	if err != nil {
		fmt.Printf("✗ Failed to save configuration: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Configuration saved to: %s\n", updatedConfigPath)
}

func demoRealWorldUsage(manager *security.SecurityManager) {
	fmt.Println("Testing real-world usage scenario...")
	
	// Simulate loading a production configuration
	productionConfig := `
name: codeauthority_first
version: 1.0.0
environment: production

database:
  host: prod-db.example.com
  port: 5432
  name: codeauthority
  username: app_user
  # This should be detected as a secret
  password: ProdAppPassword123!
  pool_size: 20
  ssl_mode: require

api:
  # This should be detected as an API key
  key: AKIAIOSFODNN7EXAMPLE
  timeout: 30
  rate_limit: 1000
  cors_enabled: true

security:
  # This should be detected as a JWT token
  jwt_secret: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
  session_timeout: 3600
  max_login_attempts: 5
  lockout_duration: 15

features:
  - user_authentication
  - role_based_access
  - audit_logging
  - encryption
  - backup

monitoring:
  enabled: true
  metrics_port: 8080
  health_check_interval: 30
`
	
	fmt.Println("1. Loading and validating production configuration...")
	
	// Parse configuration with security checks
	var config map[string]interface{}
	err := manager.ParseYAML([]byte(productionConfig), &config, "production_config.yml")
	if err != nil {
		fmt.Printf("✗ Configuration validation failed: %v\n", err)
		return
	}
	
	fmt.Println("✓ Configuration parsed successfully")
	
	fmt.Println("2. Scanning for secrets...")
	
	// Scan for secrets
	result, err := manager.ScanForSecrets([]byte(productionConfig), "production_config.yml")
	if err != nil {
		fmt.Printf("✗ Secret scanning failed: %v\n", err)
		return
	}
	
	if result.HasSecrets {
		fmt.Printf("⚠️  Found %d potential secrets in production configuration:\n", len(result.Secrets))
		for i, secret := range result.Secrets {
			fmt.Printf("  %d. %s (Line %d)\n", i+1, secret.Type, secret.Line)
		}
		
		// Ask if user wants to proceed
		fmt.Println("\n⚠️  WARNING: Production configuration contains potential secrets!")
		fmt.Println("   Consider moving secrets to environment variables or secure vault.")
	} else {
		fmt.Println("✓ No secrets detected in configuration")
	}
	
	fmt.Println("3. Validating configuration file...")
	
	// Validate configuration file
	configFile := "production_config.yml"
	err = os.WriteFile(configFile, []byte(productionConfig), 0644)
	if err != nil {
		fmt.Printf("✗ Failed to write config file: %v\n", err)
		return
	}
	defer os.Remove(configFile)
	
	err = manager.ValidateYAMLFile(configFile)
	if err != nil {
		fmt.Printf("✗ Configuration file validation failed: %v\n", err)
		return
	}
	
	fmt.Println("✓ Configuration file validated successfully")
	
	fmt.Println("4. Performing security audit...")
	
	// Perform security audit
	auditResult, err := manager.LintSecurity(".", false)
	if err != nil {
		fmt.Printf("✗ Security audit failed: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Security audit completed. Found %d issues\n", len(auditResult.Issues))
	
	if auditResult.IsCritical() {
		fmt.Println("⚠️  CRITICAL: Security audit found critical issues!")
	} else if auditResult.HasHighSeverity() {
		fmt.Println("⚠️  WARNING: Security audit found high severity issues!")
	} else {
		fmt.Println("✓ Security audit passed with minor issues")
	}
	
	fmt.Println("5. Generating security report...")
	
	// Generate security report
	report := auditResult.FormatResult()
	fmt.Println("\nSecurity Report:")
	fmt.Println(report)
}

func showFinalStatistics(manager *security.SecurityManager) {
	stats := manager.GetStats()
	
	fmt.Printf("Security Manager Statistics:\n")
	fmt.Printf("  Initialized: %v\n", stats["initialized"])
	
	if configPath, ok := stats["config_path"].(string); ok {
		fmt.Printf("  Configuration Path: %s\n", configPath)
	}
	
	if lastUpdated, ok := stats["last_updated"]; ok {
		fmt.Printf("  Last Updated: %v\n", lastUpdated)
	}
	
	if auditStats, ok := stats["audit_logger"].(map[string]interface{}); ok {
		fmt.Printf("  Audit Events: %v\n", auditStats["total_events"])
		fmt.Printf("  Audit Enabled: %v\n", auditStats["enabled"])
	}
	
	if config, ok := stats["security_config"].(*security.SecurityConfiguration); ok {
		fmt.Printf("  Current Configuration:\n")
		fmt.Printf("    YAML Sandbox: %v\n", config.YAML.SandboxParser)
		fmt.Printf("    Secret Scanning: %v\n", config.Secrets.Enabled)
		fmt.Printf("    Audit Logging: %v\n", config.Audit.Enabled)
		fmt.Printf("    Severity Threshold: %.2f\n", config.Secrets.SeverityThreshold)
	}
}
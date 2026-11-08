package security

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExampleSecurityUsage demonstrates how to use the security components
func ExampleSecurityUsage() {
	// Create a temporary directory for demonstration
	tempDir := "security_demo"
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create demo directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create security configuration
	configPath := filepath.Join(tempDir, "security.yml")
	config := DefaultSecurityConfiguration()
	config.YAML.SandboxParser = true
	config.Secrets.Enabled = true
	config.Audit.Enabled = true
	config.Audit.LogToConsole = true
	config.Audit.OutputDir = filepath.Join(tempDir, "audit_logs")

	// Save configuration
	err = config.SaveSecurityConfiguration(configPath)
	if err != nil {
		fmt.Printf("Failed to save configuration: %v\n", err)
		return
	}

	// Initialize security manager
	manager := NewSecurityManager()
	err = manager.Initialize(configPath)
	if err != nil {
		fmt.Printf("Failed to initialize security manager: %v\n", err)
		return
	}
	defer manager.Shutdown()

	fmt.Println("=== Security Manager Demo ===")
	fmt.Println("1. YAML Parsing with Security Checks")
	
	// Example 1: YAML parsing with security checks
	yamlContent := `
name: my_application
version: 1.0
database:
  url: postgresql://user:password@localhost:5432/mydb
  ssl_enabled: true
api:
  key: AKIA1234567890123456
  timeout: 30
features:
  - authentication
  - logging
  - encryption
`

	var configData map[string]interface{}
	err = manager.ParseYAML([]byte(yamlContent), &configData, "app_config.yml")
	if err != nil {
		fmt.Printf("YAML parsing error: %v\n", err)
	} else {
		fmt.Println("✓ YAML parsed successfully with security checks")
	}

	fmt.Println("\n2. Secret Scanning")
	
	// Example 2: Secret scanning
	secretsYAML := `
# Application configuration
database_url: postgresql://admin:supersecret@localhost:5432/prod_db
api_key: AKIAIOSFODNN7EXAMPLE
jwt_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
password: mysecretpassword123
oauth_secret: client_secret_abcdef123456
`

	secretResult, err := manager.ScanForSecrets([]byte(secretsYAML), "secrets_config.yml")
	if err != nil {
		fmt.Printf("Secret scanning error: %v\n", err)
	} else {
		fmt.Printf("✓ Secret scanning completed. Found %d potential secrets:\n", len(secretResult.Secrets))
		for i, secret := range secretResult.Secrets {
			fmt.Printf("  %d. %s (Line %d, Confidence: %.2f)\n", 
				i+1, secret.Type, secret.Line, secret.Confidence)
		}
	}

	fmt.Println("\n3. Security Linting")
	
	// Example 3: Security linting
	goCode := `
package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	// Hardcoded credential (security issue)
	password := "supersecretpassword123"
	
	// Insecure random number generation (security issue)
	token := rand.Intn(1000)
	
	// Potential path traversal (security issue)
	file := os.Open("/tmp/" + userInput)
	
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Token: %d\n", token)
	fmt.Printf("File: %s\n", file)
}
`

	goFile := filepath.Join(tempDir, "example.go")
	err = os.WriteFile(goFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Printf("Failed to write Go file: %v\n", err)
		return
	}

	issues, err := manager.LintFile(goFile)
	if err != nil {
		fmt.Printf("Security linting error: %v\n", err)
	} else {
		fmt.Printf("✓ Security linting completed. Found %d issues:\n", len(issues))
		for i, issue := range issues {
			fmt.Printf("  %d. %s (%s)\n", i+1, issue.Message, issue.Severity)
		}
	}

	fmt.Println("\n4. YAML File Validation")
	
	// Example 4: YAML file validation
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
`

	validFile := filepath.Join(tempDir, "valid_config.yml")
	err = os.WriteFile(validFile, []byte(validYAML), 0644)
	if err != nil {
		fmt.Printf("Failed to write valid YAML file: %v\n", err)
		return
	}

	err = manager.ValidateYAMLFile(validFile)
	if err != nil {
		fmt.Printf("✗ Valid YAML file failed validation: %v\n", err)
	} else {
		fmt.Println("✓ Valid YAML file passed validation")
	}

	// Create invalid YAML file with secrets
	invalidYAML := `
name: production_app
version: 2.0
config:
  database_url: postgresql://admin:password123@localhost:5432/prod_db
  api_key: AKIA1234567890123456
  admin_password: supersecretpassword
`

	invalidFile := filepath.Join(tempDir, "invalid_config.yml")
	err = os.WriteFile(invalidFile, []byte(invalidYAML), 0644)
	if err != nil {
		fmt.Printf("Failed to write invalid YAML file: %v\n", err)
		return
	}

	err = manager.ValidateYAMLFile(invalidFile)
	if err != nil {
		fmt.Printf("✓ Invalid YAML file correctly failed validation: %v\n", err)
	} else {
		fmt.Println("✗ Invalid YAML file should have failed validation")
	}

	fmt.Println("\n5. Security Analysis")
	
	// Example 5: Comprehensive security analysis
	result, err := manager.LintSecurity(tempDir, true)
	if err != nil {
		fmt.Printf("Security analysis error: %v\n", err)
	} else {
		fmt.Printf("✓ Security analysis completed. Found %d total issues:\n", len(result.Issues))
		
		// Show summary
		fmt.Println("Summary:")
		for severity, count := range result.Summary {
			fmt.Printf("  %s: %d\n", severity, count)
		}
		
		// Show critical issues
		criticalIssues := result.GetCriticalIssues()
		if len(criticalIssues) > 0 {
			fmt.Printf("\nCritical Issues:\n")
			for i, issue := range criticalIssues {
				fmt.Printf("  %d. %s\n", i+1, issue.Message)
			}
		}
	}

	fmt.Println("\n6. Configuration Management")
	
	// Example 6: Configuration management
	fmt.Println("Current configuration:")
	currentConfig := manager.GetConfiguration()
	fmt.Printf("  YAML Sandbox Parser: %v\n", currentConfig.YAML.SandboxParser)
	fmt.Printf("  Secret Scanning Enabled: %v\n", currentConfig.Secrets.Enabled)
	fmt.Printf("  Audit Logging Enabled: %v\n", currentConfig.Audit.Enabled)
	fmt.Printf("  Secret Severity Threshold: %.2f\n", currentConfig.Secrets.SeverityThreshold)

	// Update configuration
	err = manager.UpdateConfiguration(func(c *SecurityConfiguration) {
		c.Secrets.SeverityThreshold = 0.9
		c.Audit.LogLevel = "warning"
	})
	if err != nil {
		fmt.Printf("Failed to update configuration: %v\n", err)
	} else {
		fmt.Println("✓ Configuration updated successfully")
	}

	// Verify update
	updatedConfig := manager.GetConfiguration()
	fmt.Printf("  Updated Severity Threshold: %.2f\n", updatedConfig.Secrets.SeverityThreshold)
	fmt.Printf("  Updated Log Level: %s\n", updatedConfig.Audit.LogLevel)

	fmt.Println("\n7. Statistics")
	
	// Example 7: Get statistics
	stats := manager.GetStats()
	fmt.Printf("  Initialized: %v\n", stats["initialized"])
	fmt.Printf("  Config Path: %s\n", stats["config_path"])
	
	if auditStats, ok := stats["audit_logger"].(map[string]interface{}); ok {
		fmt.Printf("  Total Events: %v\n", auditStats["total_events"])
	}

	fmt.Println("\n=== Demo Complete ===")
}

// ExampleYAMLParser demonstrates the sandboxed YAML parser
func ExampleYAMLParser() {
	fmt.Println("=== YAML Parser Demo ===")
	
	parser := NewSandboxYAMLParser()
	
	// Valid YAML
	validYAML := []byte(`
name: test_application
version: 1.0
config:
  database_url: postgresql://localhost:5432/mydb
  features:
    - authentication
    - logging
`)

	var result map[string]interface{}
	err := parser.ParseYAML(validYAML, &result)
	if err != nil {
		fmt.Printf("Error parsing valid YAML: %v\n", err)
	} else {
		fmt.Println("✓ Valid YAML parsed successfully")
		fmt.Printf("  Name: %v\n", result["name"])
		fmt.Printf("  Version: %v\n", result["version"])
	}

	// Invalid YAML with potential security issues
	invalidYAML := []byte(`
name: !script alert('xss')
content: <script>dangerous()</script>
`)

	err = parser.ParseYAML(invalidYAML, &result)
	if err != nil {
		fmt.Printf("✓ Invalid YAML correctly blocked: %v\n", err)
	} else {
		fmt.Println("✗ Invalid YAML should have been blocked")
	}
}

// ExampleSecretScanner demonstrates the secret scanner
func ExampleSecretScanner() {
	fmt.Println("=== Secret Scanner Demo ===")
	
	scanner := NewSecretScanner()
	
	// Test content with secrets
	testContent := `
# Application Configuration
database:
  url: postgresql://admin:supersecret@localhost:5432/prod_db
  ssl_enabled: true
  
api:
  key: AKIAIOSFODNN7EXAMPLE
  timeout: 30
  
security:
  jwt_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
  password: mysecretpassword123
`

	result := scanner.ScanYAML([]byte(testContent), "test_config.yml")
	
	fmt.Printf("Scanning completed. Found %d secrets:\n", len(result.Secrets))
	
	for i, secret := range result.Secrets {
		fmt.Printf("%d. %s\n", i+1, secret.Message)
		fmt.Printf("   Type: %s\n", secret.Type)
		fmt.Printf("   Line: %d\n", secret.Line)
		fmt.Printf("   Confidence: %.2f\n", secret.Confidence)
		fmt.Printf("   Context: %s\n", secret.Context)
		fmt.Println()
	}
}

// ExampleAuditLogger demonstrates the audit logger
func ExampleAuditLogger() {
	fmt.Println("=== Audit Logger Demo ===")
	
	logger := NewAuditLogger()
	logger.Enabled = true
	logger.LogToConsole = true
	logger.LogToFile = false
	
	// Log various types of events
	logger.LogYAMLParse("config.yml", true, 100, map[string]interface{}{"size": 1024})
	logger.LogSecretScan("config.yml", 1, 200, map[string]interface{}{"secrets": []string{"api_key"}})
	logger.LogConfigLoad("security.yml", true, 50, map[string]interface{}{"format": "yaml"})
	logger.LogAuthCheck("user1", "resource1", "read", true, map[string]interface{}{"ip": "127.0.0.1"})
	
	logger.LogError("Test error", map[string]interface{}{"error_code": 500})
	logger.LogWarning("Test warning", map[string]interface{}{"warning_code": 100})
	logger.LogInfo("Test info", map[string]interface{}{"info_code": 200})
	
	// Get statistics
	stats := logger.GetStats()
	fmt.Printf("Total events logged: %v\n", stats["total_events"])
	fmt.Printf("Logger enabled: %v\n", stats["enabled"])
}

// ExampleSecurityConfiguration demonstrates security configuration
func ExampleSecurityConfiguration() {
	fmt.Println("=== Security Configuration Demo ===")
	
	// Create default configuration
	config := DefaultSecurityConfiguration()
	
	fmt.Println("Default Configuration:")
	fmt.Printf("  YAML Sandbox Parser: %v\n", config.YAML.SandboxParser)
	fmt.Printf("  Secret Scanning Enabled: %v\n", config.Secrets.Enabled)
	fmt.Printf("  Audit Logging Enabled: %v\n", config.Audit.Enabled)
	fmt.Printf("  Validation Mode: %s\n", config.YAML.ValidationMode)
	fmt.Printf("  Severity Threshold: %.2f\n", config.Secrets.SeverityThreshold)
	
	// Update configuration
	config.YAML.ValidationMode = "strict"
	config.Secrets.SeverityThreshold = 0.8
	config.Audit.LogLevel = "warning"
	
	fmt.Println("\nUpdated Configuration:")
	fmt.Printf("  Validation Mode: %s\n", config.YAML.ValidationMode)
	fmt.Printf("  Severity Threshold: %.2f\n", config.Secrets.SeverityThreshold)
	fmt.Printf("  Log Level: %s\n", config.Audit.LogLevel)
	
	// Validate configuration
	err := config.Validate()
	if err != nil {
		fmt.Printf("Configuration validation error: %v\n", err)
	} else {
		fmt.Println("✓ Configuration is valid")
	}
	
	// Convert to JSON
	jsonData, err := config.ToJSON()
	if err != nil {
		fmt.Printf("JSON conversion error: %v\n", err)
	} else {
		fmt.Printf("JSON configuration length: %d bytes\n", len(jsonData))
	}
}
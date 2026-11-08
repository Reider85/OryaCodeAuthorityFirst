package security

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// TestSandboxYAMLParser tests the sandboxed YAML parser
func TestSandboxYAMLParser(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		target      interface{}
		expectError bool
		errorType   error
	}{
		{
			name:    "Valid YAML",
			content: "name: test\nvalue: 123",
			target:  map[string]interface{}{},
			expectError: false,
		},
		{
			name:    "YAML with custom tags (blocked)",
			content: "name: !tag test",
			target:  map[string]interface{}{},
			expectError: true,
			errorType: &YAMLValidationError{},
		},
		{
			name:    "YAML with script content (blocked)",
			content: "script: <script>alert('xss')</script>",
			target:  map[string]interface{}{},
			expectError: true,
		},
		{
			name:    "YAML with template expressions (blocked)",
			content: "template: ${dangerous}",
			target:  map[string]interface{}{},
			expectError: true,
		},
		{
			name:    "YAML with XML processing instructions (blocked)",
			content: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
			target:  map[string]interface{}{},
			expectError: true,
		},
		{
			name:    "Empty YAML",
			content: "",
			target:  map[string]interface{}{},
			expectError: true,
		},
		{
			name:    "YAML with excessive nesting",
			content: "level1:\n  level2:\n    level3:\n      level4:\n        level5:\n          level6:\n            level7:\n              level8:\n                level9:\n                  level10:\n                    level11: value",
			target:  map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewSandboxYAMLParser()
			err := parser.ParseYAML([]byte(tt.content), tt.target)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got none")
				} else if tt.errorType != nil && !isErrorOfType(err, tt.errorType) {
					t.Errorf("Expected error type %T, got %T", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestSafeParseYAML tests the convenience function
func TestSafeParseYAML(t *testing.T) {
	validYAML := []byte("name: test\nvalue: 123")
	var result map[string]interface{}
	
	err := SafeParseYAML(validYAML, &result)
	if err != nil {
		t.Errorf("SafeParseYAML failed: %v", err)
	}
	
	if result["name"] != "test" {
		t.Errorf("Expected name to be 'test', got %v", result["name"])
	}
}

// TestYAMLValidationError tests YAML validation error
func TestYAMLValidationError(t *testing.T) {
	err := &YAMLValidationError{
		Field:   "test_field",
		Message: "test message",
		Line:    10,
	}
	
	expected := "YAML validation error at line 10: test message (field: test_field)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestSecretScanner tests the secret scanner
func TestSecretScanner(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectCount int
	}{
		{
			name:        "No secrets",
			content:     "name: test\nvalue: 123",
			expectCount: 0,
		},
		{
			name:        "API key",
			content:     "api_key: AKIA1234567890123456",
			expectCount: 1,
		},
		{
			name:        "JWT token",
			content:     "token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectCount: 1,
		},
		{
			name:        "Password",
			content:     "password: mysecretpassword123",
			expectCount: 1,
		},
		{
			name:        "RSA private key",
			content:     "private_key: -----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA1234567890",
			expectCount: 1,
		},
		{
			name:        "Multiple secrets",
			content:     "api_key: AKIA1234567890123456\ntoken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\npassword: secret123",
			expectCount: 3,
		},
	}

	scanner := NewSecretScanner()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.ScanYAML([]byte(tt.content), "test.yml")
			
			if result.HasSecrets != (tt.expectCount > 0) {
				t.Errorf("Expected HasSecrets to be %v, got %v", tt.expectCount > 0, result.HasSecrets)
			}
			
			if len(result.Secrets) != tt.expectCount {
				t.Errorf("Expected %d secrets, got %d", tt.expectCount, len(result.Secrets))
			}
		})
	}
}

// TestSecretScannerCustomRules tests custom secret detection rules
func TestSecretSecretScannerCustomRules(t *testing.T) {
	scanner := NewSecretScanner()
	
	// Add custom rule for database URLs
	scanner.AddRule(SecretRule{
		ID:          "custom_db_url",
		Name:        "Custom Database URL",
		Description: "Detects custom database URLs",
		Pattern:     regexp.MustCompile(`db_url:\s*["']?(postgresql|mysql|mongodb)://[^"'\s]+["']?`),
		Type:        SecretTypeDatabaseURL,
		Severity:    "high",
		Confidence:  0.8,
	})
	
	content := `name: test
db_url: postgresql://user:pass@localhost:5432/db`
	
	result := scanner.ScanYAML([]byte(content), "test.yml")
	
	if len(result.Secrets) != 1 {
		t.Errorf("Expected 1 secret, got %d", len(result.Secrets))
	}
	
	if result.Secrets[0].Type != SecretTypeDatabaseURL {
		t.Errorf("Expected secret type %s, got %s", SecretTypeDatabaseURL, result.Secrets[0].Type)
	}
}

// TestSecretScannerFilter tests filtering by severity
func TestSecretScannerFilter(t *testing.T) {
	scanner := NewSecretScanner()
	
	// Create test secrets with different severities
	secrets := []Secret{
		{Type: SecretTypeAPIKey, Confidence: 0.9},
		{Type: SecretTypeToken, Confidence: 0.5},
		{Type: SecretTypePassword, Confidence: 0.8},
	}
	
	filtered := scanner.FilterBySeverity(secrets, 0.7)
	
	if len(filtered) != 2 {
		t.Errorf("Expected 2 secrets with confidence >= 0.7, got %d", len(filtered))
	}
}

// TestSecretScannerFormat tests formatting secrets
func TestSecretScannerFormat(t *testing.T) {
	scanner := NewSecretScanner()
	
	secrets := []Secret{
		{
			Type:        SecretTypeAPIKey,
			Value:       "AKIA1234567890123456",
			Line:        1,
			Column:      10,
			Context:     "api_key: ***",
			Confidence:  0.9,
			Description: "AWS Access Key",
			RuleID:      "aws_key",
		},
	}
	
	formatted := scanner.FormatSecrets(secrets)
	
	if !strings.Contains(formatted, "AKIA1234567890123456") {
		t.Errorf("Formatted output should contain the secret value")
	}
	
	if !strings.Contains(formatted, "1 secrets found") {
		t.Errorf("Formatted output should indicate the number of secrets found")
	}
}

// TestAuditLogger tests the audit logger
func TestAuditLogger(t *testing.T) {
	logger := NewAuditLogger()
	logger.Enabled = true
	logger.LogToConsole = false
	logger.LogToFile = false
	
	// Test logging various event types
	events := []struct {
		eventType EventType
		level     EventLevel
		message   string
	}{
		{EventTypeYAMLParse, LevelInfo, "YAML parse test"},
		{EventTypeSecretScan, LevelWarning, "Secret scan test"},
		{EventTypeConfigLoad, LevelError, "Config load test"},
		{EventTypeAuthCheck, LevelCritical, "Auth check test"},
	}
	
	for _, event := range events {
		err := logger.LogEvent(SecurityEvent{
			EventType: event.eventType,
			Level:    event.level,
			Message:  event.message,
			Details:  map[string]interface{}{"test": true},
		})
		
		if err != nil {
			t.Errorf("Failed to log event: %v", err)
		}
	}
}

// TestAuditLoggerSpecificMethods tests specific audit logging methods
func TestAuditLoggerSpecificMethods(t *testing.T) {
	logger := NewAuditLogger()
	logger.Enabled = true
	logger.LogToConsole = false
	logger.LogToFile = false
	
	// Test YAML parsing log
	err := logger.LogYAMLParse("test.yaml", true, 100*time.Millisecond, map[string]interface{}{"size": 1024})
	if err != nil {
		t.Errorf("Failed to log YAML parse event: %v", err)
	}
	
	// Test secret scan log
	err = logger.LogSecretScan("test.yaml", 1, 200*time.Millisecond, map[string]interface{}{"secrets": []string{"test"}})
	if err != nil {
		t.Errorf("Failed to log secret scan event: %v", err)
	}
	
	// Test config load log
	err = logger.LogConfigLoad("config.yml", true, 50*time.Millisecond, map[string]interface{}{"format": "yaml"})
	if err != nil {
		t.Errorf("Failed to log config load event: %v", err)
	}
	
	// Test auth check log
	err = logger.LogAuthCheck("user1", "resource1", "read", true, map[string]interface{}{"ip": "127.0.0.1"})
	if err != nil {
		t.Errorf("Failed to log auth check event: %v", err)
	}
}

// TestAuditLoggerStats tests audit logger statistics
func TestAuditLoggerStats(t *testing.T) {
	logger := NewAuditLogger()
	logger.Enabled = true
	logger.LogToConsole = false
	logger.LogToFile = false
	
	// Log some events
	logger.LogInfo("Test message 1", map[string]interface{}{"test": true})
	logger.LogError("Test error", map[string]interface{}{"test": true})
	
	stats := logger.GetStats()
	
	if stats["total_events"] != 2 {
		t.Errorf("Expected 2 total events, got %v", stats["total_events"])
	}
	
	if stats["enabled"] != true {
		t.Errorf("Expected enabled to be true, got %v", stats["enabled"])
	}
}

// TestSecurityConfiguration tests the security configuration
func TestSecurityConfiguration(t *testing.T) {
	config := DefaultSecurityConfiguration()
	
	// Test default values
	if !config.YAML.SandboxParser {
		t.Error("Expected sandbox parser to be enabled by default")
	}
	
	if !config.Secrets.Enabled {
		t.Error("Expected secret scanning to be enabled by default")
	}
	
	if !config.Audit.Enabled {
		t.Error("Expected audit logging to be enabled by default")
	}
	
	// Test validation
	err := config.Validate()
	if err != nil {
		t.Errorf("Configuration validation failed: %v", err)
	}
	
	// Test invalid configuration
	config.YAML.MaxNestingDepth = 0
	err = config.Validate()
	if err == nil {
		t.Error("Expected validation to fail with invalid nesting depth")
	}
}

// TestSecurityConfigurationLoadSave tests loading and saving configuration
func TestSecurityConfigurationLoadSave(t *testing.T) {
	// Create test configuration
	config := DefaultSecurityConfiguration()
	config.YAML.ValidationMode = "strict"
	config.Secrets.SeverityThreshold = 0.8
	
	// Test saving to temporary file
	tempFile := "test_security_config.yml"
	err := config.SaveSecurityConfiguration(tempFile)
	if err != nil {
		t.Errorf("Failed to save configuration: %v", err)
	}
	defer os.Remove(tempFile)
	
	// Test loading from file
	loadedConfig, err := LoadSecurityConfiguration(tempFile)
	if err != nil {
		t.Errorf("Failed to load configuration: %v", err)
	}
	
	if loadedConfig.YAML.ValidationMode != "strict" {
		t.Errorf("Expected validation mode 'strict', got %s", loadedConfig.YAML.ValidationMode)
	}
	
	if loadedConfig.Secrets.SeverityThreshold != 0.8 {
		t.Errorf("Expected severity threshold 0.8, got %f", loadedConfig.Secrets.SeverityThreshold)
	}
}

// TestSecurityConfigurationValidation tests configuration validation
func TestSecurityConfigurationValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      func() *SecurityConfiguration
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: func() *SecurityConfiguration {
				return DefaultSecurityConfiguration()
			},
			expectError: false,
		},
		{
			name: "Invalid nesting depth",
			config: func() *SecurityConfiguration {
				config := DefaultSecurityConfiguration()
				config.YAML.MaxNestingDepth = 0
				return config
			},
			expectError: true,
		},
		{
			name: "Invalid severity threshold",
			config: func() *SecurityConfiguration {
				config := DefaultSecurityConfiguration()
				config.Secrets.SeverityThreshold = 1.5
				return config
			},
			expectError: true,
		},
		{
			name: "Invalid log level",
			config: func() *SecurityConfiguration {
				config := DefaultSecurityConfiguration()
				config.Audit.LogLevel = "invalid"
				return config
			},
			expectError: true,
		},
		{
			name: "Invalid token length",
			config: func() *SecurityConfiguration {
				config := DefaultSecurityConfiguration()
				config.Authentication.Token.MaxLength = 5
				config.Authentication.Token.MinLength = 10
				return config
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config()
			err := config.Validate()
			
			if tt.expectError && err == nil {
				t.Errorf("Expected validation error, got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

// TestSecurityLinter tests the security linter
func TestSecurityLinter(t *testing.T) {
	// Create test Go file with security issues
	testCode := `
package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	// Hardcoded credential (should be detected)
	password := "supersecretpassword123"
	
	// Insecure crypto usage (should be detected)
	hash := md5.Sum([]byte("data"))
	
	// SQL injection risk (should be detected)
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userId)
	
	// Path traversal risk (should be detected)
	file := os.Open("/tmp/" + userInput)
	
	// Insecure random (should be detected)
	token := rand.Intn(1000)
	
	fmt.Println("Running with potential security issues")
}
`

	// Write test file
	tempFile := "test_security.go"
	err := os.WriteFile(tempFile, []byte(testCode), 0644)
	if err != nil {
		t.Errorf("Failed to write test file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test linting
	linter := NewSecurityLinter()
	issues, err := linter.LintFile(tempFile)
	if err != nil {
		t.Errorf("Failed to lint file: %v", err)
	}

	// Should detect security issues
	if len(issues) == 0 {
		t.Error("Expected to detect security issues, but found none")
	}

	// Test formatting
	formatted := linter.FormatIssues(issues)
	if formatted == "No security issues found" {
		t.Error("Formatted output should show issues")
	}

	// Test statistics
	stats := linter.GetStats(issues)
	if stats["total_issues"] != len(issues) {
		t.Errorf("Statistics total should match issue count")
	}
}

// TestSecurityLinterDirectory tests linting a directory
func TestSecurityLinterDirectory(t *testing.T) {
	// Create temporary directory with test files
	tempDir := "test_security_dir"
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Errorf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test Go file
	testCode := `
package main

func main() {
	// Hardcoded credential
	api_key := "AKIA1234567890123456"
}
`
	testFile := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Errorf("Failed to write test file: %v", err)
	}

	// Test directory linting
	linter := NewSecurityLinter()
	issues, err := linter.LintDirectory(tempDir)
	if err != nil {
		t.Errorf("Failed to lint directory: %v", err)
	}

	// Should detect security issues
	if len(issues) == 0 {
		t.Error("Expected to detect security issues, but found none")
	}
}

// Helper functions

func isErrorOfType(err, expectedType error) bool {
	if _, ok := err.(*YAMLValidationError); ok {
		return true
	}
	return false
}

// Benchmark functions

func BenchmarkSandboxYAMLParser(b *testing.B) {
	parser := NewSandboxYAMLParser()
	content := []byte("name: test\nvalue: 123\nnested:\n  item1: abc\n  item2: 456")
	
	var result map[string]interface{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseYAML(content, &result)
	}
}

func BenchmarkSecretScanner(b *testing.B) {
	scanner := NewSecretScanner()
	content := []byte("api_key: AKIA1234567890123456\ntoken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\npassword: secret123")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scanner.ScanYAML(content, "test.yml")
	}
}

func BenchmarkAuditLogger(b *testing.B) {
	logger := NewAuditLogger()
	logger.Enabled = true
	logger.LogToConsole = false
	logger.LogToFile = false
	
	event := SecurityEvent{
		EventType: EventTypeYAMLParse,
		Level:     LevelInfo,
		Message:   "Benchmark test",
		Details:   map[string]interface{}{"benchmark": true},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogEvent(event)
	}
}

func BenchmarkSecurityConfiguration(b *testing.B) {
	config := DefaultSecurityConfiguration()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}
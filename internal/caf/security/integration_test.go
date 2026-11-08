package security

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestSecurityManagerIntegration tests the complete security manager workflow
func TestSecurityManagerIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	
	// Create test configuration
	configPath := filepath.Join(tempDir, "security.yml")
	config := DefaultSecurityConfiguration()
	config.YAML.SandboxParser = true
	config.Secrets.Enabled = true
	config.Audit.Enabled = true
	config.Audit.LogToConsole = false
	config.Audit.LogToFile = false
	
	// Save configuration
	err := config.SaveSecurityConfiguration(configPath)
	if err != nil {
		t.Fatalf("Failed to save test configuration: %v", err)
	}
	
	// Create security manager
	manager := NewSecurityManager()
	
	// Initialize security manager
	err = manager.Initialize(configPath)
	if err != nil {
		t.Fatalf("Failed to initialize security manager: %v", err)
	}
	defer manager.Shutdown()
	
	// Test 1: YAML parsing with security checks
	t.Run("YAML Parsing", func(t *testing.T) {
		validYAML := []byte(`
name: test_application
version: 1.0
config:
  database_url: postgresql://user:pass@localhost:5432/db
  api_key: AKIA1234567890123456
  features:
    - feature1
    - feature2
`)
		
		var result map[string]interface{}
		err = manager.ParseYAML(validYAML, &result, "test.yaml")
		if err != nil {
			t.Errorf("Failed to parse valid YAML: %v", err)
		}
		
		// Test invalid YAML
		invalidYAML := []byte(`name: !script alert('xss')`)
		err = manager.ParseYAML(invalidYAML, &result, "invalid.yaml")
		if err == nil {
			t.Error("Expected parsing to fail with invalid YAML")
		}
	})
	
	// Test 2: Secret scanning
	t.Run("Secret Scanning", func(t *testing.T) {
		testContent := []byte(`
# Configuration file
database_url: postgresql://admin:supersecret@localhost:5432/mydb
api_key: AKIAIOSFODNN7EXAMPLE
jwt_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
password: mysecretpassword123
`)
		
		result, err := manager.ScanForSecrets(testContent, "test_config.yml")
		if err != nil {
			t.Errorf("Failed to scan for secrets: %v", err)
		}
		
		if !result.HasSecrets {
			t.Error("Expected to find secrets in test content")
		}
		
		if len(result.Secrets) == 0 {
			t.Error("Expected to find at least one secret")
		}
		
		// Verify specific secrets were found
		foundSecrets := 0
		for _, secret := range result.Secrets {
			if secret.Type == SecretTypeDatabaseURL || secret.Type == SecretTypeAPIKey || 
			   secret.Type == SecretTypeToken || secret.Type == SecretTypePassword {
				foundSecrets++
			}
		}
		
		if foundSecrets < 4 {
			t.Errorf("Expected to find at least 4 secrets, found %d", foundSecrets)
		}
	})
	
	// Test 3: Security linting
	t.Run("Security Linting", func(t *testing.T) {
		// Create a test Go file with security issues
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
	// hash := md5.Sum([]byte("data"))
	
	// SQL injection risk (should be detected)
	// query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userId)
	
	// Path traversal risk (should be detected)
	// file := os.Open("/tmp/" + userInput)
	
	// Insecure random (should be detected)
	token := rand.Intn(1000)
	
	fmt.Println("Running with potential security issues")
}
`
		
		testFile := filepath.Join(tempDir, "test_security.go")
		err := os.WriteFile(testFile, []byte(testCode), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		
		issues, err := manager.LintFile(testFile)
		if err != nil {
			t.Errorf("Failed to lint file: %v", err)
		}
		
		if len(issues) == 0 {
			t.Error("Expected to find security issues in test file")
		}
		
		// Test directory linting
		subDir := filepath.Join(tempDir, "subdir")
		err = os.Mkdir(subDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}
		
		subTestFile := filepath.Join(subDir, "sub_test.go")
		err = os.WriteFile(subTestFile, []byte(testCode), 0644)
		if err != nil {
			t.Fatalf("Failed to write sub test file: %v", err)
		}
		
		dirIssues, err := manager.LintDirectory(tempDir)
		if err != nil {
			t.Errorf("Failed to lint directory: %v", err)
		}
		
		if len(dirIssues) == 0 {
			t.Error("Expected to find security issues in directory lint")
		}
		
		if len(dirIssues) < len(issues) {
			t.Error("Expected directory lint to find at least as many issues as file lint")
		}
	})
	
	// Test 4: YAML file validation
	t.Run("YAML File Validation", func(t *testing.T) {
		// Create valid YAML file
		validYAML := []byte(`name: test
version: 1.0
config:
  setting1: value1
  setting2: value2
`)
		
		validFile := filepath.Join(tempDir, "valid.yaml")
		err = os.WriteFile(validFile, validYAML, 0644)
		if err != nil {
			t.Fatalf("Failed to write valid YAML file: %v", err)
		}
		
		err = manager.ValidateYAMLFile(validFile)
		if err != nil {
			t.Errorf("Valid YAML file should pass validation: %v", err)
		}
		
		// Create YAML file with secrets
		secretsYAML := []byte(`name: test
version: 1.0
config:
  api_key: AKIA1234567890123456
  password: secret123
`)
		
		secretsFile := filepath.Join(tempDir, "secrets.yaml")
		err = os.WriteFile(secretsFile, secretsYAML, 0644)
		if err != nil {
			t.Fatalf("Failed to write secrets YAML file: %v", err)
		}
		
		err = manager.ValidateYAMLFile(secretsFile)
		if err == nil {
			t.Error("YAML file with secrets should fail validation")
		}
	})
	
	// Test 5: Configuration management
	t.Run("Configuration Management", func(t *testing.T) {
		// Test getting configuration
		config := manager.GetConfiguration()
		if config == nil {
			t.Error("Expected to get configuration")
		}
		
		// Test updating configuration
		err = manager.UpdateConfiguration(func(c *SecurityConfiguration) {
			c.Secrets.SeverityThreshold = 0.9
		})
		if err != nil {
			t.Errorf("Failed to update configuration: %v", err)
		}
		
		// Verify update
		updatedConfig := manager.GetConfiguration()
		if updatedConfig.Secrets.SeverityThreshold != 0.9 {
			t.Error("Configuration update not reflected")
		}
		
		// Test saving configuration
		newConfigPath := filepath.Join(tempDir, "updated_security.yml")
		err = manager.SaveConfiguration(newConfigPath)
		if err != nil {
			t.Errorf("Failed to save configuration: %v", err)
		}
		
		// Test loading saved configuration
		newManager := NewSecurityManager()
		err = newManager.Initialize(newConfigPath)
		if err != nil {
			t.Errorf("Failed to initialize with saved configuration: %v", err)
		}
		
		loadedConfig := newManager.GetConfiguration()
		if loadedConfig.Secrets.SeverityThreshold != 0.9 {
			t.Error("Loaded configuration should reflect updates")
		}
		
		newManager.Shutdown()
	})
	
	// Test 6: Security analysis
	t.Run("Security Analysis", func(t *testing.T) {
		// Create test file
		testContent := []byte(`api_key: AKIA1234567890123456
password: secret123
`)
		testFile := filepath.Join(tempDir, "analysis_test.yml")
		err = os.WriteFile(testFile, testContent, 0644)
		if err != nil {
			t.Fatalf("Failed to write analysis test file: %v", err)
		}
		
		// Test file analysis
		result, err := manager.LintSecurity(testFile, false)
		if err != nil {
			t.Errorf("Failed to analyze file: %v", err)
		}
		
		if len(result.Issues) == 0 {
			t.Error("Expected to find security issues")
		}
		
		if !result.HasHighSeverity() {
			t.Error("Expected to find high severity issues")
		}
		
		// Test directory analysis
		result, err = manager.LintSecurity(tempDir, true)
		if err != nil {
			t.Errorf("Failed to analyze directory: %v", err)
		}
		
		if len(result.Issues) == 0 {
			t.Error("Expected to find security issues in directory analysis")
		}
		
		// Test result formatting
		formatted := result.FormatResult()
		if formatted == "" {
			t.Error("Formatted result should not be empty")
		}
		
		if !strings.Contains(formatted, "Security Analysis Report") {
			t.Error("Formatted result should contain header")
		}
		
		// Test critical issue detection
		criticalIssues := result.GetCriticalIssues()
		if len(criticalIssues) == 0 {
			t.Error("Expected to find critical issues")
		}
	})
	
	// Test 7: Statistics
	t.Run("Statistics", func(t *testing.T) {
		stats := manager.GetStats()
		if stats == nil {
			t.Error("Expected to get statistics")
		}
		
		if stats["initialized"] != true {
			t.Error("Expected initialized to be true")
		}
		
		if stats["config_path"] != configPath {
			t.Error("Expected config path to match")
		}
	})
}

// TestSecurityManagerErrorHandling tests error handling in security manager
func TestSecurityManagerErrorHandling(t *testing.T) {
	// Create security manager without initialization
	manager := NewSecurityManager()
	
	// Test operations on uninitialized manager
	var target interface{}
	err := manager.ParseYAML([]byte("test: value"), &target, "test.yaml")
	if err == nil {
		t.Error("Expected error on uninitialized manager")
	}
	
	_, err = manager.ScanForSecrets([]byte("test: value"), "test.yaml")
	if err == nil {
		t.Error("Expected error on uninitialized manager")
	}
	
	_, err = manager.LintFile("test.go")
	if err == nil {
		t.Error("Expected error on uninitialized manager")
	}
	
	// Test initialization with invalid configuration
	invalidConfig := filepath.Join(t.TempDir(), "invalid_config.yml")
	err = os.WriteFile(invalidConfig, []byte("invalid: config"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}
	
	err = manager.Initialize(invalidConfig)
	if err == nil {
		t.Error("Expected error with invalid configuration")
	}
}

// TestSecurityManagerConcurrentAccess tests concurrent access to security manager
func TestSecurityManagerConcurrentAccess(t *testing.T) {
	// Create temporary directory and configuration
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "security.yml")
	config := DefaultSecurityConfiguration()
	config.Audit.LogToConsole = false
	config.Audit.LogToFile = false
	
	err := config.SaveSecurityConfiguration(configPath)
	if err != nil {
		t.Fatalf("Failed to save test configuration: %v", err)
	}
	
	// Initialize security manager
	manager := NewSecurityManager()
	err = manager.Initialize(configPath)
	if err != nil {
		t.Fatalf("Failed to initialize security manager: %v", err)
	}
	defer manager.Shutdown()
	
	// Test concurrent operations
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			// Perform various operations
			yamlContent := []byte(fmt.Sprintf(`name: test_%d
value: %d
`, id, id))
			
			var target interface{}
			err := manager.ParseYAML(yamlContent, &target, fmt.Sprintf("test_%d.yaml", id))
			if err != nil {
				t.Errorf("Goroutine %d: YAML parsing failed: %v", id, err)
				return
			}
			
			result, err := manager.ScanForSecrets(yamlContent, fmt.Sprintf("test_%d.yaml", id))
			if err != nil {
				t.Errorf("Goroutine %d: Secret scanning failed: %v", id, err)
				return
			}
			
			if result.HasSecrets {
				t.Errorf("Goroutine %d: Unexpected secrets found", id)
			}
			
			issues, err := manager.LintFile(fmt.Sprintf("test_%d.go", id))
			if err != nil {
				t.Errorf("Goroutine %d: Linting failed: %v", id, err)
				return
			}
			
			// Should not find issues in simple test code
			if len(issues) > 0 {
				t.Errorf("Goroutine %d: Unexpected lint issues found", id)
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// Goroutine completed
		case <-time.After(10 * time.Second):
			t.Error("Timeout waiting for goroutines to complete")
			return
		}
	}
	
	// Verify manager is still in valid state
	stats := manager.GetStats()
	if stats["initialized"] != true {
		t.Error("Security manager should remain initialized after concurrent operations")
	}
}

// BenchmarkSecurityManager benchmarks the security manager
func BenchmarkSecurityManager(b *testing.B) {
	// Create test configuration
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "security.yml")
	config := DefaultSecurityConfiguration()
	config.Audit.LogToConsole = false
	config.Audit.LogToFile = false
	
	err := config.SaveSecurityConfiguration(configPath)
	if err != nil {
		b.Fatalf("Failed to save test configuration: %v", err)
	}
	
	// Initialize security manager
	manager := NewSecurityManager()
	err = manager.Initialize(configPath)
	if err != nil {
		b.Fatalf("Failed to initialize security manager: %v", err)
	}
	defer manager.Shutdown()
	
	// Test data
	yamlContent := []byte(`
name: test_application
version: 1.0
config:
  database_url: postgresql://user:pass@localhost:5432/db
  api_key: AKIA1234567890123456
  features:
    - feature1
    - feature2
`)
	
	b.ResetTimer()
	
	// Benchmark YAML parsing
	b.Run("YAML Parsing", func(b *testing.B) {
		var target interface{}
		for i := 0; i < b.N; i++ {
			manager.ParseYAML(yamlContent, &target, "benchmark.yaml")
		}
	})
	
	// Benchmark secret scanning
	b.Run("Secret Scanning", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.ScanForSecrets(yamlContent, "benchmark.yaml")
		}
	})
	
	// Benchmark security linting
	testCode := `
package main

func main() {
	// Test code
}
`
	testFile := filepath.Join(tempDir, "benchmark.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		b.Fatalf("Failed to write benchmark test file: %v", err)
	}
	
	b.Run("Security Linting", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.LintFile(testFile)
		}
	})
	
	// Benchmark configuration operations
	b.Run("Configuration Operations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.GetConfiguration()
		}
	})
}

// TestSecurityManagerLifecycle tests the complete lifecycle of security manager
func TestSecurityManagerLifecycle(t *testing.T) {
	// Create temporary directory and configuration
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "security.yml")
	config := DefaultSecurityConfiguration()
	config.Audit.LogToConsole = false
	config.Audit.LogToFile = false
	
	err := config.SaveSecurityConfiguration(configPath)
	if err != nil {
		t.Fatalf("Failed to save test configuration: %v", err)
	}
	
	// Test 1: Create and initialize
	manager := NewSecurityManager()
	err = manager.Initialize(configPath)
	if err != nil {
		t.Fatalf("Failed to initialize security manager: %v", err)
	}
	
	// Test 2: Verify initialization
	stats := manager.GetStats()
	if stats["initialized"] != true {
		t.Error("Security manager should be initialized")
	}
	
	// Test 3: Perform operations
	yamlContent := []byte(`name: test
value: 123`)
	
	var target interface{}
	err = manager.ParseYAML(yamlContent, &target, "test.yaml")
	if err != nil {
		t.Errorf("Failed to parse YAML: %v", err)
	}
	
	// Test 4: Update configuration
	err = manager.UpdateConfiguration(func(c *SecurityConfiguration) {
		c.Secrets.SeverityThreshold = 0.8
	})
	if err != nil {
		t.Errorf("Failed to update configuration: %v", err)
	}
	
	// Test 5: Save updated configuration
	newConfigPath := filepath.Join(tempDir, "updated_security.yml")
	err = manager.SaveConfiguration(newConfigPath)
	if err != nil {
		t.Errorf("Failed to save configuration: %v", err)
	}
	
	// Test 6: Shutdown
	err = manager.Shutdown()
	if err != nil {
		t.Errorf("Failed to shutdown security manager: %v", err)
	}
	
	// Test 7: Verify shutdown
	stats = manager.GetStats()
	if stats["initialized"] != false {
		t.Error("Security manager should be shutdown")
	}
	
	// Test 8: Reinitialize
	newManager := NewSecurityManager()
	err = newManager.Initialize(newConfigPath)
	if err != nil {
		t.Errorf("Failed to reinitialize security manager: %v", err)
	}
	
	// Test 9: Verify reinitialization
	config = newManager.GetConfiguration()
	if config.Secrets.SeverityThreshold != 0.8 {
		t.Error("Reinitialized manager should have updated configuration")
	}
	
	// Test 10: Final shutdown
	err = newManager.Shutdown()
	if err != nil {
		t.Errorf("Failed to final shutdown: %v", err)
	}
}
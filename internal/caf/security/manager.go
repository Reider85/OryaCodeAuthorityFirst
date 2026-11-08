package security

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SecurityManager coordinates all security components
type SecurityManager struct {
	// Components
	sandboxParser *SandboxYAMLParser
	secretScanner *SecretScanner
	auditLogger   *AuditLogger
	linter        *SecurityLinter
	config        *SecurityConfiguration
	
	// State
	mu            sync.RWMutex
	initialized   bool
	lastConfig    *SecurityConfiguration
	configPath    string
}

// NewSecurityManager creates a new security manager
func NewSecurityManager() *SecurityManager {
	return &SecurityManager{
		sandboxParser: NewSandboxYAMLParser(),
		secretScanner: NewSecretScanner(),
		auditLogger:   NewAuditLogger(),
		linter:        NewSecurityLinter(),
		config:        DefaultSecurityConfiguration(),
	}
}

// Initialize initializes the security manager with configuration
func (sm *SecurityManager) Initialize(configPath string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Load configuration
	if configPath != "" {
		config, err := LoadSecurityConfiguration(configPath)
		if err != nil {
			return fmt.Errorf("failed to load security configuration: %w", err)
		}
		sm.config = config
		sm.configPath = configPath
	} else {
		// Use default configuration
		sm.configPath = GetSecurityConfigPath()
	}

	// Configure components based on settings
	sm.configureComponents()

	// Initialize audit logger
	if sm.config.Audit.Enabled {
		sm.auditLogger.Enabled = true
		sm.auditLogger.LogToConsole = sm.config.Audit.LogToConsole
		sm.auditLogger.LogToFile = sm.config.Audit.LogToFile
		sm.auditLogger.LogLevel = sm.config.GetAuditLogLevel()
		sm.auditLogger.MaxFileSize = sm.config.Audit.MaxFileSize
		sm.auditLogger.MaxFiles = sm.config.Audit.MaxFiles
		sm.auditLogger.RetainDays = sm.config.Audit.RetainDays
		sm.auditLogger.OutputDir = sm.config.Audit.OutputDir
		sm.auditLogger.EnableJSON = sm.config.Audit.EnableJSON
		sm.auditLogger.EnableStructured = sm.config.Audit.EnableStructured
	}

	// Configure secret scanner
	if sm.config.Secrets.Enabled {
		sm.secretScanner.Enabled = true
		sm.secretScanner.MaxLineLength = sm.config.Secrets.MaxLineLength
		sm.secretScanner.MinSecretLength = sm.config.Secrets.MinSecretLength
		sm.secretScanner.MaxSecretLength = sm.config.Secrets.MaxSecretLength
		sm.secretScanner.SeverityThreshold = sm.config.Secrets.SeverityThreshold
		sm.secretScanner.EnableHeuristics = sm.config.Secrets.EnableHeuristics
		sm.secretScanner.EnableEntropyCheck = sm.config.Secrets.EnableEntropyCheck
		sm.secretScanner.EnableContextCheck = sm.config.Secrets.EnableContextCheck
		sm.secretScanner.IgnoreFiles = sm.config.Secrets.IgnoreFiles
		
		// Add custom rules
		for _, rule := range sm.config.Secrets.CustomRules {
			sm.secretScanner.AddRule(rule)
		}
	}

	// Configure sandbox parser
	sm.sandboxParser.AllowCustomTags = sm.config.YAML.AllowCustomTags
	sm.sandboxParser.MaxNestingDepth = sm.config.YAML.MaxNestingDepth
	sm.sandboxParser.MaxDocumentSize = sm.config.YAML.MaxDocumentSize
	sm.sandboxParser.ValidationMode = sm.config.GetYAMLValidationMode()
	sm.sandboxParser.RequiredFields = sm.config.YAML.RequiredFields
	
	// Configure blocked patterns
	for _, pattern := range sm.config.YAML.BlockedPatterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			sm.sandboxParser.BlockedPatterns = append(sm.sandboxParser.BlockedPatterns, re)
		}
	}

	sm.initialized = true
	sm.lastConfig = sm.config

	// Log initialization
	sm.auditLogger.LogInfo("Security manager initialized", map[string]interface{}{
		"config_path": configPath,
		"features": map[string]bool{
			"yaml_sandbox":   sm.config.YAML.SandboxParser,
			"secret_scanning": sm.config.Secrets.Enabled,
			"audit_logging":   sm.config.Audit.Enabled,
			"security_linter": true,
		},
	})

	return nil
}

// configureComponents configures security components based on settings
func (sm *SecurityManager) configureComponents() {
	// Configure sandbox parser
	sm.sandboxParser.AllowCustomTags = sm.config.YAML.AllowCustomTags
	sm.sandboxParser.MaxNestingDepth = sm.config.YAML.MaxNestingDepth
	sm.sandboxParser.MaxDocumentSize = sm.config.YAML.MaxDocumentSize
	sm.sandboxParser.ValidationMode = sm.config.GetYAMLValidationMode()
	sm.sandboxParser.RequiredFields = sm.config.YAML.RequiredFields
	
	// Configure blocked patterns
	sm.sandboxParser.BlockedPatterns = []*regexp.Regexp{}
	for _, pattern := range sm.config.YAML.BlockedPatterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			sm.sandboxParser.BlockedPatterns = append(sm.sandboxParser.BlockedPatterns, re)
		}
	}

	// Configure secret scanner
	sm.secretScanner.Enabled = sm.config.Secrets.Enabled
	sm.secretScanner.MaxLineLength = sm.config.Secrets.MaxLineLength
	sm.secretScanner.MinSecretLength = sm.config.Secrets.MinSecretLength
	sm.secretScanner.MaxSecretLength = sm.config.Secrets.MaxSecretLength
	sm.secretScanner.SeverityThreshold = sm.config.Secrets.SeverityThreshold
	sm.secretScanner.EnableHeuristics = sm.config.Secrets.EnableHeuristics
	sm.secretScanner.EnableEntropyCheck = sm.config.Secrets.EnableEntropyCheck
	sm.secretScanner.EnableContextCheck = sm.config.Secrets.EnableContextCheck
	sm.secretScanner.IgnoreFiles = sm.config.Secrets.IgnoreFiles
	
	// Add custom rules
	for _, rule := range sm.config.Secrets.CustomRules {
		sm.secretScanner.AddRule(rule)
	}

	// Configure audit logger
	sm.auditLogger.Enabled = sm.config.Audit.Enabled
	sm.auditLogger.LogToConsole = sm.config.Audit.LogToConsole
	sm.auditLogger.LogToFile = sm.config.Audit.LogToFile
	sm.auditLogger.LogLevel = sm.config.GetAuditLogLevel()
	sm.auditLogger.MaxFileSize = sm.config.Audit.MaxFileSize
	sm.auditLogger.MaxFiles = sm.config.Audit.MaxFiles
	sm.auditLogger.RetainDays = sm.config.Audit.RetainDays
	sm.auditLogger.OutputDir = sm.config.Audit.OutputDir
	sm.auditLogger.EnableJSON = sm.config.Audit.EnableJSON
	sm.auditLogger.EnableStructured = sm.config.Audit.EnableStructured
}

// ParseYAML parses YAML content with security checks
func (sm *SecurityManager) ParseYAML(content []byte, target interface{}, filePath string) error {
	if !sm.initialized {
		return fmt.Errorf("security manager not initialized")
	}

	startTime := time.Now()
	
	// Log YAML parsing attempt
	sm.auditLogger.LogYAMLParse(filePath, true, 0, map[string]interface{}{"size": len(content)})
	
	// Parse with sandboxed parser
	err := sm.sandboxParser.ParseYAML(content, target)
	
	// Log result
	duration := time.Since(startTime)
	if err != nil {
		sm.auditLogger.LogError("YAML parsing failed", map[string]interface{}{
			"file_path": filePath,
			"error":     err.Error(),
			"duration":  duration.Milliseconds(),
		})
	} else {
		sm.auditLogger.LogYAMLParse(filePath, true, duration, map[string]interface{}{
			"size": len(content),
		})
	}
	
	return err
}

// ScanForSecrets scans content for secrets
func (sm *SecurityManager) ScanForSecrets(content []byte, filePath string) (*ScanResult, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("security manager not initialized")
	}

	if !sm.config.Secrets.Enabled {
		return &ScanResult{FilePath: filePath, HasSecrets: false}, nil
	}

	startTime := time.Now()
	result := sm.secretScanner.ScanYAML(content, filePath)
	
	duration := time.Since(startTime)
	sm.auditLogger.LogSecretScan(filePath, len(result.Secrets), duration, map[string]interface{}{
		"secrets_found": len(result.Secrets),
		"scan_duration": duration.Milliseconds(),
	})

	return result, nil
}

// LintFile performs security linting on a file
func (sm *SecurityManager) LintFile(filePath string) ([]SecurityIssue, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("security manager not initialized")
	}

	return sm.linter.LintFile(filePath)
}

// LintDirectory performs security linting on a directory
func (sm *SecurityManager) LintDirectory(dirPath string) ([]SecurityIssue, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("security manager not initialized")
	}

	return sm.linter.LintDirectory(dirPath)
}

// LoadConfiguration loads security configuration from file
func (sm *SecurityManager) LoadConfiguration(configPath string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	config, err := LoadSecurityConfiguration(configPath)
	if err != nil {
		return fmt.Errorf("failed to load security configuration: %w", err)
	}

	sm.config = config
	sm.configPath = configPath
	sm.configureComponents()

	// Log configuration change
	sm.auditLogger.LogInfo("Security configuration updated", map[string]interface{}{
		"config_path": configPath,
		"features": map[string]bool{
			"yaml_sandbox":   config.YAML.SandboxParser,
			"secret_scanning": config.Secrets.Enabled,
			"audit_logging":   config.Audit.Enabled,
		},
	})

	return nil
}

// SaveConfiguration saves current configuration to file
func (sm *SecurityManager) SaveConfiguration(configPath string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	err := sm.config.SaveSecurityConfiguration(configPath)
	if err != nil {
		return fmt.Errorf("failed to save security configuration: %w", err)
	}

	sm.configPath = configPath

	// Log configuration save
	sm.auditLogger.LogInfo("Security configuration saved", map[string]interface{}{
		"config_path": configPath,
	})

	return nil
}

// GetConfiguration returns current security configuration
func (sm *SecurityManager) GetConfiguration() *SecurityConfiguration {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.config
}

// UpdateConfiguration updates security configuration
func (sm *SecurityManager) UpdateConfiguration(updater func(*SecurityConfiguration)) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.config == nil {
		return fmt.Errorf("no configuration to update")
	}

	// Create a copy of the configuration
	newConfig := *sm.config
	updater(&newConfig)

	// Validate the new configuration
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Update configuration
	sm.config = &newConfig
	sm.configureComponents()

	// Log configuration update
	sm.auditLogger.LogInfo("Security configuration updated", map[string]interface{}{
		"config_path": sm.configPath,
	})

	return nil
}

// GetStats returns statistics about security operations
func (sm *SecurityManager) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"initialized":    sm.initialized,
		"config_path":    sm.configPath,
		"last_updated":  sm.lastConfig,
	}

	// Add component statistics
	if sm.initialized {
		stats["audit_logger"] = sm.auditLogger.GetStats()
		stats["security_config"] = sm.config
	}

	return stats
}

// Shutdown shuts down the security manager
func (sm *SecurityManager) Shutdown() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.initialized {
		// Shutdown audit logger
		if err := sm.auditLogger.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown audit logger: %w", err)
		}

		sm.initialized = false

		// Log shutdown
		sm.auditLogger.LogInfo("Security manager shutdown", map[string]interface{}{
			"shutdown_time": time.Now(),
		})
	}

	return nil
}

// ValidateYAMLFile validates a YAML file with security checks
func (sm *SecurityManager) ValidateYAMLFile(filePath string) error {
	if !sm.initialized {
		return fmt.Errorf("security manager not initialized")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Parse YAML with security checks
	var target interface{}
	if err := sm.ParseYAML(content, &target, filePath); err != nil {
		return fmt.Errorf("YAML validation failed: %w", err)
	}

	// Scan for secrets
	result, err := sm.ScanForSecrets(content, filePath)
	if err != nil {
		return fmt.Errorf("secret scanning failed: %w", err)
	}

	// Check if secrets were found
	if result.HasSecrets {
		return fmt.Errorf("YAML file contains %d potential secrets", len(result.Secrets))
	}

	return nil
}

// ScanFileForSecrets scans a file for secrets
func (sm *SecurityManager) ScanFileForSecrets(filePath string) (*ScanResult, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("security manager not initialized")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return sm.ScanForSecrets(content, filePath)
}

// LintSecurity performs comprehensive security analysis
func (sm *SecurityManager) LintSecurity(targetPath string, isDirectory bool) (*SecurityLintResult, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("security manager not initialized")
	}

	result := &SecurityLintResult{
		FilePath:  targetPath,
		Timestamp: time.Now(),
		Issues:    []SecurityIssue{},
		Summary:   map[string]int{},
	}

	var issues []SecurityIssue
	var err error

	if isDirectory {
		issues, err = sm.LintDirectory(targetPath)
	} else {
		issues, err = sm.LintFile(targetPath)
	}

	if err != nil {
		return nil, err
	}

	result.Issues = issues

	// Build summary
	for _, issue := range issues {
		result.Summary[issue.Severity]++
		result.Summary[issue.Category]++
	}

	// Log the result
	sm.auditLogger.LogInfo("Security linting completed", map[string]interface{}{
		"target_path": targetPath,
		"is_directory": isDirectory,
		"issues_found": len(issues),
		"severity_breakdown": result.Summary,
	})

	return result, nil
}

// SecurityLintResult represents the result of security linting
type SecurityLintResult struct {
	FilePath  string         `json:"file_path"`
	Timestamp time.Time      `json:"timestamp"`
	Issues    []SecurityIssue `json:"issues"`
	Summary   map[string]int `json:"summary"`
}

// GetIssuesBySeverity returns issues filtered by severity
func (r *SecurityLintResult) GetIssuesBySeverity(severity string) []SecurityIssue {
	var filtered []SecurityIssue
	for _, issue := range r.Issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// GetIssuesByCategory returns issues filtered by category
func (r *SecurityLintResult) GetIssuesByCategory(category string) []SecurityIssue {
	var filtered []SecurityIssue
	for _, issue := range r.Issues {
		if issue.Category == category {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// FormatResult formats the linting result for display
func (r *SecurityLintResult) FormatResult() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Security Analysis Report\n"))
	builder.WriteString(fmt.Sprintf("======================\n"))
	builder.WriteString(fmt.Sprintf("Target: %s\n", r.FilePath))
	builder.WriteString(fmt.Sprintf("Timestamp: %s\n", r.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("Total Issues: %d\n\n", len(r.Issues)))

	// Severity breakdown
	builder.WriteString("Severity Breakdown:\n")
	for severity, count := range r.Summary {
		if strings.Contains(severity, "low") || strings.Contains(severity, "medium") || 
		   strings.Contains(severity, "high") || strings.Contains(severity, "critical") {
			builder.WriteString(fmt.Sprintf("  %s: %d\n", severity, count))
		}
	}

	// Category breakdown
	builder.WriteString("\nCategory Breakdown:\n")
	for category, count := range r.Summary {
		if category != "low" && category != "medium" && category != "high" && category != "critical" {
			builder.WriteString(fmt.Sprintf("  %s: %d\n", category, count))
		}
	}

	// Detailed issues
	if len(r.Issues) > 0 {
		builder.WriteString("\nDetailed Issues:\n")
		for i, issue := range r.Issues {
			builder.WriteString(fmt.Sprintf("\n%d. %s\n", i+1, issue.Message))
			builder.WriteString(fmt.Sprintf("   Severity: %s\n", issue.Severity))
			builder.WriteString(fmt.Sprintf("   Category: %s\n", issue.Category))
			builder.WriteString(fmt.Sprintf("   File: %s:%d:%d\n", issue.File, issue.Line, issue.Column))
			builder.WriteString(fmt.Sprintf("   Context: %s\n", issue.Context))
			if issue.Suggestion != "" {
				builder.WriteString(fmt.Sprintf("   Suggestion: %s\n", issue.Suggestion))
			}
		}
	}

	return builder.String()
}

// IsCritical returns true if critical issues were found
func (r *SecurityLintResult) IsCritical() bool {
	return r.Summary["critical"] > 0
}

// HasHighSeverity returns true if high severity issues were found
func (r *SecurityLintResult) HasHighSeverity() bool {
	return r.Summary["high"] > 0 || r.IsCritical()
}

// GetCriticalIssues returns all critical issues
func (r *SecurityLintResult) GetCriticalIssues() []SecurityIssue {
	return r.GetIssuesBySeverity("critical")
}

// GetHighSeverityIssues returns all high severity issues
func (r *SecurityLintResult) GetHighSeverityIssues() []SecurityIssue {
	return r.GetIssuesBySeverity("high")
}
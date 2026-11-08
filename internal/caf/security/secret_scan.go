package security

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// SecretType represents the type of secret detected
type SecretType string

const (
	SecretTypeAPIKey     SecretType = "api_key"
	SecretTypeToken      SecretType = "token"
	SecretTypePassword   SecretType = "password"
	SecretTypePrivateKey SecretType = "private_key"
	SecretTypeCertificate SecretType = "certificate"
	SecretTypeDatabaseURL SecretType = "database_url"
	SecretTypeAWSKey     SecretType = "aws_key"
	SecretTypeGeneric    SecretType = "generic"
)

// Secret represents a detected secret
type Secret struct {
	Type        SecretType
	Value       string
	Line        int
	Column      int
	Context     string
	Confidence  float64
	Description string
	RuleID      string
}

// SecretScanner scans YAML content for secrets and sensitive information
type SecretScanner struct {
	// Configuration
	Enabled             bool
	MaxLineLength       int
	MinSecretLength     int
	MaxSecretLength     int
	IgnoreFiles         []string
	AllowedPatterns     []*regexp.Regexp
	BlockedPatterns     []*regexp.Regexp
	SeverityThreshold   float64
	EnableHeuristics    bool
	EnableEntropyCheck  bool
	EnableContextCheck  bool

	// Secret detection rules
	rules []SecretRule
}

// SecretRule defines a rule for detecting specific types of secrets
type SecretRule struct {
	ID          string
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Type        SecretType
	Severity    string
	Confidence  float64
	Context     []string
}

// ScanResult contains the results of a secret scan
type ScanResult struct {
	FilePath    string
	Secrets     []Secret
	TotalLines  int
	HasSecrets  bool
	Errors      []error
	Warnings    []string
}

// NewSecretScanner creates a new secret scanner with default rules
func NewSecretScanner() *SecretScanner {
	scanner := &SecretScanner{
		Enabled:           true,
		MaxLineLength:     1000,
		MinSecretLength:   8,
		MaxSecretLength:   512,
		SeverityThreshold: 0.7,
		EnableHeuristics:  true,
		EnableEntropyCheck: true,
		EnableContextCheck: true,
	}

	// Initialize default secret detection rules
	scanner.initializeDefaultRules()

	return scanner
}

// initializeDefaultRules sets up the default secret detection rules
func (s *SecretScanner) initializeDefaultRules() {
	// API Key patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "api_key_generic",
		Name:        "Generic API Key",
		Description: "Detects generic API keys",
		Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|apikey|access[_-]?key|secret[_-]?key)[\s]*[:=][\s]*["']?([a-zA-Z0-9+/]{16,})["']?`),
		Type:        SecretTypeAPIKey,
		Severity:    "high",
		Confidence:  0.8,
	})

	// Token patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "jwt_token",
		Name:        "JWT Token",
		Description: "Detects JSON Web Tokens",
		Pattern:     regexp.MustCompile(`eyJ[a-zA-Z0-9_-]+\.eyJ[a-zA-Z0-9_-]+\.?[a-zA-Z0-9_-]*`),
		Type:        SecretTypeToken,
		Severity:    "high",
		Confidence:  0.9,
	})

	// Password patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "password",
		Name:        "Password",
		Description: "Detects password fields",
		Pattern:     regexp.MustCompile(`(?i)(password|pwd|pass)[\s]*[:=][\s]*["']?([^"'\s]{8,})["']?`),
		Type:        SecretTypePassword,
		Severity:    "high",
		Confidence:  0.7,
	})

	// Private Key patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "private_key_rsa",
		Name:        "RSA Private Key",
		Description: "Detects RSA private keys",
		Pattern:     regexp.MustCompile(`-----BEGIN RSA PRIVATE KEY-----`),
		Type:        SecretTypePrivateKey,
		Severity:    "critical",
		Confidence:  0.95,
	})

	s.rules = append(s.rules, SecretRule{
		ID:          "private_key_ec",
		Name:        "EC Private Key",
		Description: "Detects Elliptic Curve private keys",
		Pattern:     regexp.MustCompile(`-----BEGIN EC PRIVATE KEY-----`),
		Type:        SecretTypePrivateKey,
		Severity:    "critical",
		Confidence:  0.95,
	})

	// Certificate patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "certificate",
		Name:        "Certificate",
		Description: "Detects certificates",
		Pattern:     regexp.MustCompile(`-----BEGIN (?:CERTIFICATE|X.509)-----`),
		Type:        SecretTypeCertificate,
		Severity:    "high",
		Confidence:  0.9,
	})

	// Database URL patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "database_url",
		Name:        "Database URL",
		Description: "Detects database connection URLs",
		Pattern:     regexp.MustCompile(`(?i)(database_url|db_url|connection_string)[\s]*[:=][\s]*["']?(postgresql://|mysql://|mongodb://|sqlite://)[^"'\s]{10,}["']?`),
		Type:        SecretTypeDatabaseURL,
		Severity:    "high",
		Confidence:  0.8,
	})

	// AWS Key patterns
	s.rules = append(s.rules, SecretRule{
		ID:          "aws_key",
		Name:        "AWS Key",
		Description: "Detects AWS access keys",
		Pattern:     regexp.MustCompile(`(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key)[\s]*[:=][\s]*["']?(AKIA[0-9A-Z]{16})["']?`),
		Type:        SecretTypeAWSKey,
		Severity:    "critical",
		Confidence:  0.95,
	})
}

// ScanYAML scans YAML content for secrets
func (s *SecretScanner) ScanYAML(content []byte, filePath string) *ScanResult {
	result := &ScanResult{
		FilePath:   filePath,
		Secrets:    []Secret{},
		TotalLines: 0,
	}

	if !s.Enabled {
		return result
	}

	// Count lines
	result.TotalLines = bytes.Count(content, []byte("\n"))

	// Parse YAML to get line information
	var yamlNode yaml.Node
	if err := yaml.Unmarshal(content, &yamlNode); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse YAML: %w", err))
		return result
	}

	// Scan line by line
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Apply line length limit
		if len(line) > s.MaxLineLength {
			continue
		}

		// Scan for secrets in this line
		secrets := s.scanLine(line, lineNum)
		result.Secrets = append(result.Secrets, secrets...)
	}

	// Check if any secrets were found
	result.HasSecrets = len(result.Secrets) > 0

	return result
}

// scanLine scans a single line for secrets
func (s *SecretScanner) scanLine(line string, lineNum int) []Secret {
	var secrets []Secret

	// Apply each rule
	for _, rule := range s.rules {
		matches := rule.Pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			secretValue := match[1]
			if !s.isValidSecret(secretValue) {
				continue
			}

			// Calculate confidence
			confidence := rule.Confidence
			if s.EnableEntropyCheck {
				confidence *= s.calculateEntropyScore(secretValue)
			}

			if s.EnableContextCheck {
				confidence *= s.calculateContextScore(line, secretValue)
			}

			// Apply severity threshold
			if confidence < s.SeverityThreshold {
				continue
			}

			// Find column position
			column := strings.Index(line, secretValue) + 1

			secret := Secret{
				Type:        rule.Type,
				Value:       secretValue,
				Line:        lineNum,
				Column:      column,
				Context:     s.extractContext(line, secretValue),
				Confidence:  confidence,
				Description: rule.Description,
				RuleID:      rule.ID,
			}

			secrets = append(secrets, secret)
		}
	}

	// Apply heuristic scanning if enabled
	if s.EnableHeuristics {
		secrets = append(secrets, s.scanHeuristics(line, lineNum)...)
	}

	return secrets
}

// isValidSecret checks if a secret value meets basic criteria
func (s *SecretScanner) isValidSecret(value string) bool {
	// Check length
	if len(value) < s.MinSecretLength || len(value) > s.MaxSecretLength {
		return false
	}

	// Check for common placeholders
	lowerValue := strings.ToLower(value)
	placeholders := []string{
		"changeme", "password", "secret", "token", "key", "test", "demo", "sample",
		"example", "placeholder", "default", "dummy", "fake", "mock",
	}

	for _, placeholder := range placeholders {
		if lowerValue == placeholder {
			return false
		}
	}

	return true
}

// calculateEntropyScore calculates the entropy of a secret value
func (s *SecretScanner) calculateEntropyScore(value string) float64 {
	entropy := 0.0
	charCounts := make(map[rune]int)

	// Count character frequencies
	for _, char := range value {
		charCounts[char]++
	}

	// Calculate entropy
	totalChars := len(value)
	for _, count := range charCounts {
		probability := float64(count) / float64(totalChars)
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}

	// Normalize entropy (max entropy for given length is log2(n))
	maxEntropy := math.Log2(float64(len(value)))
	if maxEntropy > 0 {
		return min(1.0, entropy/maxEntropy)
	}

	return 0.0
}

// calculateContextScore analyzes the context around a secret
func (s *SecretScanner) calculateContextScore(line, secret string) float64 {
	contextScore := 1.0

	// Check for context indicators
	contextIndicators := map[string]float64{
		"api":     0.9,
		"key":     0.9,
		"secret":  0.9,
		"token":   0.9,
		"auth":    0.8,
		"config":  0.7,
		"setting": 0.7,
		"env":     0.8,
		"envvar":  0.8,
	}

	lowerLine := strings.ToLower(line)
	for indicator, score := range contextIndicators {
		if strings.Contains(lowerLine, indicator) {
			contextScore *= score
		}
	}

	return contextScore
}

// extractContext extracts contextual information around a secret
func (s *SecretScanner) extractContext(line, secret string) string {
	// Extract the key-value pair containing the secret
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		return fmt.Sprintf("%s:***", key)
	}
	return strings.TrimSpace(line)
}

// scanHeuristics performs heuristic-based secret detection
func (s *SecretScanner) scanHeuristics(line string, lineNum int) []Secret {
	var secrets []Secret

	// Look for long strings of random characters
	longRandomPattern := regexp.MustCompile(`[a-zA-Z0-9+/]{32,}`)
	matches := longRandomPattern.FindAllString(line, -1)
	
	for _, match := range matches {
		if len(match) >= 32 && len(match) <= 512 {
			secrets = append(secrets, Secret{
				Type:        SecretTypeGeneric,
				Value:       match,
				Line:        lineNum,
				Column:      strings.Index(line, match) + 1,
				Context:     s.extractContext(line, match),
				Confidence:  0.6,
				Description: "Long random string (potential secret)",
				RuleID:      "heuristic_long_random",
			})
		}
	}

	return secrets
}

// AddRule adds a custom secret detection rule
func (s *SecretScanner) AddRule(rule SecretRule) {
	s.rules = append(s.rules, rule)
}

// RemoveRule removes a rule by ID
func (s *SecretScanner) RemoveRule(id string) {
	var newRules []SecretRule
	for _, rule := range s.rules {
		if rule.ID != id {
			newRules = append(newRules, rule)
		}
	}
	s.rules = newRules
}

// GetRules returns all current rules
func (s *SecretScanner) GetRules() []SecretRule {
	return s.rules
}

// ScanFile scans a file for secrets (wrapper for ScanYAML)
func (s *SecretScanner) ScanFile(filePath string, content []byte) *ScanResult {
	return s.ScanYAML(content, filePath)
}

// FilterBySeverity filters secrets by severity
func (s *SecretScanner) FilterBySeverity(secrets []Secret, minSeverity float64) []Secret {
	var filtered []Secret
	for _, secret := range secrets {
		if secret.Confidence >= minSeverity {
			filtered = append(filtered, secret)
		}
	}
	return filtered
}

// FormatSecrets formats secrets for display
func (s *SecretScanner) FormatSecrets(secrets []Secret) string {
	if len(secrets) == 0 {
		return "No secrets found"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d secret(s):\n", len(secrets)))
	
	for i, secret := range secrets {
		builder.WriteString(fmt.Sprintf("\n%d. %s (Line %d, Column %d)\n", 
			i+1, secret.Type, secret.Line, secret.Column))
		builder.WriteString(fmt.Sprintf("   Value: %s\n", secret.Value))
		builder.WriteString(fmt.Sprintf("   Confidence: %.2f\n", secret.Confidence))
		builder.WriteString(fmt.Sprintf("   Description: %s\n", secret.Description))
		builder.WriteString(fmt.Sprintf("   Context: %s\n", secret.Context))
	}
	
	return builder.String()
}

// Helper function to get min of two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// SecurityConfiguration defines the security configuration for the application
type SecurityConfiguration struct {
	// YAML Security Settings
	YAML struct {
		// Enable sandboxed YAML parser
		SandboxParser bool `yaml:"sandbox_parser" json:"sandbox_parser"`
		
		// Maximum nesting depth for YAML documents
		MaxNestingDepth int `yaml:"max_nesting_depth" json:"max_nesting_depth"`
		
		// Maximum document size in bytes
		MaxDocumentSize int64 `yaml:"max_document_size" json:"max_document_size"`
		
		// Allow custom YAML tags
		AllowCustomTags bool `yaml:"allow_custom_tags" json:"allow_custom_tags"`
		
		// Validation mode: basic, strict, paranoid
		ValidationMode string `yaml:"validation_mode" json:"validation_mode"`
		
		// Required fields in YAML documents
		RequiredFields []string `yaml:"required_fields" json:"required_fields"`
		
		// Blocked patterns in YAML content
		BlockedPatterns []string `yaml:"blocked_patterns" json:"blocked_patterns"`
	} `yaml:"yaml" json:"yaml"`

	// Secret Scanning Settings
	Secrets struct {
		// Enable secret scanning
		Enabled bool `yaml:"enabled" json:"enabled"`
		
		// Maximum line length to scan
		MaxLineLength int `yaml:"max_line_length" json:"max_line_length"`
		
		// Minimum secret length
		MinSecretLength int `yaml:"min_secret_length" json:"min_secret_length"`
		
		// Maximum secret length
		MaxSecretLength int `yaml:"max_secret_length" json:"max_secret_length"`
		
		// Severity threshold (0.0 to 1.0)
		SeverityThreshold float64 `yaml:"severity_threshold" json:"severity_threshold"`
		
		// Enable heuristic scanning
		EnableHeuristics bool `yaml:"enable_heuristics" json:"enable_heuristics"`
		
		// Enable entropy checking
		EnableEntropyCheck bool `yaml:"enable_entropy_check" json:"enable_entropy_check"`
		
		// Enable context checking
		EnableContextCheck bool `yaml:"enable_context_check" json:"enable_context_check"`
		
		// Files to ignore during scanning
		IgnoreFiles []string `yaml:"ignore_files" json:"ignore_files"`
		
		// Custom secret detection rules
		CustomRules []SecretRule `yaml:"custom_rules" json:"custom_rules"`
	} `yaml:"secrets" json:"secrets"`

	// Audit Logging Settings
	Audit struct {
		// Enable audit logging
		Enabled bool `yaml:"enabled" json:"enabled"`
		
		// Log to console
		LogToConsole bool `yaml:"log_to_console" json:"log_to_console"`
		
		// Log to file
		LogToFile bool `yaml:"log_to_file" json:"log_to_file"`
		
		// Minimum log level: debug, info, warning, error, critical
		LogLevel string `yaml:"log_level" json:"log_level"`
		
		// Maximum file size before rotation
		MaxFileSize int64 `yaml:"max_file_size" json:"max_file_size"`
		
		// Maximum number of log files to keep
		MaxFiles int `yaml:"max_files" json:"max_files"`
		
		// Number of days to retain log files
		RetainDays int `yaml:"retain_days" json:"retain_days"`
		
		// Output directory for log files
		OutputDir string `yaml:"output_dir" json:"output_dir"`
		
		// Enable JSON logging format
		EnableJSON bool `yaml:"enable_json" json:"enable_json"`
		
		// Enable structured logging
		EnableStructured bool `yaml:"enable_structured" json:"enable_structured"`
	} `yaml:"audit" json:"audit"`

	// Authentication Settings
	Authentication struct {
		// Require signed adapters
		SignedAdaptersRequired bool `yaml:"signed_adapters_required" json:"signed_adapters_required"`
		
	// Token validation settings
	Token struct {
		// Minimum token length
		MinLength int `yaml:"min_length" json:"min_length"`
		
		// Maximum token length
		MaxLength int `yaml:"max_length" json:"max_length"`
		
		// Allowed characters pattern
		AllowedCharacters string `yaml:"allowed_characters" json:"allowed_characters"`
		
		// Token expiration time
		Expiration time.Duration `yaml:"expiration" json:"expiration"`
	} `yaml:"token" json:"token"`
	
	// Session security
	Session struct {
		// Session timeout
		Timeout time.Duration `yaml:"timeout" json:"timeout"`
		
		// Maximum concurrent sessions per user
		MaxConcurrent int `yaml:"max_concurrent" json:"max_concurrent"`
		
		// Enable session encryption
		Encryption bool `yaml:"encryption" json:"encryption"`
		
		// Session cleanup interval
		CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"`
	} `yaml:"session" json:"session"`
} `yaml:"security" json:"security"`

// DefaultSecurityConfiguration returns a security configuration with default values
func DefaultSecurityConfiguration() *SecurityConfiguration {
	config := &SecurityConfiguration{
		YAML: struct {
			SandboxParser     bool
			MaxNestingDepth  int
			MaxDocumentSize  int64
			AllowCustomTags  bool
			ValidationMode   string
			RequiredFields   []string
			BlockedPatterns  []string
		}{
			SandboxParser:    true,
			MaxNestingDepth:  10,
			MaxDocumentSize:   1024 * 1024, // 1MB
			AllowCustomTags:  false,
			ValidationMode:   "strict",
			RequiredFields:   []string{},
			BlockedPatterns:  []string{},
		},
		Secrets: struct {
			Enabled            bool
			MaxLineLength      int
			MinSecretLength    int
			MaxSecretLength    int
			SeverityThreshold  float64
			EnableHeuristics   bool
			EnableEntropyCheck bool
			EnableContextCheck bool
			IgnoreFiles        []string
			CustomRules        []SecretRule
		}{
			Enabled:            true,
			MaxLineLength:      1000,
			MinSecretLength:    8,
			MaxSecretLength:    512,
			SeverityThreshold:  0.7,
			EnableHeuristics:   true,
			EnableEntropyCheck: true,
			EnableContextCheck: true,
			IgnoreFiles:        []string{".git", ".gitignore", "node_modules"},
			CustomRules:        []SecretRule{},
		},
		Audit: struct {
			Enabled           bool
			LogToConsole      bool
			LogToFile         bool
			LogLevel          string
			MaxFileSize       int64
			MaxFiles          int
			RetainDays        int
			OutputDir         string
			EnableJSON        bool
			EnableStructured bool
		}{
			Enabled:           true,
			LogToConsole:      true,
			LogToFile:         true,
			LogLevel:          "info",
			MaxFileSize:       10 * 1024 * 1024, // 10MB
			MaxFiles:          5,
			RetainDays:        30,
			OutputDir:         "audit_logs",
			EnableJSON:        true,
			EnableStructured: true,
		},
		Authentication: struct {
			SignedAdaptersRequired bool
			Token                   struct {
				MinLength          int
				MaxLength          int
				AllowedCharacters  string
				Expiration         time.Duration
			}
			Session struct {
				Timeout           time.Duration
				MaxConcurrent     int
				Encryption        bool
				CleanupInterval   time.Duration
			}
		}{
			SignedAdaptersRequired: false,
			Token: struct {
				MinLength:         16,
				MaxLength:         256,
				AllowedCharacters: "a-zA-Z0-9-_.",
				Expiration:        24 * time.Hour,
			},
			Session: struct {
				Timeout:           24 * time.Hour,
				MaxConcurrent:    5,
				Encryption:       true,
				CleanupInterval:  time.Hour,
			},
		},
	}

	return config
}

// LoadSecurityConfiguration loads security configuration from a file
func LoadSecurityConfiguration(filePath string) (*SecurityConfiguration, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("security configuration file not found: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read security configuration file: %w", err)
	}

	// Parse YAML
	var config SecurityConfiguration
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse security configuration: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid security configuration: %w", err)
	}

	return &config, nil
}

// SaveSecurityConfiguration saves security configuration to a file
func (c *SecurityConfiguration) SaveSecurityConfiguration(filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to YAML
	content, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal security configuration: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write security configuration file: %w", err)
	}

	return nil
}

// Validate validates the security configuration
func (c *SecurityConfiguration) Validate() error {
	// Validate YAML settings
	if c.YAML.MaxNestingDepth < 1 || c.YAML.MaxNestingDepth > 100 {
		return fmt.Errorf("yaml.max_nesting_depth must be between 1 and 100")
	}

	if c.YAML.MaxDocumentSize < 1024 || c.YAML.MaxDocumentSize > 100*1024*1024 {
		return fmt.Errorf("yaml.max_document_size must be between 1024 and 104857600 bytes")
	}

	validValidationModes := []string{"basic", "strict", "paranoid"}
	valid := false
	for _, mode := range validValidationModes {
		if c.YAML.ValidationMode == mode {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("yaml.validation_mode must be one of: basic, strict, paranoid")
	}

	// Validate secret scanning settings
	if c.Secrets.Enabled {
		if c.Secrets.MaxLineLength < 100 || c.Secrets.MaxLineLength > 10000 {
			return fmt.Errorf("secrets.max_line_length must be between 100 and 10000")
		}

		if c.Secrets.MinSecretLength < 4 || c.Secrets.MinSecretLength > 1024 {
			return fmt.Errorf("secrets.min_secret_length must be between 4 and 1024")
		}

		if c.Secrets.MaxSecretLength < 8 || c.Secrets.MaxSecretLength > 2048 {
			return fmt.Errorf("secrets.max_secret_length must be between 8 and 2048")
		}

		if c.Secrets.SeverityThreshold < 0.0 || c.Secrets.SeverityThreshold > 1.0 {
			return fmt.Errorf("secrets.severity_threshold must be between 0.0 and 1.0")
		}
	}

	// Validate audit logging settings
	if c.Audit.Enabled {
		validLogLevels := []string{"debug", "info", "warning", "error", "critical"}
		valid = false
		for _, level := range validLogLevels {
			if c.Audit.LogLevel == level {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("audit.log_level must be one of: debug, info, warning, error, critical")
		}

		if c.Audit.MaxFileSize < 1024 || c.Audit.MaxFileSize > 100*1024*1024 {
			return fmt.Errorf("audit.max_file_size must be between 1024 and 104857600 bytes")
		}

		if c.Audit.MaxFiles < 1 || c.Audit.MaxFiles > 100 {
			return fmt.Errorf("audit.max_files must be between 1 and 100")
		}

		if c.Audit.RetainDays < 1 || c.Audit.RetainDays > 365 {
			return fmt.Errorf("audit.retain_days must be between 1 and 365")
		}
	}

	// Validate authentication settings
	if c.Authentication.Token.MinLength < 8 || c.Authentication.Token.MaxLength > 1024 {
		return fmt.Errorf("authentication.token.min_length must be between 8 and 1024")
	}

	if c.Authentication.Token.MaxLength < c.Authentication.Token.MinLength {
		return fmt.Errorf("authentication.token.maxLength must be greater than or equal to minLength")
	}

	if c.Authentication.Session.Timeout < time.Minute || c.Authentication.Session.Timeout > 30*24*time.Hour {
		return fmt.Errorf("authentication.session.timeout must be between 1 minute and 30 days")
	}

	if c.Authentication.Session.MaxConcurrent < 1 || c.Authentication.Session.MaxConcurrent > 100 {
		return fmt.Errorf("authentication.session.max_concurrent must be between 1 and 100")
	}

	if c.Authentication.Session.CleanupInterval < time.Minute || c.Authentication.Session.CleanupInterval > 24*time.Hour {
		return fmt.Errorf("authentication.session.cleanup_interval must be between 1 minute and 24 hours")
	}

	return nil
}

// ToJSON converts the security configuration to JSON
func (c *SecurityConfiguration) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// FromJSON creates a security configuration from JSON
func (c *SecurityConfiguration) FromJSON(jsonData []byte) error {
	return json.Unmarshal(jsonData, c)
}

// GetYAMLValidationMode returns the validation mode as an enum
func (c *SecurityConfiguration) GetYAMLValidationMode() ValidationMode {
	switch c.YAML.ValidationMode {
	case "strict":
		return ValidationModeStrict
	case "paranoid":
		return ValidationModeParanoid
	default:
		return ValidationModeBasic
	}
}

// GetAuditLogLevel returns the log level as an enum
func (c *SecurityConfiguration) GetAuditLogLevel() EventLevel {
	switch c.Audit.LogLevel {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warning":
		return LevelWarning
	case "error":
		return LevelError
	case "critical":
		return LevelCritical
	default:
		return LevelInfo
	}
}

// CreateDefaultSecurityFile creates a default security configuration file
func CreateDefaultSecurityFile(filePath string) error {
	config := DefaultSecurityConfiguration()
	return config.SaveSecurityConfiguration(filePath)
}

// GetSecurityConfigPath returns the default security configuration file path
func GetSecurityConfigPath() string {
	// Check for user-specific config first
	userConfig := filepath.Join(os.Getenv("HOME"), ".config", "orya", "security.yml")
	if _, err := os.Stat(userConfig); err == nil {
		return userConfig
	}

	// Check for system-wide config
	systemConfig := filepath.Join("/etc", "orya", "security.yml")
	if _, err := os.Stat(systemConfig); err == nil {
		return systemConfig
	}

	// Default to current directory
	return filepath.Join(".", "security.yml")
}
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

// TraceabilityConfig contains configuration for traceability validation.
type TraceabilityConfig struct {
	// File patterns to include in validation
	IncludePatterns []string `yaml:"include_patterns"`
	// File patterns to exclude from validation
	ExcludePatterns []string `yaml:"exclude_patterns"`
	// Requirement prefixes to validate
	RequirementPrefixes []string `yaml:"requirement_prefixes"`
	// Number of commits to look back for traceability validation
	LookbackCommits int `yaml:"lookback_commits"`
	// Grace commits for new files (before requiring traceability)
	NewFileGraceCommits int `yaml:"new_file_grace_commits"`
	// Reporting level: advisory, warning, error
	ReportingLevel string `yaml:"reporting_level"`
	// Output format: cli, json, markdown
	OutputFormat string `yaml:"output_format"`
	// Enable traceability validation (can be disabled for debugging)
	Enabled bool `yaml:"enabled"`
	// Log level: debug, info, warn, error
	LogLevel string `yaml:"log_level"`
}

// DefaultTraceabilityConfig returns the default configuration.
func DefaultTraceabilityConfig() *TraceabilityConfig {
	return &TraceabilityConfig{
		IncludePatterns:     []string{"**/*.go", "**/*.md", "**/*.yaml", "**/*.yml"},
		ExcludePatterns:     []string{"**/vendor/**", "**/node_modules/**", "**/*.test.go", "**/test/**"},
		RequirementPrefixes: []string{"FR-", "ADR-", "INV-"},
		LookbackCommits:     10,
		NewFileGraceCommits: 5,
		ReportingLevel:      "advisory",
		OutputFormat:        "cli",
		Enabled:            true,
		LogLevel:           "info",
	}
}

// LoadTraceabilityConfig loads configuration from a file.
func LoadTraceabilityConfig(configPath string) (*TraceabilityConfig, error) {
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultTraceabilityConfig(), nil
	}

	// Load config from file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config TraceabilityConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Set defaults for any missing fields
	if len(config.IncludePatterns) == 0 {
		config.IncludePatterns = DefaultTraceabilityConfig().IncludePatterns
	}
	if len(config.ExcludePatterns) == 0 {
		config.ExcludePatterns = DefaultTraceabilityConfig().ExcludePatterns
	}
	if len(config.RequirementPrefixes) == 0 {
		config.RequirementPrefixes = DefaultTraceabilityConfig().RequirementPrefixes
	}
	if config.LookbackCommits == 0 {
		config.LookbackCommits = DefaultTraceabilityConfig().LookbackCommits
	}
	if config.NewFileGraceCommits == 0 {
		config.NewFileGraceCommits = DefaultTraceabilityConfig().NewFileGraceCommits
	}
	if config.ReportingLevel == "" {
		config.ReportingLevel = DefaultTraceabilityConfig().ReportingLevel
	}
	if config.OutputFormat == "" {
		config.OutputFormat = DefaultTraceabilityConfig().OutputFormat
	}
	if config.LogLevel == "" {
		config.LogLevel = DefaultTraceabilityConfig().LogLevel
	}

	return &config, nil
}

// SaveTraceabilityConfig saves configuration to a file.
func SaveTraceabilityConfig(config *TraceabilityConfig, configPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetTraceabilityConfigPath returns the path to the traceability config file.
func GetTraceabilityConfigPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".caf", "traceability.yaml")
}

// ValidateConfig validates the configuration.
func (c *TraceabilityConfig) Validate() error {
	// Validate reporting level
	if c.ReportingLevel != "advisory" && c.ReportingLevel != "warning" && c.ReportingLevel != "error" {
		return fmt.Errorf("invalid reporting level: %s", c.ReportingLevel)
	}

	// Validate output format
	if c.OutputFormat != "cli" && c.OutputFormat != "json" && c.OutputFormat != "markdown" {
		return fmt.Errorf("invalid output format: %s", c.OutputFormat)
	}

	// Validate log level
	if c.LogLevel != "debug" && c.LogLevel != "info" && c.LogLevel != "warn" && c.LogLevel != "error" {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Validate lookback commits
	if c.LookbackCommits <= 0 {
		return fmt.Errorf("lookback commits must be positive")
	}

	// Validate grace commits
	if c.NewFileGraceCommits < 0 {
		return fmt.Errorf("grace commits cannot be negative")
	}

	// Validate requirement prefixes
	if len(c.RequirementPrefixes) == 0 {
		return fmt.Errorf("at least one requirement prefix must be specified")
	}

	return nil
}

// MergeWithDefaults merges the config with default values.
func (c *TraceabilityConfig) MergeWithDefaults() *TraceabilityConfig {
	defaultConfig := DefaultTraceabilityConfig()
	
	// Only use non-default values from the current config
	if len(c.IncludePatterns) > 0 {
		defaultConfig.IncludePatterns = c.IncludePatterns
	}
	if len(c.ExcludePatterns) > 0 {
		defaultConfig.ExcludePatterns = c.ExcludePatterns
	}
	if len(c.RequirementPrefixes) > 0 {
		defaultConfig.RequirementPrefixes = c.RequirementPrefixes
	}
	if c.LookbackCommits > 0 {
		defaultConfig.LookbackCommits = c.LookbackCommits
	}
	if c.NewFileGraceCommits >= 0 {
		defaultConfig.NewFileGraceCommits = c.NewFileGraceCommits
	}
	if c.ReportingLevel != "" {
		defaultConfig.ReportingLevel = c.ReportingLevel
	}
	if c.OutputFormat != "" {
		defaultConfig.OutputFormat = c.OutputFormat
	}
	if c.LogLevel != "" {
		defaultConfig.LogLevel = c.LogLevel
	}
	defaultConfig.Enabled = c.Enabled

	return defaultConfig
}
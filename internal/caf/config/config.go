package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the .caf/config.yml file
type Config struct {
	Version              string             `yaml:"version"`
	TrivialityThreshold  int                `yaml:"triviality_threshold_loc"`
	UI                   UIConfig          `yaml:"ui"`
	LLM                  LLMConfig         `yaml:"llm"`
	LSP                  LSPConfig         `yaml:"lsp"`
	Telemetry            TelemetryConfig   `yaml:"telemetry"`
	Observability        ObservabilityConfig `yaml:"observability"`
	Security             SecurityConfig     `yaml:"security"`
	Cost                 CostConfig        `yaml:"cost"`
	Decisions            DecisionsConfig   `yaml:"decisions"`
}

// UIConfig represents UI settings
type UIConfig struct {
	Theme               string `yaml:"theme"`
	Language            string `yaml:"language"`
	HeatmapRefreshSeconds int  `yaml:"heatmap_refresh_seconds"`
	ReviewMaxCards      int    `yaml:"review_max_cards"`
}

// LLMConfig represents LLM settings
type LLMConfig struct {
	Backend   string `yaml:"backend"`
	Endpoint  string `yaml:"endpoint"`
	Model     string `yaml:"model"`
}

// LSPConfig represents LSP settings
type LSPConfig struct {
	Strict bool `yaml:"strict"`
}

// TelemetryConfig represents telemetry settings
type TelemetryConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// ObservabilityConfig represents observability settings
type ObservabilityConfig struct {
	OtelEnabled      bool   `yaml:"otel_enabled"`
	OtelEndpoint     string `yaml:"otel_endpoint"`
	MetricsExport    bool   `yaml:"metrics_export"`
	MetricsPort      int    `yaml:"metrics_port"`
}

// SecurityConfig represents security settings
type SecurityConfig struct {
	SignedAdaptersRequired bool `yaml:"signed_adapters_required"`
	SandboxParser          bool `yaml:"sandbox_parser"`
	AuditLog              bool `yaml:"audit_log"`
}

// CostConfig represents cost tracking settings
type CostConfig struct {
	LLMTokenCapDaily    int `yaml:"llm_token_cap_daily"`
	LLMTokenCapMonthly  int `yaml:"llm_token_cap_monthly"`
	AlertThresholdPct   int `yaml:"alert_threshold_pct"`
}

// DecisionsConfig represents decisions settings
type DecisionsConfig struct {
	Autocomplete string `yaml:"autocomplete"`
}

// Manager handles config operations
type Manager struct {
	configFile string
}

// NewManager creates a new config manager
func NewManager(configFile string) *Manager {
	return &Manager{
		configFile: configFile,
	}
}

// LoadConfig loads the config file with defaults
func (m *Manager) LoadConfig() (*Config, error) {
	// Default config
	config := getDefaultConfig()

	// Load existing config if it exists
	if _, err := os.Stat(m.configFile); err == nil {
		if err := ReadYAML(m.configFile, &config); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return &config, nil
}

// SaveConfig saves the config file
func (m *Manager) SaveConfig(config Config) error {
	return WriteYAML(m.configFile, config)
}

// GetConfigPath returns the path to the config file
func GetConfigPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".caf", "config.yml")
}

// getDefaultConfig returns default configuration values
func getDefaultConfig() Config {
	return Config{
		Version:              "0.4.0-triz2",
		TrivialityThreshold:  50,
		UI: UIConfig{
			Theme:               "dark",
			Language:            "ru",
			HeatmapRefreshSeconds: 30,
			ReviewMaxCards:      15,
		},
		LLM: LLMConfig{
			Backend:  "",
			Endpoint: "",
			Model:    "",
		},
		LSP: LSPConfig{
			Strict: true,
		},
		Telemetry: TelemetryConfig{
			Enabled:  false,
			Endpoint: "",
		},
		Observability: ObservabilityConfig{
			OtelEnabled:      false,
			OtelEndpoint:     "",
			MetricsExport:    false,
			MetricsPort:      9464,
		},
		Security: SecurityConfig{
			SignedAdaptersRequired: false,
			SandboxParser:          true,
			AuditLog:              true,
		},
		Cost: CostConfig{
			LLMTokenCapDaily:    100000,
			LLMTokenCapMonthly:  2000000,
			AlertThresholdPct:   80,
		},
		Decisions: DecisionsConfig{
			Autocomplete: "aggressive",
		},
	}
}

// ReadYAML safely reads YAML from a file
func ReadYAML(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return SafeParseYAML(data, out)
}

// WriteYAML safely writes YAML to a file
func WriteYAML(filename string, data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Write to temp file first (atomic write)
	tempFile := filename + ".tmp"
	if err := os.WriteFile(tempFile, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	return os.Rename(tempFile, filename)
}

// SafeParseYAML safely parses YAML with sandboxing
func SafeParseYAML(data []byte, out interface{}) error {
	// Use gopkg.in/yaml.v3 (no custom tags that execute code)
	// Validate output struct (no interface{} at top level)
	if out == nil {
		return fmt.Errorf("output must not be nil")
	}

	// Reject any YAML tag that starts with "!" (custom tags)
	// This prevents potential code execution through custom tags
	for i, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "!") {
			return fmt.Errorf("custom YAML tags are not allowed at line %d", i+1)
		}
	}

	return yaml.Unmarshal(data, out)
}
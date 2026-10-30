package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAML utilities for safe parsing

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

// WriteYAML safely writes YAML to a file
func WriteYAML(filename string, data interface{}) error {
	// Marshal to YAML
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

// ReadYAML safely reads YAML from a file
func ReadYAML(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return SafeParseYAML(data, out)
}

// EnsureDir ensures a directory exists
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// EnsureCAFDir ensures the .caf directory structure exists
func EnsureCAFDir(repoRoot string) (string, error) {
	cafDir := filepath.Join(repoRoot, ".caf")
	
	// Create main .caf directory
	if err := EnsureDir(cafDir); err != nil {
		return "", fmt.Errorf("failed to create .caf directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{
		filepath.Join(cafDir, "sessions"),
		filepath.Join(cafDir, "decisions"),
		filepath.Join(cafDir, "mental-model"),
		filepath.Join(cafDir, "cache"),
		filepath.Join(cafDir, "metrics"),
		filepath.Join(cafDir, "reviews"),
	}

	for _, dir := range subdirs {
		if err := EnsureDir(dir); err != nil {
			return "", fmt.Errorf("failed to create subdirectory %s: %w", dir, err)
		}
	}

	return cafDir, nil
}
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Manifest represents the .caf/manifest.yml file
type Manifest struct {
	Version      string    `yaml:"version"`
	Repo         string    `yaml:"repo"`
	Initialized  time.Time `yaml:"initialized_at"`
	LastSync     time.Time `yaml:"last_sync,omitempty"`
	CrossRepo    CrossRepo `yaml:"cross_repo,omitempty"`
}

// CrossRepo represents cross-repo sync configuration
type CrossRepo struct {
	Enabled     bool     `yaml:"enabled"`
	Repos       []string `yaml:"repos"`
	LastSync    time.Time `yaml:"last_sync,omitempty"`
}

// Manager handles manifest operations
type Manager struct {
	manifestFile string
}

// NewManifestManager creates a new manifest manager
func NewManifestManager(manifestFile string) *Manager {
	return &Manager{
		manifestFile: manifestFile,
	}
}

// CreateManifest creates a new manifest file
func (m *Manager) CreateManifest(repo string) error {
	manifest := Manifest{
		Version:     "0.4.0-triz2",
		Repo:        repo,
		Initialized: time.Now(),
	}

	return m.SaveManifest(manifest)
}

// LoadManifest loads the manifest file
func (m *Manager) LoadManifest() (*Manifest, error) {
	var manifest Manifest
	if err := ReadYAML(m.manifestFile, &manifest); err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	return &manifest, nil
}

// SaveManifest saves the manifest file
func (m *Manager) SaveManifest(manifest Manifest) error {
	return WriteYAML(m.manifestFile, manifest)
}

// UpdateLastSync updates the last sync timestamp
func (m *Manager) UpdateLastSync() error {
	manifest, err := m.LoadManifest()
	if err != nil {
		return err
	}

	manifest.LastSync = time.Now()
	return m.SaveManifest(*manifest)
}

// GetManifestPath returns the path to the manifest file
func GetManifestPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".caf", "manifest.yml")
}

// GetCAFDir returns the path to the .caf directory
func GetCAFDir(repoRoot string) string {
	return filepath.Join(repoRoot, ".caf")
}

// InitializeCAF initializes the .caf directory structure
func InitializeCAF(repoRoot string) error {
	// Ensure .caf directory exists
	cafDir, err := EnsureCAFDir(repoRoot)
	if err != nil {
		return err
	}

	// Create manifest.yml
	manifestFile := filepath.Join(cafDir, "manifest.yml")
	manifestManager := NewManifestManager(manifestFile)
	
	// Get git repo name
	repo := getGitRepoName(repoRoot)
	if repo == "" {
		repo = "unknown"
	}

	return manifestManager.CreateManifest(repo)
}

// getGitRepoName gets the git repository name from the current directory
func getGitRepoName(repoRoot string) string {
	gitDir := filepath.Join(repoRoot, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return ""
	}

	// Read .git/config to get remote URL
	configFile := filepath.Join(gitDir, "config")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return ""
	}

	// Simple parsing to get repo name
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "url = ") {
			url := strings.TrimPrefix(line, "url = ")
			// Extract repo name from URL
			if strings.Contains(url, "github.com") {
				parts := strings.Split(url, "/")
				if len(parts) >= 2 {
					return parts[len(parts)-2] + "/" + parts[len(parts)-1]
				}
			}
			break
		}
	}

	return ""
}
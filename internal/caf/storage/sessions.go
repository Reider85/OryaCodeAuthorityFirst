package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Session represents a single AI agent session
type Session struct {
	SessionID     string                 `yaml:"session_id"`
	Repo          string                 `yaml:"repo"`
	StartedAt     time.Time              `yaml:"started_at"`
	EndedAt       time.Time              `yaml:"ended_at"`
	Agent         string                 `yaml:"agent"`
	Model         string                 `yaml:"model"`
	ToolCall      string                 `yaml:"tool_call"`
	InitialPrompt string                 `yaml:"initial_prompt"`
	Iterations    []Iteration            `yaml:"iterations"`
	FilesTouched  []string               `yaml:"files_touched"`
	CommitSHA     string                 `yaml:"commit_sha"`
	DecisionsMade []Decision             `yaml:"decisions_made"`
	VoiceSignature string                `yaml:"voice_signature"`
}

// Iteration represents a single iteration in a session
type Iteration struct {
	ID      string `yaml:"id"`
	Action  string `yaml:"action"`
	Target  string `yaml:"target"`
	Diff    string `yaml:"diff_summary,omitempty"`
}

// Decision represents a decision made during a session
type Decision struct {
	ID          string    `yaml:"id"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Timestamp   time.Time `yaml:"timestamp"`
}

// Manager handles session storage and operations
type Manager struct {
	sessionsDir string
}

// NewManager creates a new storage manager
func NewManager(sessionsDir string) *Manager {
	return &Manager{
		sessionsDir: sessionsDir,
	}
}

// SaveSession saves a session to disk
func (m *Manager) SaveSession(session *Session) error {
	// Ensure sessions directory exists
	if err := os.MkdirAll(m.sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Create session file path
	sessionFile := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.yml", session.SessionID))
	
	// Marshal to YAML
	data, err := yaml.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to file (atomic write via temp file)
	tempFile := sessionFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write session temp file: %w", err)
	}
	
	// Atomic rename
	if err := os.Rename(tempFile, sessionFile); err != nil {
		return fmt.Errorf("failed to rename session file: %w", err)
	}

	return nil
}

// LoadSession loads a session from disk
func (m *Manager) LoadSession(sessionID string) (*Session, error) {
	sessionFile := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.yml", sessionID))
	
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := yaml.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// ListSessions returns all session IDs
func (m *Manager) ListSessions() ([]string, error) {
	var sessionIDs []string
	
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return sessionIDs, nil
		}
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yml" {
			sessionIDs = append(sessionIDs, entry.Name()[:len(entry.Name())-4])
		}
	}

	return sessionIDs, nil
}

// DeleteSession deletes a session file
func (m *Manager) DeleteSession(sessionID string) error {
	sessionFile := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.yml", sessionID))
	return os.Remove(sessionFile)
}
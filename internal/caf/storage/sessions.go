package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Session represents a single AI agent session
type Session struct {
	SessionID     string                 `yaml:"session_id"`
	Repo          string                 `yaml:"repo"`
	StartedAt     time.Time              `yaml:"started_at"`
	EndedAt       time.Time              `yaml:"ended_at,omitempty"`
	Agent         string                 `yaml:"agent"`
	Model         string                 `yaml:"model"`
	ToolCall      string                 `yaml:"tool_call"`
	InitialPrompt string                 `yaml:"initial_prompt"`
	Iterations    []Iteration            `yaml:"iterations"`
	FilesTouched  []string               `yaml:"files_touched"`
	CommitSHA     string                 `yaml:"commit_sha,omitempty"`
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

// validateSessionID validates a session ID to prevent path traversal
func validateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if sessionID == "." || sessionID == ".." {
		return fmt.Errorf("invalid session ID: %q", sessionID)
	}
	if strings.Contains(sessionID, "/") || strings.Contains(sessionID, "\\") {
		return fmt.Errorf("invalid session ID: contains path separators")
	}
	return nil
}

// SaveSession saves a session to disk
func (m *Manager) SaveSession(session *Session) error {
	// Validate session ID
	if err := validateSessionID(session.SessionID); err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	// Ensure sessions directory exists
	if err := EnsureDir(m.sessionsDir); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	// Use shared save function
	return saveSessionFile(filepath.Join(m.sessionsDir, session.SessionID), session)
}

// LoadSession loads a session from disk
func (m *Manager) LoadSession(sessionID string) (*Session, error) {
	// Validate session ID
	if err := validateSessionID(sessionID); err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	// Use shared load function
	return loadSessionFile(filepath.Join(m.sessionsDir, sessionID))
}

// ListSessions returns all session IDs
func (m *Manager) ListSessions() ([]string, error) {
	return ListSessions(m.sessionsDir)
}

// DeleteSession deletes a session file
func (m *Manager) DeleteSession(sessionID string) error {
	// Validate session ID
	if err := validateSessionID(sessionID); err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	// Use shared delete function
	return deleteSessionFile(filepath.Join(m.sessionsDir, sessionID))
}

// LoadSession загружает сессию по пути без расширения .yml
// (пакетная функция для ContextGenerator).
func LoadSession(path string) (*Session, error) {
	return loadSessionFile(path)
}

// ListSessions возвращает ID всех сессий в каталоге (пакетная функция
// для ContextGenerator).
func ListSessions(dir string) ([]string, error) {
	return listYMLFiles(dir)
}

// Internal helper functions

// loadSessionFile loads a session from a file path (without .yml extension)
func loadSessionFile(path string) (*Session, error) {
	sessionFile := path + ".yml"
	
	var session Session
	if err := ReadYAML(sessionFile, &session); err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	return &session, nil
}

// saveSessionFile saves a session to a file path (without .yml extension)
func saveSessionFile(path string, session *Session) error {
	return WriteYAML(path+".yml", session)
}

// deleteSessionFile deletes a session file (without .yml extension)
func deleteSessionFile(path string) error {
	sessionFile := path + ".yml"
	return os.Remove(sessionFile)
}

// listYMLFiles returns all .yml files in a directory (without extension)
func listYMLFiles(dir string) ([]string, error) {
	var sessionIDs []string
	
	entries, err := os.ReadDir(dir)
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
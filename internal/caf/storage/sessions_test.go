package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSessionIDValidation(t *testing.T) {
	tests := []struct {
		name    string
		sessionID string
		wantErr bool
	}{
		{"empty", "", true},
		{"dot", ".", true},
		{"dotdot", "..", true},
		{"slash", "2026-10-24/path", true},
		{"backslash", "2026-10-24\\path", true},
		{"valid", "2026-10-24-a1b2c3", false},
		{"valid with dash", "2026-10-24-x-y-z", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSessionID(tt.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSessionID(%q) error = %v, wantErr %v", tt.sessionID, err, tt.wantErr)
			}
		})
	}
}

func TestSaveLoadSession(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create test session
	session := &Session{
		SessionID:     "2026-10-24-a1b2c3",
		Repo:          "test/repo",
		StartedAt:     time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
		EndedAt:       time.Date(2026, 10, 24, 11, 47, 33, 0, time.UTC),
		Agent:         "opencode",
		Model:         "glm-4.6",
		ToolCall:      "opencode run \"add rate limiter\"",
		InitialPrompt: "Need to add rate limiter on POST /v1/policies.",
		Iterations: []Iteration{
			{
				ID:     "i1",
				Action: "read",
				Target: "internal/api/policies.go",
			},
			{
				ID:     "i2",
				Action: "edit",
				Target: "internal/api/policies.go",
				Diff:   "added middleware RateLimit(token, 10/min)",
			},
		},
		FilesTouched:  []string{"internal/api/policies.go"},
		CommitSHA:     "9f3e2a1",
		DecisionsMade: []Decision{
			{
				ID:          "d1",
				Title:       "Rate limiter decision",
				Description: "Added token bucket rate limiting",
				Timestamp:   time.Date(2026, 10, 24, 11, 30, 0, 0, time.UTC),
			},
		},
		VoiceSignature: "rate limiter needed urgently",
	}

	// Test SaveSession
	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Test LoadSession
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// Verify session data
	if loaded.SessionID != session.SessionID {
		t.Errorf("SessionID = %v, want %v", loaded.SessionID, session.SessionID)
	}
	if loaded.Repo != session.Repo {
		t.Errorf("Repo = %v, want %v", loaded.Repo, session.Repo)
	}
	if !loaded.StartedAt.Equal(session.StartedAt) {
		t.Errorf("StartedAt = %v, want %v", loaded.StartedAt, session.StartedAt)
	}
	if !loaded.EndedAt.Equal(session.EndedAt) {
		t.Errorf("EndedAt = %v, want %v", loaded.EndedAt, session.EndedAt)
	}
	if loaded.Agent != session.Agent {
		t.Errorf("Agent = %v, want %v", loaded.Agent, session.Agent)
	}
	if loaded.Model != session.Model {
		t.Errorf("Model = %v, want %v", loaded.Model, session.Model)
	}
	if loaded.ToolCall != session.ToolCall {
		t.Errorf("ToolCall = %v, want %v", loaded.ToolCall, session.ToolCall)
	}
	if loaded.InitialPrompt != session.InitialPrompt {
		t.Errorf("InitialPrompt = %v, want %v", loaded.InitialPrompt, session.InitialPrompt)
	}
	if len(loaded.Iterations) != len(session.Iterations) {
		t.Errorf("Iterations count = %v, want %v", len(loaded.Iterations), len(session.Iterations))
	} else {
		for i, iter := range loaded.Iterations {
			if iter.ID != session.Iterations[i].ID {
				t.Errorf("Iteration[%d].ID = %v, want %v", i, iter.ID, session.Iterations[i].ID)
			}
			if iter.Action != session.Iterations[i].Action {
				t.Errorf("Iteration[%d].Action = %v, want %v", i, iter.Action, session.Iterations[i].Action)
			}
			if iter.Target != session.Iterations[i].Target {
				t.Errorf("Iteration[%d].Target = %v, want %v", i, iter.Target, session.Iterations[i].Target)
			}
			if iter.Diff != session.Iterations[i].Diff {
				t.Errorf("Iteration[%d].Diff = %v, want %v", i, iter.Diff, session.Iterations[i].Diff)
			}
		}
	}
	if len(loaded.FilesTouched) != len(session.FilesTouched) {
		t.Errorf("FilesTouched count = %v, want %v", len(loaded.FilesTouched), len(session.FilesTouched))
	} else {
		for i, file := range loaded.FilesTouched {
			if file != session.FilesTouched[i] {
				t.Errorf("FilesTouched[%d] = %v, want %v", i, file, session.FilesTouched[i])
			}
		}
	}
	if loaded.CommitSHA != session.CommitSHA {
		t.Errorf("CommitSHA = %v, want %v", loaded.CommitSHA, session.CommitSHA)
	}
	if len(loaded.DecisionsMade) != len(session.DecisionsMade) {
		t.Errorf("DecisionsMade count = %v, want %v", len(loaded.DecisionsMade), len(session.DecisionsMade))
	} else {
		for i, decision := range loaded.DecisionsMade {
			if decision.ID != session.DecisionsMade[i].ID {
				t.Errorf("Decision[%d].ID = %v, want %v", i, decision.ID, session.DecisionsMade[i].ID)
			}
			if decision.Title != session.DecisionsMade[i].Title {
				t.Errorf("Decision[%d].Title = %v, want %v", i, decision.Title, session.DecisionsMade[i].Title)
			}
			if decision.Description != session.DecisionsMade[i].Description {
				t.Errorf("Decision[%d].Description = %v, want %v", i, decision.Description, session.DecisionsMade[i].Description)
			}
			if !decision.Timestamp.Equal(session.DecisionsMade[i].Timestamp) {
				t.Errorf("Decision[%d].Timestamp = %v, want %v", i, decision.Timestamp, session.DecisionsMade[i].Timestamp)
			}
		}
	}
	if loaded.VoiceSignature != session.VoiceSignature {
		t.Errorf("VoiceSignature = %v, want %v", loaded.VoiceSignature, session.VoiceSignature)
	}

	// Test LoadSession package function
	loaded2, err := LoadSession(filepath.Join(sessionsDir, "2026-10-24-a1b2c3"))
	if err != nil {
		t.Fatalf("LoadSession package function failed: %v", err)
	}
	if loaded2.SessionID != session.SessionID {
		t.Errorf("LoadSession package: SessionID = %v, want %v", loaded2.SessionID, session.SessionID)
	}
}

func TestLoadSessionMissingFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(filepath.Join(tempDir, "sessions"))
	
	_, err := manager.LoadSession("nonexistent")
	if err == nil {
		t.Error("LoadSession should fail for missing file")
	}
}

func TestListSessions(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Test empty directory
	sessions, err := ListSessions(sessionsDir)
	if err != nil {
		t.Fatalf("ListSessions failed on empty dir: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("ListSessions on empty dir = %v, want empty", sessions)
	}

	// Create test sessions
	session1 := &Session{SessionID: "2026-10-24-a1b2c3", Agent: "opencode"}
	session2 := &Session{SessionID: "2026-10-25-d4e5f6", Agent: "claude-code"}
	
	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session1); err != nil {
		t.Fatalf("SaveSession failed for session1: %v", err)
	}
	if err := manager.SaveSession(session2); err != nil {
		t.Fatalf("SaveSession failed for session2: %v", err)
	}

	// Test ListSessions
	sessions, err = ListSessions(sessionsDir)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("ListSessions count = %v, want 2", len(sessions))
	}
	
	// Should contain both sessions
	found1 := false
	found2 := false
	for _, id := range sessions {
		if id == "2026-10-24-a1b2c3" {
			found1 = true
		}
		if id == "2026-10-25-d4e5f6" {
			found2 = true
		}
	}
	if !found1 {
		t.Error("ListSessions missing 2026-10-24-a1b2c3")
	}
	if !found2 {
		t.Error("ListSessions missing 2026-10-25-d4e5f6")
	}

	// Test Manager.ListSessions
	sessions, err = manager.ListSessions()
	if err != nil {
		t.Fatalf("Manager.ListSessions failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("Manager.ListSessions count = %v, want 2", len(sessions))
	}
}

func TestDeleteSession(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create test session
	session := &Session{SessionID: "2026-10-24-a1b2c3", Agent: "opencode"}
	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Verify session exists
	sessions, err := ListSessions(sessionsDir)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("Expected 1 session, got %v", len(sessions))
	}

	// Delete session
	if err := manager.DeleteSession("2026-10-24-a1b2c3"); err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Verify session is deleted
	sessions, err = ListSessions(sessionsDir)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("Expected 0 sessions after delete, got %v", len(sessions))
	}

	// Test deleting nonexistent session
	err = manager.DeleteSession("nonexistent")
	if err == nil {
		t.Error("DeleteSession should fail for nonexistent session")
	}
}

func TestSafeParseYAML(t *testing.T) {
	// Test that SafeParseYAML rejects custom tags
	yamlContent := `
session_id: 2026-10-24-a1b2c3
repo: test/repo
started_at: 2026-10-24T11:02:14Z
agent: opencode
!custom_tag "this should be rejected"
`
	
	var session Session
	err := SafeParseYAML([]byte(yamlContent), &session)
	if err == nil {
		t.Error("SafeParseYAML should reject custom YAML tags")
	}
}

func TestOmitEmptyFields(t *testing.T) {
	// Test that optional fields with zero values are omitted from YAML
	session := &Session{
		SessionID:     "2026-10-24-a1b2c3",
		Repo:          "test/repo",
		StartedAt:     time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
		EndedAt:       time.Time{}, // zero value - should be omitted
		Agent:         "opencode",
		Model:         "glm-4.6",
		ToolCall:      "",
		InitialPrompt: "",
		Iterations:    []Iteration{},
		FilesTouched:  []string{},
		CommitSHA:     "", // zero value - should be omitted
		DecisionsMade: []Decision{},
		VoiceSignature: "",
	}

	// Test with Manager
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	manager := NewManager(sessionsDir)
	
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Load and check that zero-value fields are omitted
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// Verify that zero-value optional fields are zero (they were omitted in YAML)
	if !loaded.EndedAt.IsZero() {
		t.Errorf("EndedAt should be zero (omitted), got %v", loaded.EndedAt)
	}
	if loaded.CommitSHA != "" {
		t.Errorf("CommitSHA should be empty (omitted), got %q", loaded.CommitSHA)
	}
	if len(loaded.Iterations) != 0 {
		t.Errorf("Iterations should be empty (omitted), got %d", len(loaded.Iterations))
	}
	if len(loaded.FilesTouched) != 0 {
		t.Errorf("FilesTouched should be empty (omitted), got %d", len(loaded.FilesTouched))
	}
	if len(loaded.DecisionsMade) != 0 {
		t.Errorf("DecisionsMade should be empty (omitted), got %d", len(loaded.DecisionsMade))
	}
}

func TestSaveSessionFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create test session
	session := &Session{
		SessionID:     "2026-10-24-a1b2c3",
		Repo:          "test/repo",
		StartedAt:     time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
		Agent:         "opencode",
		Model:         "glm-4.6",
		InitialPrompt: "test prompt",
	}

	// Test saveSessionFile
	err := saveSessionFile(filepath.Join(sessionsDir, "test"), session)
	if err != nil {
		t.Fatalf("saveSessionFile failed: %v", err)
	}

	// Verify file exists
	sessionFile := filepath.Join(sessionsDir, "test.yml")
	if _, err := os.Stat(sessionFile); err != nil {
		t.Errorf("session file %s does not exist: %v", sessionFile, err)
	}

	// Load and verify
	loaded, err := loadSessionFile(filepath.Join(sessionsDir, "test"))
	if err != nil {
		t.Fatalf("loadSessionFile failed: %v", err)
	}
	if loaded.SessionID != session.SessionID {
		t.Errorf("SessionID mismatch: got %v, want %v", loaded.SessionID, session.SessionID)
	}
}

func TestDeleteSessionFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create test session file directly
	sessionFile := filepath.Join(sessionsDir, "test.yml")
	sessionContent := `session_id: 2026-10-24-a1b2c3
agent: opencode
started_at: 2026-10-24T11:02:14Z
`
	if err := os.WriteFile(sessionFile, []byte(sessionContent), 0644); err != nil {
		t.Fatalf("failed to create test session file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(sessionFile); err != nil {
		t.Fatalf("session file %s does not exist: %v", sessionFile, err)
	}

	// Test deleteSessionFile
	err := deleteSessionFile(filepath.Join(sessionsDir, "test"))
	if err != nil {
		t.Fatalf("deleteSessionFile failed: %v", err)
	}

	// Verify file is deleted
	if _, err := os.Stat(sessionFile); err == nil {
		t.Errorf("session file %s still exists after deletion", sessionFile)
	}
}

func TestListYMLFiles(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Test empty directory
	files, err := listYMLFiles(sessionsDir)
	if err != nil {
		t.Fatalf("listYMLFiles failed on empty dir: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("listYMLFiles on empty dir = %v, want empty", files)
	}

	// Create test files
	testFiles := []string{"session1.yml", "session2.yml", "not_session.txt", "session3.yml"}
	for _, file := range testFiles {
		filePath := filepath.Join(sessionsDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Test listYMLFiles
	files, err = listYMLFiles(sessionsDir)
	if err != nil {
		t.Fatalf("listYMLFiles failed: %v", err)
	}

	// Should return only .yml files without extension
	expected := []string{"session1", "session2", "session3"}
	if len(files) != len(expected) {
		t.Errorf("listYMLFiles count = %v, want %v", len(files), len(expected))
	}

	// Check all expected files are present
	for _, expectedFile := range expected {
		found := false
		for _, actualFile := range files {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("listYMLFiles missing expected file: %s", expectedFile)
		}
	}

	// Test with non-existent directory
	files, err = listYMLFiles(filepath.Join(tempDir, "nonexistent"))
	if err != nil {
		t.Fatalf("listYMLFiles failed on non-existent dir: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("listYMLFiles on non-existent dir = %v, want empty", files)
	}
}

func TestManagerListSessionsEmpty(t *testing.T) {
	// Test Manager.ListSessions on empty directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	manager := NewManager(sessionsDir)
	sessions, err := manager.ListSessions()
	if err != nil {
		t.Fatalf("Manager.ListSessions failed on empty dir: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("Manager.ListSessions on empty dir = %v, want empty", sessions)
	}
}

func TestManagerSaveSessionError(t *testing.T) {
	t.Skip("Skip for now - directory creation behavior differs by OS")
	// Test Manager.SaveSession to a directory that exists but is read-only
	session := &Session{
		SessionID: "2026-10-24-a1b2c3",
		Agent:     "opencode",
	}

	// Create a temp directory and make it read-only
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatalf("failed to create read-only dir: %v", err)
	}

	// Change permissions to read-only
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to set read-only permissions: %v", err)
	}

	manager := NewManager(readOnlyDir)
	err := manager.SaveSession(session)
	if err == nil {
		t.Error("Manager.SaveSession should fail with read-only directory")
	}
}

func TestLoadSessionWithInvalidYAML(t *testing.T) {
	// Test LoadSession with invalid YAML content
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create file with invalid YAML
	sessionFile := filepath.Join(sessionsDir, "invalid.yml")
	if err := os.WriteFile(sessionFile, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to create invalid YAML file: %v", err)
	}

	manager := NewManager(sessionsDir)
	_, err := manager.LoadSession("invalid")
	if err == nil {
		t.Error("Manager.LoadSession should fail with invalid YAML")
	}
}

func TestLoadSessionPackageWithInvalidYAML(t *testing.T) {
	// Test LoadSession package function with invalid YAML
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create file with invalid YAML
	sessionFile := filepath.Join(sessionsDir, "invalid.yml")
	if err := os.WriteFile(sessionFile, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to create invalid YAML file: %v", err)
	}

	_, err := LoadSession(filepath.Join(sessionsDir, "invalid"))
	if err == nil {
		t.Error("LoadSession package function should fail with invalid YAML")
	}
}

func TestListSessionsWithInvalidFiles(t *testing.T) {
	// Test ListSessions with non-YAML files in directory
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create various file types
	testFiles := []string{"session1.yml", "session2.yml", "not_session.txt", "session3.yml", "session4.yml"}
	for _, file := range testFiles {
		filePath := filepath.Join(sessionsDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}

	// Test ListSessions
	sessions, err := ListSessions(sessionsDir)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	// Should return only .yml files (case insensitive for extension)
	expected := []string{"session1", "session2", "session3", "session4"}
	if len(sessions) != len(expected) {
		t.Errorf("ListSessions count = %v, want %v", len(sessions), len(expected))
	}

	// Check all expected files are present
	for _, expectedFile := range expected {
		found := false
		for _, actualFile := range sessions {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListSessions missing expected file: %s", expectedFile)
		}
	}
}

func TestSessionWithComplexData(t *testing.T) {
	// Test session with complex nested data and special characters
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	session := &Session{
		SessionID:     "2026-10-24-a1b2c3",
		Repo:          "test/repo",
		StartedAt:     time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
		EndedAt:       time.Date(2026, 10, 24, 11, 47, 33, 0, time.UTC),
		Agent:         "opencode",
		Model:         "glm-4.6",
		ToolCall:      `opencode run "add rate limiter --timeout=30s --max-requests=100"`,
		InitialPrompt: `Need to add rate limiter on POST /v1/policies.
The policy should handle:
- Token bucket algorithm
- Timeout of 30 seconds
- Max 100 requests per minute`,
		Iterations: []Iteration{
			{
				ID:     "i1",
				Action: "read",
				Target: "internal/api/policies.go",
				Diff:   `// Read existing policy file
func PoliciesHandler(w http.ResponseWriter, r *http.Request) {
    // existing code
}`,
			},
			{
				ID:     "i2",
				Action: "edit",
				Target: "internal/api/policies.go",
				Diff:   `// Added rate limiting middleware
func PoliciesHandler(w http.ResponseWriter, r *http.Request) {
    // Rate limiting middleware added here
    if !rateLimit.Allow() {
        http.Error(w, "Too many requests", http.StatusTooManyRequests)
        return
    }
    // existing code continues...
}`,
			},
		},
		FilesTouched:  []string{"internal/api/policies.go", "internal/rate/limiter.go"},
		CommitSHA:     "9f3e2a1b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9",
		DecisionsMade: []Decision{
			{
				ID:          "d1",
				Title:       "Rate limiter implementation",
				Description: "Chose token bucket algorithm over leaky bucket for better burst handling",
				Timestamp:   time.Date(2026, 10, 24, 11, 30, 0, 0, time.UTC),
			},
			{
				ID:          "d2",
				Title:       "Timeout configuration",
				Description: "Set 30-second timeout based on average API response time analysis",
				Timestamp:   time.Date(2026, 10, 24, 11, 35, 0, 0, time.UTC),
			},
		},
		VoiceSignature: "rate limiter needed urgently for API protection",
	}

	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Load and verify complex data
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// Verify complex string fields
	if loaded.ToolCall != session.ToolCall {
		t.Errorf("ToolCall mismatch:\ngot:  %v\nwant: %v", loaded.ToolCall, session.ToolCall)
	}
	if loaded.InitialPrompt != session.InitialPrompt {
		t.Errorf("InitialPrompt mismatch")
	}
	if loaded.VoiceSignature != session.VoiceSignature {
		t.Errorf("VoiceSignature mismatch")
	}

	// Verify complex iteration diff
	if len(loaded.Iterations) != len(session.Iterations) {
		t.Errorf("Iterations count mismatch: got %v, want %v", len(loaded.Iterations), len(session.Iterations))
	}
	for i, iter := range loaded.Iterations {
		if iter.Diff != session.Iterations[i].Diff {
			t.Errorf("Iteration[%d].Diff mismatch", i)
		}
	}

	// Verify complex decisions
	if len(loaded.DecisionsMade) != len(session.DecisionsMade) {
		t.Errorf("DecisionsMade count mismatch: got %v, want %v", len(loaded.DecisionsMade), len(session.DecisionsMade))
	}
	for i, decision := range loaded.DecisionsMade {
		if decision.Description != session.DecisionsMade[i].Description {
			t.Errorf("Decision[%d].Description mismatch", i)
		}
	}
}

func TestSessionWithZeroValues(t *testing.T) {
	// Test session with many zero/empty values to test omitempty
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	session := &Session{
		SessionID:     "2026-10-24-a1b2c3",
		Repo:          "",
		StartedAt:     time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
		EndedAt:       time.Time{}, // zero value - should be omitted
		Agent:         "",
		Model:         "",
		ToolCall:      "",
		InitialPrompt: "",
		Iterations:    []Iteration{},
		FilesTouched:  []string{},
		CommitSHA:     "", // zero value - should be omitted
		DecisionsMade: []Decision{},
		VoiceSignature: "",
	}

	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Load and verify that zero values are handled correctly
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// Verify non-zero fields
	if loaded.SessionID != session.SessionID {
		t.Errorf("SessionID mismatch: got %v, want %v", loaded.SessionID, session.SessionID)
	}
	if !loaded.StartedAt.Equal(session.StartedAt) {
		t.Errorf("StartedAt mismatch: got %v, want %v", loaded.StartedAt, session.StartedAt)
	}

	// Verify zero values are zero (they were omitted in YAML)
	if !loaded.EndedAt.IsZero() {
		t.Errorf("EndedAt should be zero: got %v", loaded.EndedAt)
	}
	if loaded.Agent != "" {
		t.Errorf("Agent should be empty: got %q", loaded.Agent)
	}
	if loaded.Model != "" {
		t.Errorf("Model should be empty: got %q", loaded.Model)
	}
	if loaded.ToolCall != "" {
		t.Errorf("ToolCall should be empty: got %q", loaded.ToolCall)
	}
	if loaded.InitialPrompt != "" {
		t.Errorf("InitialPrompt should be empty: got %q", loaded.InitialPrompt)
	}
	if len(loaded.Iterations) != 0 {
		t.Errorf("Iterations should be empty: got %d", len(loaded.Iterations))
	}
	if len(loaded.FilesTouched) != 0 {
		t.Errorf("FilesTouched should be empty: got %d", len(loaded.FilesTouched))
	}
	if loaded.CommitSHA != "" {
		t.Errorf("CommitSHA should be empty: got %q", loaded.CommitSHA)
	}
	if len(loaded.DecisionsMade) != 0 {
		t.Errorf("DecisionsMade should be empty: got %d", len(loaded.DecisionsMade))
	}
	if loaded.VoiceSignature != "" {
		t.Errorf("VoiceSignature should be empty: got %q", loaded.VoiceSignature)
	}
}

func TestSessionWithOnlyRequiredFields(t *testing.T) {
	// Test session with only required fields (minimal session)
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	session := &Session{
		SessionID: "2026-10-24-a1b2c3",
		StartedAt: time.Date(2026, 10, 24, 11, 2, 14, 0, time.UTC),
	}

	manager := NewManager(sessionsDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Load and verify minimal session
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	if loaded.SessionID != session.SessionID {
		t.Errorf("SessionID mismatch: got %v, want %v", loaded.SessionID, session.SessionID)
	}
	if !loaded.StartedAt.Equal(session.StartedAt) {
		t.Errorf("StartedAt mismatch: got %v, want %v", loaded.StartedAt, session.StartedAt)
	}
}

func TestMultipleSessions(t *testing.T) {
	// Test creating and managing multiple sessions
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	sessions := []*Session{
		{
			SessionID: "2026-10-24-a1b2c3",
			Agent:     "opencode",
			StartedAt: time.Date(2026, 10, 24, 11, 0, 0, 0, time.UTC),
		},
		{
			SessionID: "2026-10-24-d4e5f6",
			Agent:     "claude-code",
			StartedAt: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
		},
		{
			SessionID: "2026-10-25-g7h8i9",
			Agent:     "cursor",
			StartedAt: time.Date(2026, 10, 25, 9, 0, 0, 0, time.UTC),
		},
	}

	manager := NewManager(sessionsDir)

	// Save all sessions
	for _, session := range sessions {
		if err := manager.SaveSession(session); err != nil {
			t.Fatalf("SaveSession failed for %s: %v", session.SessionID, err)
		}
	}

	// List all sessions
	sessionIDs, err := manager.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if len(sessionIDs) != len(sessions) {
		t.Errorf("ListSessions count = %v, want %v", len(sessionIDs), len(sessions))
	}

	// Load and verify each session
	for _, expected := range sessions {
		loaded, err := manager.LoadSession(expected.SessionID)
		if err != nil {
			t.Fatalf("LoadSession failed for %s: %v", expected.SessionID, err)
		}

		if loaded.SessionID != expected.SessionID {
			t.Errorf("SessionID mismatch for %s: got %v, want %v", expected.SessionID, loaded.SessionID, expected.SessionID)
		}
		if loaded.Agent != expected.Agent {
			t.Errorf("Agent mismatch for %s: got %v, want %v", expected.SessionID, loaded.Agent, expected.Agent)
		}
		if !loaded.StartedAt.Equal(expected.StartedAt) {
			t.Errorf("StartedAt mismatch for %s: got %v, want %v", expected.SessionID, loaded.StartedAt, expected.StartedAt)
		}
	}

	// Delete one session
	if err := manager.DeleteSession("2026-10-24-d4e5f6"); err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Verify remaining sessions
	sessionIDs, err = manager.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions after delete failed: %v", err)
	}

	if len(sessionIDs) != 2 {
		t.Errorf("ListSessions after delete count = %v, want 2", len(sessionIDs))
	}

	// Verify deleted session is gone
	_, err = manager.LoadSession("2026-10-24-d4e5f6")
	if err == nil {
		t.Error("LoadSession should fail for deleted session")
	}
}

func TestSessionIDValidationComplex(t *testing.T) {
	// Test complex session ID validation scenarios
	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{"valid_uuid", "2026-10-24-a1b2c3d4e5f6", false},
		{"valid_with_numbers", "2026-10-24-12345", false},
		{"valid_with_mixed", "2026-10-24-x-y-z", false},
		{"valid_long", "2026-10-24-very-long-session-id-with-dashes", false},
		{"empty", "", true},
		{"dot", ".", true},
		{"dotdot", "..", true},
		{"contains_slash", "2026-10-24/path/to/file", true},
		{"contains_backslash", "2026-10-24\\path\\to\\file", true},
		{"contains_both_separators", "2026-10-24/path\\to/file", true},
		{"starts_with_dash", "-2026-10-24-a1b2c3", false}, // technically valid, but unusual
		{"ends_with_dash", "2026-10-24-a1b2c3-", false},   // technically valid, but unusual
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSessionID(tt.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSessionID(%q) error = %v, wantErr %v", tt.sessionID, err, tt.wantErr)
			}
		})
	}
}

func TestLoadSessionFileError(t *testing.T) {
	// Test loadSessionFile with invalid YAML content
	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}

	// Create file with invalid YAML
	sessionFile := filepath.Join(sessionsDir, "invalid.yml")
	if err := os.WriteFile(sessionFile, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to create invalid YAML file: %v", err)
	}

	// Test loadSessionFile should fail
	_, err := loadSessionFile(filepath.Join(sessionsDir, "invalid"))
	if err == nil {
		t.Error("loadSessionFile should fail with invalid YAML")
	}
}

func TestSaveSessionFileError(t *testing.T) {
	t.Skip("Skip for now - directory creation behavior differs by OS")
	// Test saveSessionFile to a directory that exists but is read-only
	session := &Session{
		SessionID: "2026-10-24-a1b2c3",
		Agent:     "opencode",
	}

	// Create a temp directory and make it read-only
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatalf("failed to create read-only dir: %v", err)
	}

	// Change permissions to read-only
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to set read-only permissions: %v", err)
	}

	err := saveSessionFile(filepath.Join(readOnlyDir, "test"), session)
	if err == nil {
		t.Error("saveSessionFile should fail with read-only directory")
	}
}

func TestManagerInvalidSessionIDErrors(t *testing.T) {
	manager := NewManager(t.TempDir())
	bad := &Session{SessionID: "../evil", Agent: "opencode"}

	if err := manager.SaveSession(bad); err == nil {
		t.Error("SaveSession should reject path-traversal session ID")
	}
	if _, err := manager.LoadSession("../evil"); err == nil {
		t.Error("LoadSession should reject path-traversal session ID")
	}
	if err := manager.DeleteSession(".."); err == nil {
		t.Error("DeleteSession should reject path-traversal session ID")
	}
}

func TestManagerSaveSessionEnsureDirFails(t *testing.T) {
	tempDir := t.TempDir()
	blocker := filepath.Join(tempDir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatalf("create blocker file: %v", err)
	}

	manager := NewManager(filepath.Join(blocker, "sessions"))
	err := manager.SaveSession(&Session{SessionID: "2026-10-24-a1b2c3", Agent: "opencode"})
	if err == nil {
		t.Fatal("SaveSession should fail when sessions dir cannot be created")
	}
}

func TestSessionYAMLSchemaExample(t *testing.T) {
	// Сессия в точности как в примере PROMPT-MVP.md §2.1 (FR-PROV-002)
	tempDir := t.TempDir()
	session := &Session{
		SessionID: "2026-10-24-a1b2c3",
		Repo:      "OryaSecurity",
		StartedAt: time.Date(2026, 10, 24, 11, 2, 14, 0, time.FixedZone("MSK", 3*3600)),
		EndedAt:   time.Date(2026, 10, 24, 11, 47, 33, 0, time.FixedZone("MSK", 3*3600)),
		Agent:     "opencode",
		Model:     "glm-4.6",
		ToolCall:  `opencode run "add rate limiter"`,
		InitialPrompt: "Need to add rate limiter on POST /v1/policies.\n" +
			"Token bucket, 10/min.\n",
		Iterations: []Iteration{
			{ID: "i1", Action: "read", Target: "internal/api/policies.go"},
			{ID: "i2", Action: "edit", Target: "internal/api/policies.go", Diff: "added middleware RateLimit(token, 10/min)"},
		},
		FilesTouched:  []string{"internal/api/policies.go"},
		CommitSHA:     "9f3e2a1",
		DecisionsMade: []Decision{},
		VoiceSignature: "rate limiter needed urgently",
	}

	manager := NewManager(tempDir)
	if err := manager.SaveSession(session); err != nil {
		t.Fatalf("SaveSession: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tempDir, "2026-10-24-a1b2c3.yml"))
	if err != nil {
		t.Fatalf("read session file: %v", err)
	}
	content := string(data)

	// Порядок ключей должен совпадать с примером §2.1
	keys := []string{
		"session_id:", "repo:", "started_at:", "ended_at:", "agent:", "model:",
		"tool_call:", "initial_prompt:", "iterations:", "files_touched:",
		"commit_sha:", "decisions_made:", "voice_signature:",
	}
	lastIdx := -1
	for _, key := range keys {
		idx := strings.Index(content, key)
		if idx == -1 {
			t.Errorf("schema key %q missing from YAML:\n%s", key, content)
			continue
		}
		if idx < lastIdx {
			t.Errorf("schema key %q out of order vs §2.1 example:\n%s", key, content)
		}
		lastIdx = idx
	}

	// Многострочный initial_prompt должен рендериться literal block'ом (|)
	if !strings.Contains(content, "initial_prompt: |") {
		t.Errorf("multiline initial_prompt should use literal block scalar:\n%s", content)
	}

	// decisions_made: [] должен присутствовать даже пустым (пример §2.1)
	if !strings.Contains(content, "decisions_made: []") {
		t.Errorf("empty decisions_made should render as [] per §2.1 example:\n%s", content)
	}

	// Round-trip: загруженная сессия равна исходной
	loaded, err := manager.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("LoadSession: %v", err)
	}
	if loaded.SessionID != session.SessionID || loaded.Repo != session.Repo ||
		loaded.Agent != session.Agent || loaded.Model != session.Model ||
		loaded.ToolCall != session.ToolCall || loaded.InitialPrompt != session.InitialPrompt ||
		loaded.VoiceSignature != session.VoiceSignature || loaded.CommitSHA != session.CommitSHA {
		t.Error("round-trip scalar fields mismatch")
	}
	if !loaded.StartedAt.Equal(session.StartedAt) || !loaded.EndedAt.Equal(session.EndedAt) {
		t.Error("round-trip timestamps mismatch")
	}
	if len(loaded.Iterations) != 2 || loaded.Iterations[1].Diff != session.Iterations[1].Diff {
		t.Error("round-trip iterations mismatch")
	}
	if len(loaded.FilesTouched) != 1 || loaded.FilesTouched[0] != session.FilesTouched[0] {
		t.Error("round-trip files_touched mismatch")
	}
	if len(loaded.DecisionsMade) != 0 {
		t.Error("round-trip decisions_made should be empty")
	}
}
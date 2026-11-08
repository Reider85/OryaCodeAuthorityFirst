// Package monitoring provides comprehensive monitoring infrastructure for real-world testing
// of CodeAuthorityFirst exit criteria 1,2,3,9,10,11,12 §12.
package monitoring

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionMonitor tracks session creation and metadata collection
type SessionMonitor struct {
	sessionDir     string
	logFile        *os.File
	mu             sync.Mutex
	sessions       map[string]*SessionMetadata
	startTime      time.Time
	totalSessions  int
	errorCount     int
	performance    PerformanceMetrics
}

// SessionMetadata contains session tracking metadata
type SessionMetadata struct {
	ID           string            `json:"id"`
	Agent        string            `json:"agent"`
	Model        string            `json:"model"`
	Started      time.Time         `json:"started"`
	Duration     time.Duration     `json:"duration"`
	FilesTouched []string          `json:"files_touched"`
	Iterations   []Iteration       `json:"iterations"`
	Status       string            `json:"status"`
	Error        string            `json:"error,omitempty"`
	Metadata     map[string]string `json:"metadata"`
}

// Iteration represents a single session iteration
type Iteration struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Files     []string  `json:"files"`
	Message   string    `json:"message"`
	Duration  time.Duration `json:"duration"`
}

// PerformanceMetrics tracks performance data
type PerformanceMetrics struct {
	AvgSessionCreation time.Duration `json:"avg_session_creation"`
	MaxSessionCreation time.Duration `json:"max_session_creation"`
	MinSessionCreation time.Duration `json:"min_session_creation"`
	AvgHookExecution  time.Duration `json:"avg_hook_execution"`
	MaxHookExecution  time.Duration `json:"max_hook_execution"`
	MinHookExecution  time.Duration `json:"min_hook_execution"`
	TotalSessions     int           `json:"total_sessions"`
	SuccessRate       float64       `json:"success_rate"`
	ErrorRate         float64       `json:"error_rate"`
}

// NewSessionMonitor creates a new session monitor
func NewSessionMonitor(sessionDir string) (*SessionMonitor, error) {
	// Ensure session directory exists
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	// Create log file
	logFile, err := os.Create(filepath.Join(sessionDir, "session-monitor.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	monitor := &SessionMonitor{
		sessionDir: sessionDir,
		logFile:    logFile,
		sessions:   make(map[string]*SessionMetadata),
		startTime:  time.Now(),
	}

	// Load existing sessions
	if err := monitor.loadExistingSessions(); err != nil {
		log.Printf("Warning: failed to load existing sessions: %v", err)
	}

	return monitor, nil
}

// loadExistingSessions loads existing session data
func (m *SessionMonitor) loadExistingSessions() error {
	files, err := os.ReadDir(m.sessionDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		if filepath.Ext(file.Name()) == ".yml" {
			sessionPath := filepath.Join(m.sessionDir, file.Name())
			metadata, err := m.parseSessionFile(sessionPath)
			if err != nil {
				log.Printf("Warning: failed to parse session file %s: %v", sessionPath, err)
				continue
			}
			m.sessions[metadata.ID] = metadata
			m.totalSessions++
		}
	}

	return nil
}

// parseSessionFile parses a session YAML file
func (m *SessionMonitor) parseSessionFile(filePath string) (*SessionMetadata, error) {
	// This would normally parse YAML files
	// For now, create a basic metadata structure
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &SessionMetadata{
		ID:        filepath.Base(filePath),
		Started:   fileInfo.ModTime(),
		Duration:  time.Since(fileInfo.ModTime()),
		Status:    "completed",
		Metadata:  make(map[string]string),
	}, nil
}

// TrackSessionStart tracks the beginning of a new session
func (m *SessionMonitor) TrackSessionStart(sessionID, agent, model string) *SessionMetadata {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &SessionMetadata{
		ID:        sessionID,
		Agent:     agent,
		Model:     model,
		Started:  time.Now(),
		Status:    "active",
		Metadata:  make(map[string]string),
	}

	m.sessions[sessionID] = session
	m.totalSessions++

	log.Printf("Session started: %s (agent: %s, model: %s)", sessionID, agent, model)
	m.logEvent("session_start", sessionID, map[string]interface{}{
		"agent": agent,
		"model": model,
	})

	return session
}

// TrackSessionEnd tracks the end of a session
func (m *SessionMonitor) TrackSessionEnd(sessionID string, error error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		log.Printf("Warning: session %s not found for end tracking", sessionID)
		return
	}

	session.Status = "completed"
	if error != nil {
		session.Status = "error"
		session.Error = error.Error()
		m.errorCount++
	}

	session.Duration = time.Since(session.Started)

	log.Printf("Session ended: %s (status: %s, duration: %v)", sessionID, session.Status, session.Duration)
	m.logEvent("session_end", sessionID, map[string]interface{}{
		"status":   session.Status,
		"duration": session.Duration.String(),
		"error":    error,
	})

	// Update performance metrics
	m.updatePerformanceMetrics(session)
}

// TrackIteration tracks a session iteration
func (m *SessionMonitor) TrackIteration(sessionID, action string, files []string, message string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		log.Printf("Warning: session %s not found for iteration tracking", sessionID)
		return
	}

	iteration := Iteration{
		ID:        fmt.Sprintf("i%d", len(session.Iterations)+1),
		Action:    action,
		Timestamp: time.Now(),
		Files:     files,
		Message:   message,
		Duration:  duration,
	}

	session.Iterations = append(session.Iterations, iteration)

	log.Printf("Iteration tracked: %s.%s (action: %s, files: %v)", sessionID, iteration.ID, action, files)
	m.logEvent("iteration", sessionID, map[string]interface{}{
		"action":   action,
		"files":    files,
		"message":  message,
		"duration": duration.String(),
	})
}

// TrackFileTouch tracks when a file is touched during a session
func (m *SessionMonitor) TrackFileTouch(sessionID, filePath string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[sessionID]
	if !exists {
		log.Printf("Warning: session %s not found for file touch tracking", sessionID)
		return
	}

	// Check if file already tracked
	for _, file := range session.FilesTouched {
		if file == filePath {
			return
		}
	}

	session.FilesTouched = append(session.FilesTouched, filePath)
	m.logEvent("file_touch", sessionID, map[string]interface{}{
		"file": filePath,
	})
}

// GetSessionStats returns session statistics
func (m *SessionMonitor) GetSessionStats() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	activeSessions := 0
	completedSessions := 0
	errorSessions := 0

	for _, session := range m.sessions {
		switch session.Status {
		case "active":
			activeSessions++
		case "completed":
			completedSessions++
		case "error":
			errorSessions++
		}
	}

	return map[string]interface{}{
		"total_sessions":      m.totalSessions,
		"active_sessions":    activeSessions,
		"completed_sessions": completedSessions,
		"error_sessions":      errorSessions,
		"error_rate":         float64(m.errorCount) / float64(m.totalSessions),
		"uptime":             time.Since(m.startTime).String(),
	}
}

// GetPerformanceMetrics returns performance metrics
func (m *SessionMonitor) GetPerformanceMetrics() PerformanceMetrics {
	return m.performance
}

// GenerateReport generates a comprehensive session monitoring report
func (m *SessionMonitor) GenerateReport() map[string]interface{} {
	stats := m.GetSessionStats()
	
	return map[string]interface{}{
		"timestamp":           time.Now().Format(time.RFC3339),
		"session_stats":       stats,
		"performance_metrics": m.performance,
		"monitoring_uptime":   time.Since(m.startTime).String(),
		"session_details":     m.getSessionDetails(),
	}
}

// getSessionDetails returns detailed session information
func (m *SessionMonitor) getSessionDetails() []map[string]interface{} {
	details := make([]map[string]interface{}, 0, len(m.sessions))
	
	for _, session := range m.sessions {
		sessionDetail := map[string]interface{}{
			"id":         session.ID,
			"agent":      session.Agent,
			"model":      session.Model,
			"started":    session.Started.Format(time.RFC3339),
			"duration":   session.Duration.String(),
			"status":     session.Status,
			"files_touched": len(session.FilesTouched),
			"iterations": len(session.Iterations),
			"error":      session.Error,
		}
		details = append(details, sessionDetail)
	}
	
	return details
}

// updatePerformanceMetrics updates performance metrics
func (m *SessionMonitor) updatePerformanceMetrics(session *SessionMetadata) {
	// Update session creation metrics (simplified for now)
	if m.performance.TotalSessions == 0 {
		m.performance.AvgSessionCreation = session.Duration
		m.performance.MaxSessionCreation = session.Duration
		m.performance.MinSessionCreation = session.Duration
	} else {
		m.performance.AvgSessionCreation = (m.performance.AvgSessionCreation*time.Duration(m.performance.TotalSessions) + session.Duration) / time.Duration(m.performance.TotalSessions+1)
		if session.Duration > m.performance.MaxSessionCreation {
			m.performance.MaxSessionCreation = session.Duration
		}
		if session.Duration < m.performance.MinSessionCreation {
			m.performance.MinSessionCreation = session.Duration
		}
	}

	m.performance.TotalSessions = m.totalSessions
	m.performance.SuccessRate = float64(m.totalSessions-m.errorCount) / float64(m.totalSessions)
	m.performance.ErrorRate = float64(m.errorCount) / float64(m.totalSessions)
}

// logEvent logs a monitoring event
func (m *SessionMonitor) logEvent(eventType, sessionID string, data map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"event_type":  eventType,
		"session_id":  sessionID,
		"data":        data,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	// Write to log file
	if _, err := m.logFile.WriteString(string(jsonData) + "\n"); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}
}

// Close closes the session monitor and cleans up resources
func (m *SessionMonitor) Close() error {
	if m.logFile != nil {
		return m.logFile.Close()
	}
	return nil
}
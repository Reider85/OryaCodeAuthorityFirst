package security

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// EventType represents the type of security event
type EventType string

const (
	EventTypeYAMLParse      EventType = "yaml_parse"
	EventTypeSecretScan     EventType = "secret_scan"
	EventTypeConfigLoad    EventType = "config_load"
	EventTypeAuthCheck      EventType = "auth_check"
	EventTypePolicyCheck    EventType = "policy_check"
	EventTypeDataAccess     EventType = "data_access"
	EventTypeError          EventType = "error"
	EventTypeWarning        EventType = "warning"
	EventTypeInfo           EventType = "info"
)

// EventLevel represents the severity level of an event
type EventLevel string

const (
	LevelDebug   EventLevel = "debug"
	LevelInfo    EventLevel = "info"
	LevelWarning EventLevel = "warning"
	LevelError   EventLevel = "error"
	LevelCritical EventLevel = "critical"
)

// SecurityEvent represents a security-related event for audit logging
type SecurityEvent struct {
	Timestamp    time.Time     `json:"timestamp"`
	EventType    EventType     `json:"event_type"`
	Level        EventLevel    `json:"level"`
	Message      string        `json:"message"`
	Details      map[string]interface{} `json:"details"`
	UserID       string        `json:"user_id,omitempty"`
	SessionID    string        `json:"session_id,omitempty"`
	Resource     string        `json:"resource,omitempty"`
	Action       string        `json:"action,omitempty"`
	Result       string        `json:"result,omitempty"`
	SourceIP     string        `json:"source_ip,omitempty"`
	UserAgent    string        `json:"user_agent,omitempty"`
	CorrelationID string       `json:"correlation_id,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
}

// AuditLogger handles security event logging and storage
type AuditLogger struct {
	// Configuration
	Enabled          bool
	LogToFile        bool
	LogToConsole     bool
	LogLevel         EventLevel
	MaxFileSize      int64
	MaxFiles         int
	RetainDays       int
	OutputDir        string
	EnableJSON       bool
	EnableStructured bool
	
	// Internal state
	file          *os.File
	fileMutex     sync.Mutex
	currentSize   int64
	eventCount    int64
	startTime     time.Time
	shutdownChan  chan struct{}
}

// NewAuditLogger creates a new audit logger with default settings
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		Enabled:          true,
		LogToFile:        true,
		LogToConsole:     true,
		LogLevel:         LevelInfo,
		MaxFileSize:      10 * 1024 * 1024, // 10MB
		MaxFiles:         5,
		RetainDays:       30,
		OutputDir:        "audit_logs",
		EnableJSON:       true,
		EnableStructured: true,
		startTime:        time.Now(),
		shutdownChan:     make(chan struct{}),
	}
}

// LogEvent logs a security event
func (l *AuditLogger) LogEvent(event SecurityEvent) error {
	if !l.Enabled {
		return nil
	}

	// Check event level against minimum level
	if !l.shouldLog(event.Level) {
		return nil
	}

	// Add timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Format the event message
	message, err := l.formatEvent(event)
	if err != nil {
		return fmt.Errorf("failed to format event: %w", err)
	}

	// Log to console if enabled
	if l.LogToConsole {
		l.logToConsole(message, event.Level)
	}

	// Log to file if enabled
	if l.LogToFile {
		if err := l.logToFile(message, event); err != nil {
			return fmt.Errorf("failed to log to file: %w", err)
		}
	}

	// Update event count
	l.eventCount++

	return nil
}

// LogYAMLParse logs YAML parsing events
func (l *AuditLogger) LogYAMLParse(filePath string, success bool, duration time.Duration, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeYAMLParse,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("YAML parsing %s", map[bool]string{true: "succeeded", false: "failed"}[success]),
		Details: map[string]interface{}{
			"file_path": filePath,
			"success":   success,
			"duration":  duration.Milliseconds(),
		},
	}

	if success {
		event.Level = LevelDebug
	} else {
		event.Level = LevelError
	}

	return l.LogEvent(event)
}

// LogSecretScan logs secret scanning events
func (l *AuditLogger) LogSecretScan(filePath string, secretsFound int, duration time.Duration, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeSecretScan,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("Secret scan completed: %d secrets found", secretsFound),
		Details: map[string]interface{}{
			"file_path":    filePath,
			"secrets_found": secretsFound,
			"duration":     duration.Milliseconds(),
		},
	}

	if secretsFound > 0 {
		event.Level = LevelWarning
		event.Tags = []string{"secret_detected", "security_issue"}
	} else {
		event.Level = LevelDebug
	}

	return l.LogEvent(event)
}

// LogConfigLoad logs configuration loading events
func (l *AuditLogger) LogConfigLoad(configPath string, success bool, duration time.Duration, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeConfigLoad,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("Configuration loading %s", map[bool]string{true: "succeeded", false: "failed"}[success]),
		Details: map[string]interface{}{
			"config_path": configPath,
			"success":     success,
			"duration":    duration.Milliseconds(),
		},
	}

	if success {
		event.Level = LevelDebug
	} else {
		event.Level = LevelError
	}

	return l.LogEvent(event)
}

// LogAuthCheck logs authentication/authorization events
func (l *AuditLogger) LogAuthCheck(userID, resource, action string, success bool, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeAuthCheck,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("Auth check for %s %s %s", action, resource, map[bool]string{true: "succeeded", false: "failed"}[success]),
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Result:    map[bool]string{true: "success", false: "failure"}[success],
		Details:   details,
	}

	if success {
		event.Level = LevelDebug
	} else {
		event.Level = LevelWarning
		event.Tags = []string{"auth_failure", "security_issue"}
	}

	return l.LogEvent(event)
}

// LogPolicyCheck logs policy enforcement events
func (l *AuditLogger) LogPolicyCheck(policyName, resource string, success bool, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypePolicyCheck,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("Policy check for %s on %s %s", policyName, resource, map[bool]string{true: "succeeded", false: "failed"}[success]),
		Resource:  resource,
		Result:    map[bool]string{true: "success", false: "failure"}[success],
		Details:   details,
	}

	if success {
		event.Level = LevelDebug
	} else {
		event.Level = LevelWarning
		event.Tags = []string{"policy_violation", "security_issue"}
	}

	return l.LogEvent(event)
}

// LogDataAccess logs data access events
func (l *AuditLogger) LogDataAccess(userID, resource, action string, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeDataAccess,
		Level:     LevelInfo,
		Message:   fmt.Sprintf("Data access: %s %s by %s", action, resource, userID),
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Details:   details,
	}

	return l.LogEvent(event)
}

// LogError logs error events
func (l *AuditLogger) LogError(message string, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeError,
		Level:     LevelError,
		Message:   message,
		Details:   details,
		Tags:      []string{"error", "security_issue"},
	}

	return l.LogEvent(event)
}

// LogWarning logs warning events
func (l *AuditLogger) LogWarning(message string, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeWarning,
		Level:     LevelWarning,
		Message:   message,
		Details:   details,
		Tags:      []string{"warning", "security_issue"},
	}

	return l.LogEvent(event)
}

// LogInfo logs informational events
func (l *AuditLogger) LogInfo(message string, details map[string]interface{}) error {
	event := SecurityEvent{
		EventType: EventTypeInfo,
		Level:     LevelInfo,
		Message:   message,
		Details:   details,
	}

	return l.LogEvent(event)
}

// shouldLog determines if an event should be logged based on level
func (l *AuditLogger) shouldLog(level EventLevel) bool {
	switch l.LogLevel {
	case LevelDebug:
		return true
	case LevelInfo:
		return level != LevelDebug
	case LevelWarning:
		return level == LevelWarning || level == LevelError || level == LevelCritical
	case LevelError:
		return level == LevelError || level == LevelCritical
	case LevelCritical:
		return level == LevelCritical
	default:
		return true
	}
}

// formatEvent formats a security event for output
func (l *AuditLogger) formatEvent(event SecurityEvent) (string, error) {
	if l.EnableJSON {
		return l.formatJSON(event)
	}
	return l.formatText(event)
}

// formatJSON formats an event as JSON
func (l *AuditLogger) formatJSON(event SecurityEvent) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// formatText formats an event as human-readable text
func (l *AuditLogger) formatText(event SecurityEvent) (string, error) {
	return fmt.Sprintf("[%s] %s - %s: %s", 
		event.Timestamp.Format(time.RFC3339), 
		event.Level, 
		event.EventType, 
		event.Message), nil
}

// logToConsole logs an event to the console
func (l *AuditLogger) logToConsole(message string, level EventLevel) {
	log.Printf("%s", message)
}

// logToFile logs an event to the audit log file
func (l *AuditLogger) logToFile(message string, event SecurityEvent) error {
	l.fileMutex.Lock()
	defer l.fileMutex.Unlock()

	// Ensure output directory exists
	if err := os.MkdirAll(l.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open or rotate log file if needed
	if l.file == nil || l.currentSize >= l.MaxFileSize {
		if err := l.rotateLogFile(); err != nil {
			return err
		}
	}

	// Write to file
	if _, err := l.file.WriteString(message + "\n"); err != nil {
		return err
	}

	// Flush to ensure immediate write
	if err := l.file.Sync(); err != nil {
		return err
	}

	l.currentSize += int64(len(message)) + 1

	return nil
}

// rotateLogFile rotates the log file when it reaches size limit
func (l *AuditLogger) rotateLogFile() error {
	// Close existing file if open
	if l.file != nil {
		l.file.Close()
	}

	// Create new log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("audit_%s.log", timestamp)
	logFilePath := filepath.Join(l.OutputDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.file = file
	l.currentSize = 0

	// Clean up old log files
	go l.cleanupOldLogs()

	return nil
}

// cleanupOldLogs removes log files older than RetainDays
func (l *AuditLogger) cleanupOldLogs() {
	files, err := os.ReadDir(l.OutputDir)
	if err != nil {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -l.RetainDays)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			os.Remove(filepath.Join(l.OutputDir, file.Name()))
		}
	}
}

// GetStats returns audit logging statistics
func (l *AuditLogger) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"enabled":     l.Enabled,
		"total_events": l.eventCount,
		"start_time":  l.startTime,
		"uptime":      time.Since(l.startTime),
		"log_level":   l.LogLevel,
		"log_to_file": l.LogToFile,
		"log_to_console": l.LogToConsole,
		"max_file_size": l.MaxFileSize,
		"max_files":   l.MaxFiles,
		"retain_days": l.RetainDays,
	}
}

// Shutdown closes the audit logger and performs cleanup
func (l *AuditLogger) Shutdown() error {
	l.fileMutex.Lock()
	defer l.fileMutex.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLogLevel sets the minimum log level
func (l *AuditLogger) SetLogLevel(level EventLevel) {
	l.LogLevel = level
}

// Enable enables or disables audit logging
func (l *AuditLogger) Enable(enabled bool) {
	l.Enabled = enabled
}

// SetOutputDir sets the output directory for log files
func (l *AuditLogger) SetOutputDir(dir string) {
	l.OutputDir = dir
}

// SetMaxFileSize sets the maximum file size before rotation
func (l *AuditLogger) SetMaxFileSize(size int64) {
	l.MaxFileSize = size
}

// SetRetainDays sets the number of days to retain log files
func (l *AuditLogger) SetRetainDays(days int) {
	l.RetainDays = days
}
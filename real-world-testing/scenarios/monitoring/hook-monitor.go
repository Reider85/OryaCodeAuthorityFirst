// Package monitoring provides hook execution monitoring for real-world testing
// of CodeAuthorityFirst universal pre-commit hook across different AI agents.
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

// HookMonitor tracks hook execution across different agents and events
type HookMonitor struct {
	hookDir      string
	logFile      *os.File
	mu           sync.Mutex
	hookExecs    map[string]*HookExecution
	performance  HookPerformanceMetrics
	config       HookMonitorConfig
}

// HookExecution represents a single hook execution
type HookExecution struct {
	ID           string                 `json:"id"`
	Event        string                 `json:"event"`        // pre-session, post-edit, pre-commit
	Agent        string                 `json:"agent"`        // opencode, claude-code, etc.
	SessionID    string                 `json:"session_id,omitempty"`
	Files        []string               `json:"files,omitempty"`
	Message      string                 `json:"message,omitempty"`
	Started      time.Time              `json:"started"`
	Completed    time.Time              `json:"completed"`
	Duration     time.Duration           `json:"duration"`
	Status       string                 `json:"status"`        // success, error, timeout
	Error        string                 `json:"error,omitempty"`
	ExitCode     int                    `json:"exit_code"`
	Output       string                 `json:"output,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// HookPerformanceMetrics tracks hook performance data
type HookPerformanceMetrics struct {
	TotalExecutions     int                    `json:"total_executions"`
	SuccessExecutions   int                    `json:"success_executions"`
	ErrorExecutions     int                    `json:"error_executions"`
	TimeoutExecutions   int                    `json:"timeout_executions"`
	
	// Timing metrics
	AvgExecutionTime    time.Duration          `json:"avg_execution_time"`
	MaxExecutionTime    time.Duration          `json:"max_execution_time"`
	MinExecutionTime    time.Duration          `json:"min_execution_time"`
	
	// Event-specific metrics
	PreSessionMetrics   EventMetrics           `json:"pre_session_metrics"`
	PostEditMetrics     EventMetrics           `json:"post_edit_metrics"`
	PreCommitMetrics    EventMetrics           `json:"pre_commit_metrics"`
	
	// Agent-specific metrics
	AgentMetrics        map[string]AgentMetrics `json:"agent_metrics"`
	
	// Performance thresholds
	Violations          []PerformanceViolation `json:"violations"`
}

// EventMetrics tracks metrics for specific hook events
type EventMetrics struct {
	TotalExecutions    int           `json:"total_executions"`
	SuccessRate        float64       `json:"success_rate"`
	AvgExecutionTime   time.Duration `json:"avg_execution_time"`
	MaxExecutionTime   time.Duration `json:"max_execution_time"`
	MinExecutionTime   time.Duration `json:"min_execution_time"`
}

// AgentMetrics tracks metrics for specific agents
type AgentMetrics struct {
	TotalExecutions    int           `json:"total_executions"`
	SuccessRate        float64       `json:"success_rate"`
	AvgExecutionTime   time.Duration `json:"avg_execution_time"`
	MaxExecutionTime   time.Duration `json:"max_execution_time"`
	MinExecutionTime   time.Duration `json:"min_execution_time"`
}

// PerformanceViolation tracks performance threshold violations
type PerformanceViolation struct {
	Type        string        `json:"type"`        // execution_time, success_rate, etc.
	Threshold   time.Duration  `json:"threshold"`
	Actual      time.Duration  `json:"actual"`
	Agent       string        `json:"agent"`
	Event       string        `json:"event"`
	Timestamp   time.Time     `json:"timestamp"`
	Severity    string        `json:"severity"` // warning, error, critical
}

// HookMonitorConfig contains configuration for hook monitoring
type HookMonitorConfig struct {
	LogDir                string        `json:"log_dir"`
	MaxExecutionTime      time.Duration `json:"max_execution_time"`      // Default: 5 seconds
	MinSuccessRate        float64       `json:"min_success_rate"`        // Default: 0.95
	MaxMemoryUsage        int64         `json:"max_memory_usage"`        // Default: 100MB
	SampleRate            float64       `json:"sample_rate"`             // Default: 1.0 (100%)
	RetentionPeriod       time.Duration `json:"retention_period"`        // Default: 7 days
	EnableDetailedLogging bool          `json:"enable_detailed_logging"`
}

// NewHookMonitor creates a new hook monitor
func NewHookMonitor(hookDir string, config HookMonitorConfig) (*HookMonitor, error) {
	// Ensure hook directory exists
	if err := os.MkdirAll(hookDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create hook directory: %w", err)
	}

	// Create log file
	logFile, err := os.Create(filepath.Join(hookDir, "hook-monitor.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	monitor := &HookMonitor{
		hookDir:      hookDir,
		logFile:      logFile,
		hookExecs:    make(map[string]*HookExecution),
		performance:  HookPerformanceMetrics{},
		config:       config,
		AgentMetrics: make(map[string]AgentMetrics),
	}

	// Load existing hook executions
	if err := monitor.loadExistingExecutions(); err != nil {
		log.Printf("Warning: failed to load existing hook executions: %v", err)
	}

	return monitor, nil
}

// loadExistingExecutions loads existing hook execution data
func (m *HookMonitor) loadExistingExecutions() error {
	files, err := os.ReadDir(m.hookDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		if filepath.Ext(file.Name()) == ".json" {
			executionPath := filepath.Join(m.hookDir, file.Name())
			execution, err := m.parseExecutionFile(executionPath)
			if err != nil {
				log.Printf("Warning: failed to parse execution file %s: %v", executionPath, err)
				continue
			}
			m.hookExecs[execution.ID] = execution
			m.updatePerformanceMetrics(execution)
		}
	}

	return nil
}

// parseExecutionFile parses a hook execution JSON file
func (m *HookMonitor) parseExecutionFile(filePath string) (*HookExecution, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var execution HookExecution
	if err := json.Unmarshal(data, &execution); err != nil {
		return nil, err
	}

	return &execution, nil
}

// TrackHookStart tracks the start of a hook execution
func (m *HookMonitor) TrackHookStart(event, agent, sessionID string, files []string, message string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	executionID := fmt.Sprintf("hook_%s_%s_%d", event, agent, time.Now().UnixNano())
	
	execution := &HookExecution{
		ID:        executionID,
		Event:     event,
		Agent:     agent,
		SessionID: sessionID,
		Files:     files,
		Message:   message,
		Started:   time.Now(),
		Status:    "running",
		Metadata:  make(map[string]interface{}),
	}

	m.hookExecs[executionID] = execution

	log.Printf("Hook execution started: %s (event: %s, agent: %s)", executionID, event, agent)
	m.logEvent("hook_start", executionID, map[string]interface{}{
		"event":     event,
		"agent":     agent,
		"session_id": sessionID,
		"files":     files,
		"message":   message,
	})

	return executionID
}

// TrackHookEnd tracks the end of a hook execution
func (m *HookMonitor) TrackHookEnd(executionID string, exitCode int, output string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	execution, exists := m.hookExecs[executionID]
	if !exists {
		log.Printf("Warning: hook execution %s not found for end tracking", executionID)
		return
	}

	execution.Completed = time.Now()
	execution.Duration = execution.Completed.Sub(execution.Started)
	execution.ExitCode = exitCode
	execution.Output = output

	if err != nil {
		execution.Status = "error"
		execution.Error = err.Error()
	} else {
		execution.Status = "success"
	}

	log.Printf("Hook execution ended: %s (status: %s, duration: %v, exit_code: %d)", 
		executionID, execution.Status, execution.Duration, exitCode)
	m.logEvent("hook_end", executionID, map[string]interface{}{
		"status":    execution.Status,
		"duration":  execution.Duration.String(),
		"exit_code": exitCode,
		"error":     err,
	})

	// Update performance metrics
	m.updatePerformanceMetrics(execution)
	
	// Check for performance violations
	m.checkPerformanceViolations(execution)
}

// TrackHookTimeout tracks a hook timeout
func (m *HookMonitor) TrackHookTimeout(executionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	execution, exists := m.hookExecs[executionID]
	if !exists {
		log.Printf("Warning: hook execution %s not found for timeout tracking", executionID)
		return
	}

	execution.Status = "timeout"
	execution.Duration = m.config.MaxExecutionTime

	log.Printf("Hook execution timeout: %s (duration: %v)", executionID, execution.Duration)
	m.logEvent("hook_timeout", executionID, map[string]interface{}{
		"duration": execution.Duration.String(),
	})

	// Update performance metrics
	m.updatePerformanceMetrics(execution)
}

// GetHookStats returns hook execution statistics
func (m *HookMonitor) GetHookStats() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalExecutions := len(m.hookExecs)
	successExecutions := 0
	errorExecutions := 0
	timeoutExecutions := 0

	for _, execution := range m.hookExecs {
		switch execution.Status {
		case "success":
			successExecutions++
		case "error":
			errorExecutions++
		case "timeout":
			timeoutExecutions++
		}
	}

	successRate := 0.0
	if totalExecutions > 0 {
		successRate = float64(successExecutions) / float64(totalExecutions)
	}

	return map[string]interface{}{
		"total_executions":    totalExecutions,
		"success_executions":  successExecutions,
		"error_executions":    errorExecutions,
		"timeout_executions":  timeoutExecutions,
		"success_rate":        successRate,
		"violations":          len(m.performance.Violations),
	}
}

// GetPerformanceMetrics returns hook performance metrics
func (m *HookMonitor) GetPerformanceMetrics() HookPerformanceMetrics {
	return m.performance
}

// GetAgentPerformance returns performance metrics for specific agents
func (m *HookMonitor) GetAgentPerformance(agent string) (AgentMetrics, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, exists := m.AgentMetrics[agent]
	return metrics, exists
}

// GetEventPerformance returns performance metrics for specific events
func (m *HookMonitor) GetEventPerformance(event string) EventMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch event {
	case "pre-session":
		return m.performance.PreSessionMetrics
	case "post-edit":
		return m.performance.PostEditMetrics
	case "pre-commit":
		return m.performance.PreCommitMetrics
	default:
		return EventMetrics{}
	}
}

// GenerateReport generates a comprehensive hook monitoring report
func (m *HookMonitor) GenerateReport() map[string]interface{} {
	stats := m.GetHookStats()
	
	return map[string]interface{}{
		"timestamp":            time.Now().Format(time.RFC3339),
		"hook_stats":           stats,
		"performance_metrics":  m.performance,
		"agent_performance":   m.AgentMetrics,
		"event_performance":   map[string]EventMetrics{
			"pre-session":  m.performance.PreSessionMetrics,
			"post-edit":    m.performance.PostEditMetrics,
			"pre-commit":   m.performance.PreCommitMetrics,
		},
		"violations":           m.performance.Violations,
		"monitoring_uptime":    time.Since(m.startTime).String(),
	}
}

// updatePerformanceMetrics updates performance metrics based on a hook execution
func (m *HookMonitor) updatePerformanceMetrics(execution *HookExecution) {
	// Update total metrics
	m.performance.TotalExecutions++
	
	switch execution.Status {
	case "success":
		m.performance.SuccessExecutions++
	case "error":
		m.performance.ErrorExecutions++
	case "timeout":
		m.performance.TimeoutExecutions++
	}

	// Update timing metrics
	if m.performance.TotalExecutions == 1 {
		m.performance.AvgExecutionTime = execution.Duration
		m.performance.MaxExecutionTime = execution.Duration
		m.performance.MinExecutionTime = execution.Duration
	} else {
		totalDuration := m.performance.AvgExecutionTime*time.Duration(m.performance.TotalExecutions-1) + execution.Duration
		m.performance.AvgExecutionTime = totalDuration / time.Duration(m.performance.TotalExecutions)
		
		if execution.Duration > m.performance.MaxExecutionTime {
			m.performance.MaxExecutionTime = execution.Duration
		}
		if execution.Duration < m.performance.MinExecutionTime {
			m.performance.MinExecutionTime = execution.Duration
		}
	}

	// Update event-specific metrics
	switch execution.Event {
	case "pre-session":
		m.updateEventMetrics(&m.performance.PreSessionMetrics, execution)
	case "post-edit":
		m.updateEventMetrics(&m.performance.PostEditMetrics, execution)
	case "pre-commit":
		m.updateEventMetrics(&m.performance.PreCommitMetrics, execution)
	}

	// Update agent-specific metrics
	if _, exists := m.AgentMetrics[execution.Agent]; !exists {
		m.AgentMetrics[execution.Agent] = AgentMetrics{}
	}
	m.updateAgentMetrics(&m.AgentMetrics[execution.Agent], execution)
}

// updateEventMetrics updates event-specific metrics
func (m *HookMonitor) updateEventMetrics(metrics *EventMetrics, execution *HookExecution) {
	metrics.TotalExecutions++
	
	if metrics.TotalExecutions == 1 {
		metrics.AvgExecutionTime = execution.Duration
		metrics.MaxExecutionTime = execution.Duration
		metrics.MinExecutionTime = execution.Duration
	} else {
		totalDuration := metrics.AvgExecutionTime*time.Duration(metrics.TotalExecutions-1) + execution.Duration
		metrics.AvgExecutionTime = totalDuration / time.Duration(metrics.TotalExecutions)
		
		if execution.Duration > metrics.MaxExecutionTime {
			metrics.MaxExecutionTime = execution.Duration
		}
		if execution.Duration < metrics.MinExecutionTime {
			metrics.MinExecutionTime = execution.Duration
		}
	}

	if metrics.TotalExecutions > 0 {
		metrics.SuccessRate = float64(m.performance.SuccessExecutions) / float64(metrics.TotalExecutions)
	}
}

// updateAgentMetrics updates agent-specific metrics
func (m *HookMonitor) updateAgentMetrics(metrics *AgentMetrics, execution *HookExecution) {
	metrics.TotalExecutions++
	
	if metrics.TotalExecutions == 1 {
		metrics.AvgExecutionTime = execution.Duration
		metrics.MaxExecutionTime = execution.Duration
		metrics.MinExecutionTime = execution.Duration
	} else {
		totalDuration := metrics.AvgExecutionTime*time.Duration(metrics.TotalExecutions-1) + execution.Duration
		metrics.AvgExecutionTime = totalDuration / time.Duration(metrics.TotalExecutions)
		
		if execution.Duration > metrics.MaxExecutionTime {
			metrics.MaxExecutionTime = execution.Duration
		}
		if execution.Duration < metrics.MinExecutionTime {
			metrics.MinExecutionTime = execution.Duration
		}
	}

	if metrics.TotalExecutions > 0 {
		metrics.SuccessRate = float64(m.performance.SuccessExecutions) / float64(metrics.TotalExecutions)
	}
}

// checkPerformanceViolations checks for performance threshold violations
func (m *HookMonitor) checkPerformanceViolations(execution *HookExecution) {
	// Check execution time violations
	if execution.Duration > m.config.MaxExecutionTime {
		violation := PerformanceViolation{
			Type:      "execution_time",
			Threshold: m.config.MaxExecutionTime,
			Actual:    execution.Duration,
			Agent:     execution.Agent,
			Event:     execution.Event,
			Timestamp: time.Now(),
			Severity:  "warning",
		}
		m.performance.Violations = append(m.performance.Violations, violation)
	}

	// Check success rate violations (simplified - would need historical data)
	if execution.Status == "error" {
		violation := PerformanceViolation{
			Type:      "success_rate",
			Threshold: time.Duration(0), // Not applicable for success rate
			Actual:    time.Duration(0),
			Agent:     execution.Agent,
			Event:     execution.Event,
			Timestamp: time.Now(),
			Severity:  "error",
		}
		m.performance.Violations = append(m.performance.Violations, violation)
	}
}

// logEvent logs a monitoring event
func (m *HookMonitor) logEvent(eventType, executionID string, data map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"event_type":  eventType,
		"execution_id": executionID,
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

// Close closes the hook monitor and cleans up resources
func (m *HookMonitor) Close() error {
	if m.logFile != nil {
		return m.logFile.Close()
	}
	return nil
}
// Package monitoring provides LSP performance monitoring for real-world testing
// of CodeAuthorityFirst LSP fallback with large repositories (up to 1000 files).
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

// LSPMonitor tracks LSP performance and symbol graph generation
type LSPMonitor struct {
	lspDir       string
	logFile      *os.File
	mu           sync.Mutex
	symbolGraphs map[string]*SymbolGraphMetrics
	performance  LSPPerformanceMetrics
	config       LSPMonitorConfig
}

// SymbolGraphMetrics contains metrics for symbol graph generation
type SymbolGraphMetrics struct {
	ID              string                 `json:"id"`
	RepositoryPath  string                 `json:"repository_path"`
	FileCount       int                    `json:"file_count"`
	SymbolCount     int                    `json:"symbol_count"`
	GeneratedAt     time.Time              `json:"generated_at"`
	Duration        time.Duration          `json:"duration"`
	Status          string                 `json:"status"`        // success, error, timeout
	Error           string                 `json:"error,omitempty"`
	MemoryUsage     int64                  `json:"memory_usage"`
	Symbols         []SymbolInfo          `json:"symbols"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// SymbolInfo contains information about a single symbol
type SymbolInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	FilePath  string `json:"file_path"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Signature string `json:"signature,omitempty"`
}

// LSPPerformanceMetrics tracks LSP performance data
type LSPPerformanceMetrics struct {
	TotalGraphs      int                     `json:"total_graphs"`
	SuccessGraphs    int                     `json:"success_graphs"`
	ErrorGraphs      int                     `json:"error_graphs"`
	TimeoutGraphs    int                     `json:"timeout_graphs"`
	
	// Performance metrics
	AvgGenerationTime time.Duration          `json:"avg_generation_time"`
	MaxGenerationTime time.Duration          `json:"max_generation_time"`
	MinGenerationTime time.Duration          `json:"min_generation_time"`
	
	// Memory usage metrics
	AvgMemoryUsage    int64                   `json:"avg_memory_usage"`
	MaxMemoryUsage    int64                   `json:"max_memory_usage"`
	MinMemoryUsage    int64                   `json:"min_memory_usage"`
	
	// Repository size metrics
	SmallRepoMetrics  RepoSizeMetrics         `json:"small_repo_metrics"`
	MediumRepoMetrics RepoSizeMetrics         `json:"medium_repo_metrics"`
	LargeRepoMetrics  RepoSizeMetrics         `json:"large_repo_metrics"`
	
	// Accuracy metrics
	SymbolAccuracy    float64                 `json:"symbol_accuracy"`
	AccuracyChecks    []AccuracyCheck         `json:"accuracy_checks"`
	
	// Violations
	Violations       []LSPViolation          `json:"violations"`
}

// RepoSizeMetrics tracks metrics for different repository sizes
type RepoSizeMetrics struct {
	Count            int           `json:"count"`
	AvgGenerationTime time.Duration `json:"avg_generation_time"`
	AvgMemoryUsage   int64         `json:"avg_memory_usage"`
	SuccessRate      float64       `json:"success_rate"`
}

// AccuracyCheck tracks symbol accuracy validation
type AccuracyCheck struct {
	Timestamp    time.Time  `json:"timestamp"`
	Expected     int        `json:"expected"`
	Actual       int        `json:"actual"`
	Accuracy     float64    `json:"accuracy"`
	Repository   string     `json:"repository"`
	SampleSize   int        `json:"sample_size"`
	Errors       []string   `json:"errors"`
}

// LSPViolation tracks LSP performance violations
type LSPViolation struct {
	Type        string        `json:"type"`        // generation_time, memory_usage, symbol_accuracy
	Threshold   interface{}   `json:"threshold"`  // time.Duration, int64, float64
	Actual      interface{}   `json:"actual"`
	Repository  string        `json:"repository"`
	Timestamp   time.Time     `json:"timestamp"`
	Severity    string        `json:"severity"`    // warning, error, critical
}

// LSPMonitorConfig contains configuration for LSP monitoring
type LSPMonitorConfig struct {
	LogDir                string        `json:"log_dir"`
	MaxGenerationTime     time.Duration `json:"max_generation_time"`    // Default: 10 seconds
	MaxMemoryUsage        int64         `json:"max_memory_usage"`       // Default: 100MB
	MinSymbolAccuracy     float64       `json:"min_symbol_accuracy"`    // Default: 0.95
	SampleRate            float64       `json:"sample_rate"`            // Default: 1.0 (100%)
	RetentionPeriod       time.Duration `json:"retention_period"`        // Default: 7 days
	EnableDetailedLogging bool          `json:"enable_detailed_logging"`
	AccuracySampleSize    int           `json:"accuracy_sample_size"`    // Default: 100
}

// NewLSPMonitor creates a new LSP monitor
func NewLSPMonitor(lspDir string, config LSPMonitorConfig) (*LSPMonitor, error) {
	// Ensure LSP directory exists
	if err := os.MkdirAll(lspDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create LSP directory: %w", err)
	}

	// Create log file
	logFile, err := os.Create(filepath.Join(lspDir, "lsp-monitor.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	monitor := &LSPMonitor{
		lspDir:       lspDir,
		logFile:      logFile,
		symbolGraphs: make(map[string]*SymbolGraphMetrics),
		performance:  LSPPerformanceMetrics{},
		config:       config,
	}

	// Load existing symbol graphs
	if err := monitor.loadExistingGraphs(); err != nil {
		log.Printf("Warning: failed to load existing symbol graphs: %v", err)
	}

	return monitor, nil
}

// loadExistingGraphs loads existing symbol graph data
func (m *LSPMonitor) loadExistingGraphs() error {
	files, err := os.ReadDir(m.lspDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		if filepath.Ext(file.Name()) == ".json" {
			graphPath := filepath.Join(m.lspDir, file.Name())
			graph, err := m.parseGraphFile(graphPath)
			if err != nil {
				log.Printf("Warning: failed to parse graph file %s: %v", graphPath, err)
				continue
			}
			m.symbolGraphs[graph.ID] = graph
			m.updatePerformanceMetrics(graph)
		}
	}

	return nil
}

// parseGraphFile parses a symbol graph JSON file
func (m *LSPMonitor) parseGraphFile(filePath string) (*SymbolGraphMetrics, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var graph SymbolGraphMetrics
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, err
	}

	return &graph, nil
}

// TrackSymbolGraphStart tracks the start of symbol graph generation
func (m *LSPMonitor) TrackSymbolGraphStart(repositoryPath string, fileCount int) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	graphID := fmt.Sprintf("graph_%s_%d", filepath.Base(repositoryPath), time.Now().UnixNano())
	
	graph := &SymbolGraphMetrics{
		ID:             graphID,
		RepositoryPath: repositoryPath,
		FileCount:      fileCount,
		GeneratedAt:    time.Now(),
		Status:         "generating",
		Metadata:       make(map[string]interface{}),
	}

	m.symbolGraphs[graphID] = graph

	log.Printf("Symbol graph generation started: %s (repository: %s, files: %d)", graphID, repositoryPath, fileCount)
	m.logEvent("graph_start", graphID, map[string]interface{}{
		"repository_path": repositoryPath,
		"file_count":     fileCount,
	})

	return graphID
}

// TrackSymbolGraphComplete tracks the completion of symbol graph generation
func (m *LSPMonitor) TrackSymbolGraphComplete(graphID string, symbolCount int, memoryUsage int64, symbols []SymbolInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	graph, exists := m.symbolGraphs[graphID]
	if !exists {
		log.Printf("Warning: symbol graph %s not found for completion tracking", graphID)
		return
	}

	graph.Status = "success"
	graph.SymbolCount = symbolCount
	graph.MemoryUsage = memoryUsage
	graph.Symbols = symbols
	graph.Duration = time.Since(graph.GeneratedAt)

	log.Printf("Symbol graph generation completed: %s (duration: %v, symbols: %d, memory: %d bytes)", 
		graphID, graph.Duration, symbolCount, memoryUsage)
	m.logEvent("graph_complete", graphID, map[string]interface{}{
		"duration":     graph.Duration.String(),
		"symbol_count": symbolCount,
		"memory_usage": memoryUsage,
	})

	// Update performance metrics
	m.updatePerformanceMetrics(graph)
	
	// Check for performance violations
	m.checkPerformanceViolations(graph)
	
	// Validate symbol accuracy
	m.validateSymbolAccuracy(graph)
}

// TrackSymbolGraphError tracks an error in symbol graph generation
func (m *LSPMonitor) TrackSymbolGraphError(graphID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	graph, exists := m.symbolGraphs[graphID]
	if !exists {
		log.Printf("Warning: symbol graph %s not found for error tracking", graphID)
		return
	}

	graph.Status = "error"
	graph.Error = err.Error()
	graph.Duration = time.Since(graph.GeneratedAt)

	log.Printf("Symbol graph generation error: %s (duration: %v, error: %v)", graphID, graph.Duration, err)
	m.logEvent("graph_error", graphID, map[string]interface{}{
		"duration": graph.Duration.String(),
		"error":    err.Error(),
	})

	// Update performance metrics
	m.updatePerformanceMetrics(graph)
}

// TrackSymbolGraphTimeout tracks a timeout in symbol graph generation
func (m *LSPMonitor) TrackSymbolGraphTimeout(graphID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	graph, exists := m.symbolGraphs[graphID]
	if !exists {
		log.Printf("Warning: symbol graph %s not found for timeout tracking", graphID)
		return
	}

	graph.Status = "timeout"
	graph.Duration = m.config.MaxGenerationTime

	log.Printf("Symbol graph generation timeout: %s (duration: %v)", graphID, graph.Duration)
	m.logEvent("graph_timeout", graphID, map[string]interface{}{
		"duration": graph.Duration.String(),
	})

	// Update performance metrics
	m.updatePerformanceMetrics(graph)
}

// GetGraphStats returns symbol graph statistics
func (m *LSPMonitor) GetGraphStats() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalGraphs := len(m.symbolGraphs)
	successGraphs := 0
	errorGraphs := 0
	timeoutGraphs := 0

	for _, graph := range m.symbolGraphs {
		switch graph.Status {
		case "success":
			successGraphs++
		case "error":
			errorGraphs++
		case "timeout":
			timeoutGraphs++
		}
	}

	successRate := 0.0
	if totalGraphs > 0 {
		successRate = float64(successGraphs) / float64(totalGraphs)
	}

	return map[string]interface{}{
		"total_graphs":     totalGraphs,
		"success_graphs":   successGraphs,
		"error_graphs":     errorGraphs,
		"timeout_graphs":   timeoutGraphs,
		"success_rate":     successRate,
		"violations":       len(m.performance.Violations),
		"symbol_accuracy":  m.performance.SymbolAccuracy,
	}
}

// GetPerformanceMetrics returns LSP performance metrics
func (m *LSPMonitor) GetPerformanceMetrics() LSPPerformanceMetrics {
	return m.performance
}

// GetRepositorySizeMetrics returns metrics for specific repository sizes
func (m *LSPMonitor) GetRepositorySizeMetrics(sizeCategory string) RepoSizeMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch sizeCategory {
	case "small":
		return m.performance.SmallRepoMetrics
	case "medium":
		return m.performance.MediumRepoMetrics
	case "large":
		return m.performance.LargeRepoMetrics
	default:
		return RepoSizeMetrics{}
	}
}

// GenerateReport generates a comprehensive LSP monitoring report
func (m *LSPMonitor) GenerateReport() map[string]interface{} {
	stats := m.GetGraphStats()
	
	return map[string]interface{}{
		"timestamp":            time.Now().Format(time.RFC3339),
		"graph_stats":          stats,
		"performance_metrics":  m.performance,
		"repository_size_metrics": map[string]RepoSizeMetrics{
			"small":  m.performance.SmallRepoMetrics,
			"medium": m.performance.MediumRepoMetrics,
			"large":  m.performance.LargeRepoMetrics,
		},
		"violations":           m.performance.Violations,
		"monitoring_uptime":    time.Since(m.startTime).String(),
	}
}

// updatePerformanceMetrics updates performance metrics based on a symbol graph
func (m *LSPMonitor) updatePerformanceMetrics(graph *SymbolGraphMetrics) {
	// Update total metrics
	m.performance.TotalGraphs++
	
	switch graph.Status {
	case "success":
		m.performance.SuccessGraphs++
	case "error":
		m.performance.ErrorGraphs++
	case "timeout":
		m.performance.TimeoutGraphs++
	}

	// Update timing metrics
	if m.performance.TotalGraphs == 1 {
		m.performance.AvgGenerationTime = graph.Duration
		m.performance.MaxGenerationTime = graph.Duration
		m.performance.MinGenerationTime = graph.Duration
	} else {
		totalDuration := m.performance.AvgGenerationTime*time.Duration(m.performance.TotalGraphs-1) + graph.Duration
		m.performance.AvgGenerationTime = totalDuration / time.Duration(m.performance.TotalGraphs)
		
		if graph.Duration > m.performance.MaxGenerationTime {
			m.performance.MaxGenerationTime = graph.Duration
		}
		if graph.Duration < m.performance.MinGenerationTime {
			m.performance.MinGenerationTime = graph.Duration
		}
	}

	// Update memory usage metrics
	if graph.Status == "success" {
		if m.performance.SuccessGraphs == 1 {
			m.performance.AvgMemoryUsage = graph.MemoryUsage
			m.performance.MaxMemoryUsage = graph.MemoryUsage
			m.performance.MinMemoryUsage = graph.MemoryUsage
		} else {
			totalMemory := m.performance.AvgMemoryUsage*(m.performance.SuccessGraphs-1) + graph.MemoryUsage
			m.performance.AvgMemoryUsage = totalMemory / int64(m.performance.SuccessGraphs)
			
			if graph.MemoryUsage > m.performance.MaxMemoryUsage {
				m.performance.MaxMemoryUsage = graph.MemoryUsage
			}
			if graph.MemoryUsage < m.performance.MinMemoryUsage {
				m.performance.MinMemoryUsage = graph.MemoryUsage
			}
		}
	}

	// Update repository size metrics
	m.updateRepoSizeMetrics(graph)
}

// updateRepoSizeMetrics updates repository size-specific metrics
func (m *LSPMonitor) updateRepoSizeMetrics(graph *SymbolGraphMetrics) {
	var metrics *RepoSizeMetrics
	sizeCategory := m.getRepositorySizeCategory(graph.FileCount)

	switch sizeCategory {
	case "small":
		metrics = &m.performance.SmallRepoMetrics
	case "medium":
		metrics = &m.performance.MediumRepoMetrics
	case "large":
		metrics = &m.performance.LargeRepoMetrics
	}

	if metrics != nil {
		metrics.Count++
		
		if metrics.Count == 1 {
			metrics.AvgGenerationTime = graph.Duration
			metrics.AvgMemoryUsage = graph.MemoryUsage
		} else {
			totalDuration := metrics.AvgGenerationTime*time.Duration(metrics.Count-1) + graph.Duration
			metrics.AvgGenerationTime = totalDuration / time.Duration(metrics.Count)
			
			totalMemory := metrics.AvgMemoryUsage*(metrics.Count-1) + graph.MemoryUsage
			metrics.AvgMemoryUsage = totalMemory / int64(metrics.Count)
		}

		if metrics.Count > 0 {
			metrics.SuccessRate = float64(m.performance.SuccessGraphs) / float64(metrics.Count)
		}
	}
}

// getRepositorySizeCategory categorizes repository by size
func (m *LSPMonitor) getRepositorySizeCategory(fileCount int) string {
	switch {
	case fileCount < 100:
		return "small"
	case fileCount < 500:
		return "medium"
	default:
		return "large"
	}
}

// validateSymbolAccuracy validates symbol accuracy for a generated graph
func (m *LSPMonitor) validateSymbolAccuracy(graph *SymbolGraphMetrics) {
	if graph.Status != "success" || len(graph.Symbols) == 0 {
		return
	}

	// For demonstration, we'll use a simplified accuracy check
	// In a real implementation, this would compare against ground truth
	expectedSymbols := graph.FileCount * 10 // Expected ~10 symbols per file
	actualSymbols := len(graph.Symbols)
	
	accuracy := float64(actualSymbols) / float64(expectedSymbols)
	if accuracy > 1.0 {
		accuracy = 1.0
	}

	check := AccuracyCheck{
		Timestamp:  time.Now(),
		Expected:   expectedSymbols,
		Actual:    actualSymbols,
		Accuracy:   accuracy,
		Repository: graph.RepositoryPath,
		SampleSize: m.config.AccuracySampleSize,
	}

	m.performance.AccuracyChecks = append(m.performance.AccuracyChecks, check)
	
	// Update overall accuracy
	if len(m.performance.AccuracyChecks) > 0 {
		totalAccuracy := 0.0
		for _, check := range m.performance.AccuracyChecks {
			totalAccuracy += check.Accuracy
		}
		m.performance.SymbolAccuracy = totalAccuracy / float64(len(m.performance.AccuracyChecks))
	}

	// Check for accuracy violations
	if accuracy < m.config.MinSymbolAccuracy {
		violation := LSPViolation{
			Type:       "symbol_accuracy",
			Threshold:  m.config.MinSymbolAccuracy,
			Actual:     accuracy,
			Repository: graph.RepositoryPath,
			Timestamp:  time.Now(),
			Severity:   "error",
		}
		m.performance.Violations = append(m.performance.Violations, violation)
	}
}

// checkPerformanceViolations checks for performance threshold violations
func (m *LSPMonitor) checkPerformanceViolations(graph *SymbolGraphMetrics) {
	// Check generation time violations
	if graph.Duration > m.config.MaxGenerationTime {
		violation := LSPViolation{
			Type:       "generation_time",
			Threshold:  m.config.MaxGenerationTime,
			Actual:     graph.Duration,
			Repository: graph.RepositoryPath,
			Timestamp:  time.Now(),
			Severity:   "warning",
		}
		m.performance.Violations = append(m.performance.Violations, violation)
	}

	// Check memory usage violations
	if graph.MemoryUsage > m.config.MaxMemoryUsage {
		violation := LSPViolation{
			Type:       "memory_usage",
			Threshold:  m.config.MaxMemoryUsage,
			Actual:     graph.MemoryUsage,
			Repository: graph.RepositoryPath,
			Timestamp:  time.Now(),
			Severity:   "error",
		}
		m.performance.Violations = append(m.performance.Violations, violation)
	}
}

// logEvent logs a monitoring event
func (m *LSPMonitor) logEvent(eventType, graphID string, data map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"event_type":  eventType,
		"graph_id":    graphID,
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

// Close closes the LSP monitor and cleans up resources
func (m *LSPMonitor) Close() error {
	if m.logFile != nil {
		return m.logFile.Close()
	}
	return nil
}
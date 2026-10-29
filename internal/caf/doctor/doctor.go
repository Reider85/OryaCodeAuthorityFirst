package doctor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/config"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/hooks"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/lsp"
)

// HealthCheck represents a health check result
type HealthCheck struct {
	Name        string `json:"name"`
	Status      string `json:"status"` // "healthy", "warning", "error", "unknown"
	Message     string `json:"message"`
	Details     string `json:"details,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// HealthReport represents the complete health report
type HealthReport struct {
	Overall     string        `json:"overall"`
	Checks      []HealthCheck `json:"checks"`
	Summary     string        `json:"summary"`
	Timestamp   string        `json:"timestamp"`
}

// Doctor performs health checks on the system
type Doctor struct {
	repoRoot string
}

// NewDoctor creates a new doctor instance
func NewDoctor(repoRoot string) *Doctor {
	return &Doctor{
		repoRoot: repoRoot,
	}
}

// Run runs all health checks and returns a report
func (d *Doctor) Run() (*HealthReport, error) {
	log.Printf("🔧 Running health checks for %s", d.repoRoot)

	var checks []HealthCheck
	var healthy, warnings, errors, unknown int

	// Check 1: .caf/ directory validity
	check := d.checkCAFValidity()
	checks = append(checks, check)
	d.updateStatusCount(check.Status, &healthy, &warnings, &errors, &unknown)

	// Check 2: Config schema
	check = d.checkConfigSchema()
	checks = append(checks, check)
	d.updateStatusCount(check.Status, &healthy, &warnings, &errors, &unknown)

	// Check 3: Sessions present
	check = d.checkSessionsPresent()
	checks = append(checks, check)
	d.updateStatusCount(check.Status, &healthy, &warnings, &errors, &unknown)

	// Check 4: Hook installed
	check = d.checkHookInstalled()
	checks = append(checks, check)
	d.updateStatusCount(check.Status, &healthy, &warnings, &errors, &unknown)

	// Check 5: LSP responds
	check = d.checkLSPResponds()
	checks = append(checks, check)
	d.updateStatusCount(check.Status, &healthy, &warnings, &errors, &unknown)

	// Determine overall status
	var overall string
	if errors == 0 && warnings == 0 {
		overall = "healthy"
	} else if errors == 0 {
		overall = "warning"
	} else {
		overall = "error"
	}

	// Create summary
	summary := fmt.Sprintf("%d passed, %d failed, %d warnings, %d unknown", 
		healthy, errors, warnings, unknown)

	report := &HealthReport{
		Overall:   overall,
		Checks:    checks,
		Summary:   summary,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	log.Printf("✓ Health check completed: %s", summary)
	return report, nil
}

// checkCAFValidity checks if the .caf/ directory is valid
func (d *Doctor) checkCAFValidity() HealthCheck {
	cafDir := storage.GetCAFDir(d.repoRoot)
	
	// Check if .caf directory exists
	if _, err := os.Stat(cafDir); err != nil {
		if os.IsNotExist(err) {
			return HealthCheck{
				Name:    "CAF directory validity",
				Status:  "error",
				Message: ".caf/ directory not found",
				Suggestion: "Run 'caf init' in this repository",
			}
		}
		return HealthCheck{
			Name:    "CAF directory validity",
			Status:  "error",
			Message: "Failed to access .caf/ directory",
			Details: err.Error(),
		}
	}

	// Check required subdirectories
	requiredDirs := []string{
		filepath.Join(cafDir, "sessions"),
		filepath.Join(cafDir, "decisions"),
		filepath.Join(cafDir, "mental-model"),
		filepath.Join(cafDir, "cache"),
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				return HealthCheck{
					Name:    "CAF directory validity",
					Status:  "warning",
					Message: "Missing subdirectory: " + dir,
					Suggestion: "Run 'caf init' to fix directory structure",
				}
			}
		}
	}

	// Check manifest.yml
	manifestFile := filepath.Join(cafDir, "manifest.yml")
	if _, err := os.Stat(manifestFile); err != nil {
		return HealthCheck{
			Name:    "CAF directory validity",
			Status:  "warning",
			Message: "manifest.yml not found",
			Suggestion: "Run 'caf init' to create manifest",
		}
	}

	return HealthCheck{
		Name:    "CAF directory validity",
		Status:  "healthy",
		Message: ".caf/ directory is valid",
	}
}

// checkConfigSchema checks if the config schema is valid
func (d *Doctor) checkConfigSchema() HealthCheck {
	configFile := storage.GetConfigPath(d.repoRoot)
	
	// Check if config file exists
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			return HealthCheck{
				Name:    "Config schema",
				Status:  "warning",
				Message: "config.yml not found, using defaults",
				Suggestion: "Run 'caf init' to create config file",
			}
		}
		return HealthCheck{
			Name:    "Config schema",
			Status:  "error",
			Message: "Failed to access config.yml",
			Details: err.Error(),
		}
	}

	// Try to load config
	configManager := config.NewManager(configFile)
	cfg, err := configManager.LoadConfig()
	if err != nil {
		return HealthCheck{
			Name:    "Config schema",
			Status:  "error",
			Message: "Invalid config.yml schema",
			Details: err.Error(),
			Suggestion: "Check config.yml syntax or delete it to regenerate",
		}
	}

	// Check required fields
	if cfg.Version == "" {
		return HealthCheck{
			Name:    "Config schema",
			Status:  "warning",
			Message: "Missing version field",
			Suggestion: "Update config.yml with version: 0.4.0-triz2",
		}
	}

	return HealthCheck{
		Name:    "Config schema",
		Status:  "healthy",
		Message: "config.yml is valid",
	}
}

// checkSessionsPresent checks if sessions are present
func (d *Doctor) checkSessionsPresent() HealthCheck {
	sessionsDir := filepath.Join(d.repoRoot, ".caf", "sessions")
	
	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDir); err != nil {
		if os.IsNotExist(err) {
			return HealthCheck{
				Name:    "Sessions present",
				Status:  "unknown",
				Message: "No sessions yet",
				Suggestion: "Run an AI agent in this repository to start tracking",
			}
		}
		return HealthCheck{
			Name:    "Sessions present",
			Status:  "error",
			Message: "Failed to access sessions directory",
			Details: err.Error(),
		}
	}

	// List sessions
	sessionManager := storage.NewManager(sessionsDir)
	sessionIDs, err := sessionManager.ListSessions()
	if err != nil {
		return HealthCheck{
			Name:    "Sessions present",
			Status:  "error",
			Message: "Failed to list sessions",
			Details: err.Error(),
		}
	}

	if len(sessionIDs) == 0 {
		return HealthCheck{
			Name:    "Sessions present",
			Status:  "unknown",
			Message: "No sessions yet",
			Suggestion: "Run an AI agent in this repository to start tracking",
		}
	}

	return HealthCheck{
		Name:    "Sessions present",
		Status:  "healthy",
		Message: fmt.Sprintf("Found %d sessions", len(sessionIDs)),
	}
}

// checkHookInstalled checks if the pre-commit hook is installed
func (d *Doctor) checkHookInstalled() HealthCheck {
	installed, err := hooks.IsHookInstalled(d.repoRoot)
	if err != nil {
		return HealthCheck{
			Name:    "Hook installed",
			Status:  "error",
			Message: "Failed to check hook status",
			Details: err.Error(),
		}
	}

	if !installed {
		return HealthCheck{
			Name:    "Hook installed",
			Status:  "warning",
			Message: "Universal pre-commit hook not installed",
			Suggestion: "Run 'caf init' to install the hook",
		}
	}

	return HealthCheck{
		Name:    "Hook installed",
		Status:  "healthy",
		Message: "Universal pre-commit hook is installed",
	}
}

// checkLSPResponds checks if the LSP server responds
func (d *Doctor) checkLSPResponds() HealthCheck {
	// Check if gopls is available
	if _, err := exec.LookPath("gopls"); err != nil {
		return HealthCheck{
			Name:    "LSP responds",
			Status:  "unknown",
			Message: "gopls not available",
			Suggestion: "Install Go and gopls for symbol graph support",
		}
	}

	// Try to create LSP manager
	manager, err := lsp.GetManagerForRepo(d.repoRoot)
	if err != nil {
		return HealthCheck{
			Name:    "LSP responds",
			Status:  "warning",
			Message: "LSP server not responding",
			Details: err.Error(),
			Suggestion: "Run 'caf sync' to start LSP server",
		}
	}

	// Check if symbols can be loaded
	symbols := manager.GetSymbols()
	if len(symbols) == 0 {
		return HealthCheck{
			Name:    "LSP responds",
			Status:  "warning",
			Message: "LSP server running but no symbols loaded",
			Suggestion: "Run 'caf sync' to refresh symbol graph",
		}
	}

	return HealthCheck{
		Name:    "LSP responds",
		Status:  "healthy",
		Message: fmt.Sprintf("LSP server responding, %d symbols loaded", len(symbols)),
	}
}

// updateStatusCount updates the status counters
func (d *Doctor) updateStatusCount(status string, healthy, warnings, errors, unknown *int) {
	switch status {
	case "healthy":
		*healthy++
	case "warning":
		warnings++
	case "error":
		*errors++
	case "unknown":
		unknown++
	}
}

// RunCommand handles the caf doctor command
func RunCommand(repoRoot string) (*HealthReport, error) {
	doctor := NewDoctor(repoRoot)
	return doctor.Run()
}

// FormatReport formats a health report for display
func FormatReport(report *HealthReport) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("CodeAuthorityFirst Health Report\n"))
	builder.WriteString(fmt.Sprintf("Timestamp: %s\n", report.Timestamp))
	builder.WriteString(fmt.Sprintf("Overall: %s\n", report.Overall))
	builder.WriteString(fmt.Sprintf("Summary: %s\n\n", report.Summary))
	
	builder.WriteString("Health Checks:\n")
	for _, check := range report.Checks {
		statusIcon := "✓"
		switch check.Status {
		case "warning":
			statusIcon = "⚠"
		case "error":
			statusIcon = "✗"
		case "unknown":
			statusIcon = "?"
		}
		
		builder.WriteString(fmt.Sprintf("  %s %s: %s\n", statusIcon, check.Name, check.Message))
		if check.Suggestion != "" {
			builder.WriteString(fmt.Sprintf("    Suggestion: %s\n", check.Suggestion))
		}
		if check.Details != "" {
			builder.WriteString(fmt.Sprintf("    Details: %s\n", check.Details))
		}
	}
	
	return builder.String()
}
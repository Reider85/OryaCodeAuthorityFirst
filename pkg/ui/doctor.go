package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// DoctorScreen represents the health check screen
// FR-TUI-002: Doctor shows 5 health checks with detailed diagnostics
type DoctorScreen struct {
	*Screen
	Components []Component
	Title      string
	
	// Health check results
	Results []HealthCheckResult
	
	// UI components
	CheckList     *List
	RunButton     *Button
	RefreshButton *Button
	StatusSpinner *Spinner
	
	// State
	Loading       bool
	LastRun       time.Time
	Error         error
}

// HealthCheckResult represents a health check result
type HealthCheckResult struct {
	ID          string
	Name        string
	Status      string // "healthy", "warning", "error", "unknown"
	Message     string
	Details     string
	Recommendations []string
	Timestamp   time.Time
	Duration    time.Duration
}

// NewDoctorScreen creates a new doctor screen
func NewDoctorScreen(theme *Theme, keys *Keys, a11y *A11y, perfTracker *PerfTracker) *DoctorScreen {
	screen := NewScreen(theme, keys, perfTracker)
	
	return &DoctorScreen{
		Screen:        screen,
		Components:    []Component{},
		Title:         "CodeAuthorityFirst Health Check",
		Results:       []HealthCheckResult{},
		CheckList:     NewList("checks", 10, theme, keys),
		RunButton:     NewButton("run", "Run Health Checks", func() tea.Cmd {
			return func() tea.Msg {
				return DoctorRunMsg{Force: true}
			}
		}, theme, keys),
		RefreshButton: NewButton("refresh", "Refresh", func() tea.Cmd {
			return func() tea.Msg {
				return RefreshMsg{Force: true}
			}
		}, theme, keys),
		StatusSpinner: NewSpinner("spinner", theme, keys),
		Loading:       false,
		LastRun:       time.Now(),
	}
}

// Init initializes the doctor screen
func (d *DoctorScreen) Init() tea.Cmd {
	// Set up screen reader announcement
	d.Screen.ScreenReader.Announce(AnnouncementImportant, "Doctor screen activated", "Application", false)
	
	// Run initial health checks
	return d.runHealthChecks()
}

// runHealthChecks runs all health checks
func (d *DoctorScreen) runHealthChecks() tea.Cmd {
	d.Loading = true
	d.StatusSpinner.SetAnimating(true)
	d.Error = nil
	
	return func() tea.Msg {
		// Simulate health check execution
		time.Sleep(2 * time.Second)
		
		// Generate mock results
		results := d.generateMockResults()
		
		return DoctorResultsLoaded{
			Results: results,
		}
	}
}

// generateMockResults generates mock health check results
func (d *DoctorScreen) generateMockResults() []HealthCheckResult {
	return []HealthCheckResult{
		{
			ID:          "caf-dir",
			Name:        ".caf/ Directory",
			Status:      "healthy",
			Message:     "Directory structure is valid",
			Details:     "Found all required files and directories",
			Recommendations: []string{},
			Timestamp:   time.Now(),
			Duration:    150 * time.Millisecond,
		},
		{
			ID:          "config",
			Name:        "Configuration",
			Status:      "healthy",
			Message:     "Configuration file is valid",
			Details:     "All required settings are present",
			Recommendations: []string{},
			Timestamp:   time.Now(),
			Duration:    200 * time.Millisecond,
		},
		{
			ID:          "sessions",
			Name:        "Session Storage",
			Status:      "warning",
			Message:     "Some sessions have incomplete data",
			Details:     "Found 3 sessions missing metadata",
			Recommendations: []string{
				"Run 'caf doctor --fix-sessions' to repair",
				"Check session file permissions",
			},
			Timestamp:   time.Now(),
			Duration:    500 * time.Millisecond,
		},
		{
			ID:          "hook",
			Name:        "Pre-commit Hook",
			Status:      "healthy",
			Message:     "Hook is installed and executable",
			Details:     "Universal pre-commit hook found in .git/hooks/",
			Recommendations: []string{},
			Timestamp:   time.Now(),
			Duration:    100 * time.Millisecond,
		},
		{
			ID:          "lsp",
			Name:        "LSP Connection",
			Status:      "error",
			Message:     "LSP server not responding",
			Details:     "gopls failed to start after 3 attempts",
			Recommendations: []string{
				"Install Go 1.22+ to use LSP features",
				"Check PATH for gopls binary",
				"Run 'caf doctor --install-lsp' to install",
			},
			Timestamp:   time.Now(),
			Duration:    3 * time.Second,
		},
	}
}

// Update updates the doctor screen
func (d *DoctorScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.Type {
		case tea.Key1:
			// Switch to dashboard
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDashboard} }
		case tea.Key2:
			// Stay on doctor
			return d, nil
		case tea.Key3:
			// Switch to metrics
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenMetrics} }
		case tea.Key4:
			// Switch to cost
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenCost} }
		case tea.Key5:
			// Switch to settings
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenSettings} }
		case tea.KeyR:
			// Run health checks
			return d, d.runHealthChecks()
		case tea.KeyQuestion:
			// Show help
			return d, d.showHelp()
		}
		
		// Handle list selection
		if d.CheckList.IsFocused() {
			switch msg.Type {
			case tea.KeyEnter:
				if selectedItem, ok := d.CheckList.GetSelectedItem(); ok {
					return d, func() tea.Msg {
						return DoctorItemSelected{
							Item: selectedItem,
							Result: d.getResultByID(selectedItem.ID),
						}
					}
				}
			}
		}
		
		// Update components
		for i, component := range d.Components {
			updatedComponent, cmd := component.Update(msg)
			d.Components[i] = updatedComponent
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case DoctorResultsLoaded:
		// Handle results loaded
		d.Loading = false
		d.StatusSpinner.SetAnimating(false)
		d.Results = msg.Results
		d.LastRun = time.Now()
		
		// Update list
		d.updateCheckList()
		
	case DoctorItemSelected:
		// Handle item selection
		d.showItemDetails(msg.Result)
		
	case DoctorRunMsg:
		// Handle run command
		if msg.Force {
			return d, d.runHealthChecks()
		}
		
	case ErrorMsg:
		// Handle error
		d.Error = msg.Err
		d.Loading = false
		d.StatusSpinner.SetAnimating(false)
	}
	
	return d, tea.Batch(cmds...)
}

// updateCheckList updates the check list with results
func (d *DoctorScreen) updateCheckList() {
	items := make([]ListItem, len(d.Results))
	for i, result := range d.Results {
		statusIcon, statusColor := d.getStatusIconAndColor(result.Status)
		
		title := fmt.Sprintf("%s %s", statusIcon, result.Name)
		subtitle := fmt.Sprintf("%s • %s", result.Message, d.formatDuration(result.Duration))
		
		items[i] = ListItem{
			ID:       result.ID,
			Title:    title,
			Subtitle: subtitle,
			Data:     result,
		}
	}
	
	d.CheckList.SetItems(items)
}

// getStatusIconAndColor returns status icon and color
func (d *DoctorScreen) getStatusIconAndColor(status string) (string, string) {
	switch status {
	case "healthy":
		return "✓", "secondary"
	case "warning":
		return "⚠", "warning"
	case "error":
		return "✗", "error"
	default:
		return "?", "muted"
	}
}

// getResultByID returns a result by ID
func (d *DoctorScreen) getResultByID(id string) *HealthCheckResult {
	for _, result := range d.Results {
		if result.ID == id {
			return &result
		}
	}
	return nil
}

// showItemDetails shows details for a selected item
func (d *DoctorScreen) showItemDetails(result *HealthCheckResult) {
	if result == nil {
		return
	}
	
	// Announce to screen reader
	d.Screen.ScreenReader.Announce(AnnouncementImportant, 
		fmt.Sprintf("%s: %s", result.Name, result.Message), "Health Check", false)
	
	// In a real implementation, this would show a modal or expand the item
	// For now, we'll just log it
	fmt.Printf("Details for %s:\n%s\n\nRecommendations:\n", result.Name, result.Details)
	for _, rec := range result.Recommendations {
		fmt.Printf("- %s\n", rec)
	}
}

// formatDuration formats a duration for display
func (d *DoctorScreen) formatDuration(duration time.Duration) string {
	if duration < time.Millisecond {
		return fmt.Sprintf("%dns", duration.Nanoseconds())
	}
	if duration < time.Second {
		return fmt.Sprintf("%dms", duration.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", duration.Seconds())
}

// showHelp shows help for the doctor screen
func (d *DoctorScreen) showHelp() tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{
			ScreenID: "doctor",
			Message:  d.getDoctorHelp(),
		}
	}
}

// getDoctorHelp returns doctor help text
func (d *DoctorScreen) getDoctorHelp() string {
	return `Doctor Controls:
1-5: Switch screens
r: Run health checks
?: Show this help
q: Quit

Health Checks:
• .caf/ Directory: Validates directory structure
• Configuration: Checks config file validity
• Session Storage: Verifies session data integrity
• Pre-commit Hook: Ensures hook is installed
• LSP Connection: Tests LSP server connectivity

Press Enter on a check to see details`
}

// View returns the doctor view
func (d *DoctorScreen) View() string {
	if d.Loading {
		return d.renderLoading()
	}
	
	if d.Error != nil {
		return d.renderError()
	}
	
	return d.renderContent()
}

// renderLoading renders the loading state
func (d *DoctorScreen) renderLoading() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	header += d.Theme.Muted.Render("Running health checks...")
	
	content := fmt.Sprintf("%s\n\n%s", header, d.StatusSpinner.View())
	
	return d.renderLayout(content)
}

// renderError renders the error state
func (d *DoctorScreen) renderError() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	errorBox := d.Theme.ErrorBox.Render(fmt.Sprintf("Error: %s", d.Error.Error()))
	
	content := fmt.Sprintf("%s\n\n%s", header, errorBox)
	
	return d.renderLayout(content)
}

// renderContent renders the main content
func (d *DoctorScreen) renderContent() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	header += d.Theme.Muted.Render(fmt.Sprintf("Last run: %s", d.LastRun.Format("15:04:05"))) + "\n"
	
	// Overall status
	overallStatus := d.getOverallStatus()
	statusColor := d.getStatusColor(overallStatus)
	statusText := fmt.Sprintf("Overall Status: %s", statusColor.Render(overallStatus))
	
	// Create layout
	layout := fmt.Sprintf("%s\n%s\n%s\n%s",
		header,
		statusText,
	 d.Theme.Title.Render("Health Checks"),
	 d.CheckList.View(),
	)
	
	// Add buttons
	layout += fmt.Sprintf("\n\n%s", d.RunButton.View())
	layout += fmt.Sprintf("  %s", d.RefreshButton.View())
	
	return d.renderLayout(layout)
}

// getOverallStatus returns the overall health status
func (d *DoctorScreen) getOverallStatus() string {
	if len(d.Results) == 0 {
		return "unknown"
	}
	
	hasError := false
	hasWarning := false
	
	for _, result := range d.Results {
		if result.Status == "error" {
			hasError = true
		} else if result.Status == "warning" {
			hasWarning = true
		}
	}
	
	if hasError {
		return "error"
	}
	if hasWarning {
		return "warning"
	}
	return "healthy"
}

// getStatusColor returns the color for a status
func (d *DoctorScreen) getStatusColor(status string) string {
	switch status {
	case "healthy":
		return d.Theme.Secondary.Render(status)
	case "warning":
		return d.Theme.Warning.Render(status)
	case "error":
		return d.Theme.Error.Render(status)
	default:
		return d.Theme.Muted.Render(status)
	}
}

// renderLayout renders the layout with proper styling
func (d *DoctorScreen) renderLayout(content string) string {
	// Add status bar
	statusBar := d.Theme.StatusBar.Render("Doctor | Press r to run checks, ? for help")
	
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

// GetID returns the screen ID
func (d *DoctorScreen) GetID() string {
	return "doctor"
}

// IsActive returns true if the screen is active
func (d *DoctorScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (d *DoctorScreen) SetActive(active bool) {
	if active {
		d.Screen.ScreenReader.Announce(AnnouncementImportant, "Doctor screen activated", "Application", false)
	} else {
		d.Screen.ScreenReader.Announce(AnnouncementNone, "Doctor screen deactivated", "Application", false)
	}
}

// DoctorResultsLoaded represents doctor results loaded message
type DoctorResultsLoaded struct {
	Results []HealthCheckResult
}

// DoctorRunMsg represents a doctor run message
type DoctorRunMsg struct {
	Force bool
}

// DoctorItemSelected represents a doctor item selection message
type DoctorItemSelected struct {
	Item   ListItem
	Result *HealthCheckResult
}

// ErrorMsg represents an error message
type ErrorMsg struct {
	Err error
}

// RefreshMsg represents a refresh message
type RefreshMsg struct {
	Force bool
}

// HelpMsg represents a help message
type HelpMsg struct {
	ScreenID string
	Message  string
}
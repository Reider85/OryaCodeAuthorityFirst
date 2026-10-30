package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// DashboardScreen represents the main dashboard screen
// FR-TUI-001: Dashboard shows recent sessions, decisions, and agent activity
type DashboardScreen struct {
	*ScreenBase
	Components []Component
	Title      string
	
	// Dashboard data
	RecentSessions []SessionItem
	Decisions      []DecisionItem
	AgentActivity  []AgentActivityItem
	
	// UI components
	SessionList    *List
	DecisionList   *List
	ActivityList   *List
	RefreshButton  *Button
	StatusSpinner  *Spinner
	
	// State
	Loading       bool
	LastRefresh   time.Time
	Error         error
}

// SessionItem represents a session item
type SessionItem struct {
	ID          string
	Agent       string
	Prompt      string
	Duration    time.Duration
	Status      string
	Timestamp   time.Time
}

// DecisionItem represents a decision item
type DecisionItem struct {
	ID          string
	Title       string
	Description string
	Reason      string
	Outcome     string
	Timestamp   time.Time
}

// AgentActivityItem represents an agent activity item
type AgentActivityItem struct {
	Agent       string
	Action      string
	Target      string
	Timestamp   time.Time
	Duration    time.Duration
}

// NewDashboardScreen creates a new dashboard screen
func NewDashboardScreen(theme *Theme, keys *Keys, a11y *A11y, perfTracker *PerfTracker) *DashboardScreen {
	screen := NewScreen(theme, keys, a11y, perfTracker)
	
	return &DashboardScreen{
		ScreenBase:        screen,
		Components:    []Component{},
		Title:         "CodeAuthorityFirst Dashboard",
		RecentSessions: []SessionItem{},
		Decisions:     []DecisionItem{},
		AgentActivity: []AgentActivityItem{},
		SessionList:   NewList("sessions", 8, theme, keys),
		DecisionList:  NewList("decisions", 6, theme, keys),
		ActivityList:  NewList("activity", 10, theme, keys),
		RefreshButton: NewButton("refresh", "Refresh", func() tea.Cmd {
			return func() tea.Msg {
				return RefreshMsg{Force: true}
			}
		}, theme, keys),
		StatusSpinner: NewSpinner("spinner", theme, keys),
		Loading:       false,
		LastRefresh:   time.Now(),
	}
}

// Init initializes the dashboard screen
func (d *DashboardScreen) Init() tea.Cmd {
	// Set up screen reader announcement
	d.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Dashboard screen activated", "Application", false)
	
	// Start loading
	d.Loading = true
	d.StatusSpinner.SetAnimating(true)
	
	// Load initial data
	return d.loadData()
}

// loadData loads dashboard data
func (d *DashboardScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		// Simulate data loading
		time.Sleep(500 * time.Millisecond)
		
		// Generate mock data
		sessions := d.generateMockSessions()
		decisions := d.generateMockDecisions()
		activity := d.generateMockActivity()
		
		return DashboardDataLoaded{
			Sessions: sessions,
			Decisions: decisions,
			Activity: activity,
		}
	}
}

// generateMockSessions generates mock session data
func (d *DashboardScreen) generateMockSessions() []SessionItem {
	return []SessionItem{
		{
			ID:        "2026-10-30-a1b2c3",
			Agent:     "opencode",
			Prompt:    "Add rate limiter to API endpoints",
			Duration:  2*time.Minute + 30*time.Second,
			Status:    "completed",
			Timestamp: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "2026-10-30-d4e5f6",
			Agent:     "claude-code",
			Prompt:    "Refactor authentication module",
			Duration:  5*time.Minute + 15*time.Second,
			Status:    "completed",
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "2026-10-30-g7h8i9",
			Agent:     "universal-lsp",
			Prompt:    "Fix type errors in main.go",
			Duration:  1*time.Minute + 45*time.Second,
			Status:    "completed",
			Timestamp: time.Now().Add(-3 * time.Hour),
		},
	}
}

// generateMockDecisions generates mock decision data
func (d *DashboardScreen) generateMockDecisions() []DecisionItem {
	return []DecisionItem{
		{
			ID:          "dec-001",
			Title:       "API Rate Limiting Strategy",
			Description: "Implement token bucket algorithm for API endpoints",
			Reason:      "Prevent abuse and ensure fair usage",
			Outcome:     "Successfully implemented with configurable limits",
			Timestamp:   time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "dec-002",
			Title:       "Database Connection Pool",
			Description: "Configure connection pool for PostgreSQL",
			Reason:      "Improve performance and resource utilization",
			Outcome:     "Pool size optimized for current load",
			Timestamp:   time.Now().Add(-2 * time.Hour),
		},
		{
			ID:          "dec-003",
			Title:       "Error Handling Strategy",
			Description: "Standardize error responses across all endpoints",
			Reason:      "Improve API consistency and debugging",
			Outcome:     "Unified error response format implemented",
			Timestamp:   time.Now().Add(-3 * time.Hour),
		},
	}
}

// generateMockActivity generates mock activity data
func (d *DashboardScreen) generateMockActivity() []AgentActivityItem {
	return []AgentActivityItem{
		{
			Agent:     "opencode",
			Action:    "modified",
			Target:    "internal/api/rate_limiter.go",
			Timestamp: time.Now().Add(-30 * time.Minute),
			Duration:  2*time.Minute + 15*time.Second,
		},
		{
			Agent:     "claude-code",
			Action:    "created",
			Target:    "internal/auth/middleware.go",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Duration:  4*time.Minute + 30*time.Second,
		},
		{
			Agent:     "universal-lsp",
			Action:    "fixed",
			Target:    "cmd/server/main.go",
			Timestamp: time.Now().Add(-1 * time.Hour - 30 * time.Minute),
			Duration:  1*time.Minute + 20*time.Second,
		},
		{
			Agent:     "opencode",
			Action:    "refactored",
			Target:    "internal/database/connection.go",
			Timestamp: time.Now().Add(-2 * time.Hour),
			Duration:  3*time.Minute + 45*time.Second,
		},
	}
}

// Update updates the dashboard screen
func (d *DashboardScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.String() {
		case "1":
			// Stay on dashboard
			return d, nil
		case "2":
			// Switch to doctor
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDoctor} }
		case "3":
			// Switch to metrics
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenMetrics} }
		case "4":
			// Switch to cost
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenCost} }
		case "5":
			// Switch to settings
			return d, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenSettings} }
		case "r":
			// Refresh dashboard
			return d, d.refresh()
		case "?":
			// Show help
			return d, d.showHelp()
		}
		
		// Update components
		for i, component := range d.Components {
			updatedComponent, cmd := component.Update(msg)
			d.Components[i] = updatedComponent
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case DashboardDataLoaded:
		// Handle data loaded
		d.Loading = false
		d.StatusSpinner.SetAnimating(false)
		d.RecentSessions = msg.Sessions
		d.Decisions = msg.Decisions
		d.AgentActivity = msg.Activity
		
		// Update lists
		d.updateLists()
		
	case RefreshMsg:
		// Handle refresh
		if msg.Force {
			return d, d.refresh()
		}
		
	case ErrorMsg:
		// Handle error
		d.Error = msg.Err
		d.Loading = false
		d.StatusSpinner.SetAnimating(false)
	}
	
	return d, tea.Batch(cmds...)
}

// refresh refreshes the dashboard
func (d *DashboardScreen) refresh() tea.Cmd {
	d.Loading = true
	d.StatusSpinner.SetAnimating(true)
	d.Error = nil
	return d.loadData()
}

// updateLists updates the list components with new data
func (d *DashboardScreen) updateLists() {
	// Update session list
	sessionItems := make([]ListItem, len(d.RecentSessions))
	for i, session := range d.RecentSessions {
		statusIcon := "✓"
		if session.Status != "completed" {
			statusIcon = "○"
		}
		
		sessionItems[i] = ListItem{
			ID:       session.ID,
			Title:    fmt.Sprintf("%s %s", statusIcon, session.Agent),
			Subtitle: fmt.Sprintf("%s • %s", session.Prompt, d.formatDuration(session.Duration)),
			Data:     session,
		}
	}
	d.SessionList.SetItems(sessionItems)
	
	// Update decision list
	decisionItems := make([]ListItem, len(d.Decisions))
	for i, decision := range d.Decisions {
		decisionItems[i] = ListItem{
			ID:       decision.ID,
			Title:    decision.Title,
			Subtitle: fmt.Sprintf("%s • %s", decision.Reason, decision.Outcome),
			Data:     decision,
		}
	}
	d.DecisionList.SetItems(decisionItems)
	
	// Update activity list
	activityItems := make([]ListItem, len(d.AgentActivity))
	for i, activity := range d.AgentActivity {
		activityItems[i] = ListItem{
			ID:       activity.Agent,
			Title:    fmt.Sprintf("%s %s", activity.Agent, activity.Action),
			Subtitle: fmt.Sprintf("%s • %s", activity.Target, d.formatDuration(activity.Duration)),
			Data:     activity,
		}
	}
	d.ActivityList.SetItems(activityItems)
}

// formatDuration formats a duration for display
func (d *DashboardScreen) formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	return fmt.Sprintf("%dh", int(duration.Hours()))
}

// showHelp shows help for the dashboard
func (d *DashboardScreen) showHelp() tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{
			ScreenID: "dashboard",
			Message:  d.getDashboardHelp(),
		}
	}
}

// getDashboardHelp returns dashboard help text
func (d *DashboardScreen) getDashboardHelp() string {
	return `Dashboard Controls:
1-5: Switch screens
r: Refresh dashboard
?: Show this help
q: Quit

Dashboard Sections:
Recent Sessions: AI sessions with prompts and duration
Decisions: Technical decisions made
Agent Activity: Recent agent actions and targets`
}

// View returns the dashboard view
func (d *DashboardScreen) View() string {
	if d.Loading {
		return d.renderLoading()
	}
	
	if d.Error != nil {
		return d.renderError()
	}
	
	return d.renderContent()
}

// renderLoading renders the loading state
func (d *DashboardScreen) renderLoading() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	header += d.Theme.Muted.Render("Loading dashboard data...")
	
	content := fmt.Sprintf("%s\n\n%s", header, d.StatusSpinner.View())
	
	return d.renderLayout(content)
}

// renderError renders the error state
func (d *DashboardScreen) renderError() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	errorBox := d.Theme.ErrorBox.Render(fmt.Sprintf("Error: %s", d.Error.Error()))
	
	content := fmt.Sprintf("%s\n\n%s", header, errorBox)
	
	return d.renderLayout(content)
}

// renderContent renders the main content
func (d *DashboardScreen) renderContent() string {
	header := d.Theme.Title.Render(d.Title) + "\n"
	header += d.Theme.Muted.Render(fmt.Sprintf("Last updated: %s", d.LastRefresh.Format("15:04:05"))) + "\n"
	
	// Create layout
	layout := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		header,
	 d.Theme.Title.Render("Recent Sessions"),
	 d.SessionList.View(),
	 d.Theme.Title.Render("Recent Decisions"),
	 d.DecisionList.View(),
	 d.Theme.Title.Render("Agent Activity"),
	)
	
	// Add refresh button
	layout += fmt.Sprintf("\n%s", d.RefreshButton.View())
	
	return d.renderLayout(layout)
}

// renderLayout renders the layout with proper styling
func (d *DashboardScreen) renderLayout(content string) string {
	// Add status bar
	statusBar := d.Theme.StatusBar.Render("Dashboard | Press r to refresh, ? for help")
	
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

// GetID returns the screen ID
func (d *DashboardScreen) GetID() string {
	return "dashboard"
}

// IsActive returns true if the screen is active
func (d *DashboardScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (d *DashboardScreen) SetActive(active bool) {
	if active {
		d.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Dashboard screen activated", "Application", false)
	} else {
		d.ScreenBase.ScreenReader.Announce(AnnouncementNone, "Dashboard screen deactivated", "Application", false)
	}
}

// DashboardDataLoaded represents dashboard data loaded message
type DashboardDataLoaded struct {
	Sessions  []SessionItem
	Decisions []DecisionItem
	Activity  []AgentActivityItem
}
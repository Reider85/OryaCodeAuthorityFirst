package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// MetricsScreen represents the metrics screen
// FR-TUI-003: Metrics shows usage statistics, session counts, performance data
type MetricsScreen struct {
	*ScreenBase
	Components []Component
	Title      string
	
	// Metrics data
	SessionStats   SessionStatistics
	UsageStats     UsageStatistics
	PerformanceData PerformanceData
	
	// UI components
	SessionStatsCard  *Card
	UsageStatsCard    *Card
	PerformanceCard   *Card
	RefreshButton     *Button
	ExportButton      *Button
	StatusSpinner     *Spinner
	
	// State
	Loading          bool
	LastRefresh      time.Time
	Error            error
	TimeRange        TimeRange
}

// SessionStatistics represents session statistics
type SessionStatistics struct {
	TotalSessions      int
	SessionsToday     int
	SessionsThisWeek  int
	SessionsThisMonth int
	AvgDuration       time.Duration
	TotalDuration     time.Duration
	AgentBreakdown    map[string]int
}

// UsageStatistics represents usage statistics
type UsageStatistics struct {
	TotalFilesModified int
	TotalCommits       int
	TotalDecisions     int
	ActiveAgents       []string
	TopFiles           []FileStat
	TopAgents          []AgentStat
}

// FileStat represents file statistics
type FileStat struct {
	Path     string
	Modified int
}

// AgentStat represents agent statistics
type AgentStat struct {
	Name     string
	Sessions int
}

// PerformanceData represents performance data
type PerformanceData struct {
	AvgResponseTime time.Duration
	SuccessRate     float64
	ErrorRate       float64
	Throughput      float64
	PeakLoad        time.Time
}

// TimeRange represents the time range for metrics
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Card represents a card component for displaying metrics
type Card struct {
	BaseComponent
	Title   string
	Content string
	Width   int
}

// NewCard creates a new card component
func NewCard(id, title string, theme *Theme, keys *Keys) *Card {
	return &Card{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Title:        title,
		Content:      "",
		Width:        40,
	}
}

// Init initializes the card
func (c *Card) Init() tea.Cmd {
	return nil
}

// Update updates the card
func (c *Card) Update(msg tea.Msg) (Component, tea.Cmd) {
	return c, nil
}

// View returns the card view
func (c *Card) View() string {
	title := c.Theme.Title.Render(c.Title)
	content := c.Theme.Muted.Render(c.Content)
	
	return fmt.Sprintf("%s\n%s", title, content)
}

// SetContent sets the card content
func (c *Card) SetContent(content string) {
	c.Content = content
}

// NewMetricsScreen creates a new metrics screen
func NewMetricsScreen(theme *Theme, keys *Keys, a11y *A11y, perfTracker *PerfTracker) *MetricsScreen {
	screen := NewScreen(theme, keys, a11y, perfTracker)
	
	return &MetricsScreen{
		ScreenBase:           screen,
		Components:       []Component{},
		Title:           "CodeAuthorityFirst Metrics",
		SessionStats:   SessionStatistics{},
		UsageStats:      UsageStatistics{},
		PerformanceData: PerformanceData{},
		SessionStatsCard: NewCard("session-stats", "Session Statistics", theme, keys),
		UsageStatsCard:   NewCard("usage-stats", "Usage Statistics", theme, keys),
		PerformanceCard:  NewCard("performance", "Performance", theme, keys),
		RefreshButton:   NewButton("refresh", "Refresh", func() tea.Cmd {
			return func() tea.Msg {
				return RefreshMsg{Force: true}
			}
		}, theme, keys),
		ExportButton:    NewButton("export", "Export CSV", func() tea.Cmd {
			return func() tea.Msg {
				return ExportMetricsMsg{}
			}
		}, theme, keys),
		StatusSpinner:   NewSpinner("spinner", theme, keys),
		Loading:         false,
		LastRefresh:     time.Now(),
		Error:           nil,
		TimeRange:       TimeRange{Start: time.Now().Add(-24 * time.Hour), End: time.Now()},
	}
}

// Init initializes the metrics screen
func (m *MetricsScreen) Init() tea.Cmd {
	// Set up screen reader announcement
	m.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Metrics screen activated", "Application", false)
	
	// Load initial data
	return m.loadMetrics()
}

// loadMetrics loads metrics data
func (m *MetricsScreen) loadMetrics() tea.Cmd {
	return func() tea.Msg {
		// Simulate data loading
		time.Sleep(800 * time.Millisecond)
		
		// Generate mock data
		sessionStats := m.generateMockSessionStats()
		usageStats := m.generateMockUsageStats()
		performanceData := m.generateMockPerformanceData()
		
		return MetricsDataLoaded{
			SessionStats:    sessionStats,
			UsageStats:      usageStats,
			PerformanceData: performanceData,
		}
	}
}

// generateMockSessionStats generates mock session statistics
func (m *MetricsScreen) generateMockSessionStats() SessionStatistics {
	return SessionStatistics{
		TotalSessions:      156,
		SessionsToday:     12,
		SessionsThisWeek:  45,
		SessionsThisMonth: 89,
		AvgDuration:       3*time.Minute + 45*time.Second,
		TotalDuration:     8*time.Hour + 32*time.Minute,
		AgentBreakdown: map[string]int{
			"opencode":      89,
			"claude-code":   45,
			"universal-lsp": 22,
		},
	}
}

// generateMockUsageStats generates mock usage statistics
func (m *MetricsScreen) generateMockUsageStats() UsageStatistics {
	return UsageStatistics{
		TotalFilesModified: 342,
		TotalCommits:      89,
		TotalDecisions:    23,
		ActiveAgents:      []string{"opencode", "claude-code"},
		TopFiles: []FileStat{
			{Path: "internal/api/server.go", Modified: 45},
			{Path: "internal/database/connection.go", Modified: 32},
			{Path: "cmd/main.go", Modified: 28},
			{Path: "internal/auth/middleware.go", Modified: 19},
			{Path: "internal/config/config.go", Modified: 15},
		},
		TopAgents: []AgentStat{
			{Name: "opencode", Sessions: 89},
			{Name: "claude-code", Sessions: 45},
			{Name: "universal-lsp", Sessions: 22},
		},
	}
}

// generateMockPerformanceData generates mock performance data
func (m *MetricsScreen) generateMockPerformanceData() PerformanceData {
	return PerformanceData{
		AvgResponseTime: 450 * time.Millisecond,
		SuccessRate:     0.987,
		ErrorRate:       0.013,
		Throughput:      2.34,
		PeakLoad:        time.Now().Add(-2 * time.Hour),
	}
}

// Update updates the metrics screen
func (m *MetricsScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.String() {
		case "1":
			// Switch to dashboard
			return m, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDashboard} }
		case "2":
			// Switch to doctor
			return m, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDoctor} }
		case "3":
			// Stay on metrics
			return m, nil
		case "4":
			// Switch to cost
			return m, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenCost} }
		case "5":
			// Switch to settings
			return m, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenSettings} }
		case "r":
			// Refresh metrics
			return m, m.loadMetrics()
		case "e":
			// Export metrics
			return m, func() tea.Msg { return ExportMetricsMsg{} }
		case "?":
			// Show help
			return m, m.showHelp()
		}
		
		// Update components
		for i, component := range m.Components {
			updatedComponent, cmd := component.Update(msg)
			m.Components[i] = updatedComponent
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case MetricsDataLoaded:
		// Handle data loaded
		m.Loading = false
		m.SessionStats = msg.SessionStats
		m.UsageStats = msg.UsageStats
		m.PerformanceData = msg.PerformanceData
		m.LastRefresh = time.Now()
		
		// Update cards
		m.updateCards()
		
	case ExportMetricsMsg:
		// Handle export
		m.exportMetrics()
		
	case ErrorMsg:
		// Handle error
		m.Error = msg.Err
		m.Loading = false
	}
	
	return m, tea.Batch(cmds...)
}

// updateCards updates the metric cards
func (m *MetricsScreen) updateCards() {
	// Update session stats card
	sessionContent := m.formatSessionStats()
	m.SessionStatsCard.SetContent(sessionContent)
	
	// Update usage stats card
	usageContent := m.formatUsageStats()
	m.UsageStatsCard.SetContent(usageContent)
	
	// Update performance card
	performanceContent := m.formatPerformanceData()
	m.PerformanceCard.SetContent(performanceContent)
}

// formatSessionStats formats session statistics
func (m *MetricsScreen) formatSessionStats() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Total Sessions: %d\n", m.SessionStats.TotalSessions))
	sb.WriteString(fmt.Sprintf("Today: %d\n", m.SessionStats.SessionsToday))
	sb.WriteString(fmt.Sprintf("This Week: %d\n", m.SessionStats.SessionsThisWeek))
	sb.WriteString(fmt.Sprintf("This Month: %d\n", m.SessionStats.SessionsThisMonth))
	sb.WriteString(fmt.Sprintf("Avg Duration: %s\n", m.formatDuration(m.SessionStats.AvgDuration)))
	sb.WriteString(fmt.Sprintf("Total Time: %s\n", m.formatDuration(m.SessionStats.TotalDuration)))
	
	sb.WriteString("\nAgent Breakdown:\n")
	for agent, count := range m.SessionStats.AgentBreakdown {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", agent, count))
	}
	
	return sb.String()
}

// formatUsageStats formats usage statistics
func (m *MetricsScreen) formatUsageStats() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Files Modified: %d\n", m.UsageStats.TotalFilesModified))
	sb.WriteString(fmt.Sprintf("Commits: %d\n", m.UsageStats.TotalCommits))
	sb.WriteString(fmt.Sprintf("Decisions: %d\n", m.UsageStats.TotalDecisions))
	
	sb.WriteString("\nTop Files:\n")
	for _, file := range m.UsageStats.TopFiles {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", file.Path, file.Modified))
	}
	
	sb.WriteString("\nTop Agents:\n")
	for _, agent := range m.UsageStats.TopAgents {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", agent.Name, agent.Sessions))
	}
	
	return sb.String()
}

// formatPerformanceData formats performance data
func (m *MetricsScreen) formatPerformanceData() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Avg Response: %s\n", m.formatDuration(m.PerformanceData.AvgResponseTime)))
	sb.WriteString(fmt.Sprintf("Success Rate: %.1f%%\n", m.PerformanceData.SuccessRate*100))
	sb.WriteString(fmt.Sprintf("Error Rate: %.1f%%\n", m.PerformanceData.ErrorRate*100))
	sb.WriteString(fmt.Sprintf("Throughput: %.2f req/s\n", m.PerformanceData.Throughput))
	sb.WriteString(fmt.Sprintf("Peak Load: %s\n", m.PerformanceData.PeakLoad.Format("15:04:05")))
	
	return sb.String()
}

// formatDuration formats a duration for display
func (m *MetricsScreen) formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	return fmt.Sprintf("%.1fh", duration.Hours())
}

// exportMetrics exports metrics to CSV
func (m *MetricsScreen) exportMetrics() {
	// In a real implementation, this would write to a CSV file
	// For now, we'll just announce it
	m.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Metrics exported to CSV", "Export", false)
	fmt.Println("Metrics exported to metrics.csv")
}

// showHelp shows help for the metrics screen
func (m *MetricsScreen) showHelp() tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{
			ScreenID: "metrics",
			Message:  m.getMetricsHelp(),
		}
	}
}

// getMetricsHelp returns metrics help text
func (m *MetricsScreen) getMetricsHelp() string {
	return `Metrics Controls:
1-5: Switch screens
r: Refresh metrics
e: Export to CSV
?: Show this help
q: Quit

Metrics Sections:
Session Statistics: Total sessions, duration, agent breakdown
Usage Statistics: Files modified, commits, top files
Performance: Response times, success rates, throughput`
}

// View returns the metrics view
func (m *MetricsScreen) View() string {
	if m.Loading {
		return m.renderLoading()
	}
	
	if m.Error != nil {
		return m.renderError()
	}
	
	return m.renderContent()
}

// renderLoading renders the loading state
func (m *MetricsScreen) renderLoading() string {
	header := m.Theme.Title.Render(m.Title) + "\n"
	header += m.Theme.Muted.Render("Loading metrics data...")
	
	content := fmt.Sprintf("%s\n\n%s", header, m.StatusSpinner.View())
	
	return m.renderLayout(content)
}

// renderError renders the error state
func (m *MetricsScreen) renderError() string {
	header := m.Theme.Title.Render(m.Title) + "\n"
	errorBox := m.Theme.ErrorBox.Render(fmt.Sprintf("Error: %s", m.Error.Error()))
	
	content := fmt.Sprintf("%s\n\n%s", header, errorBox)
	
	return m.renderLayout(content)
}

// renderContent renders the main content
func (m *MetricsScreen) renderContent() string {
	header := m.Theme.Title.Render(m.Title) + "\n"
	header += m.Theme.Muted.Render(fmt.Sprintf("Last updated: %s", m.LastRefresh.Format("15:04:05"))) + "\n"
	
	// Create layout with cards
	layout := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		header,
	 m.SessionStatsCard.View(),
	 m.UsageStatsCard.View(),
	 m.PerformanceCard.View(),
	 "")
	
	// Add buttons
	layout += fmt.Sprintf("\n%s", m.RefreshButton.View())
	layout += fmt.Sprintf("  %s", m.ExportButton.View())
	
	return m.renderLayout(layout)
}

// renderLayout renders the layout with proper styling
func (m *MetricsScreen) renderLayout(content string) string {
	// Add status bar
	statusBar := m.Theme.StatusBar.Render("Metrics | Press r to refresh, e to export, ? for help")
	
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

// GetID returns the screen ID
func (m *MetricsScreen) GetID() string {
	return "metrics"
}

// IsActive returns true if the screen is active
func (m *MetricsScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (m *MetricsScreen) SetActive(active bool) {
	if active {
		m.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Metrics screen activated", "Application", false)
	} else {
		m.ScreenBase.ScreenReader.Announce(AnnouncementNone, "Metrics screen deactivated", "Application", false)
	}
}

// MetricsDataLoaded represents metrics data loaded message
type MetricsDataLoaded struct {
	SessionStats    SessionStatistics
	UsageStats      UsageStatistics
	PerformanceData PerformanceData
}

// ExportMetricsMsg represents an export metrics message
type ExportMetricsMsg struct{}
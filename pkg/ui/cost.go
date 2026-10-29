package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// CostScreen represents the cost screen
// FR-TUI-004: Cost shows cost tracking, usage costs, budget information
type CostScreen struct {
	*Screen
	Components []Component
	Title      string
	
	// Cost data
	TotalCost         float64
	MonthlyBudget     float64
	RemainingBudget   float64
	UsageByAgent      map[string]float64
	CostByFile       map[string]float64
	RecentCosts      []CostEntry
	CostTrend        []CostPoint
	
	// UI components
	OverviewCard      *Card
	BudgetCard        *Card
	AgentCostCard     *Card
	FileCostCard      *Card
	TrendCard        *Card
	RefreshButton    *Button
	ResetButton      *Button
	
	// State
	Loading          bool
	LastRefresh      time.Time
	Error            error
	TimeRange        TimeRange
	Alerts           []CostAlert
}

// CostEntry represents a cost entry
type CostEntry struct {
	ID        string
	Agent     string
	File      string
	Action    string
	Cost      float64
	Timestamp time.Time
}

// CostPoint represents a cost trend point
type CostPoint struct {
	Date  time.Time
	Cost  float64
}

// CostAlert represents a cost alert
type CostAlert struct {
	Type     string
	Message  string
	Level    AnnouncementLevel
}

// NewCostScreen creates a new cost screen
func NewCostScreen(theme *Theme, keys *Keys, a11y *A11y, perfTracker *PerfTracker) *CostScreen {
	screen := NewScreen(theme, keys, perfTracker)
	
	return &CostScreen{
		Screen:           screen,
		Components:       []Component{},
		Title:           "CodeAuthorityFirst Cost Tracking",
		TotalCost:       0.0,
		MonthlyBudget:   100.0,
		RemainingBudget: 100.0,
		UsageByAgent:    make(map[string]float64),
		CostByFile:     make(map[string]float64),
		RecentCosts:    []CostEntry{},
		CostTrend:      []CostPoint{},
		OverviewCard:   NewCard("overview", "Cost Overview", theme, keys),
		BudgetCard:     NewCard("budget", "Budget Status", theme, keys),
		AgentCostCard:  NewCard("agent-costs", "Cost by Agent", theme, keys),
		FileCostCard:   NewCard("file-costs", "Cost by File", theme, keys),
		TrendCard:      NewCard("trend", "Cost Trend", theme, keys),
		RefreshButton: NewButton("refresh", "Refresh", func() tea.Cmd {
			return func() tea.Msg {
				return RefreshMsg{Force: true}
			}
		}, theme, keys),
		ResetButton:    NewButton("reset", "Reset Budget", func() tea.Cmd {
			return func() tea.Msg {
				return ResetBudgetMsg{}
			}
		}, theme, keys),
		Loading:         false,
		LastRefresh:     time.Now(),
		Error:           nil,
		TimeRange:       TimeRange{Start: time.Now().Add(-30 * 24 * time.Hour), End: time.Now()},
		Alerts:         []CostAlert{},
	}
}

// Init initializes the cost screen
func (c *CostScreen) Init() tea.Cmd {
	// Set up screen reader announcement
	c.Screen.ScreenReader.Announce(AnnouncementImportant, "Cost screen activated", "Application", false)
	
	// Load initial data
	return c.loadCostData()
}

// loadCostData loads cost data
func (c *CostScreen) loadCostData() tea.Cmd {
	return func() tea.Msg {
		// Simulate data loading
		time.Sleep(600 * time.Millisecond)
		
		// Generate mock data
		totalCost := 45.67
		monthlyBudget := 100.0
		remainingBudget := monthlyBudget - totalCost
		usageByAgent := map[string]float64{
			"opencode":      25.34,
			"claude-code":   15.89,
			"universal-lsp": 4.44,
		}
		costByFile := map[string]float64{
			"internal/api/server.go":      12.45,
			"internal/database/connection.go": 8.92,
			"cmd/main.go":                 6.78,
			"internal/auth/middleware.go": 5.34,
			"internal/config/config.go":   4.23,
		}
		recentCosts := c.generateMockRecentCosts()
		costTrend := c.generateMockCostTrend()
		alerts := c.generateMockAlerts(totalCost, monthlyBudget)
		
		return CostDataLoaded{
			TotalCost:       totalCost,
			MonthlyBudget:   monthlyBudget,
			RemainingBudget: remainingBudget,
			UsageByAgent:    usageByAgent,
			CostByFile:     costByFile,
			RecentCosts:    recentCosts,
			CostTrend:      costTrend,
			Alerts:         alerts,
		}
	}
}

// generateMockRecentCosts generates mock recent costs
func (c *CostScreen) generateMockRecentCosts() []CostEntry {
	return []CostEntry{
		{
			ID:        "cost-001",
			Agent:     "opencode",
			File:      "internal/api/server.go",
			Action:    "modified",
			Cost:      2.34,
			Timestamp: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:        "cost-002",
			Agent:     "claude-code",
			File:      "internal/auth/middleware.go",
			Action:    "created",
			Cost:      1.89,
			Timestamp: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "cost-003",
			Agent:     "universal-lsp",
			File:      "cmd/main.go",
			Action:    "fixed",
			Cost:      0.67,
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "cost-004",
			Agent:     "opencode",
			File:      "internal/database/connection.go",
			Action:    "refactored",
			Cost:      3.45,
			Timestamp: time.Now().Add(-3 * time.Hour),
		},
		{
			ID:        "cost-005",
			Agent:     "claude-code",
			File:      "internal/config/config.go",
			Action:    "modified",
			Cost:      1.23,
			Timestamp: time.Now().Add(-4 * time.Hour),
		},
	}
}

// generateMockCostTrend generates mock cost trend
func (c *CostScreen) generateMockCostTrend() []CostPoint {
	trend := []CostPoint{}
	now := time.Now()
	
	for i := 29; i >= 0; i-- {
		date := now.Add(-time.Duration(i) * 24 * time.Hour)
		cost := 1.0 + float64(i)*0.1 // Increasing trend
		trend = append(trend, CostPoint{Date: date, Cost: cost})
	}
	
	return trend
}

// generateMockAlerts generates mock cost alerts
func (c *CostScreen) generateMockAlerts(totalCost, budget float64) []CostAlert {
	alerts := []CostAlert{}
	
	// Budget alert
	if totalCost > budget*0.8 {
		alerts = append(alerts, CostAlert{
			Type:    "budget",
			Message: fmt.Sprintf("Cost at %.0f%% of budget", (totalCost/budget)*100),
			Level:   AnnouncementWarning,
		})
	}
	
	// High cost alert
	if totalCost > budget*0.9 {
		alerts = append(alerts, CostAlert{
			Type:    "budget",
			Message: "Cost exceeding 90% of budget!",
			Level:   AnnouncementCritical,
		})
	}
	
	return alerts
}

// Update updates the cost screen
func (c *CostScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.Type {
		case tea.Key1:
			// Switch to dashboard
			return c, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDashboard} }
		case tea.Key2:
			// Switch to doctor
			return c, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDoctor} }
		case tea.Key3:
			// Switch to metrics
			return c, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenMetrics} }
		case tea.Key4:
			// Stay on cost
			return c, nil
		case tea.Key5:
			// Switch to settings
			return c, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenSettings} }
		case tea.KeyR:
			// Refresh cost data
			return c, c.loadCostData()
		case tea.KeyQuestion:
			// Show help
			return c, c.showHelp()
		}
		
		// Update components
		for i, component := range c.Components {
			updatedComponent, cmd := component.Update(msg)
			c.Components[i] = updatedComponent
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case CostDataLoaded:
		// Handle data loaded
		c.Loading = false
		c.TotalCost = msg.TotalCost
		c.MonthlyBudget = msg.MonthlyBudget
		c.RemainingBudget = msg.RemainingBudget
		c.UsageByAgent = msg.UsageByAgent
		c.CostByFile = msg.CostByFile
		c.RecentCosts = msg.RecentCosts
		c.CostTrend = msg.CostTrend
		c.Alerts = msg.Alerts
		c.LastRefresh = time.Now()
		
		// Update cards
		c.updateCards()
		
	case ResetBudgetMsg:
		// Handle budget reset
		c.MonthlyBudget = 100.0
		c.RemainingBudget = 100.0 - c.TotalCost
		c.Screen.ScreenReader.Announce(AnnouncementImportant, "Budget reset to $100.00", "Budget", false)
		
	case ErrorMsg:
		// Handle error
		c.Error = msg.Err
		c.Loading = false
	}
	
	return c, tea.Batch(cmds...)
}

// updateCards updates the cost cards
func (c *CostScreen) updateCards() {
	// Update overview card
	overviewContent := c.formatOverview()
	c.OverviewCard.SetContent(overviewContent)
	
	// Update budget card
	budgetContent := c.formatBudget()
	c.BudgetCard.SetContent(budgetContent)
	
	// Update agent cost card
	agentContent := c.formatAgentCosts()
	c.AgentCostCard.SetContent(agentContent)
	
	// Update file cost card
	fileContent := c.formatFileCosts()
	c.FileCostCard.SetContent(fileContent)
	
	// Update trend card
	trendContent := c.formatTrend()
	c.TrendCard.SetContent(trendContent)
}

// formatOverview formats cost overview
func (c *CostScreen) formatOverview() string {
	budgetPercentage := (c.TotalCost / c.MonthlyBudget) * 100
	budgetStatus := "Within Budget"
	if budgetPercentage > 90 {
		budgetStatus = "Over Budget"
	} else if budgetPercentage > 80 {
		budgetStatus = "Near Budget"
	}
	
	return fmt.Sprintf("Total Cost: $%.2f\nBudget: $%.2f\nUsed: %.1f%%\nStatus: %s",
		c.TotalCost, c.MonthlyBudget, budgetPercentage, budgetStatus)
}

// formatBudget formats budget information
func (c *CostScreen) formatBudget() string {
	budgetPercentage := (c.TotalCost / c.MonthlyBudget) * 100
	
	var progressBar strings.Builder
	progressBar.WriteString("Budget: [")
	for i := 0; i < 20; i++ {
		if i < int(budgetPercentage/5) {
			progressBar.WriteString("=")
		} else {
			progressBar.WriteString(" ")
		}
	}
	progressBar.WriteString("]")
	
	return fmt.Sprintf("%s\nRemaining: $%.2f\nDays Left: %d",
		progressBar.String(), c.RemainingBudget, 25)
}

// formatAgentCosts formats costs by agent
func (c *CostScreen) formatAgentCosts() string {
	var sb strings.Builder
	
	for agent, cost := range c.UsageByAgent {
		percentage := (cost / c.TotalCost) * 100
		sb.WriteString(fmt.Sprintf("%s: $%.2f (%.1f%%)\n", agent, cost, percentage))
	}
	
	return sb.String()
}

// formatFileCosts formats costs by file
func (c *CostScreen) formatFileCosts() string {
	var sb strings.Builder
	
	for file, cost := range c.CostByFile {
		sb.WriteString(fmt.Sprintf("%s: $%.2f\n", file, cost))
	}
	
	return sb.String()
}

// formatTrend formats cost trend
func (c *CostScreen) formatTrend() string {
	if len(c.CostTrend) == 0 {
		return "No trend data available"
	}
	
	var sb strings.Builder
	
	// Show last 7 days
	start := len(c.CostTrend) - 7
	if start < 0 {
		start = 0
	}
	
	for i := start; i < len(c.CostTrend); i++ {
		point := c.CostTrend[i]
		sb.WriteString(fmt.Sprintf("%s: $%.2f\n", point.Date.Format("01/02"), point.Cost))
	}
	
	return sb.String()
}

// showHelp shows help for the cost screen
func (c *CostScreen) showHelp() tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{
			ScreenID: "cost",
			Message:  c.getCostHelp(),
		}
	}
}

// getCostHelp returns cost help text
func (c *CostScreen) getCostHelp() string {
	return `Cost Controls:
1-5: Switch screens
r: Refresh cost data
?: Show this help
q: Quit

Cost Sections:
Overview: Total cost, budget usage, status
Budget: Progress bar, remaining budget
Agent Costs: Cost breakdown by agent
File Costs: Cost breakdown by file
Trend: Daily cost trends over time

Alerts:
- Yellow: Cost near budget (80%)
- Red: Cost over budget (90%)`
}

// View returns the cost view
func (c *CostScreen) View() string {
	if c.Loading {
		return c.renderLoading()
	}
	
	if c.Error != nil {
		return c.renderError()
	}
	
	return c.renderContent()
}

// renderLoading renders the loading state
func (c *CostScreen) renderLoading() string {
	header := c.Theme.Title.Render(c.Title) + "\n"
	header += c.Theme.Muted.Render("Loading cost data...")
	
	content := fmt.Sprintf("%s\n\n%s", header, c.StatusSpinner.View())
	
	return c.renderLayout(content)
}

// renderError renders the error state
func (c *CostScreen) renderError() string {
	header := c.Theme.Title.Render(c.Title) + "\n"
	errorBox := c.Theme.ErrorBox.Render(fmt.Sprintf("Error: %s", c.Error.Error()))
	
	content := fmt.Sprintf("%s\n\n%s", header, errorBox)
	
	return c.renderLayout(content)
}

// renderContent renders the main content
func (c *CostScreen) renderContent() string {
	header := c.Theme.Title.Render(c.Title) + "\n"
	header += c.Theme.Muted.Render(fmt.Sprintf("Last updated: %s", c.LastRefresh.Format("15:04:05"))) + "\n"
	
	// Show alerts
	if len(c.Alerts) > 0 {
		header += "\n" + c.Theme.Title.Render("Alerts:") + "\n"
		for _, alert := range c.Alerts {
			alertStyle := c.Theme.GetColor(alert.Type)
			header += fmt.Sprintf("%s %s\n", alertStyle.Render("⚠"), alert.Message)
		}
	}
	
	// Create layout with cards
	layout := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		header,
	 c.OverviewCard.View(),
	 c.BudgetCard.View(),
	 c.AgentCostCard.View(),
	 c.FileCostCard.View(),
	 c.TrendCard.View())
	
	// Add buttons
	layout += fmt.Sprintf("\n\n%s", c.RefreshButton.View())
	layout += fmt.Sprintf("  %s", c.ResetButton.View())
	
	return c.renderLayout(layout)
}

// renderLayout renders the layout with proper styling
func (c *CostScreen) renderLayout(content string) string {
	// Add status bar
	statusBar := c.Theme.StatusBar.Render("Cost | Press r to refresh, ? for help")
	
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

// GetID returns the screen ID
func (c *CostScreen) GetID() string {
	return "cost"
}

// IsActive returns true if the screen is active
func (c *CostScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (c *CostScreen) SetActive(active bool) {
	if active {
		c.Screen.ScreenReader.Announce(AnnouncementImportant, "Cost screen activated", "Application", false)
	} else {
		c.Screen.ScreenReader.Announce(AnnouncementNone, "Cost screen deactivated", "Application", false)
	}
}

// CostDataLoaded represents cost data loaded message
type CostDataLoaded struct {
	TotalCost       float64
	MonthlyBudget   float64
	RemainingBudget float64
	UsageByAgent    map[string]float64
	CostByFile     map[string]float64
	RecentCosts    []CostEntry
	CostTrend      []CostPoint
	Alerts         []CostAlert
}

// ResetBudgetMsg represents a reset budget message
type ResetBudgetMsg struct{}

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
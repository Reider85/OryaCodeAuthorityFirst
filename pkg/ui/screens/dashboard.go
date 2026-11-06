// dashboard.go — TUI-дашборд, хаб-экран (PROMPT-MVP §4, §2.6).
package screens

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui/components"
)

// widgetTitleOrder — заголовки виджетов в порядке WidgetIDs (фокус-маркер
// сравнивает заголовок с позицией, а не с ID: INV-PROMPT-MVP-007).
var widgetTitleOrder = []string{"Today", "Authorship Score", "Active Agents", "Recent Decisions", "Cross-repo"}

// DashboardScreen represents the dashboard screen
type DashboardScreen struct {
	ui.StateMachine
	
	Version      string
	RepoRoot     string
	Theme        *ui.Theme
	Keys         *ui.Keys
	A11y         *ui.A11y
	StatusBar    *components.StatusBar
	
	// Data
	Sessions     []*storage.Session
	Manifest     *storage.Manifest
	Adapters     []caf.Adapter
	
	// Widgets
	TodayWidget        *TodayWidget
	AuthorshipWidget   *AuthorshipWidget
	AgentsWidget       *AgentsWidget
	DecisionsWidget    *DecisionsWidget
	CrossRepoWidget    *CrossRepoWidget
	
	// Loading state
	Loading      bool
	LoadError    error
	RefreshTime  time.Time
	
	// Focus
	CurrentWidget int
	WidgetIDs     []string

	// Help overlay (§18.4: ? показывает help на каждом экране)
	ShowHelp bool
	
	// Widget data caching for performance optimization
	widgetDataCache     *WidgetDataCache
	widgetCacheValid    bool
	lastSessionsHash    string
	lastManifestHash    string
	lastAdaptersHash    string
	
	// Memory pooling for performance optimization
	stringPool          *StringPool
	lastPoolReset       time.Time
}

// WidgetDataCache holds pre-computed widget data to avoid recomputation
type WidgetDataCache struct {
	TodayData        TodayWidgetData
	AuthorshipData   AuthorshipWidgetData
	AgentsData       AgentsWidgetData
	DecisionsData    DecisionsWidgetData
	CrossRepoData    CrossRepoWidgetData
	GridLayout       GridLayoutData
	BorderStrings    BorderStringCache
	WidgetPriority   []string // Order of widget rendering priority
	LastRenderTime   time.Time
	WidgetRenderStats map[string]WidgetRenderStats
}

// WidgetRenderStats tracks performance metrics for individual widgets
type WidgetRenderStats struct {
	LastRenderTime  time.Time
	RenderCount     int
	AvgRenderTime   time.Duration
	IsVisible       bool
}

// Widget data structures
type TodayWidgetData struct {
	CardsDue        int
	PendingDecisions int
	NewBlindspots   int
}

type AuthorshipWidgetData struct {
	ACS       float64
	Recall    float64
	Coverage  float64
	Trend     string
}

type AgentsWidgetData struct {
	Agents []AgentStatus
}

type DecisionsWidgetData struct {
	Decisions []storage.Decision
	Limit     int
}

type CrossRepoWidgetData struct {
	Enabled   bool
	Repos     []string
	LastSync  time.Time
	Stale     bool
}

type GridLayoutData struct {
	WidgetWidth  int
	WidgetHeight int
	WidgetPositions map[string][2]int // widget name -> [row, col]
}

type BorderStringCache struct {
	Top    string
	Bottom string
	Left   string
	Right  string
}

// StringPool provides a pool of reusable strings to reduce allocations
type StringPool struct {
	strings     chan string
	poolSize   int
	maxLength  int
}

// NewStringPool creates a new string pool
func NewStringPool(poolSize, maxLength int) *StringPool {
	return &StringPool{
		strings:    make(chan string, poolSize),
		poolSize:   poolSize,
		maxLength:  maxLength,
	}
}

// GetString gets a string from the pool or creates a new one
func (sp *StringPool) GetString(s string) string {
	select {
	case pooled := <-sp.strings:
		// Reuse existing string
		if len(pooled) >= len(s) {
			copy([]byte(pooled), s)
			return pooled[:len(s)]
		}
		return s
	default:
		// Pool is empty, create new string
		if len(s) <= sp.maxLength {
			return s
		}
		return s
	}
}

// ReleaseString returns a string to the pool if it fits
func (sp *StringPool) ReleaseString(s string) {
	if len(s) <= sp.maxLength {
		select {
		case sp.strings <- s:
			// String returned to pool
		default:
			// Pool is full, discard string
		}
	}
}

// ResetPool clears the pool and frees memory
func (sp *StringPool) ResetPool() {
	for {
		select {
		case <-sp.strings:
			// Clear all strings from pool
		default:
			return
		}
	}
}

// Widget interface for dashboard widgets
type Widget interface {
	ID() string
	Render(width, height int, theme *ui.Theme, a11y *ui.A11y) string
	Update(data interface{}) error
	IsFocusable() bool
}

// TodayWidget shows today's activity
type TodayWidget struct {
	CardsDue      int
	PendingDecisions int
	NewBlindspots int
}

// AuthorshipWidget shows ACS and related metrics
type AuthorshipWidget struct {
	ACS       float64
	Recall    float64
	Coverage  float64
	Trend     string // 30-day trend
}

// AgentsWidget shows active agents
type AgentsWidget struct {
	Agents []AgentStatus
}

// AgentStatus represents an agent's status
type AgentStatus struct {
	Name     string
	ID       string
	Active   bool
	Health   string
}

// DecisionsWidget shows recent decisions
type DecisionsWidget struct {
	Decisions []storage.Decision
	Limit     int
}

// CrossRepoWidget shows cross-repo sync status
type CrossRepoWidget struct {
	Enabled   bool
	Repos     []string
	LastSync  time.Time
	Stale     bool
}

// NewDashboardScreen creates a new dashboard screen
func NewDashboardScreen(repoRoot string, version string) *DashboardScreen {
	theme := ui.NewTheme()
	keys := ui.NewKeys()
	a11y := ui.NewA11y(ui.NewAccessibilityConfig(), theme, keys)
	
	screen := &DashboardScreen{
		StateMachine: *ui.NewStateMachine(),
		Version:      version,
		RepoRoot:     repoRoot,
		Theme:        theme,
		Keys:         keys,
		A11y:         a11y,
		Loading:      true,
		CurrentWidget: 0,
		WidgetIDs:     []string{"today", "authorship", "agents", "decisions", "crossrepo"},
		widgetDataCache: &WidgetDataCache{
			GridLayout: GridLayoutData{
				WidgetPositions: make(map[string][2]int),
			},
			BorderStrings: BorderStringCache{
				Top:    "┌" + strings.Repeat("─", 19) + "┐",
				Bottom: "└" + strings.Repeat("─", 19) + "┘",
				Left:   "│",
				Right:  "│",
			},
			WidgetPriority: []string{"Today", "Authorship Score", "Active Agents", "Recent Decisions", "Cross-repo"},
			WidgetRenderStats: make(map[string]WidgetRenderStats),
		},
		widgetCacheValid: false,
		stringPool:       NewStringPool(10, 256), // Pool of 10 strings, max 256 chars each
		lastPoolReset:    time.Now(),
	}
	
	// Initialize status bar
	screen.StatusBar = components.NewStatusBar("dashboard")
	screen.StatusBar.SetKeyHints([]string{
		"[Tab] widgets", "[r] refresh", "[?] help", "[q] quit",
	})
	
	// Initialize widgets
	screen.TodayWidget = &TodayWidget{}
	screen.AuthorshipWidget = &AuthorshipWidget{}
	screen.AgentsWidget = &AgentsWidget{}
	screen.DecisionsWidget = &DecisionsWidget{Limit: 5}
	screen.CrossRepoWidget = &CrossRepoWidget{}
	
	// Initialize widget render stats
	screen.initWidgetRenderStats()
	
	return screen
}

// Init initializes the dashboard screen
func (d *DashboardScreen) Init() tea.Cmd {
	return d.LoadData()
}

// LoadData loads dashboard data
func (d *DashboardScreen) LoadData() tea.Cmd {
	return func() tea.Msg {
	// Load sessions (параллельный bulk-load — бюджет §6 file ops <1s на 10k)
	sessionsDir := filepath.Join(d.RepoRoot, ".caf", "sessions")
	sessionManager := storage.NewManager(sessionsDir)
	// errors.Is, а НЕ os.IsNotExist: ошибки storage обёрнуты через %w,
	// os.IsNotExist их не раскапывает — свежий репо падал в Error state
	// вместо Empty-CTA (§2.6, INV-PROMPT-MVP-006).
	sessions, err := sessionManager.LoadAllSessions()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return DashboardLoadError{Err: fmt.Errorf("failed to load sessions: %w", err)}
	}
		
		// Load manifest
		manifestFile := storage.GetManifestPath(d.RepoRoot)
		manifestManager := storage.NewManifestManager(manifestFile)
		manifest, err := manifestManager.LoadManifest()
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return DashboardLoadError{Err: fmt.Errorf("failed to load manifest: %w", err)}
		}
		
	// Get adapters
	adapters := caf.ListAdapters()
	var adapterList []caf.Adapter
	for _, id := range adapters {
		adapter, ok := caf.GetAdapter(id)
		if ok {
			adapterList = append(adapterList, adapter)
		}
	}
		
		return DashboardDataLoaded{
			Sessions: sessions,
			Manifest: manifest,
			Adapters: adapterList,
			LoadTime: time.Now(),
		}
	}
}

// DashboardDataLoaded is the message for data loading
type DashboardDataLoaded struct {
	Sessions []*storage.Session
	Manifest *storage.Manifest
	Adapters []caf.Adapter
	LoadTime time.Time
}

// DashboardLoadError is the error message for data loading
type DashboardLoadError struct {
	Err error
}

// Update handles messages for the dashboard screen
func (d *DashboardScreen) Update(msg tea.Msg) (ui.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return d.handleKey(msg)
	case DashboardDataLoaded:
		d.handleDataLoaded(msg)
	case DashboardLoadError:
		d.handleLoadError(msg)
	case tea.WindowSizeMsg:
		// Handle resize if needed
	}
	
	// Update state machine
	d.StateMachine.Transition(nil)
	
	return d, nil
}

// handleKey handles keyboard input
func (d *DashboardScreen) handleKey(msg tea.KeyMsg) (ui.Screen, tea.Cmd) {
	// Help-open режим: ?/Esc закрывают оверлей, q/Ctrl+C выходят, остальные
	// клавиши не действуют, пока открыт help (нет магии под оверлеем).
	if d.ShowHelp {
		switch {
		case runeOf(msg) == '?', msg.Type == tea.KeyEsc:
			d.ShowHelp = false
			return d, nil
		case runeOf(msg) == 'q', msg.Type == tea.KeyCtrlC:
			return d, tea.Quit
		}
		return d, nil
	}

	// Печатные клавиши — через runeOf: bubbletea кодирует их как KeyRunes (§18.4).
	switch runeOf(msg) {
	case 'r':
		d.Loading = true
		d.LoadError = nil
		// §16.4: Success → Loading на refresh (скелет вместо устаревших данных).
		d.StateMachine.SetState(ui.StateLoading)
		return d, d.LoadData()
	case 'q':
		// Статбар обещает "[q] quit" — q обязан выходить (§18.4, INV-UI-009).
		return d, tea.Quit
	case '?':
		d.ShowHelp = true
		return d, nil
	}

	switch msg.Type {
	case tea.KeyTab:
		d.CurrentWidget = (d.CurrentWidget + 1) % len(d.WidgetIDs)
		d.StatusBar.SetFocus(d.WidgetIDs[d.CurrentWidget])
	case tea.KeyShiftTab:
		d.CurrentWidget = (d.CurrentWidget - 1 + len(d.WidgetIDs)) % len(d.WidgetIDs)
		d.StatusBar.SetFocus(d.WidgetIDs[d.CurrentWidget])
	case tea.KeyEsc, tea.KeyCtrlC:
		return d, tea.Quit
	}
	
	return d, nil
}

// handleDataLoaded handles successful data loading
func (d *DashboardScreen) handleDataLoaded(msg DashboardDataLoaded) {
	d.Sessions = msg.Sessions
	d.Manifest = msg.Manifest
	d.Adapters = msg.Adapters
	d.Loading = false
	d.LoadError = nil
	d.RefreshTime = msg.LoadTime
	
	// Invalidate cache and update widgets with data
	d.widgetCacheValid = false
	d.updateWidgets()
	
	// 0 сессий → empty-with-CTA «Run caf init» (§2.6), не happy-path.
	if len(msg.Sessions) == 0 {
		d.StateMachine.SetState(ui.StateEmpty)
		return
	}
	// Transition to success state
	d.StateMachine.SetState(ui.StateSuccess)
}

// handleLoadError handles data loading errors
func (d *DashboardScreen) handleLoadError(msg DashboardLoadError) {
	d.Loading = false
	d.LoadError = msg.Err
	d.StateMachine.SetState(ui.StateError)
}

// updateWidgets updates all widgets with current data (cached version)
func (d *DashboardScreen) updateWidgets() {
	// Check if we need to recompute widget data
	currentSessionsHash := d.computeDataHash(d.Sessions)
	currentManifestHash := d.computeDataHash(d.Manifest)
	currentAdaptersHash := d.computeDataHash(d.Adapters)
	
	if d.widgetCacheValid && 
	   currentSessionsHash == d.lastSessionsHash && 
	   currentManifestHash == d.lastManifestHash && 
	   currentAdaptersHash == d.lastAdaptersHash {
		return // Cache is valid, skip recomputation
	}
	
	// Recompute widget data
	d.computeWidgetData(currentSessionsHash, currentManifestHash, currentAdaptersHash)
	d.widgetCacheValid = true
	d.lastSessionsHash = currentSessionsHash
	d.lastManifestHash = currentManifestHash
	d.lastAdaptersHash = currentAdaptersHash
}

// computeDataHash creates a simple hash for data comparison
func (d *DashboardScreen) computeDataHash(data interface{}) string {
	switch v := data.(type) {
	case []*storage.Session:
		return fmt.Sprintf("%d", len(v))
	case *storage.Manifest:
		if v == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", v.LastSync)
	case []caf.Adapter:
		return fmt.Sprintf("%d", len(v))
	default:
		return "unknown"
	}
}

// computeWidgetData computes all widget data and updates cache
func (d *DashboardScreen) computeWidgetData(sessionsHash, manifestHash, adaptersHash string) {
	// Update Today Widget
	today := time.Now()
	var cardsDue, pendingDecisions, newBlindspots int
	for _, session := range d.Sessions {
		sessionStart := session.StartedAt
		if sessionStart.Year() == today.Year() && sessionStart.Month() == today.Month() && sessionStart.Day() == today.Day() {
			cardsDue += len(session.FilesTouched)
			pendingDecisions += len(session.DecisionsMade)
			// Simple blindspot detection: files touched but not committed
			if session.CommitSHA == "" {
				newBlindspots++
			}
		}
	}
	
	// Update cache
	d.widgetDataCache.TodayData = TodayWidgetData{
		CardsDue:        cardsDue,
		PendingDecisions: pendingDecisions,
		NewBlindspots:   newBlindspots,
	}
	
	// Update Authorship Widget - honest empty/partial for MVP
	if len(d.Sessions) > 0 {
		// Simple proxy: average files per session
		var totalFiles int
		for _, session := range d.Sessions {
			totalFiles += len(session.FilesTouched)
		}
		d.widgetDataCache.AuthorshipData = AuthorshipWidgetData{
			ACS:       float64(totalFiles) / float64(len(d.Sessions)) * 0.1, // Very naive proxy
			Recall:    0.0, // Not tracked in MVP
			Coverage:  0.0, // Not tracked in MVP
			Trend:     "→", // No trend calculation in MVP
		}
	}
	
	// Update Agents Widget
	var agents []AgentStatus
	for _, adapter := range d.Adapters {
		status := "unknown"
		// Simple health check: if adapter has methods, assume healthy
		if adapter.ID() != "" {
			status = "healthy"
		}
		agents = append(agents, AgentStatus{
			Name:   adapter.ID(),
			ID:     adapter.ID(),
			Active: true,
			Health: status,
		})
	}
	d.widgetDataCache.AgentsData = AgentsWidgetData{Agents: agents}
	
	// Update Decisions Widget
	var recentDecisions []storage.Decision
	for _, session := range d.Sessions {
		recentDecisions = append(recentDecisions, session.DecisionsMade...)
	}
	// Sort by timestamp and take most recent
	for i := 0; i < len(recentDecisions)-1; i++ {
		for j := i + 1; j < len(recentDecisions); j++ {
			if recentDecisions[i].Timestamp.Before(recentDecisions[j].Timestamp) {
				recentDecisions[i], recentDecisions[j] = recentDecisions[j], recentDecisions[i]
			}
		}
	}
	if len(recentDecisions) > d.DecisionsWidget.Limit {
		recentDecisions = recentDecisions[:d.DecisionsWidget.Limit]
	}
	d.widgetDataCache.DecisionsData = DecisionsWidgetData{
		Decisions: recentDecisions,
		Limit:     d.DecisionsWidget.Limit,
	}
	
	// Update Cross Repo Widget
	if d.Manifest != nil {
		d.widgetDataCache.CrossRepoData = CrossRepoWidgetData{
			Enabled:   d.Manifest.CrossRepo.Enabled,
			Repos:     d.Manifest.CrossRepo.Repos,
			LastSync:  d.Manifest.LastSync,
			Stale:     time.Since(d.Manifest.LastSync) > 30*time.Minute,
		}
	}
	
	// Update layout cache
	d.updateGridLayoutCache()
}

// updateGridLayoutCache updates cached layout calculations
func (d *DashboardScreen) updateGridLayoutCache() {
	width := 80 // Default width
	height := 24 // Default height
	
	widgetWidth := width / 4
	widgetHeight := height / 3
	
	d.widgetDataCache.GridLayout = GridLayoutData{
		WidgetWidth: widgetWidth,
		WidgetHeight: widgetHeight,
		WidgetPositions: map[string][2]int{
			"Today":          {0, 0},
			"Authorship":     {0, 1},
			"Agents":         {0, 2},
			"Recent Decisions": {1, 0},
			"Cross-repo":    {1, 1},
		},
	}
}

// initWidgetRenderStats initializes performance tracking for each widget
func (d *DashboardScreen) initWidgetRenderStats() {
	for _, widgetName := range d.widgetDataCache.WidgetPriority {
		d.widgetDataCache.WidgetRenderStats[widgetName] = WidgetRenderStats{
			LastRenderTime: time.Now(),
			RenderCount:    0,
			AvgRenderTime:  0,
			IsVisible:      false,
		}
	}
}

// getVisibleWidgets determines which widgets should be rendered based on priority and visibility
func (d *DashboardScreen) getVisibleWidgets() []string {
	var visibleWidgets []string
	
	// Always render Today widget - it's lightweight and important
	if d.widgetDataCache.TodayData.CardsDue > 0 || d.widgetDataCache.TodayData.PendingDecisions > 0 {
		visibleWidgets = append(visibleWidgets, "Today")
	}
	
	// Render Authorship if we have data
	if d.widgetDataCache.AuthorshipData.ACS > 0 {
		visibleWidgets = append(visibleWidgets, "Authorship Score")
	}
	
	// Always render Agents - lightweight and important for user awareness
	if len(d.widgetDataCache.AgentsData.Agents) > 0 {
		visibleWidgets = append(visibleWidgets, "Active Agents")
	}
	
	// Render Decisions if we have recent decisions
	if len(d.widgetDataCache.DecisionsData.Decisions) > 0 {
		visibleWidgets = append(visibleWidgets, "Recent Decisions")
	}
	
	// Only render Cross-repo if enabled
	if d.widgetDataCache.CrossRepoData.Enabled {
		visibleWidgets = append(visibleWidgets, "Cross-repo")
	}
	
	// If no widgets are visible, show Today as default
	if len(visibleWidgets) == 0 {
		visibleWidgets = append(visibleWidgets, "Today")
	}
	
	return visibleWidgets
}

// updateWidgetRenderStats updates performance tracking for a widget
func (d *DashboardScreen) updateWidgetRenderStats(widgetName string, renderTime time.Duration) {
	stats, exists := d.widgetDataCache.WidgetRenderStats[widgetName]
	if !exists {
		return
	}
	
	stats.RenderCount++
	stats.LastRenderTime = time.Now()
	
	// Update average render time using exponential moving average
	if stats.AvgRenderTime == 0 {
		stats.AvgRenderTime = renderTime
	} else {
		stats.AvgRenderTime = (stats.AvgRenderTime + renderTime) / 2
	}
	
	d.widgetDataCache.WidgetRenderStats[widgetName] = stats
}

// shouldSkipWidget determines if a widget should be skipped based on performance metrics
func (d *DashboardScreen) shouldSkipWidget(widgetName string) bool {
	stats, exists := d.widgetDataCache.WidgetRenderStats[widgetName]
	if !exists {
		return false
	}
	
	// Skip widget if:
	// 1. It's taking too long to render (> 20ms)
	// 2. It hasn't been updated recently (> 1 second)
	// 3. It's not in the current visible widgets list
	visibleWidgets := d.getVisibleWidgets()
	isVisible := false
	for _, visible := range visibleWidgets {
		if visible == widgetName {
			isVisible = true
			break
		}
	}
	
	if !isVisible {
		return true
	}
	
	if stats.AvgRenderTime > 20*time.Millisecond {
		// Log performance warning
		return false // Still render but log performance issue
	}
	
	return false
}

// View returns the dashboard view (optimized version)
func (d *DashboardScreen) View() string {
	// Update status bar
	d.StatusBar.SetRefreshCompleted()
	
	// Help overlay (§18.4): ? на каждом экране показывает все клавиши.
	// Header и StatusBar — «священные» (§15.2), оверлей заменяет только контент.
	if d.ShowHelp {
		overlay := components.NewHelpOverlay("dashboard", d.helpSections(), d.Theme)
		return components.Header(d.Version, "dashboard") + "\n" + overlay.View() + "\n" + d.StatusBar.Render()
	}
	
	// Render based on current state
	switch d.StateMachine.Current() {
	case ui.StateLoading:
		return d.renderSkeleton()
	case ui.StateError:
		return d.renderError()
	case ui.StateEmpty:
		return d.renderEmpty()
	case ui.StatePartial:
		return d.renderPartial()
	case ui.StateSuccess:
		return d.renderSuccess()
	default:
		return d.renderSkeleton()
	}
}

// renderSkeleton renders the loading skeleton
func (d *DashboardScreen) renderSkeleton() string {
	width := 80 // Default width
	height := 24 // Default height
	
	// Render skeleton widgets
	skeleton := d.renderWidgetSkeleton("Today", width/4-2, height/2-1)
	skeleton += d.renderWidgetSkeleton("Authorship Score", width/4-2, height/2-1)
	skeleton += d.renderWidgetSkeleton("Active Agents", width/4-2, height/2-1)
	skeleton += d.renderWidgetSkeleton("Recent Decisions", width/4-2, height/2-1)
	skeleton += d.renderWidgetSkeleton("Cross-repo", width/4-2, height/2-1)
	
	return components.Header(d.Version, "dashboard") + "\n" + skeleton + "\n" + d.StatusBar.Render()
}

// renderWidgetSkeleton renders a skeleton widget (optimized version)
func (d *DashboardScreen) renderWidgetSkeleton(title string, width, height int) string {
	var builder strings.Builder
	builder.Grow(width * height) // Pre-allocate for skeleton content
	
	// Use cached border strings with adjusted width
	border := "┌" + strings.Repeat("─", width-2) + "┐"
	footer := "└" + strings.Repeat("─", width-2) + "┘"
	content := strings.Repeat("│" + strings.Repeat(" ", width-2) + "│\n", height-2)
	
	builder.WriteString("  ") // Focus prefix for skeleton
	builder.WriteString(title)
	builder.WriteString("\n")
	builder.WriteString(border)
	builder.WriteString("\n")
	builder.WriteString(content)
	builder.WriteString(footer)
	builder.WriteString("\n")
	
	return builder.String()
}

// renderError renders the error state (optimized version)
func (d *DashboardScreen) renderError() string {
	var builder strings.Builder
	builder.Grow(1000) // Pre-allocate for error content
	
	builder.WriteString(components.Header(d.Version, "dashboard"))
	builder.WriteString("\n")
	
	errorBox := components.NewErrorBox(fmt.Sprintf("Failed to load dashboard: %v", d.LoadError), d.Theme, d.A11y)
	builder.WriteString(errorBox.View())
	builder.WriteString("\n")
	builder.WriteString(d.StatusBar.Render())
	
	return builder.String()
}

// renderEmpty renders the empty state (optimized version)
func (d *DashboardScreen) renderEmpty() string {
	var builder strings.Builder
	builder.Grow(1000) // Pre-allocate for empty state content
	
	builder.WriteString(components.Header(d.Version, "dashboard"))
	builder.WriteString("\n")
	
	// Подсказка обязана называть живую клавишу: [i] на дашборде не обрабатывается,
	// честный путь — запустить caf init и освежить [r] (INV-UI-009, §18.4).
	emptyCTA := components.NewEmptyCTA(
		"No sessions found",
		"caf hasn't been invoked in this repo yet",
		"Run `caf init` to start tracking",
		"Run `caf init` in your terminal, then press [r] to refresh",
		d.Theme,
	)
	builder.WriteString(emptyCTA.View())
	builder.WriteString("\n")
	builder.WriteString(d.StatusBar.Render())
	
	return builder.String()
}

// renderPartial renders the partial state (optimized version)
func (d *DashboardScreen) renderPartial() string {
	var builder strings.Builder
	builder.Grow(2000) // Pre-allocate for partial state content
	
	builder.WriteString(components.Header(d.Version, "dashboard"))
	builder.WriteString("\n")
	
	// Add partial banner if needed
	partialBanner := components.NewPartialBanner(1, 2, []string{"cross-repo"}, d.Theme)
	builder.WriteString(partialBanner.View())
	builder.WriteString("\n")
	
	// Render partial widgets with some data using optimized method
	widgets := d.renderWidgetsOptimized()
	builder.WriteString(widgets)
	builder.WriteString("\n")
	builder.WriteString(d.StatusBar.Render())
	
	return builder.String()
}

// renderSuccess renders the success state with actual data (optimized version)
func (d *DashboardScreen) renderSuccess() string {
	var builder strings.Builder
	
	// Pre-calculate sizes to avoid repeated allocations
	builder.Grow(2000) // Pre-allocate for typical dashboard size
	
	// Header
	builder.WriteString(components.Header(d.Version, "dashboard"))
	builder.WriteString("\n")
	
	// Widgets using cached data
	widgets := d.renderWidgetsOptimized()
	builder.WriteString(widgets)
	builder.WriteString("\n")
	
	// Status bar
	builder.WriteString(d.StatusBar.Render())
	
	return builder.String()
}

// renderWidgets renders all widgets in a grid (optimized version)
func (d *DashboardScreen) renderWidgetsOptimized() string {
	var builder strings.Builder
	
	// Get visible widgets based on priority and data availability
	visibleWidgets := d.getVisibleWidgets()
	
	// Use cached layout data
	layout := d.widgetDataCache.GridLayout
	widgetWidth := layout.WidgetWidth
	widgetHeight := layout.WidgetHeight
	
	// Track render start time for performance measurement
	renderStart := time.Now()
	
	// Render only visible widgets in priority order
	for _, widgetName := range visibleWidgets {
		if d.shouldSkipWidget(widgetName) {
			continue // Skip this widget based on performance metrics
		}
		
		widgetData := d.getWidgetDataByName(widgetName)
		if widgetData == nil {
			continue // Skip if widget data is not available
		}
		
		widgetContent := d.renderWidgetOptimized(widgetName, widgetData, widgetWidth, widgetHeight)
		builder.WriteString(widgetContent)
		
		// Update performance stats
		renderTime := time.Since(renderStart)
		d.updateWidgetRenderStats(widgetName, renderTime)
	}
	
	// Update overall cache render time
	d.widgetDataCache.LastRenderTime = time.Since(renderStart)
	
	return builder.String()
}

// getWidgetDataByName retrieves widget data by name from cache
func (d *DashboardScreen) getWidgetDataByName(name string) interface{} {
	switch name {
	case "Today":
		return d.widgetDataCache.TodayData
	case "Authorship Score":
		return d.widgetDataCache.AuthorshipData
	case "Active Agents":
		return d.widgetDataCache.AgentsData
	case "Recent Decisions":
		return d.widgetDataCache.DecisionsData
	case "Cross-repo":
		return d.widgetDataCache.CrossRepoData
	default:
		return nil
	}
}

// renderWidgets renders all widgets in a grid (backward compatibility)
func (d *DashboardScreen) renderWidgets(width, height int) string {
	// Grid layout: 2x3 for 5 widgets
	widgetWidth := width / 4
	widgetHeight := height / 3
	
	// Today widget (top-left)
	today := d.renderWidget("Today", d.TodayWidget, widgetWidth, widgetHeight)
	
	// Authorship widget (top-middle)  
	authorship := d.renderWidget("Authorship Score", d.AuthorshipWidget, widgetWidth, widgetHeight)
	
	// Agents widget (top-right)
	agents := d.renderWidget("Active Agents", d.AgentsWidget, widgetWidth, widgetHeight)
	
	// Decisions widget (bottom-left, spans 2 columns)
	decisions := d.renderWidget("Recent Decisions", d.DecisionsWidget, widgetWidth*2-1, widgetHeight)
	
	// Cross-repo widget (bottom-right)
	crossrepo := d.renderWidget("Cross-repo", d.CrossRepoWidget, widgetWidth, widgetHeight)
	
	// Combine widgets
	return today + authorship + agents + "\n" + decisions + crossrepo
}

// renderWidget renders a single widget
func (d *DashboardScreen) renderWidget(title string, widget interface{}, width, height int) string {
	theme := d.Theme
	
	// Add focus indicator if this widget is focused
	focusPrefix := "  "
	if d.CurrentWidget >= 0 && d.CurrentWidget < len(widgetTitleOrder) && widgetTitleOrder[d.CurrentWidget] == title {
		focusPrefix = "▸ "
	}
	
	border := theme.Border.Render("┌" + strings.Repeat("─", width-2) + "┐")
	
	// Render widget content
	var content string
	switch w := widget.(type) {
	case *TodayWidget:
		content = d.renderTodayWidget(w, width-2, height-2)
	case *AuthorshipWidget:
		content = d.renderAuthorshipWidget(w, width-2, height-2)
	case *AgentsWidget:
		content = d.renderAgentsWidget(w, width-2, height-2)
	case *DecisionsWidget:
		content = d.renderDecisionsWidget(w, width-2, height-2)
	case *CrossRepoWidget:
		content = d.renderCrossRepoWidget(w, width-2, height-2)
	}
	
	footer := theme.Border.Render("└" + strings.Repeat("─", width-2) + "┘")
	
	return focusPrefix + title + "\n" + border + "\n" + content + footer + "\n"
}

// renderWidgetOptimized renders a single widget using cached data
func (d *DashboardScreen) renderWidgetOptimized(title string, data interface{}, width, height int) string {
	var builder strings.Builder
	builder.Grow(width * height) // Pre-allocate for widget content
	
	// Add focus indicator if this widget is focused
	focusPrefix := "  "
	if d.CurrentWidget >= 0 && d.CurrentWidget < len(widgetTitleOrder) && widgetTitleOrder[d.CurrentWidget] == title {
		focusPrefix = "▸ "
	}
	
	// Use cached border strings
	border := d.widgetDataCache.BorderStrings.Top
	
	// Render widget content using cached data
	switch v := data.(type) {
	case TodayWidgetData:
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		builder.WriteString(fmt.Sprintf("  %d cards due\n", v.CardsDue))
		builder.WriteString(fmt.Sprintf("  %d pending decisions\n", v.PendingDecisions))
		builder.WriteString(fmt.Sprintf("  %d new blindspots\n", v.NewBlindspots))
		for i := 3; i < height-1; i++ {
			builder.WriteString("  \n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
		
	case AuthorshipWidgetData:
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		builder.WriteString(fmt.Sprintf("  ACS: %.2f\n", v.ACS))
		builder.WriteString(fmt.Sprintf("  Recall: %.2f\n", v.Recall))
		builder.WriteString(fmt.Sprintf("  Coverage: %.2f\n", v.Coverage))
		builder.WriteString(fmt.Sprintf("  Trend: %s\n", v.Trend))
		for i := 4; i < height-1; i++ {
			builder.WriteString("  \n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
		
	case AgentsWidgetData:
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		for _, agent := range v.Agents {
			builder.WriteString(fmt.Sprintf("  %s (%s)\n", agent.Name, agent.Health))
		}
		for i := len(v.Agents); i < height-1; i++ {
			builder.WriteString("  \n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
		
	case DecisionsWidgetData:
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		for _, decision := range v.Decisions {
			builder.WriteString(fmt.Sprintf("  %s\n", decision.Summary))
		}
		for i := len(v.Decisions); i < height-1; i++ {
			builder.WriteString("  \n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
		
	case CrossRepoWidgetData:
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		builder.WriteString(fmt.Sprintf("  Enabled: %t\n", v.Enabled))
		if v.Enabled {
			builder.WriteString(fmt.Sprintf("  Repos: %d\n", len(v.Repos)))
			builder.WriteString(fmt.Sprintf("  Last sync: %s\n", v.LastSync.Format("15:04")))
			builder.WriteString(fmt.Sprintf("  Stale: %t\n", v.Stale))
		} else {
			builder.WriteString("  Not enabled\n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
		
	default:
		// Fallback for unknown data types
		builder.WriteString(focusPrefix + title + "\n")
		builder.WriteString(border + "\n")
		for i := 0; i < height-1; i++ {
			builder.WriteString("  \n")
		}
		builder.WriteString(d.widgetDataCache.BorderStrings.Bottom + "\n")
	}
	
	return builder.String()
}

// renderTodayWidget renders the today widget
func (d *DashboardScreen) renderTodayWidget(w *TodayWidget, width, height int) string {
	lines := []string{
		fmt.Sprintf("  %d cards due", w.CardsDue),
		fmt.Sprintf("  %d pending decisions", w.PendingDecisions),
		fmt.Sprintf("  %d new blindspots", w.NewBlindspots),
	}
	
	for len(lines) < height-1 {
		lines = append(lines, "  ")
	}
	
	return strings.Join(lines, "\n")
}

// renderAuthorshipWidget renders the authorship widget
func (d *DashboardScreen) renderAuthorshipWidget(w *AuthorshipWidget, width, height int) string {
	lines := []string{}
	
	if w.ACS > 0 {
		lines = append(lines, fmt.Sprintf("  ACS %.1f / target 3.5", w.ACS))
	} else {
		lines = append(lines, "  ACS not tracked yet")
	}
	
	if w.Recall > 0 {
		lines = append(lines, fmt.Sprintf("  Recall  %.0f%%", w.Recall*100))
	} else {
		lines = append(lines, "  Recall not tracked")
	}
	
	if w.Coverage > 0 {
		lines = append(lines, fmt.Sprintf("  Coverage %.0f%%", w.Coverage*100))
	} else {
		lines = append(lines, "  Coverage not tracked")
	}
	
	if w.Trend != "" {
		lines = append(lines, fmt.Sprintf("  30-day: %s", w.Trend))
	}
	
	for len(lines) < height-1 {
		lines = append(lines, "  ")
	}
	
	return strings.Join(lines, "\n")
}

// renderAgentsWidget renders the agents widget
func (d *DashboardScreen) renderAgentsWidget(w *AgentsWidget, width, height int) string {
	lines := []string{}
	for _, agent := range w.Agents {
		status := "✓"
		if agent.Health != "healthy" {
			status = "⚠"
		}
		lines = append(lines, fmt.Sprintf("  %s %s (%s)", status, agent.Name, agent.Health))
	}
	
	for len(lines) < height-1 {
		lines = append(lines, "  ")
	}
	
	return strings.Join(lines, "\n")
}

// renderDecisionsWidget renders the decisions widget
func (d *DashboardScreen) renderDecisionsWidget(w *DecisionsWidget, width, height int) string {
	lines := []string{}
	for _, decision := range w.Decisions {
		lines = append(lines, fmt.Sprintf("  • %s", decision.Title))
	}
	
	for len(lines) < height-1 {
		lines = append(lines, "  ")
	}
	
	return strings.Join(lines, "\n")
}

// renderCrossRepoWidget renders the cross-repo widget
func (d *DashboardScreen) renderCrossRepoWidget(w *CrossRepoWidget, width, height int) string {
	lines := []string{}
	
	if w.Enabled {
		status := "✓"
		if w.Stale {
			status = "⚠"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", status, "OryaAgent ← Sec"))
		lines = append(lines, fmt.Sprintf("    stale decision"))
	} else {
		lines = append(lines, "  ◇ No cross-repo sync")
	}
	
	for len(lines) < height-1 {
		lines = append(lines, "  ")
	}
	
	return strings.Join(lines, "\n")
}

// helpSections — все живые клавиши дашборда для help-оверлея (§18.4).
// Перечисляются только реально обрабатываемые клавиши — help не врёт (INV-UI-009).
// Screen-specific-категории нет: у дашборда нет собственных клавиш вне карты
// действий — выдумывать bindings ради симметрии нельзя (§18.5).
func (d *DashboardScreen) helpSections() []components.HelpSection {
	return []components.HelpSection{
		{Title: "Navigation", Entries: []components.HelpEntry{
			{Key: "Tab", Desc: "next widget"},
			{Key: "Shift+Tab", Desc: "previous widget"},
			{Key: "Esc", Desc: "quit (root screen)"},
		}},
		{Title: "Actions", Entries: []components.HelpEntry{
			{Key: "r", Desc: "refresh data"},
			{Key: "q", Desc: "quit"},
			{Key: "?", Desc: "toggle this help"},
		}},
	}
}

// GetID returns the screen ID
func (d *DashboardScreen) GetID() string {
	return "dashboard"
}

// IsActive returns if the screen is active
func (d *DashboardScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (d *DashboardScreen) SetActive(active bool) {
	// Update status bar focus when active
	if active {
		d.StatusBar.SetFocus(d.WidgetIDs[d.CurrentWidget])
	}
}

// SetKeyKeys sets keyboard hints (for compatibility)
func (d *DashboardScreen) SetKeyKeys(hints []string) {
	d.StatusBar.SetKeyHints(hints)
}

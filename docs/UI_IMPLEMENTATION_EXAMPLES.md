# UI Self-Review Examples & Templates

**Companion to UI_SELF_REVIEW_CHECKLIST.md**

*Version: 0.4.1-mvp.1*  
*Scope: Practical examples for UI implementation*

---

## 📋 Screen Implementation Template

### Template: New Screen Implementation

```go
// pkg/ui/screens/example_screen.go
// §<N> <Title> (caf <command>)
//
// Purpose: <one sentence describing the screen's purpose>
// Layout: <ASCII mockup reference>
// Components: <list of components used>
// Keyboard map: <table of keyboard shortcuts>
// Performance budget: <table of performance targets>
// Edge cases: <list of edge cases handled>
// Accessibility: <accessibility considerations>

package screens

import (
    "time"
    
    "github.com/charmbracelet/bubbletea"
    "github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
    "github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
    "github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui/components"
)

// ExampleScreen represents the example screen
type ExampleScreen struct {
    ui.StateMachine
    
    Version    string
    RepoRoot   string
    Theme      *ui.Theme
    Keys       *ui.Keys
    A11y       *ui.A11y
    
    // Data
    Data       []DataType
    Loading    bool
    LoadError  error
    
    // Focus and navigation
    CurrentItem int
    ShowHelp    bool
    
    // Performance optimization
    lastDataHash string
    widgetCache  *WidgetDataCache
}

// WidgetDataCache for performance optimization (INV-UI-004, INV-UI-006)
type WidgetDataCache struct {
    RenderTime    time.Time
    DataHash      string
    LayoutCache   string
}

// Init implements the tea.Model interface
func (s *ExampleScreen) Init() tea.Cmd {
    // Load skeleton immediately, then real data (INV-UI-004 SWR pattern)
    return tea.Batch(
        s.loadData(),
        components.ShowSkeleton(),
    )
}

// loadData fetches data with skeleton pattern
func (s *ExampleScreen) loadData() tea.Cmd {
    return func() tea.Msg {
        // Simulate data loading
        time.Sleep(100 * time.Millisecond)
        data := fetchExampleData()
        return DataLoadedMsg{Data: data}
    }
}

// DataLoadedMsg implements SWR pattern (INV-UI-004)
type DataLoadedMsg struct {
    Data []DataType
}

// Update implements the tea.Model interface
func (s *ExampleScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrl+C, tea.KeyEsc:
            return s, tea.Quit
        case tea.KeyQuestionMark:
            s.ShowHelp = !s.ShowHelp
            return s, nil
        case tea.KeyTab, tea.KeyShiftTab:
            // Widget navigation (INV-UI-009 keyboard-only)
            s.updateFocus(msg.Type)
            return s, nil
        case tea.KeyEnter:
            // Action handling
            return s.handleAction()
        }
    case DataLoadedMsg:
        // SWR pattern: replace skeleton with real data
        s.Data = msg.Data
        s.Loading = false
        s.updateCache()
    }
    
    return s, cmd
}

// View implements the tea.Model interface
func (s *ExampleScreen) View() string {
    // High-contrast mode support (INV-UI-008)
    theme := s.Theme
    if theme.IsHighContrast {
        theme = ui.HighContrastTheme()
    }
    
    var builder strings.Builder
    
    // Header
    builder.WriteString(theme.Header.Render("# CodeAuthorityFirst — Example Screen"))
    
    // Main content with skeleton pattern (INV-UI-004)
    if s.Loading {
        return builder.WriteString(components.SkeletonView())
    }
    
    // Focus visibility marker (INV-PROMPT-MVP-007)
    focusMarker := " "
    if s.hasFocus() {
        focusMarker = "▸"
    }
    
    // Content with accessibility support
    builder.WriteString(theme.Section.Render("## Data"))
    builder.WriteString(s.renderDataList(focusMarker))
    
    // Help overlay (INV-UI-010)
    if s.ShowHelp {
        builder.WriteString(s.renderHelp())
    }
    
    return builder.String()
}

// renderDataList renders data with focus visibility
func (s *ExampleScreen) renderDataList(focusMarker string) string {
    theme := s.Theme
    
    var builder strings.Builder
    for i, item := range s.Data {
        hasFocus := i == s.CurrentItem
        prefix := "  "
        if hasFocus {
            prefix = focusMarker + " "
        }
        
        // Color-blind safe: color + symbol + text (INV-PROMPT-MVP-008)
        statusSymbol := "•"
        statusColor := theme.Muted
        if item.IsValid {
            statusSymbol = "✓"
            statusColor = theme.Secondary
        }
        
        line := fmt.Sprintf("%s%s %s %s", prefix, statusSymbol, item.Name, item.Description)
        builder.WriteString(theme.Body.Render(line))
        builder.WriteString("\n")
    }
    
    return builder.String()
}

// renderHelp shows keyboard shortcuts
func (s *ExampleScreen) renderHelp() string {
    theme := s.Theme
    
    help := []string{
        "Keyboard Shortcuts:",
        "  Tab/Shift+Tab  - Navigate items",
        "  Enter          - Select item",
        "  ?              - Toggle help",
        "  q/Ctrl+C       - Quit",
        "  Esc            - Back",
    }
    
    var builder strings.Builder
    for _, line := range help {
        builder.WriteString(theme.Muted.Render(line))
        builder.WriteString("\n")
    }
    
    return builder.String()
}

// Performance optimization methods
func (s *ExampleScreen) updateCache() {
    if s.widgetCache == nil {
        s.widgetCache = &WidgetDataCache{}
    }
    
    // Cache data hash for change detection (INV-UI-006)
    hash := s.computeDataHash()
    s.widgetCache.DataHash = hash
    s.widgetCache.RenderTime = time.Now()
}

func (s *ExampleScreen) computeDataHash() string {
    // Simple hash for performance optimization
    return fmt.Sprintf("%d", len(s.Data))
}

// Focus visibility (INV-PROMPT-MVP-007)
func (s *ExampleScreen) hasFocus() bool {
    return s.CurrentItem >= 0 && s.CurrentItem < len(s.Data)
}

// Widget navigation
func (s *ExampleScreen) updateFocus(keyType tea.KeyType) {
    if len(s.Data) == 0 {
        return
    }
    
    if keyType == tea.KeyTab {
        s.CurrentItem = (s.CurrentItem + 1) % len(s.Data)
    } else {
        s.CurrentItem = (s.CurrentItem - 1 + len(s.Data)) % len(s.Data)
    }
}

// Action handling
func (s *ExampleScreen) handleAction() (tea.Model, tea.Cmd) {
    if !s.hasFocus() {
        return s, nil
    }
    
    selected := s.Data[s.CurrentItem]
    
    // Handle action with error recovery (INV-UI-010)
    err := performAction(selected)
    if err != nil {
        return s, s.showError(err)
    }
    
    return s, nil
}

// Error handling with recovery options
func (s *ExampleScreen) showError(err error) tea.Cmd {
    return tea.Sequence(
        ui.ShowError(err.Error()),
        components.WaitForRetryOrQuit(),
    )
}

// Helper functions
func fetchExampleData() []DataType {
    // Simulate data fetching
    return []DataType{
        {Name: "Item 1", Description: "First item", IsValid: true},
        {Name: "Item 2", Description: "Second item", IsValid: false},
        {Name: "Item 3", Description: "Third item", IsValid: true},
    }
}

func performAction(item DataType) error {
    // Simulate action
    if !item.IsValid {
        return fmt.Errorf("cannot perform action on invalid item")
    }
    return nil
}

// NewExampleScreen creates a new example screen
func NewExampleScreen(repoRoot, version string) *ExampleScreen {
    return &ExampleScreen{
        Version:  version,
        RepoRoot: repoRoot,
        Theme:    ui.DefaultTheme(),
        Keys:     ui.DefaultKeys(),
        A11y:     ui.DefaultA11y(),
    }
}
```

---

## 🎨 Performance Optimization Examples

### Example 1: Widget Data Caching

```go
// Performance optimization: Widget data caching (INV-UI-004, INV-UI-006)
type WidgetDataCache struct {
    Data       interface{}
    Hash       string
    LastUpdate time.Time
    TTL        time.Duration
}

func (w *WidgetDataCache) IsExpired() bool {
    return time.Since(w.LastUpdate) > w.TTL
}

func (w *WidgetDataCache) Update(data interface{}, hash string) {
    w.Data = data
    w.Hash = hash
    w.LastUpdate = time.Now()
}

func (w *WidgetDataCache) ShouldRecompute(newHash string) bool {
    return w.Hash != newHash || w.IsExpired()
}
```

### Example 2: String Building Optimization

```go
// Performance optimization: Efficient string building (INV-UI-002)
func (s *DashboardScreen) renderView() string {
    var builder strings.Builder
    
    // Pre-compute border strings (INV-UI-006)
    border := s.Theme.Border.Render("─")
    
    // Build view efficiently
    builder.WriteString(s.renderHeader())
    builder.WriteString(border + "\n")
    builder.WriteString(s.renderContent())
    builder.WriteString(border + "\n")
    builder.WriteString(s.renderFooter())
    
    return builder.String()
}
```

### Example 3: Virtualized List Rendering

```go
// Performance optimization: Virtualized list (INV-UI-002)
func (s *ListScreen) renderVirtualizedList() string {
    const visibleItems = 10
    const itemHeight = 3
    
    totalItems := len(s.Items)
    startIndex := max(0, s.CurrentIndex-visibleItems/2)
    endIndex := min(totalItems, startIndex+visibleItems)
    
    var builder strings.Builder
    
    for i := startIndex; i < endIndex; i++ {
        item := s.Items[i]
        focused := i == s.CurrentIndex
        
        // Render only visible items
        line := s.renderListItem(item, focused)
        builder.WriteString(line)
        builder.WriteString("\n")
    }
    
    return builder.String()
}
```

---

## ♿ Accessibility Examples

### Example 1: Focus Visibility

```go
// Focus visibility implementation (INV-PROMPT-MVP-007)
func (s *Screen) renderFocusedItem(item Item) string {
    theme := s.Theme
    
    // Focus marker
    marker := "  "
    if s.hasFocus(item) {
        marker = "▸ "
    }
    
    // Color-blind safe: color + symbol + text
    status := "• " + theme.Muted.Render(item.Name)
    if item.IsActive {
        status = "✓ " + theme.Secondary.Render(item.Name)
    }
    
    return marker + status
}
```

### Example 2: Screen Reader Support

```go
// Screen reader support (Accessibility)
func (s *Screen) renderWithScreenReader() string {
    var builder strings.Builder
    
    // Status-first announcement
    if s.HasError {
        builder.WriteString("ERROR: ")
        builder.WriteString(s.ErrorMessage)
        builder.WriteString("\n")
    }
    
    // Main content
    builder.WriteString(s.renderMainContent())
    
    // Alt-text for visual elements
    if s.HasChart {
        builder.WriteString("Chart: ")
        builder.WriteString(s.getChartAltText())
        builder.WriteString("\n")
    }
    
    return builder.String()
}
```

### Example 3: High-Contrast Mode

```go
// High-contrast mode support (INV-UI-008)
func (s *Screen) renderHighContrast() string {
    theme := s.Theme
    
    if theme.IsHighContrast {
        theme = ui.HighContrastTheme()
    }
    
    return theme.Render(s.getContent())
}

// High contrast theme
func HighContrastTheme() *ui.Theme {
    return &ui.Theme{
        Primary:    lipgloss.Color("#FFFFFF"), // White on black
        Secondary:  lipgloss.Color("#00FF00"), // Bright green
        Warning:    lipgloss.Color("#FFFF00"), // Bright yellow
        Error:      lipgloss.Color("#FF0000"), // Bright red
        Muted:      lipgloss.Color("#CCCCCC"), // Light gray
        Background: lipgloss.Color("#000000"), // Black
    }
}
```

---

## ⌨️ Keyboard Navigation Examples

### Example 1: Universal Keyboard Map

```go
// Universal keyboard map implementation (INV-UI-009)
func (s *Screen) handleKey(msg tea.KeyMsg) tea.Cmd {
    switch msg.Type {
    // Global keys
    case tea.KeyQuestionMark:
        return s.toggleHelp()
    case tea.KeyCtrl+C, tea.KeyCtrl+D:
        return tea.Quit
    case tea.KeyEsc:
        return s.goBack()
    
    // Navigation keys
    case tea.KeyTab, tea.KeyShiftTab:
        return s.cycleWidgets(msg.Type)
    case tea.KeyUp, tea.KeyK:
        return s.moveUp()
    case tea.KeyDown, tea.KeyJ:
        return s.moveDown()
    case tea.KeyHome, tea.KeyG:
        return s.moveToTop()
    case tea.KeyEnd, tea.KeyShift+G:
        return s.moveToBottom()
    
    // Action keys
    case tea.KeyEnter:
        return s.activate()
    case tea.KeySpace:
        return s.toggle()
    case tea.KeyR:
        return s.refresh()
    
    // Global shortcuts
    case tea.KeyD:
        return s.openDoctor()
    case tea.KeyM:
        return s.openMetrics()
    case tea.KeyDollar:
        return s.openCost()
    }
    
    return nil
}
```

### Example 2: Widget Cycling

```go
// Widget cycling implementation
func (s *Screen) cycleWidgets(direction tea.KeyType) tea.Cmd {
    if len(s.Widgets) == 0 {
        return nil
    }
    
    if direction == tea.KeyTab {
        s.CurrentWidget = (s.CurrentWidget + 1) % len(s.Widgets)
    } else {
        s.CurrentWidget = (s.CurrentWidget - 1 + len(s.Widgets)) % len(s.Widgets)
    }
    
    return nil
}
```

---

## 🧪 Testing Examples

### Example 1: Performance Test

```go
// Performance test example (INV-UI-001, INV-UI-002)
func TestScreenPerformance(t *testing.T) {
    // Create test screen
    screen := NewTestScreen()
    
    // Test first paint time (INV-UI-001)
    t.Run("first paint <100ms", func(t *testing.T) {
        start := time.Now()
        view := screen.View()
        duration := time.Since(start)
        
        if duration >= 100*time.Millisecond {
            t.Errorf("First paint took %v, expected <100ms", duration)
        }
        
        // Verify view is not empty
        if view == "" {
            t.Error("View should not be empty")
        }
    })
    
    // Test render performance (INV-UI-002)
    t.Run("render <16ms", func(t *testing.T) {
        start := time.Now()
        for i := 0; i < 100; i++ {
            screen.View()
        }
        avgRender := time.Since(start) / 100
        
        if avgRender >= 16*time.Millisecond {
            t.Errorf("Average render time %v, expected <16ms", avgRender)
        }
    })
}
```

### Example 2: Accessibility Test

```go
// Accessibility test example
func TestScreenAccessibility(t *testing.T) {
    screen := NewTestScreen()
    
    // Test focus visibility (INV-PROMPT-MVP-007)
    t.Run("focus visibility", func(t *testing.T) {
        screen.SetFocus(true)
        view := screen.View()
        
        if !strings.Contains(view, "▸") {
            t.Error("Focused item should have marker")
        }
    })
    
    // Test symbol-text pairing (INV-PROMPT-MVP-008)
    t.Run("symbol-text pairing", func(t *testing.T) {
        view := screen.View()
        
        // Check that status symbols have text
        for _, line := range strings.Split(view, "\n") {
            if strings.ContainsAny(line, "✓⚠✗") {
                if !strings.ContainsFunc(line, func(r rune) bool {
                    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
                }) {
                    t.Errorf("Status symbol without text: %q", line)
                }
            }
        }
    })
    
    // Test help overlay (INV-UI-010)
    t.Run("help overlay", func(t *testing.T) {
        screen.ShowHelp = true
        view := screen.View()
        
        if !strings.Contains(view, "Keyboard Shortcuts") {
            t.Error("Help overlay should show keyboard shortcuts")
        }
    })
}
```

### Example 3: Integration Test

```go
// Integration test example
func TestScreenIntegration(t *testing.T) {
    // Create temporary repository
    tempDir := t.TempDir()
    createTestRepo(tempDir)
    
    // Create screen with test data
    screen := NewExampleScreen(tempDir, "0.4.1-test")
    
    // Test full flow
    t.Run("full screen flow", func(t *testing.T) {
        // Initialize screen
        cmd := screen.Init()
        screen, cmd = screen.Update(cmd())
        
        // Test data loading
        if cmd != nil {
            screen, cmd = screen.Update(cmd())
        }
        
        // Test view rendering
        view := screen.View()
        
        // Verify screen state
        if screen.Loading {
            t.Error("Screen should not be loading after initialization")
        }
        
        if view == "" {
            t.Error("View should not be empty")
        }
        
        // Test keyboard navigation
        screen, cmd = screen.Update(tea.KeyMsg{Type: tea.KeyTab})
        if screen.CurrentWidget == 0 {
            t.Error("Widget should have changed after Tab")
        }
    })
}
```

---

## 📊 Performance Budget Examples

### Example 1: Dashboard Performance

```go
// Dashboard performance budget (INV-UI-007)
const (
    DashboardFirstPaint     = 100 * time.Millisecond
    DashboardFrameRender    = 16 * time.Millisecond
    DashboardWidgetRender   = 50 * time.Millisecond
    DashboardMemoryUsage    = 50 * 1024 * 1024 // 50MB
)

func TestDashboardPerformance(t *testing.T) {
    dashboard := NewDashboardScreen("/tmp/test", "0.4.1-test")
    
    // Test first paint (INV-UI-001)
    t.Run("first paint", func(t *testing.T) {
        start := time.Now()
        view := dashboard.View()
        duration := time.Since(start)
        
        if duration >= DashboardFirstPaint {
            t.Errorf("First paint %v >= %v", duration, DashboardFirstPaint)
        }
        
        if view == "" {
            t.Error("View should not be empty")
        }
    })
    
    // Test widget rendering (INV-UI-002)
    t.Run("widget render", func(t *testing.T) {
        start := time.Now()
        for i := 0; i < 10; i++ {
            dashboard.renderWidgets()
        }
        avgRender := time.Since(start) / 10
        
        if avgRender >= DashboardWidgetRender {
            t.Errorf("Widget render %v >= %v", avgRender, DashboardWidgetRender)
        }
    })
}
```

### Example 2: Memory Usage

```go
// Memory usage monitoring
func TestMemoryUsage(t *testing.T) {
    var m1, m2 runtime.MemStats
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    screen := NewTestScreen()
    screen.renderLargeView()
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    usedMemory := m2.Alloc - m1.Alloc
    
    if usedMemory > DashboardMemoryUsage {
        t.Errorf("Memory usage %v > %v", usedMemory, DashboardMemoryUsage)
    }
}
```

---

## 🎯 Quick Reference

### Performance Budgets
- **First paint**: <100ms (INV-UI-001)
- **Frame render**: <16ms (INV-UI-002)
- **Keystroke response**: <50ms (INV-UI-003)
- **Widget render**: <50ms
- **Memory usage**: <50MB

### Accessibility Requirements
- **Focus visibility**: Always visible (INV-PROMPT-MVP-007)
- **Symbol-text pairing**: Status symbols must have text (INV-PROMPT-MVP-008)
- **Keyboard-only**: 100% coverage (INV-UI-009)
- **Help overlay**: Available on all screens (INV-UI-010)

### Keyboard Shortcuts
- **Global**: ? (help), q (quit), Esc (back)
- **Navigation**: Tab/Shift+Tab, j/k/↑/↓, g/G (top/bottom)
- **Actions**: Enter (activate), r (refresh), 1-9 (switch)
- **Global shortcuts**: D (doctor), M (metrics), $ (cost)

### UI Patterns
- **Skeleton/SWR**: Immediate skeleton + data replacement
- **Virtualization**: Only render visible items
- **Caching**: Pre-compute styles and data
- **Progressive enhancement**: 16-color → 256-color → truecolor

---

*These examples provide practical implementation guidance for CodeAuthorityFirst UI requirements.*
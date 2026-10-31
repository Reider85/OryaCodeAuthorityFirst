//go:build integration

package ui

import (
	"testing"
	"time"
)

// TestTheme tests the theme functionality
func TestTheme(t *testing.T) {
	theme := NewTheme()
	
	// Test theme creation
	if theme == nil {
		t.Fatal("Theme should not be nil")
	}
	
	// Test color styles
	if theme.Primary.GetForeground() == "" {
		t.Error("Primary color should not be empty")
	}
	
	if theme.Secondary.GetForeground() == "" {
		t.Error("Secondary color should not be empty")
	}
	
	if theme.Error.GetForeground() == "" {
		t.Error("Error color should not be empty")
	}
	
	// Test high contrast theme
	highContrastTheme := theme.HighContrast()
	if highContrastTheme == nil {
		t.Fatal("High contrast theme should not be nil")
	}
	
	// Test precompute styles
	theme.PrecomputeStyles()
	
	// Test color getter
	color := theme.GetColor("primary")
	if color == nil {
		t.Error("GetColor should return a valid style")
	}
}

// TestKeys tests the keyboard mapping functionality
func TestKeys(t *testing.T) {
	keys := NewKeys()
	
	// Test keys creation
	if keys == nil {
		t.Fatal("Keys should not be nil")
	}
	
	// Test navigation keys
	if keys.Tab == "" {
		t.Error("Tab key should not be empty")
	}
	
	if keys.Enter == "" {
		t.Error("Enter key should not be empty")
	}
	
	if keys.Esc == "" {
		t.Error("Esc key should not be empty")
	}
	
	// Test action keys
	if keys.Refresh == "" {
		t.Error("Refresh key should not be empty")
	}
	
	if keys.Quit == "" {
		t.Error("Quit key should not be empty")
	}
	
	// Test global shortcuts
	if keys.Dashboard == "" {
		t.Error("Dashboard key should not be empty")
	}
	
	if keys.Doctor == "" {
		t.Error("Doctor key should not be empty")
	}
	
	// Test key type checking
	if !keys.IsNavigationKey(keys.Up) {
		t.Error("Up should be a navigation key")
	}
	
	if keys.IsNavigationKey(keys.Refresh) {
		t.Error("Refresh should not be a navigation key")
	}
	
	if !keys.IsActionKey(keys.Enter) {
		t.Error("Enter should be an action key")
	}
	
	if keys.IsActionKey(keys.Up) {
		t.Error("Up should not be an action key")
	}
	
	if !keys.IsGlobalShortcut(keys.Dashboard) {
		t.Error("Dashboard should be a global shortcut")
	}
	
	if keys.IsGlobalShortcut(keys.Up) {
		t.Error("Up should not be a global shortcut")
	}
}

// TestPerformanceBudget tests the performance budget functionality
func TestPerformanceBudget(t *testing.T) {
	budget := NewPerformanceBudget()
	
	// Test budget creation
	if budget == nil {
		t.Fatal("Performance budget should not be nil")
	}
	
	// Test budget values
	if budget.FirstPaint <= 0 {
		t.Error("First paint budget should be positive")
	}
	
	if budget.FrameRender <= 0 {
		t.Error("Frame render budget should be positive")
	}
	
	if budget.Keystroke <= 0 {
		t.Error("Keystroke budget should be positive")
	}
}

// TestPerfTracker tests the performance tracker functionality
func TestPerfTracker(t *testing.T) {
	budget := NewPerformanceBudget()
	tracker := NewPerfTracker(budget)
	
	// Test tracker creation
	if tracker == nil {
		t.Fatal("Performance tracker should not be nil")
	}
	
	// Test timing
	tracker.Start("test_operation")
	time.Sleep(10 * time.Millisecond) // Simulate work
	duration := tracker.End("test_operation")
	
	if duration <= 0 {
		t.Error("Duration should be positive")
	}
	
	// Test budget checking
	withinBudget, actual, budget := tracker.CheckBudget("test_operation")
	if withinBudget {
		t.Log("Operation within budget")
	} else {
		t.Logf("Operation exceeded budget: %s > %s", actual, budget)
	}
}

// TestStateMachine tests the state machine functionality
func TestStateMachine(t *testing.T) {
	sm := NewStateMachine()
	
	// Test state machine creation
	if sm == nil {
		t.Fatal("State machine should not be nil")
	}
	
	// Test initial state
	if sm.Current() != StateLoading {
		t.Error("Initial state should be loading")
	}
	
	// Test state transitions
	sm.Transition(DataLoadedMsg{Empty: false})
	if sm.Current() != StateSuccess {
		t.Error("Should transition to success on data loaded")
	}
	
	sm.Transition(DataErrorMsg{Err: fmt.Errorf("test error")})
	if sm.Current() != StateError {
		t.Error("Should transition to error on data error")
	}
	
	sm.Transition(PartialDataMsg{Missing: []string{"file1"}})
	if sm.Current() != StatePartial {
		t.Error("Should transition to partial on partial data")
	}
	
	// Test error state
	if sm.Error() == "" {
		t.Error("Error state should have an error message")
	}
}

// TestAccessibilityConfig tests the accessibility configuration
func TestAccessibilityConfig(t *testing.T) {
	config := NewAccessibilityConfig()
	
	// Test config creation
	if config == nil {
		t.Fatal("Accessibility config should not be nil")
	}
	
	// Test default values
	if !config.KeyboardOnly {
		t.Error("Keyboard-only should be enabled by default")
	}
	
	if !config.TabNavigation {
		t.Error("Tab navigation should be enabled by default")
	}
	
	if !config.EscapeExits {
		t.Error("Escape exits should be enabled by default")
	}
	
	if !config.ConfirmActions {
		t.Error("Confirm actions should be enabled by default")
	}
}

// TestA11y tests the accessibility functionality
func TestA11y(t *testing.T) {
	config := NewAccessibilityConfig()
	theme := NewTheme()
	keys := NewKeys()
	a11y := NewA11y(config, theme, keys)
	
	// Test a11y creation
	if a11y == nil {
		t.Fatal("Accessibility should not be nil")
	}
	
	// Test theme getter
	testTheme := a11y.GetTheme()
	if testTheme == nil {
		t.Error("GetTheme should return a valid theme")
	}
	
	// Test announcement text
	announcement := a11y.GetAnnouncementText("test_action", "test_description")
	if announcement == "" {
		t.Error("GetAnnouncementText should return text when screen reader is enabled")
	}
	
	// Test animation check
	shouldAnimate := a11y.ShouldAnimate()
	if shouldAnimate != config.ReducedMotion {
		t.Error("ShouldAnimate should respect reduced motion setting")
	}
	
	// Test keyboard shortcuts
	shortcuts := a11y.GetKeyboardShortcuts()
	if shortcuts == "" {
		t.Error("GetKeyboardShortcuts should return help text")
	}
	
	// Test contrast validation (simplified)
	valid := a11y.ValidateColorContrast("#ffffff", "#000000")
	if !valid {
		t.Error("White on black should have good contrast")
	}
	
	// Test validation report
	report := a11y.GetAccessibilityReport()
	if report == "" {
		t.Error("GetAccessibilityReport should return a report")
	}
}

// TestFocusManager tests the focus manager functionality
func TestFocusManager(t *testing.T) {
	config := NewAccessibilityConfig()
	theme := NewTheme()
	keys := NewKeys()
	a11y := NewA11y(config, theme, keys)
	focusManager := NewFocusManager(a11y)
	
	// Test focus manager creation
	if focusManager == nil {
		t.Fatal("Focus manager should not be nil")
	}
	
	// Test adding focusable elements
	focusManager.AddFocusableElement("element1")
	focusManager.AddFocusableElement("element2")
	
	// Test focus movement
	element := focusManager.MoveForward()
	if element != "element1" {
		t.Error("MoveForward should return the first element")
	}
	
	element = focusManager.MoveForward()
	if element != "element2" {
		t.Error("MoveForward should return the second element")
	}
	
	element = focusManager.MoveForward()
	if element != "element1" {
		t.Error("MoveForward should cycle back to first element")
	}
	
	// Test focus setting
	set := focusManager.SetFocus("element2")
	if !set {
		t.Error("SetFocus should return true for existing element")
	}
	
	set = focusManager.SetFocus("nonexistent")
	if set {
		t.Error("SetFocus should return false for nonexistent element")
	}
	
	// Test focus path
	path := focusManager.GetFocusPath()
	if len(path) == 0 {
		t.Error("GetFocusPath should return a non-empty path")
	}
}

// TestScreenReaderManager tests the screen reader manager functionality
func TestScreenReaderManager(t *testing.T) {
	config := NewAccessibilityConfig()
	theme := NewTheme()
	keys := NewKeys()
	a11y := NewA11y(config, theme, keys)
	screenReader := NewScreenReaderManager(a11y)
	
	// Test screen reader manager creation
	if screenReader == nil {
		t.Fatal("Screen reader manager should not be nil")
	}
	
	// Test adding announcements
	screenReader.AddAnnouncement("test_id", "test message", AnnouncementImportant, "test context", false)
	
	// Test getting announcements
	announcements := screenReader.GetAnnouncements()
	if len(announcements) == 0 {
		t.Error("GetAnnouncements should return announcements")
	}
	
	// Test getting latest announcement
	latest := screenReader.GetLatestAnnouncement()
	if latest == nil {
		t.Error("GetLatestAnnouncement should return the latest announcement")
	}
	
	// Test clearing announcements
	screenReader.ClearAnnouncements()
	announcements = screenReader.GetAnnouncements()
	if len(announcements) != 0 {
		t.Error("ClearAnnouncements should remove all announcements")
	}
}

// TestComponents tests the UI components
func TestComponents(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	
	// Test button
	button := NewButton("test_button", "Test Button", func() tea.Cmd { return nil }, theme, keys)
	if button == nil {
		t.Fatal("Button should not be nil")
	}
	
	if button.GetID() != "test_button" {
		t.Error("Button ID should match")
	}
	
	// Test input
	input := NewInput("test_input", "Test Input", theme, keys)
	if input == nil {
		t.Fatal("Input should not be nil")
	}
	
	if input.GetID() != "test_input" {
		t.Error("Input ID should match")
	}
	
	// Test list
	list := NewList("test_list", 10, theme, keys)
	if list == nil {
		t.Fatal("List should not be nil")
	}
	
	if list.GetID() != "test_list" {
		t.Error("List ID should match")
	}
	
	// Test container
	container := NewContainer("test_container", "Test Container", theme, keys)
	if container == nil {
		t.Fatal("Container should not be nil")
	}
	
	if container.GetID() != "test_container" {
		t.Error("Container ID should match")
	}
	
	// Test text
	text := NewText("test_text", "Test Text", theme, keys)
	if text == nil {
		t.Fatal("Text should not be nil")
	}
	
	if text.GetID() != "test_text" {
		t.Error("Text ID should match")
	}
	
	// Test spinner
	spinner := NewSpinner("test_spinner", theme, keys)
	if spinner == nil {
		t.Fatal("Spinner should not be nil")
	}
	
	if spinner.GetID() != "test_spinner" {
		t.Error("Spinner ID should match")
	}
}

// TestIntegrationLayer tests the integration layer functionality
func TestIntegrationLayer(t *testing.T) {
	integration := NewIntegrationLayer()
	
	// Test integration creation
	if integration == nil {
		t.Fatal("Integration layer should not be nil")
	}
	
	// Test screen manager
	screenManager := integration.ScreenManager
	if screenManager == nil {
		t.Fatal("Screen manager should not be nil")
	}
	
	// Test data manager
	dataManager := integration.DataManager
	if dataManager == nil {
		t.Fatal("Data manager should not be nil")
	}
	
	// Test event bus
	eventBus := integration.EventBus
	if eventBus == nil {
		t.Fatal("Event bus should not be nil")
	}
	
	// Test setting and getting shared data
	dataManager.SetSharedData("test_key", "test_value")
	value, exists := dataManager.GetSharedData("test_key")
	if !exists || value != "test_value" {
		t.Error("Shared data should be set and retrieved correctly")
	}
	
	// Test caching
	dataManager.CacheData("cache_key", "cached_value")
	cachedValue, exists := dataManager.GetCachedData("cache_key")
	if !exists || cachedValue != "cached_value" {
		t.Error("Cached data should be set and retrieved correctly")
	}
	
	// Test event bus subscriptions
	channel := eventBus.Subscribe("test_event")
	eventBus.Publish(Event{Type: "test_event", Payload: "test_payload", Timestamp: time.Now(), Sender: ScreenDashboard})
	
	select {
	case event := <-channel:
		if event.Type != "test_event" {
			t.Error("Event type should match")
		}
	default:
		t.Error("Should receive published event")
	}
}

// TestTUI tests the TUI application
func TestTUI(t *testing.T) {
	tui := NewTUI()
	
	// Test TUI creation
	if tui == nil {
		t.Fatal("TUI should not be nil")
	}
	
	// Test screen manager
	screenManager := tui.GetScreenManager()
	if screenManager == nil {
		t.Fatal("Screen manager should not be nil")
	}
	
	// Test data manager
	dataManager := tui.GetDataManager()
	if dataManager == nil {
		t.Fatal("Data manager should not be nil")
	}
	
	// Test event bus
	eventBus := tui.GetEventBus()
	if eventBus == nil {
		t.Fatal("Event bus should not be nil")
	}
	
	// Test shared data operations
	tui.SetSharedData("test_key", "test_value")
	value, exists := tui.GetSharedData("test_key")
	if !exists || value != "test_value" {
		t.Error("Shared data should be set and retrieved correctly")
	}
	
	// Test event subscriptions
	channel := tui.SubscribeToEvent("test_event")
	tui.PublishEvent("test_event", "test_payload", ScreenDashboard)
	
	select {
	case event := <-channel:
		if event.Type != "test_event" {
			t.Error("Event type should match")
		}
	default:
		t.Error("Should receive published event")
	}
}
package ui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// TUI represents the complete TUI application
// FR-TUI-006: TUI integrates all screens with smooth transitions and shared state
type TUI struct {
	model       *AppModel
	program     *bubbletea.Program
	integration *IntegrationLayer
}

// IntegrationLayer manages integration between screens and shared state
// FR-TUI-007: Integration layer maintains shared state and smooth transitions
type IntegrationLayer struct {
	Theme         *Theme
	Keys          *Keys
	A11y          *A11y
	PerfTracker   *PerfTracker
	ScreenManager *ScreenManager
	DataManager   *DataManager
	EventBus      *EventBus
}

// ScreenManager manages screen transitions and state
type ScreenManager struct {
	CurrentScreen ScreenType
	Screens       map[ScreenType]Screen
	History       []ScreenType
	Transitions   []ScreenTransition
}

// ScreenTransition represents a screen transition
type ScreenTransition struct {
	From     ScreenType
	To       ScreenType
	Timestamp time.Time
	Duration time.Duration
}

// DataManager manages shared data between screens
type DataManager struct {
	SharedData    map[string]interface{}
	Cache         map[string]interface{}
	LastUpdate    time.Time
	AutoRefresh   bool
	RefreshRate   time.Duration
}

// EventBus manages events between screens
type EventBus struct {
	Subscriptions map[string][]chan Event
	Events        []Event
}

// Event represents an event in the event bus
type Event struct {
	Type      string
	Payload   interface{}
	Timestamp time.Time
	Sender    ScreenType
}

// NewIntegrationLayer creates a new integration layer
func NewIntegrationLayer() *IntegrationLayer {
	theme := NewTheme()
	keys := NewKeys()
	a11y := NewA11y(NewAccessibilityConfig(), theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	screenManager := NewScreenManager()
	dataManager := NewDataManager()
	eventBus := NewEventBus()
	
	return &IntegrationLayer{
		Theme:         theme,
		Keys:          keys,
		A11y:          a11y,
		PerfTracker:   perfTracker,
		ScreenManager: screenManager,
		DataManager:   dataManager,
		EventBus:      eventBus,
	}
}

// NewScreenManager creates a new screen manager
func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		CurrentScreen: ScreenDashboard,
		Screens:       make(map[ScreenType]Screen),
		History:       []ScreenType{},
		Transitions:   []ScreenTransition{},
	}
}

// NewDataManager creates a new data manager
func NewDataManager() *DataManager {
	return &DataManager{
		SharedData:  make(map[string]interface{}),
		Cache:       make(map[string]interface{}),
		LastUpdate:  time.Now(),
		AutoRefresh: true,
		RefreshRate: 30 * time.Second,
	}
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		Subscriptions: make(map[string][]chan Event),
		Events:       []Event{},
	}
}

// Initialize initializes the integration layer
func (il *IntegrationLayer) Initialize() error {
	// Initialize screens
	il.ScreenManager.Screens[ScreenDashboard] = NewDashboardScreen(il.Theme, il.Keys, il.A11y, il.PerfTracker)
	il.ScreenManager.Screens[ScreenDoctor] = NewDoctorScreen(il.Theme, il.Keys, il.A11y, il.PerfTracker)
	il.ScreenManager.Screens[ScreenMetrics] = NewMetricsScreen(il.Theme, il.Keys, il.A11y, il.PerfTracker)
	il.ScreenManager.Screens[ScreenCost] = NewCostScreen(il.Theme, il.Keys, il.A11y, il.PerfTracker)
	il.ScreenManager.Screens[ScreenSettings] = NewSettingsScreen(il.Theme, il.Keys, il.A11y, il.PerfTracker)
	
	// Set up auto-refresh
	if il.DataManager.AutoRefresh {
		go il.startAutoRefresh()
	}
	
	// Set up event bus
	go il.eventBusLoop()
	
	return nil
}

// startAutoRefresh starts the auto-refresh mechanism
func (il *IntegrationLayer) startAutoRefresh() {
	ticker := time.NewTicker(il.DataManager.RefreshRate)
	defer ticker.Stop()
	
	for range ticker.C {
		il.refreshAllScreens()
	}
}

// refreshAllScreens refreshes all screens
func (il *IntegrationLayer) refreshAllScreens() {
	for screenType, screen := range il.ScreenManager.Screens {
		// Send refresh event to screen
		event := Event{
			Type:      "refresh",
			Payload:   nil,
			Timestamp: time.Now(),
			Sender:    screenType,
		}
		il.EventBus.Publish(event)
	}
}

// eventBusLoop handles the event bus loop
func (il *IntegrationLayer) eventBusLoop() {
	for {
		// Process events
		for len(il.EventBus.Events) > 0 {
			event := il.EventBus.Events[0]
			il.EventBus.Events = il.EventBus.Events[1:]
			
			// Notify subscribers
			il.notifySubscribers(event)
		}
		
		// Sleep to prevent high CPU usage
		time.Sleep(100 * time.Millisecond)
	}
}

// notifySubscribers notifies event subscribers
func (il *IntegrationLayer) notifySubscribers(event Event) {
	if subscribers, exists := il.EventBus.Subscriptions[event.Type]; exists {
		for _, subscriber := range subscribers {
			select {
			case subscriber <- event:
			default:
				// Channel full, skip
			}
		}
	}
}

// Subscribe subscribes to an event type
func (eb *EventBus) Subscribe(eventType string) chan Event {
	if eb.Subscriptions[eventType] == nil {
		eb.Subscriptions[eventType] = []chan Event{}
	}
	
	channel := make(chan Event, 10)
	eb.Subscriptions[eventType] = append(eb.Subscriptions[eventType], channel)
	
	return channel
}

// Publish publishes an event
func (eb *EventBus) Publish(event Event) {
	eb.Events = append(eb.Events, event)
}

// SwitchScreen switches to a different screen
func (sm *ScreenManager) SwitchScreen(newScreen ScreenType) {
	if sm.CurrentScreen != newScreen {
		// Record transition
		transition := ScreenTransition{
			From:     sm.CurrentScreen,
			To:       newScreen,
			Timestamp: time.Now(),
			Duration: 0, // Will be calculated when transition completes
		}
		
		sm.Transitions = append(sm.Transitions, transition)
		sm.History = append(sm.History, sm.CurrentScreen)
		sm.CurrentScreen = newScreen
	}
}

// GetScreen returns a screen by type
func (sm *ScreenManager) GetScreen(screenType ScreenType) Screen {
	return sm.Screens[screenType]
}

// SetSharedData sets shared data
func (dm *DataManager) SetSharedData(key string, value interface{}) {
	dm.SharedData[key] = value
	dm.Cache[key] = value
	dm.LastUpdate = time.Now()
}

// GetSharedData gets shared data
func (dm *DataManager) GetSharedData(key string) (interface{}, bool) {
	value, exists := dm.SharedData[key]
	return value, exists
}

// CacheData caches data
func (dm *DataManager) CacheData(key string, value interface{}) {
	dm.Cache[key] = value
}

// GetCachedData gets cached data
func (dm *DataManager) GetCachedData(key string) (interface{}, bool) {
	value, exists := dm.Cache[key]
	return value, exists
}

// NewTUI creates a new TUI application
func NewTUI() *TUI {
	integration := NewIntegrationLayer()
	model := NewAppModel()
	
	return &TUI{
		model:       model,
		program:     nil,
		integration: integration,
	}
}

// Initialize initializes the TUI application
func (t *TUI) Initialize() error {
	// Initialize integration layer
	if err := t.integration.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize integration layer: %w", err)
	}
	
	// Set up model with integration
	t.model.Theme = t.integration.Theme
	t.model.Keys = t.integration.Keys
	t.model.A11y = t.integration.A11y
	t.model.PerfTracker = t.integration.PerfTracker
	t.model.Screens = t.integration.ScreenManager.Screens
	t.model.CurrentScreen = t.integration.ScreenManager.CurrentScreen
	t.model.ActiveScreen = t.integration.ScreenManager.Screens[t.integration.ScreenManager.CurrentScreen]
	
	// Create bubble tea program
	t.program = bubbletea.NewProgram(t.model, bubbletea.WithAltScreen(), bubbletea.WithMouseCellMotion())
	
	return nil
}

// Run runs the TUI application
func (t *TUI) Run() error {
	if t.program == nil {
		return fmt.Errorf("TUI not initialized")
	}
	
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Run program in a goroutine
	go func() {
		_, err := t.program.Run()
		if err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			sigChan <- syscall.SIGTERM
		}
	}()
	
	// Wait for signal
	<-sigChan
	
	// Clean up
	t.program.Quit()
	
	return nil
}

// GetScreenManager returns the screen manager
func (t *TUI) GetScreenManager() *ScreenManager {
	return t.integration.ScreenManager
}

// GetDataManager returns the data manager
func (t *TUI) GetDataManager() *DataManager {
	return t.integration.DataManager
}

// GetEventBus returns the event bus
func (t *TUI) GetEventBus() *EventBus {
	return t.integration.EventBus
}

// GetPerformanceMetrics returns performance metrics
func (t *TUI) GetPerformanceMetrics() string {
	return t.integration.PerfTracker.GetReport()
}

// GetAccessibilityReport returns accessibility report
func (t *TUI) GetAccessibilityReport() string {
	return t.integration.A11y.GetAccessibilityReport()
}

// GetScreenTransitionHistory returns screen transition history
func (t *TUI) GetScreenTransitionHistory() []ScreenTransition {
	return t.integration.ScreenManager.Transitions
}

// GetSharedData returns shared data
func (t *TUI) GetSharedData(key string) (interface{}, bool) {
	return t.integration.DataManager.GetSharedData(key)
}

// SetSharedData sets shared data
func (t *TUI) SetSharedData(key string, value interface{}) {
	t.integration.DataManager.SetSharedData(key, value)
}

// SubscribeToEvent subscribes to an event
func (t *TUI) SubscribeToEvent(eventType string) chan Event {
	return t.integration.EventBus.Subscribe(eventType)
}

// PublishEvent publishes an event
func (t *TUI) PublishEvent(eventType string, payload interface{}, sender ScreenType) {
	event := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
		Sender:    sender,
	}
	t.integration.EventBus.Publish(event)
}
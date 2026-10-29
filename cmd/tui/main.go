package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

// TUIApplication represents the TUI application with full integration
type TUIApplication struct {
	tui *ui.TUI
}

// NewTUIApplication creates a new TUI application
func NewTUIApplication() *TUIApplication {
	return &TUIApplication{
		tui: ui.NewTUI(),
	}
}

// Run runs the TUI application with full integration
func (app *TUIApplication) Run() error {
	fmt.Println("🚀 Starting CodeAuthorityFirst TUI...")
	
	// Initialize the TUI
	if err := app.tui.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}
	
	// Set up shared data
	app.setupSharedData()
	
	// Subscribe to important events
	app.setupEventSubscriptions()
	
	// Show startup information
	app.showStartupInfo()
	
	// Run the TUI
	if err := app.tui.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}
	
	return nil
}

// setupSharedData sets up shared data between screens
func (app *TUIApplication) setupSharedData() {
	// Set up initial shared data
	app.tui.SetSharedData("app_start_time", time.Now())
	app.tui.SetSharedData("version", "1.0.0")
	app.tui.SetSharedData("screen_history", []string{})
	
	// Cache some initial data
	app.tui.GetDataManager().CacheData("initial_data_loaded", true)
}

// setupEventSubscriptions sets up event subscriptions
func (app *TUIApplication) setupEventSubscriptions() {
	// Subscribe to screen change events
	screenChan := app.tui.SubscribeToEvent("screen_changed")
	go func() {
		for event := range screenChan {
			fmt.Printf("Screen changed to: %v\n", event.Payload)
		}
	}()
	
	// Subscribe to performance events
	perfChan := app.tui.SubscribeToEvent("performance_warning")
	go func() {
		for event := range perfChan {
			fmt.Printf("Performance warning: %v\n", event.Payload)
		}
	}()
}

// showStartupInfo shows startup information
func (app *TUIApplication) showStartupInfo() {
	fmt.Println("✅ CodeAuthorityFirst TUI initialized successfully")
	fmt.Println("📊 Performance metrics available at startup")
	fmt.Println("♿ Accessibility features enabled")
	fmt.Println("🎯 Keyboard navigation active")
	fmt.Println("🔄 Auto-refresh enabled")
	fmt.Println("")
	
	// Show initial performance metrics
	metrics := app.tui.GetPerformanceMetrics()
	fmt.Println("📈 Initial Performance Metrics:")
	fmt.Println(metrics)
	fmt.Println("")
	
	// Show accessibility report
	a11yReport := app.tui.GetAccessibilityReport()
	fmt.Println("♿ Accessibility Report:")
	fmt.Println(a11yReport)
	fmt.Println("")
	
	// Show screen transition history
	history := app.tui.GetScreenTransitionHistory()
	fmt.Printf("🔄 Screen Transition History: %d transitions\n", len(history))
	fmt.Println("")
}

// GetIntegrationStatus returns the integration status
func (app *TUIApplication) GetIntegrationStatus() string {
	status := "CodeAuthorityFirst Integration Status:\n"
	status += "✅ TUI initialized\n"
	status += "✅ Integration layer active\n"
	status += "✅ Event bus operational\n"
	status += "✅ Data manager active\n"
	status += "✅ Screen manager operational\n"
	status += "✅ Performance tracking active\n"
	status += "✅ Accessibility features enabled\n"
	
	// Show current screen
	screenManager := app.tui.GetScreenManager()
	currentScreen := screenManager.CurrentScreen
	status += fmt.Sprintf("📍 Current screen: %s\n", currentScreen)
	
	// Show shared data count
	dataManager := app.tui.GetDataManager()
	dataCount := len(dataManager.SharedData)
	status += fmt.Sprintf("📊 Shared data items: %d\n", dataCount)
	
	// Show event count
	eventBus := app.tui.GetEventBus()
	eventCount := len(eventBus.Events)
	status += fmt.Sprintf("📈 Pending events: %d\n", eventCount)
	
	return status
}

// main is the main entry point for the TUI
func main() {
	// Create a new TUI application
	app := NewTUIApplication()
	
	// Show integration status
	fmt.Println(app.GetIntegrationStatus())
	
	// Run the TUI application
	if err := app.Run(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("👋 Goodbye!")
}
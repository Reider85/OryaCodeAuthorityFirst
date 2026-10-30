package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// SettingsScreen represents the settings screen
// FR-TUI-005: Settings shows configuration options, accessibility settings, preferences
type SettingsScreen struct {
	*ScreenBase
	Components []Component
	Title      string
	
	// Settings data
	General      GeneralSettings
	Accessibility AccessibilitySettings
	Notifications NotificationSettings
	Advanced     AdvancedSettings
	
	// UI components
	GeneralCard      *Card
	AccessibilityCard *Card
	NotificationsCard *Card
	AdvancedCard     *Card
	SaveButton       *Button
	ResetButton      *Button
	StatusSpinner    *Spinner
	
	// State
	Loading          bool
	LastSave         time.Time
	Error            error
	HasChanges       bool
}

// GeneralSettings represents general settings
type GeneralSettings struct {
	AutoRefresh     bool
	RefreshInterval time.Duration
	Theme           string
	Language        string
	Timezone        string
}

// AccessibilitySettings represents accessibility settings
type AccessibilitySettings struct {
	HighContrast    bool
	ReducedMotion   bool
	ScreenReader    bool
	LargeText      bool
	KeyboardOnly   bool
	TabNavigation  bool
	EscapeExits    bool
	ConfirmActions bool
}

// NotificationSettings represents notification settings
type NotificationSettings struct {
	EnableEmail     bool
	EnableDesktop  bool
	EnableSound    bool
	EmailAddress   string
	AlertThreshold int
}

// AdvancedSettings represents advanced settings
type AdvancedSettings struct {
	DebugMode      bool
	LogLevel      string
	DataRetention  int // days
	AutoBackup    bool
	BackupPath    string
}

// NewSettingsScreen creates a new settings screen
func NewSettingsScreen(theme *Theme, keys *Keys, a11y *A11y, perfTracker *PerfTracker) *SettingsScreen {
	screen := NewScreen(theme, keys, a11y, perfTracker)
	
	return &SettingsScreen{
		ScreenBase:           screen,
		Components:       []Component{},
		Title:           "CodeAuthorityFirst Settings",
		General:         GeneralSettings{},
		Accessibility:  AccessibilitySettings{},
		Notifications:   NotificationSettings{},
		Advanced:        AdvancedSettings{},
		GeneralCard:     NewCard("general", "General", theme, keys),
		AccessibilityCard: NewCard("accessibility", "Accessibility", theme, keys),
		NotificationsCard: NewCard("notifications", "Notifications", theme, keys),
		AdvancedCard:   NewCard("advanced", "Advanced", theme, keys),
		SaveButton:     NewButton("save", "Save Settings", func() tea.Cmd {
			return func() tea.Msg {
				return SaveSettingsMsg{}
			}
		}, theme, keys),
		ResetButton:     NewButton("reset", "Reset to Defaults", func() tea.Cmd {
			return func() tea.Msg {
				return ResetSettingsMsg{}
			}
		}, theme, keys),
		StatusSpinner:   NewSpinner("spinner", theme, keys),
		Loading:         false,
		LastSave:        time.Now(),
		Error:           nil,
		HasChanges:      false,
	}
}

// Init initializes the settings screen
func (s *SettingsScreen) Init() tea.Cmd {
	// Set up screen reader announcement
	s.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Settings screen activated", "Application", false)
	
	// Load initial settings
	return s.loadSettings()
}

// loadSettings loads settings data
func (s *SettingsScreen) loadSettings() tea.Cmd {
	return func() tea.Msg {
		// Simulate loading
		time.Sleep(300 * time.Millisecond)
		
		// Generate mock settings
		general := GeneralSettings{
			AutoRefresh:     true,
			RefreshInterval: 30 * time.Second,
			Theme:           "default",
			Language:        "en",
			Timezone:        "UTC",
		}
		
		accessibility := AccessibilitySettings{
			HighContrast:    false,
			ReducedMotion:   false,
			ScreenReader:   false,
			LargeText:      false,
			KeyboardOnly:   true,
			TabNavigation:  true,
			EscapeExits:    true,
			ConfirmActions: true,
		}
		
		notifications := NotificationSettings{
			EnableEmail:     false,
			EnableDesktop:  true,
			EnableSound:    true,
			EmailAddress:   "",
			AlertThreshold: 80,
		}
		
		advanced := AdvancedSettings{
			DebugMode:      false,
			LogLevel:      "info",
			DataRetention: 30,
			AutoBackup:    true,
			BackupPath:    "~/.config/caf/backups/",
		}
		
		return SettingsLoaded{
			General:       general,
			Accessibility: accessibility,
			Notifications: notifications,
			Advanced:      advanced,
		}
	}
}

// Update updates the settings screen
func (s *SettingsScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.String() {
		case "1":
			// Switch to dashboard
			return s, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDashboard} }
		case "2":
			// Switch to doctor
			return s, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenDoctor} }
		case "3":
			// Switch to metrics
			return s, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenMetrics} }
		case "4":
			// Switch to cost
			return s, func() tea.Msg { return ScreenSwitchMsg{ScreenType: ScreenCost} }
		case "5":
			// Stay on settings
			return s, nil
		case "s":
			// Save settings
			return s, func() tea.Msg { return SaveSettingsMsg{} }
		case "?":
			// Show help
			return s, s.showHelp()
		}
		
		// Update components
		for i, component := range s.Components {
			updatedComponent, cmd := component.Update(msg)
			s.Components[i] = updatedComponent
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case SettingsLoaded:
		// Handle settings loaded
		s.Loading = false
		s.General = msg.General
		s.Accessibility = msg.Accessibility
		s.Notifications = msg.Notifications
		s.Advanced = msg.Advanced
		s.LastSave = time.Now()
		
		// Update cards
		s.updateCards()
		
	case SaveSettingsMsg:
		// Handle save settings
		return s, s.saveSettings()
		
	case ResetSettingsMsg:
		// Handle reset settings
		return s, s.resetSettings()
		
	case SettingsSavedMsg:
		// Handle settings saved
		s.HasChanges = false
		s.LastSave = time.Now()
		s.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Settings saved successfully", "Settings", false)
		
	case SettingsResetMsg:
		// Handle settings reset
		s.HasChanges = false
		s.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Settings reset to defaults", "Settings", false)
		
	case ErrorMsg:
		// Handle error
		s.Error = msg.Err
		s.Loading = false
	}
	
	return s, tea.Batch(cmds...)
}

// updateCards updates the settings cards
func (s *SettingsScreen) updateCards() {
	// Update general card
	generalContent := s.formatGeneral()
	s.GeneralCard.SetContent(generalContent)
	
	// Update accessibility card
	accessibilityContent := s.formatAccessibility()
	s.AccessibilityCard.SetContent(accessibilityContent)
	
	// Update notifications card
	notificationsContent := s.formatNotifications()
	s.NotificationsCard.SetContent(notificationsContent)
	
	// Update advanced card
	advancedContent := s.formatAdvanced()
	s.AdvancedCard.SetContent(advancedContent)
}

// formatGeneral formats general settings
func (s *SettingsScreen) formatGeneral() string {
	return fmt.Sprintf("Auto Refresh: %t\nInterval: %s\nTheme: %s\nLanguage: %s\nTimezone: %s",
		s.General.AutoRefresh, s.formatDuration(s.General.RefreshInterval), s.General.Theme, s.General.Language, s.General.Timezone)
}

// formatAccessibility formats accessibility settings
func (s *SettingsScreen) formatAccessibility() string {
	return fmt.Sprintf("High Contrast: %t\nReduced Motion: %t\nScreen Reader: %t\nLarge Text: %t\nKeyboard Only: %t\nTab Navigation: %t\nEscape Exits: %t\nConfirm Actions: %t",
		s.Accessibility.HighContrast, s.Accessibility.ReducedMotion, s.Accessibility.ScreenReader, s.Accessibility.LargeText,
		s.Accessibility.KeyboardOnly, s.Accessibility.TabNavigation, s.Accessibility.EscapeExits, s.Accessibility.ConfirmActions)
}

// formatNotifications formats notification settings
func (s *SettingsScreen) formatNotifications() string {
	return fmt.Sprintf("Email: %t (%s)\nDesktop: %t\nSound: %t\nAlert Threshold: %d%%",
		s.Notifications.EnableEmail, s.Notifications.EmailAddress, s.Notifications.EnableDesktop, s.Notifications.EnableSound, s.Notifications.AlertThreshold)
}

// formatAdvanced formats advanced settings
func (s *SettingsScreen) formatAdvanced() string {
	return fmt.Sprintf("Debug Mode: %t\nLog Level: %s\nData Retention: %d days\nAuto Backup: %t\nBackup Path: %s",
		s.Advanced.DebugMode, s.Advanced.LogLevel, s.Advanced.DataRetention, s.Advanced.AutoBackup, s.Advanced.BackupPath)
}

// formatDuration formats a duration for display
func (s *SettingsScreen) formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
	return fmt.Sprintf("%dm", int(duration.Minutes()))
}

// saveSettings saves settings
func (s *SettingsScreen) saveSettings() tea.Cmd {
	s.Loading = true
	s.HasChanges = false
	
	return func() tea.Msg {
		// Simulate saving
		time.Sleep(500 * time.Millisecond)
		
		return SettingsSavedMsg{}
	}
}

// resetSettings resets settings to defaults
func (s *SettingsScreen) resetSettings() tea.Cmd {
	s.Loading = true
	s.HasChanges = false
	
	return func() tea.Msg {
		// Simulate reset
		time.Sleep(300 * time.Millisecond)
		
		return SettingsResetMsg{}
	}
}

// showHelp shows help for the settings screen
func (s *SettingsScreen) showHelp() tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{
			ScreenID: "settings",
			Message:  s.getSettingsHelp(),
		}
	}
}

// getSettingsHelp returns settings help text
func (s *SettingsScreen) getSettingsHelp() string {
	return `Settings Controls:
1-5: Switch screens
s: Save settings
?: Show this help
q: Quit

Settings Sections:
General: Auto refresh, theme, language
Accessibility: High contrast, screen reader, keyboard
Notifications: Email, desktop, alerts
Advanced: Debug, logging, backup

Changes are automatically detected and can be saved with s`
}

// View returns the settings view
func (s *SettingsScreen) View() string {
	if s.Loading {
		return s.renderLoading()
	}
	
	if s.Error != nil {
		return s.renderError()
	}
	
	return s.renderContent()
}

// renderLoading renders the loading state
func (s *SettingsScreen) renderLoading() string {
	header := s.Theme.Title.Render(s.Title) + "\n"
	header += s.Theme.Muted.Render("Saving settings...")
	
	content := fmt.Sprintf("%s\n\n%s", header, s.StatusSpinner.View())
	
	return s.renderLayout(content)
}

// renderError renders the error state
func (s *SettingsScreen) renderError() string {
	header := s.Theme.Title.Render(s.Title) + "\n"
	errorBox := s.Theme.ErrorBox.Render(fmt.Sprintf("Error: %s", s.Error.Error()))
	
	content := fmt.Sprintf("%s\n\n%s", header, errorBox)
	
	return s.renderLayout(content)
}

// renderContent renders the main content
func (s *SettingsScreen) renderContent() string {
	header := s.Theme.Title.Render(s.Title) + "\n"
	header += s.Theme.Muted.Render(fmt.Sprintf("Last saved: %s", s.LastSave.Format("15:04:05"))) + "\n"
	
	if s.HasChanges {
		header += s.Theme.Warning.Render("• Unsaved changes") + "\n"
	}
	
	// Create layout with cards
	layout := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		header,
	 s.GeneralCard.View(),
	 s.AccessibilityCard.View(),
	 s.NotificationsCard.View(),
	 s.AdvancedCard.View(),
	 "")
	
	// Add buttons
	layout += fmt.Sprintf("\n%s", s.SaveButton.View())
	layout += fmt.Sprintf("  %s", s.ResetButton.View())
	
	return s.renderLayout(layout)
}

// renderLayout renders the layout with proper styling
func (s *SettingsScreen) renderLayout(content string) string {
	// Add status bar
	statusBar := s.Theme.StatusBar.Render("Settings | Press s to save, ? for help")
	
	return fmt.Sprintf("%s\n%s", content, statusBar)
}

// GetID returns the screen ID
func (s *SettingsScreen) GetID() string {
	return "settings"
}

// IsActive returns true if the screen is active
func (s *SettingsScreen) IsActive() bool {
	return true
}

// SetActive sets the active state
func (s *SettingsScreen) SetActive(active bool) {
	if active {
		s.ScreenBase.ScreenReader.Announce(AnnouncementImportant, "Settings screen activated", "Application", false)
	} else {
		s.ScreenBase.ScreenReader.Announce(AnnouncementNone, "Settings screen deactivated", "Application", false)
	}
}

// SettingsLoaded represents settings loaded message
type SettingsLoaded struct {
	General       GeneralSettings
	Accessibility AccessibilitySettings
	Notifications NotificationSettings
	Advanced      AdvancedSettings
}

// SaveSettingsMsg represents a save settings message
type SaveSettingsMsg struct{}

// ResetSettingsMsg represents a reset settings message
type ResetSettingsMsg struct{}

// SettingsSavedMsg represents settings saved message
type SettingsSavedMsg struct{}

// SettingsResetMsg represents settings reset message
type SettingsResetMsg struct{}